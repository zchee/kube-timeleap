// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package ptrace

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func ptrace(request, pid int, addr, data uintptr) error {
	_, _, errno := unix.RawSyscall6(
		unix.SYS_PTRACE,
		uintptr(request),
		uintptr(pid),
		addr,
		data,
		0, 0)
	return unix.Errno(errno)
}

// peek requests are machine-size oriented, so we wrap it to retrieve arbitrary-length data.
func peek(req, pid int, addr uintptr, out []byte) (count int, err error) {
	// The ptrace syscall differs from glibc's ptrace.
	// Peeks returns the word in *data, not as the return value.
	var buf [unix.SizeofPtr]byte

	// Leading edge.
	// PEEKTEXT/PEEKDATA don't require aligned access (PEEKUSER warns that it might),
	// but if we don't align our reads, we might straddle an unmapped page
	// boundary and not get the bytes leading up to the page boundary.
	count = 0
	if addr%unix.SizeofPtr != 0 {
		err = ptrace(req, pid, addr-addr%unix.SizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if err != nil {
			return 0, err
		}
		count += copy(out, buf[addr%unix.SizeofPtr:])
		out = out[count:]
	}

	// Remainder.
	for len(out) > 0 {
		// We use an internal buffer to guarantee alignment.
		// It's not documented if this is necessary, but we're paranoid.
		err = ptrace(req, pid, addr+uintptr(count), uintptr(unsafe.Pointer(&buf[0])))
		if err != nil {
			return count, err
		}
		copied := copy(out, buf[0:])
		count += copied
		out = out[copied:]
	}

	return count, nil
}

// PeekText reads a word at the address addr in the tracee's memory,
// returning the word as the result of the ptrace call.
func PeekText(pid int, addr uintptr, out []byte) (count int, err error) {
	return peek(unix.PTRACE_PEEKTEXT, pid, addr, out)
}

// PeekData reads a word at the address addr in the tracee's memory,
// returning the word as the result of the ptrace call.
func PeekData(pid int, addr uintptr, out []byte) (count int, err error) {
	return peek(unix.PTRACE_PEEKDATA, pid, addr, out)
}

// PeekUser reads a word at offset addr in the tracee's USER area, which
// holds the registers and other information about the process.
//
// The word is returned as the result of the ptrace call.
func PeekUser(pid int, addr uintptr, out []byte) (count int, err error) {
	return peek(unix.PTRACE_PEEKUSR, pid, addr, out)
}

// poke is as for peek, we need to align our accesses to deal
// with the possibility of straddling an invalid page.
func poke(pokeReq, peekReq, pid int, addr uintptr, data []byte) (count int, err error) {
	// Leading edge.
	count = 0
	if addr%unix.SizeofPtr != 0 {
		var buf [unix.SizeofPtr]byte
		err = ptrace(peekReq, pid, addr-addr%unix.SizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if err != nil {
			return 0, err
		}
		count += copy(buf[addr%unix.SizeofPtr:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr-addr%unix.SizeofPtr, word)
		if err != nil {
			return 0, err
		}
		data = data[count:]
	}

	// Interior.
	for len(data) > unix.SizeofPtr {
		word := *((*uintptr)(unsafe.Pointer(&data[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(count), word)
		if err != nil {
			return count, err
		}
		count += unix.SizeofPtr
		data = data[unix.SizeofPtr:]
	}

	// Trailing edge.
	if len(data) > 0 {
		var buf [unix.SizeofPtr]byte
		err = ptrace(peekReq, pid, addr+uintptr(count), uintptr(unsafe.Pointer(&buf[0])))
		if err != nil {
			return count, err
		}
		copy(buf[0:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(count), word)
		if err != nil {
			return count, err
		}
		count += len(data)
	}

	return count, nil
}

// PokeText copies the word data to the address addr in the tracee's memory.
func PokeText(pid int, addr uintptr, data []byte) (count int, err error) {
	return poke(unix.PTRACE_POKETEXT, unix.PTRACE_PEEKTEXT, pid, addr, data)
}

// PokeData copies the word data to the address addr in the tracee's memory.
func PokeData(pid int, addr uintptr, data []byte) (count int, err error) {
	return poke(unix.PTRACE_POKEDATA, unix.PTRACE_PEEKDATA, pid, addr, data)
}

// PokeUser copies the word data to offset addr in the tracee's USER area.
func PokeUser(pid int, addr uintptr, data []byte) (count int, err error) {
	return poke(unix.PTRACE_POKEUSR, unix.PTRACE_PEEKUSR, pid, addr, data)
}

// GetRegs copies the tracee's general-purpose registers, respectively, to the address data in the tracer.
func GetRegs(pid int, regsout *unix.PtraceRegs) (err error) {
	return ptrace(unix.PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}

// GetFPRegs copies the tracee's floating-point registers, respectively, to the address data in the tracer.
func GetFPRegs(pid int, regsout *unix.PtraceRegs) (err error) {
	return ptrace(unix.PTRACE_GETFPREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
}

// SetRegs modifies the tracee's general-purpose registers, respectively, from the address data in the tracer.
func SetRegs(pid int, regs *unix.PtraceRegs) (err error) {
	return ptrace(unix.PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}

// SetFPRegs modifies the tracee's floating-point registers, respectively, from the address data in the tracer.
func SetFPRegs(pid int, regs *unix.PtraceRegs) (err error) {
	return ptrace(unix.PTRACE_SETFPREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
}

// SetOptions sets ptrace options from data. data is interpreted as a bit mask of options.
func SetOptions(pid, options int) (err error) {
	return ptrace(unix.PTRACE_SETOPTIONS, pid, 0, uintptr(options))
}

// GetEventMsg retrieves a message about the ptrace event that just happened, placing it at the address data in the tracer.
func GetEventMsg(pid int) (msg uintptr, err error) {
	err = ptrace(unix.PTRACE_GETEVENTMSG, pid, 0, uintptr(unsafe.Pointer(&msg)))
	return
}

// Cont restarts the stopped tracee process.
func Cont(pid, signal int) (err error) {
	return ptrace(unix.PTRACE_CONT, pid, 0, uintptr(signal))
}

// Syscall restarts the stopped tracee as for PTRACE_CONT, but arrange for the tracee to be stopped
// at the next entry to or exit from a system call, or after execution of a single instruction, respectively.
func Syscall(pid, signal int) (err error) {
	return ptrace(unix.PTRACE_SYSCALL, pid, 0, uintptr(signal))
}

// SingleStep restarts the stopped tracee as for PTRACE_CONT, but arrange for the tracee to be stopped
// at the next entry to or exit from a system call, or after execution of a single instruction, respectively.
func SingleStep(pid int) (err error) {
	return ptrace(unix.PTRACE_SINGLESTEP, pid, 0, 0)
}

// Interrupt stops a tracee.
func Interrupt(pid int) (err error) {
	return ptrace(unix.PTRACE_INTERRUPT, pid, 0, 0)
}

// Attach attachs to the process specified in pid, making it a tracee of the calling process.
func Attach(pid int) (err error) {
	return ptrace(unix.PTRACE_ATTACH, pid, 0, 0)
}

// Seize attachs to the process specified in pid, making it a tracee of the calling process.
func Seize(pid int) (err error) {
	return ptrace(unix.PTRACE_SEIZE, pid, 0, 0)
}

// Detach restarts the stopped tracee as for PTRACE_CONT, but first detach from it.
func Detach(pid, sig int) (err error) {
	return ptrace(unix.PTRACE_DETACH, pid, 0, uintptr(sig))
}

// GetRegset returns floating point registers of the specified thread using PTRACE.
//
// See amd64_linux_fetch_inferior_registers in gdb/amd64-linux-nat.c.html
// and amd64_supply_xsave in gdb/amd64-tdep.c.html
// and Section 13.1 (and following) of Intel® 64 and IA-32 Architectures Software Developer’s Manual, Volume 1: Basic Architecture.
func GetRegset(tid int) (regset Xstate, errno error) {
	_, _, errno = unix.RawSyscall6(
		unix.SYS_PTRACE,
		unix.PTRACE_GETFPREGS,
		uintptr(tid),
		uintptr(0),
		uintptr(unsafe.Pointer(&regset.FPRegs)),
		0, 0)
	if errno == unix.Errno(0) || errno == unix.ENODEV {
		// ignore ENODEV, it just means this CPU doesn't have X87 registers
		errno = nil
	}

	var xstateArgs [X86_XSTATE_MAX_SIZE]byte
	iovec := unix.Iovec{Base: &xstateArgs[0], Len: X86_XSTATE_MAX_SIZE}
	_, _, errno = unix.RawSyscall6(
		unix.SYS_PTRACE,
		unix.PTRACE_GETREGSET,
		uintptr(tid),
		uintptr(NT_X86_XSTATE),
		uintptr(unsafe.Pointer(&iovec)),
		0, 0)
	if errno != unix.Errno(0) {
		if errno == unix.ENODEV || errno == unix.EIO {
			// ignore ENODEV, it just means this CPU or kernel doesn't support XSTATE.
			// Also ignore EIO, it means that we are running on an old kernel (pre 2.6.34) and PTRACE_GETREGSET is not implemented
			errno = nil
		}
		return
	} else {
		errno = nil
	}

	regset.Xsave = xstateArgs[:iovec.Len]
	errno = ReadXstate(regset.Xsave, false, &regset)
	return
}

// ProcessVMReadv transfers data from the remote tid process to the local process.
func ProcessVMReadv(pid int, addr *uintptr, data []byte) (int, error) {
	sz := len(data)

	localIov := []unix.Iovec{
		{Base: &data[0], Len: uint64(sz)},
	}
	remoteIov := []unix.RemoteIovec{
		{Base: uintptr(unsafe.Pointer(addr)), Len: sz},
	}

	// The flags argument is currently unused and must be set to 0.
	// See also: https://man7.org/linux/man-pages/man2/process_vm_readv.2.html
	flags := uint(0)
	n, err := unix.ProcessVMReadv(pid, localIov, remoteIov, flags)
	if err != unix.Errno(0) {
		return 0, err
	}

	return n, nil
}

// ProcessVMWritev transfers data from the local process to the remote pid process.
func ProcessVMWritev(pid int, addr *uintptr, data []byte) (int, error) {
	sz := len(data)

	localIov := []unix.Iovec{
		{Base: &data[0], Len: uint64(sz)},
	}
	remoteIov := []unix.RemoteIovec{
		{Base: uintptr(unsafe.Pointer(addr)), Len: sz},
	}

	// The flags argument is currently unused and must be set to 0.
	// See also: https://man7.org/linux/man-pages/man2/process_vm_writev.2.html
	flags := uint(0)
	n, err := unix.ProcessVMWritev(pid, localIov, remoteIov, flags)
	if err != unix.Errno(0) {
		return 0, err
	}

	return n, nil
}
