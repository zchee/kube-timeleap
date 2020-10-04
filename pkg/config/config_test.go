// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"encoding/json"
	"os"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	zapobserver "go.uber.org/zap/zaptest/observer"
)

func Test_env_MarshalLogObject(t *testing.T) {
	type fields struct {
		Debug bool
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "Omitempty",
			fields: fields{},
			want:   []byte(`{"Debug":false}`),
		},
		{
			name: "AllFields",
			fields: fields{
				Debug: true,
			},
			want: []byte(`{"Debug":true}`),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			core, obs := zapobserver.New(zapcore.DebugLevel)
			logger := zap.New(core)

			e := &env{
				Debug: tt.fields.Debug,
			}
			logger.Debug("Test_env_MarshalLogObject", zap.Object("env", e))

			for _, gotEntry := range obs.All() {
				gotIface, ok := gotEntry.Context[0].Interface.(*env)
				if ok {
					got, err := json.Marshal(gotIface)
					if err != nil {
						t.Fatal(err)
					}
					if diff := cmp.Diff(tt.want, got); diff != "" {
						t.Fatalf("(-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}

func setEnvs(tb testing.TB, keyValue map[string]string) (cleanup func()) {
	tb.Helper()

	for k, v := range keyValue {
		if err := os.Setenv(k, v); err != nil {
			tb.Fatal(err)
		}
	}

	return func() {
		for k := range keyValue {
			if err := os.Unsetenv(k); err != nil {
				tb.Fatal(err)
			}
		}
	}
}

func TestIsDebug(t *testing.T) {
	tests := []struct {
		name     string
		keyValue map[string]string
		want     bool
	}{
		{
			name: "HaveEnv",
			keyValue: map[string]string{
				"DEBUG": "true",
			},
			want: true,
		},
		{
			name:     "Empty",
			keyValue: map[string]string{},
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setEnvs(t, tt.keyValue)
			defer func() {
				cleanup()
				envValue = atomic.Value{}
			}()

			if got := IsDebug(); got != tt.want {
				t.Fatalf("IsDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}
