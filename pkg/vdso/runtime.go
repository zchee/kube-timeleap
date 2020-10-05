// +build linux

package vdso

import (
	_ "runtime" // for go:linkname
	"unsafe"
)

//go:linkname gostringnocopy runtime.gostringnocopy
//go:nosplit
func gostringnocopy(str *byte) string

//go:linkname add runtime.add
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname physPageSize runtime.physPageSize

// physPageSize is the size in bytes of the OS's physical pages.
// Mapping and unmapping operations must be done at multiples of
// physPageSize.
//
// This must be set by the OS init code (typically in osinit) before
// mallocinit.
var physPageSize uintptr

//go:linkname noescape runtime.noescape
//go:nosplit

// noescape hides a pointer from escape analysis.
//
// noescape is the identity function but escape analysis doesn't think the
// output depends on the input.
//
// noescape is inlined and currently compiles down to zero instructions.
func noescape(p unsafe.Pointer) unsafe.Pointer
