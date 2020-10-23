// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

// +build linux

package ptrace

import (
	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("ptrace")

// RegisterLogger registers a logger on ptrace pkg.
func RegisterLogger(logger logr.Logger) {
	log = logger
}
