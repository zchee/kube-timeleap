// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

// Package elfdef prodives the ELF64 structure definitions for use by the vDSO loader.
package elfdef

import (
	"debug/elf"
)

// DynTag type alias of elf.DynTag.
type DynTag = elf.DynTag

const (
	AT_SYSINFO_EHDR = 33

	// Indexes into the Header.Ident array.
	EI_NIDENT = 16 // Size of Header.Ident array.

	// Version is found in Header.Ident[EI_VERSION] and Header.Version.
	VER_FLG_BASE = 0x1 // Version definition of file itself

	// Special section indices.
	SHN_UNDEF = elf.SectionIndex(0) // Undefined section

	// Section type.
	SHT_DYNSYM = elf.SectionType(11) // Dynamic linker symbol table

	// Prog.Type
	PT_LOAD    = elf.ProgType(1) // Loadable program segment
	PT_DYNAMIC = elf.ProgType(2) // Dynamic linking information

	// Dyn.Tag
	DT_NULL     = DynTag(0)          // Marks end of dynamic section
	DT_HASH     = DynTag(4)          // Dynamic symbol hash table
	DT_STRTAB   = DynTag(5)          // Address of string table
	DT_SYMTAB   = DynTag(6)          // Address of symbol table
	DT_GNU_HASH = DynTag(0x6ffffef5) // GNU-style dynamic symbol hash table
	DT_VERSYM   = DynTag(0x6ffffff0) // Address of version symbols table
	DT_VERDEF   = DynTag(0x6ffffffc) // Address of version dependencies table

	// Symbol Binding
	STB_GLOBAL = elf.SymBind(1) // Global symbol
	STB_WEAK   = elf.SymBind(2) // Weak symbol

	// Symbol type
	STT_FUNC   = elf.SymType(2) // Symbol is a code object
	STT_NOTYPE = elf.SymType(0) // Symbol type is not specified
)

// Header represents a ELF64 file header.
type Header struct {
	Ident     [EI_NIDENT]byte // Magic number and other info
	Type      uint16          // Object file type
	Machine   uint16          // Architecture
	Version   uint32          // Object file version
	Entry     uint64          // Entry point virtual address
	Phoff     uint64          // Program header table file offset
	Shoff     uint64          // Section header table file offset
	Flags     uint32          // Processor-specific flags
	Ehsize    uint16          // ELF header size in bytes
	Phentsize uint16          // Program header table entry size
	Phnum     uint16          // Program header table entry count
	Shentsize uint16          // Section header table entry size
	Shnum     uint16          // Section header table entry count
	Shstrndx  uint16          // Section header string table index
}

// Section represents a ELF64 Section header.
type Section struct {
	Name      uint32 // Section name (string tbl index)
	Type      uint32 // Section type
	Flags     uint64 // Section flags
	Addr      uint64 // Section virtual addr at execution
	Offset    uint64 // Section file offset
	Size      uint64 // Section size in bytes
	Link      uint32 // Link to another section
	Info      uint32 // Additional section information
	Addralign uint64 // Section alignment
	Entsize   uint64 // Entry size if section holds table
}

// Prog represents a ELF64 Program header.
type Prog struct {
	Type   uint32 // Segment type
	Flags  uint32 // Segment flags
	Offset uint64 // Segment file offset
	Vaddr  uint64 // Segment virtual address
	Paddr  uint64 // Segment physical address
	Filesz uint64 // Segment size in file
	Memsz  uint64 // Segment size in memory
	Align  uint64 // Segment alignment
}

// Dyn represents a ELF64 Dynamic structure.
//
// The ".dynamic" section contains an array of them.
type Dyn struct {
	Tag int64  // Dynamic entry type
	Val uint64 // Integer value
}

// Sym represents a ELF64 symbol table entries.
type Sym struct {
	Name  uint32 // String table index of name
	Info  byte   // Type and binding information
	Other byte   // Reserved (not used)
	Shndx uint16 // Section index of symbol
	Value uint64 // Symbol value
	Size  uint64 // Size of associated object
}

// Verdef represents a ELF64 Version table entries.
type Verdef struct {
	Version uint16 // Version revision
	Flags   uint16 // Version information
	Ndx     uint16 // Version Index
	Cnt     uint16 // Number of associated aux entries
	Hash    uint32 // Version name hash value
	Aux     uint32 // Offset in bytes to verdaux array
	Next    uint32 // Offset in bytes to next verdef entry
}

// Verdaux represents a ELF64 Version dependency table entries.
type Verdaux struct {
	Name uint32 // Version or dependency names
	Next uint32 // Offset in bytes to next verdaux entry
}

// How to extract and insert information held in the st_info field.
func ST_BIND(info uint8) elf.SymBind { return elf.SymBind(info >> 4) }
func ST_TYPE(info uint8) elf.SymType { return elf.SymType(info & 0xF) }
