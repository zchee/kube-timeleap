// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package ptrace

import (
	"bytes"
	"encoding/binary"
)

const (
	X86_XSTATE_MAX_SIZE = 2688
	NT_X86_XSTATE       = 0x202

	XSAVE_HDR_OFFSET             = 512
	XSAVE_HDR_SIZE               = 64
	XSAVE_EXTENDED_REGION_OFFSET = 576
	XSAVE_SSE_REGION_LEN         = 416
)

// FPRegs represents a user_fpregs_struct in /usr/include/x86_64-linux-gnu/sys/user.h.
type FPRegs struct {
	Cwd      uint16     // Control Word
	Swd      uint16     // Status Word
	Ftw      uint16     // Tag Word
	Fop      uint16     // Last Instruction Opcode
	Rip      uint64     // Instruction Pointer
	Rdp      uint64     // Data Pointer
	Mxcsr    uint32     // MXCSR Register State
	MxcrMask uint32     // MXCR Mask
	StSpace  [32]uint32 // 8*16 bytes for each FP-reg = 128 bytes
	XMMSpace [256]byte  // 16*16 bytes for each XMM-reg = 256 bytes
	_        [24]uint32 // padding
}

// Xstate represents amd64 XSAVE area.
//
// See Section 13.1 (and following) of Intel® 64 and IA-32 Architectures Software Developer’s Manual, Volume 1: Basic Architecture.
type Xstate struct {
	FPRegs
	Xsave    []byte    // raw xsave area
	AVXState bool      // contains AVX state
	YMMSpace [256]byte // YMM register space
}

// ReadXstate reads a byte array containing an XSAVE area into register set.
//
// If readLegacy is true regset.PtraceFpRegs will be filled with the
// contents of the legacy region of the XSAVE area.
// See Section 13.1 (and following) of Intel® 64 and IA-32 Architectures
// Software Developer’s Manual, Volume 1: Basic Architecture.
func ReadXstate(xstateArgs []byte, readLegacy bool, regset *Xstate) error {
	if XSAVE_HDR_OFFSET+XSAVE_HDR_SIZE >= len(xstateArgs) {
		return nil
	}
	if readLegacy {
		rdr := bytes.NewReader(xstateArgs[:XSAVE_HDR_OFFSET])
		if err := binary.Read(rdr, binary.LittleEndian, &regset.FPRegs); err != nil {
			return err
		}
	}
	xsaveHdr := xstateArgs[XSAVE_HDR_OFFSET : XSAVE_HDR_OFFSET+XSAVE_HDR_SIZE]
	xstate_bv := binary.LittleEndian.Uint64(xsaveHdr[0:8])
	xcomp_bv := binary.LittleEndian.Uint64(xsaveHdr[8:16])

	if xcomp_bv&(1<<63) != 0 {
		// compact format not supported
		return nil
	}

	if xstate_bv&(1<<2) == 0 {
		// AVX state not present
		return nil
	}

	avxState := xstateArgs[XSAVE_EXTENDED_REGION_OFFSET:]
	regset.AVXState = true
	copy(regset.YMMSpace[:], avxState[:len(regset.YMMSpace)])

	return nil
}
