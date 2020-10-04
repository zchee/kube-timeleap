// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

// +build !windows

package signalctx

import (
	"os"

	"golang.org/x/sys/unix"
)

var shutdownSignals = []os.Signal{os.Interrupt, unix.SIGTERM}
