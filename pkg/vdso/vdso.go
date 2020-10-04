// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package vdso

import (
	"unsafe"

	"github.com/zchee/kube-timeleap/pkg/vdso/elfdef"
)

// Look up symbols in the Linux vDSO.

const (
	// ArrayMax is the byte-size of a maximally sized array on this architecture.
	ArrayMax = 1<<50 - 1

	// Maximum indices for the array types used when traversing the vDSO ELF structures.
	SymTabSize     = ArrayMax / unsafe.Sizeof(elfdef.Sym{})
	DynSize        = ArrayMax / unsafe.Sizeof(elfdef.Dyn{})
	SymStringsSize = ArrayMax     // byte
	VerSymSize     = ArrayMax / 2 // uint16
	HashSize       = ArrayMax / 4 // uint32

	// BloomSizeScale is a scaling factor for gnuhash tables which are uint32 indexed, but contain uintptrs.
	BloomSizeScale = unsafe.Sizeof(uintptr(0)) / 4 // uint32
)

// SymbolKey represents a vDSO symbol key entries.
type SymbolKey struct {
	Name    string
	SymHash uint32
	GnuHash uint32
	Ptr     *uintptr
}

// VersionKey represents a vDSO Version key entries.
type VersionKey struct {
	Version string
	VerHash uint32
}

// LinuxVersion version of Linux Kernel in vDSO object.
var LinuxVersion = VersionKey{"LINUX_2.6", 0x3ae75f6}

// initialize with vsyscall fallbacks.
var (
	GettimeofdaySym uintptr = 0xffffffffff600000
	ClockgettimeSym uintptr = 0
)

// SymbolKeys is the Linux amd64 vDSO symbol keys.
var SymbolKeys = []SymbolKey{
	{"__vdso_gettimeofday", 0x315ca59, 0xb01bca00, &GettimeofdaySym},
	{"__vdso_clock_gettime", 0xd35ec75, 0x6e43a318, &ClockgettimeSym},
}

// VDSO represents a vDSO object.
type VDSO struct {
	Valid bool

	// Load information
	LoadAddr   uintptr
	LoadOffset uintptr // LoadAddr - recorded vaddr

	// Symbol table
	Symtab     *[SymTabSize]elfdef.Sym
	Symstrings *[SymStringsSize]byte
	Chain      []uint32
	Bucket     []uint32
	SymOff     uint32
	IsGNUHash  bool

	// Version table
	Versym *[VerSymSize]uint16
	Verdef *elfdef.Verdef
}

// InitFromSysinfoELFHeader init vDSO from hdr.
//go:nocheckptr
func InitFromSysinfoELFHeader(vdso *VDSO, hdr *elfdef.Header) {
	vdso.Valid = false
	vdso.LoadAddr = uintptr(unsafe.Pointer(hdr))

	pt := unsafe.Pointer(vdso.LoadAddr + uintptr(hdr.Phoff))

	// We need two things from the segment table: the load offset
	// and the dynamic table.
	var foundVaddr bool
	var dyn *[DynSize]elfdef.Dyn
	for i := uint16(0); i < hdr.Phnum; i++ {
		pt := (*elfdef.Prog)(add(pt, uintptr(i)*unsafe.Sizeof(elfdef.Prog{})))
		switch pt.Type {
		case uint32(elfdef.PT_LOAD):
			if !foundVaddr {
				foundVaddr = true
				vdso.LoadOffset = vdso.LoadAddr + uintptr(pt.Offset-pt.Vaddr)
			}

		case uint32(elfdef.PT_DYNAMIC):
			dyn = (*[DynSize]elfdef.Dyn)(unsafe.Pointer(vdso.LoadAddr + uintptr(pt.Offset)))
		}
	}

	if !foundVaddr || dyn == nil {
		return // Failed
	}

	// Fish out the useful bits of the dynamic table.
	var hash, gnuhash *[HashSize]uint32
	vdso.Symstrings = nil
	vdso.Symtab = nil
	vdso.Versym = nil
	vdso.Verdef = nil
	for i := 0; elfdef.DynTag(dyn[i].Tag) != elfdef.DT_NULL; i++ {
		dt := &dyn[i]
		p := vdso.LoadOffset + uintptr(dt.Val)
		switch elfdef.DynTag(dt.Tag) {
		case elfdef.DT_STRTAB:
			vdso.Symstrings = (*[SymStringsSize]byte)(unsafe.Pointer(p))
		case elfdef.DT_SYMTAB:
			vdso.Symtab = (*[SymTabSize]elfdef.Sym)(unsafe.Pointer(p))
		case elfdef.DT_HASH:
			hash = (*[HashSize]uint32)(unsafe.Pointer(p))
		case elfdef.DT_GNU_HASH:
			gnuhash = (*[HashSize]uint32)(unsafe.Pointer(p))
		case elfdef.DT_VERSYM:
			vdso.Versym = (*[VerSymSize]uint16)(unsafe.Pointer(p))
		case elfdef.DT_VERDEF:
			vdso.Verdef = (*elfdef.Verdef)(unsafe.Pointer(p))
		}
	}

	if vdso.Symstrings == nil || vdso.Symtab == nil || (hash == nil && gnuhash == nil) {
		return // Failed
	}

	if vdso.Verdef == nil {
		vdso.Versym = nil
	}

	if gnuhash != nil {
		// Parse the GNU hash table header.
		nbucket := gnuhash[0]
		vdso.SymOff = gnuhash[1]
		bloomSize := gnuhash[2]
		vdso.Bucket = gnuhash[4+bloomSize*uint32(BloomSizeScale):][:nbucket]
		vdso.Chain = gnuhash[4+bloomSize*uint32(BloomSizeScale)+nbucket:]
		vdso.IsGNUHash = true
	} else {
		// Parse the hash table header.
		nbucket := hash[0]
		nchain := hash[1]
		vdso.Bucket = hash[2 : 2+nbucket]
		vdso.Chain = hash[2+nbucket : 2+nbucket+nchain]
	}

	// That's all we need.
	vdso.Valid = true
}

