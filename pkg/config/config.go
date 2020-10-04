// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"
)

// env represents a environment variabels for kube-timeleap.
type env struct {
	Debug bool `envconfig:"DEBUG"`
}

// compile time check whether the env implements zapcore.ObjectMarshaler interface.
var _ zapcore.ObjectMarshaler = (*env)(nil)

// MarshalLogObject implements zapcore.ObjectMarshaler.
func (e env) MarshalLogObject(enc zapcore.ObjectEncoder) (_ error) {
	// Debug
	enc.AddBool("debug", e.Debug)

	return
}

// envValue is the global env variable with atomically.
var envValue atomic.Value // actual type is env

// envOnce for runs the Process once.
var envOnce sync.Once

// Process populates the specified struct based on environment variables.
func Process() (e env, err error) {
	envOnce.Do(func() {
		err = envconfig.Process("", &e)
		if err != nil {
			err = handleError(err)
			return
		}
		envValue.Store(e)
	})

	return e, err
}

func handleError(err error) error {
	if errors.Is(err, envconfig.ErrInvalidSpecification) {
		return fmt.Errorf("failed to Process: %w", err)
	}

	var e *envconfig.ParseError
	if errors.As(err, &e) {
		return fmt.Errorf("failed to Process: %w", e)
	}

	return err
}

func loadEnv() env {
	e, ok := envValue.Load().(env)
	if ok {
		return e // fast pass
	}

	// fallback
	envconfig.MustProcess("", &e)
	envValue.Store(e)

	return e
}

// IsDebug reports whether the debug environment.
func IsDebug() bool {
	return loadEnv().Debug
}
