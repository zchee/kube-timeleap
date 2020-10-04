// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package ptrace

import (
	"debug/elf"
	"errors"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Thread is a traced thread; it is a thread identifier.
//
// This is a convenience type for defining ptrace operations.
type Thread struct {
	tgid int32
	tid  int32
	cpu  uint32

	// initRegs are the initial registers for the first thread.
	//
	// These are used for the register set for system calls.
	initRegs unix.PtraceRegs
}

// GetRegs gets the general purpose register set.
func (t *Thread) GetRegs(regs *unix.PtraceRegs) error {
	iovec := unix.Iovec{
		Base: (*byte)(unsafe.Pointer(regs)),
		Len:  uint64(unsafe.Sizeof(*regs)),
	}

	_, _, errno := unix.RawSyscall6(
		unix.SYS_PTRACE,
		unix.PTRACE_GETREGSET,
		uintptr(t.tid),
		uintptr(elf.NT_PRSTATUS),
		uintptr(unsafe.Pointer(&iovec)),
		0, 0)
	if errno != unix.Errno(0) {
		return unix.Errno(errno)
	}

	return nil
}

// DumpRegs dumps regs.
func DumpRegs(regs *unix.PtraceRegs) string {
	var m strings.Builder

	fmt.Fprintf(&m, "Registers:\n")
	fmt.Fprintf(&m, "\tR15\t = %016x\n", regs.R15)
	fmt.Fprintf(&m, "\tR14\t = %016x\n", regs.R14)
	fmt.Fprintf(&m, "\tR13\t = %016x\n", regs.R13)
	fmt.Fprintf(&m, "\tR12\t = %016x\n", regs.R12)
	fmt.Fprintf(&m, "\tRbp\t = %016x\n", regs.Rbp)
	fmt.Fprintf(&m, "\tRbx\t = %016x\n", regs.Rbx)
	fmt.Fprintf(&m, "\tR11\t = %016x\n", regs.R11)
	fmt.Fprintf(&m, "\tR10\t = %016x\n", regs.R10)
	fmt.Fprintf(&m, "\tR9\t = %016x\n", regs.R9)
	fmt.Fprintf(&m, "\tR8\t = %016x\n", regs.R8)
	fmt.Fprintf(&m, "\tRax\t = %016x\n", regs.Rax)
	fmt.Fprintf(&m, "\tRcx\t = %016x\n", regs.Rcx)
	fmt.Fprintf(&m, "\tRdx\t = %016x\n", regs.Rdx)
	fmt.Fprintf(&m, "\tRsi\t = %016x\n", regs.Rsi)
	fmt.Fprintf(&m, "\tRdi\t = %016x\n", regs.Rdi)
	fmt.Fprintf(&m, "\tOrig_rax = %016x\n", regs.Orig_rax)
	fmt.Fprintf(&m, "\tRip\t = %016x\n", regs.Rip)
	fmt.Fprintf(&m, "\tCs\t = %016x\n", regs.Cs)
	fmt.Fprintf(&m, "\tEflags\t = %016x\n", regs.Eflags)
	fmt.Fprintf(&m, "\tRsp\t = %016x\n", regs.Rsp)
	fmt.Fprintf(&m, "\tSs\t = %016x\n", regs.Ss)
	fmt.Fprintf(&m, "\tFs_base\t = %016x\n", regs.Fs_base)
	fmt.Fprintf(&m, "\tGs_base\t = %016x\n", regs.Gs_base)
	fmt.Fprintf(&m, "\tDs\t = %016x\n", regs.Ds)
	fmt.Fprintf(&m, "\tEs\t = %016x\n", regs.Es)
	fmt.Fprintf(&m, "\tFs\t = %016x\n", regs.Fs)
	fmt.Fprintf(&m, "\tGs\t = %016x\n", regs.Gs)

	return m.String()
}

const (
	// maximumUserAddress is the largest possible user address.
	maximumUserAddress = 0x7ffffffff000

	// stubInitAddress is the initial attempt link address for the stub.
	stubInitAddress = 0x7fffffff0000

	// initRegsRipAdjustment is the size of the syscall instruction.
	initRegsRipAdjustment = 2
)

var (
	// StubStart is the link address for our stub, and determines the
	// maximum user address. This is valid only after a call to stubInit.
	//
	// We attempt to link the stub here, and adjust downward as needed.
	StubStart uintptr = stubInitAddress

	// StubEnd is the first byte past the end of the stub, as with
	// stubStart this is valid only after a call to stubInit.
	StubEnd uintptr

	// stubInitialized controls one-time stub initialization.
	stubInitialized sync.Once
)

// WaitOutcome is used for wait below.
type WaitOutcome int

const (
	// Stopped indicates that the process was Stopped.
	Stopped WaitOutcome = iota

	// Killed indicates that the process was Killed.
	Killed
)

func (t *Thread) dumpAndPanic(message string) {
	var regs unix.PtraceRegs
	message += "\n"
	if err := t.GetRegs(&regs); err == nil {
		message += DumpRegs(&regs)
	} else {
		// log.Warningf("unable to get registers: %v", err)
	}
	message += fmt.Sprintf("stubStart\t = %016x\n", StubStart)

	panic(message)
}

func (t *Thread) unexpectedStubExit() {
	msg, err := t.GetEventMessage()
	status := unix.WaitStatus(msg)
	if status.Signaled() && status.Signal() == unix.SIGKILL {
		// SIGKILL can be only sent by a user or OOM-killer. In both
		// these cases, we don't need to panic. There is no reasons to
		// think that something wrong in gVisor.
		// log.Warningf("The ptrace stub process %v has been killed by SIGKILL.", t.tgid)
		pid := unix.Getpid()
		unix.Tgkill(pid, pid, unix.Signal(unix.SIGKILL))
	}

	t.dumpAndPanic(fmt.Sprintf("wait failed: the process %d:%d exited: %x (err %v)", t.tgid, t.tid, msg, err))
}

// GetEventMessage retrieves a message about the ptrace event that just happened.
func (t *Thread) GetEventMessage() (uintptr, error) {
	var msg uintptr
	_, _, errno := unix.RawSyscall6(
		unix.SYS_PTRACE,
		unix.PTRACE_GETEVENTMSG,
		uintptr(t.tid),
		0,
		uintptr(unsafe.Pointer(&msg)),
		0, 0)
	if errno != unix.Errno(0) {
		return msg, unix.Errno(errno)
	}
	return msg, nil
}

// Wait waits for a stop event.
func (t *Thread) Wait(outcome WaitOutcome) unix.Signal {
	var status unix.WaitStatus

	for {
		r, err := unix.Wait4(int(t.tid), &status, unix.WALL|unix.WUNTRACED, nil)

		switch {
		case errors.Is(err, unix.EINTR), errors.Is(err, unix.EAGAIN):
			// Wait was interrupted; wait again.
			continue
		case err != nil:
			panic(fmt.Sprintf("ptrace wait failed: %v", err))
		}

		if int(r) != int(t.tid) {
			panic(fmt.Sprintf("ptrace wait returned %v, expected %v", r, t.tid))
		}

		switch outcome {
		case Stopped:
			if !status.Stopped() {
				t.dumpAndPanic(fmt.Sprintf("ptrace status unexpected: got %v, wanted stopped", status))
			}
			stopSig := status.StopSignal()
			if stopSig == 0 {
				continue // Spurious stop.
			}
			if stopSig == unix.SIGTRAP {
				if status.TrapCause() == unix.PTRACE_EVENT_EXIT {
					t.unexpectedStubExit()
				}
				// Re-encode the trap cause the way it's expected.
				return stopSig | unix.Signal(status.TrapCause()<<8)
			}
			// Not a trap signal.
			return stopSig

		case Killed:
			if !status.Exited() && !status.Signaled() {
				t.dumpAndPanic(fmt.Sprintf("ptrace status unexpected: got %v, wanted exited", status))
			}
			return unix.Signal(status.ExitStatus())

		default:
			// Should not happen.
			t.dumpAndPanic(fmt.Sprintf("unknown outcome: %v", outcome))
		}
	}
}