// FindVersion finds version from vDSO.
func FindVersion(vdso *VDSO, ver *VersionKey) int32 {
	if !vdso.Valid {
		return 0
	}

	def := vdso.Verdef
	for {
		if def.Flags&elfdef.VER_FLG_BASE == 0 {
			aux := (*elfdef.Verdaux)(add(unsafe.Pointer(def), uintptr(def.Aux)))
			if def.Hash == ver.VerHash && ver.Version == gostringnocopy(&vdso.Symstrings[aux.Name]) {
				return int32(def.Ndx & 0x7fff)
			}
		}

		if def.Next == 0 {
			break
		}
		def = (*elfdef.Verdef)(add(unsafe.Pointer(def), uintptr(def.Next)))
	}

	return -1 // cannot match any version
}

// ParseSymbols parses and allocates symbols from vDSO.
func ParseSymbols(vdso *VDSO, version int32) {
	if !vdso.Valid {
		return
	}

	apply := func(symIndex uint32, k SymbolKey) bool {
		sym := &vdso.Symtab[symIndex]
		typ := elfdef.ST_TYPE(sym.Info)
		bind := elfdef.ST_BIND(sym.Info)

		// On ppc64x, VDSO functions are of type _STT_NOTYPE.
		if typ != elfdef.STT_FUNC && typ != elfdef.STT_NOTYPE || bind != elfdef.STB_GLOBAL && bind != elfdef.STB_WEAK || sym.Shndx == uint16(elfdef.SHN_UNDEF) {
			return false
		}
		if k.Name != gostringnocopy(&vdso.Symstrings[sym.Name]) {
			return false
		}
		// Check symbol version.
		if vdso.Versym != nil && version != 0 && int32(vdso.Versym[symIndex]&0x7fff) != version {
			return false
		}

		*k.Ptr = vdso.LoadOffset + uintptr(sym.Value)
		return true
	}

	if vdso.IsGNUHash {
		// New-style DT_GNU_HASH table.
		for _, k := range SymbolKeys {
			symIndex := vdso.Bucket[k.GnuHash%uint32(len(vdso.Bucket))]
			if symIndex < vdso.SymOff {
				continue
			}
			for ; ; symIndex++ {
				hash := vdso.Chain[symIndex-vdso.SymOff]
				if hash|1 == k.GnuHash|1 {
					// Found a hash match.
					if apply(symIndex, k) {
						break
					}
				}
				if hash&1 != 0 {
					// End of chain.
					break
				}
			}
		}
		return
	}

	// Old-style DT_HASH table.
	for _, k := range SymbolKeys {
		for chain := vdso.Bucket[k.SymHash%uint32(len(vdso.Bucket))]; chain != 0; chain = vdso.Chain[chain] {
			if apply(chain, k) {
				break
			}
		}
	}
	return
}

// Auxv parses auxv by tag and allocate elfdef.Header to val.
//go:nocheckptr
func Auxv(tag, val uintptr) {
	switch tag {
	case elfdef.AT_SYSINFO_EHDR:
		if val == 0 {
			// Something went wrong
			return
		}
		var info VDSO
		info1 := (*VDSO)(noescape(unsafe.Pointer(&info)))
		InitFromSysinfoELFHeader(info1, (*elfdef.Header)(unsafe.Pointer(val)))
		ParseSymbols(info1, FindVersion(info1, &LinuxVersion))
	}
}

//go:nosplit

// InVDSOPage reports whether the pc is on the VDSO page.
func InVDSOPage(pc uintptr) bool {
	for _, k := range SymbolKeys {
		if *k.Ptr != 0 {
			page := *k.Ptr &^ (physPageSize - 1)
			return pc >= page && pc < page+physPageSize
		}
	}
	return false
}
