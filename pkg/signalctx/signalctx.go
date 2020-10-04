// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package signalctx

import (
	"context"
	"errors"
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})

// setupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

// NewContext creates a new context with SetupSignalHandler() as our Done() channel.
func NewContext() context.Context {
	return &signalContext{
		stopCh: setupSignalHandler(),
	}
}

// signalContext represents a signal handling with implements context.Context interface.
type signalContext struct {
	context.Context // embedded for Deadline and Value

	stopCh <-chan struct{}
}

// Done implements context.Context.
func (sgc *signalContext) Done() <-chan struct{} { return sgc.stopCh }

// Err implements context.Context.
func (sgc *signalContext) Err() error {
	select {
	case _, ok := <-sgc.Done():
		if !ok {
			return errors.New("received a termination signal")
		}
	default:
		// non blocking
	}

	return nil
}
