// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package signalctx

import (
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testTask struct{}

func (t *testTask) Run(c chan os.Signal) {
	for {
		go func(stopChan chan os.Signal) {
			fmt.Printf("...")
			time.Sleep(time.Second)
			stopChan <- os.Interrupt
		}(c)

		for i := 0; i < 5; i++ {
			fmt.Print("#")
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println() // add newline
	}
}

func TestSetupSignalHandler(t *testing.T) {
	stop := setupSignalHandler()
	task := &testTask{}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		task.Run(c)
	}(c)

	select {
	case sig := <-c:
		fmt.Printf("\ngot %s signal\n", sig)
	case _, ok := <-stop:
		if ok {
			t.Fatal("context not stopped")
		}
	}
}

func Test_signalContext_Deadline(t *testing.T) {
	tests := []struct {
		name   string
		wantOk bool
	}{
		{
			name:   "Valid",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := &signalContext{}
			gotDeadline, gotOk := ctx.Deadline()
			if !gotDeadline.IsZero() {
				t.Fatalf("signalContext.Deadline() gotDeadline not IsZero = %v", gotDeadline)
			}
			if gotOk != tt.wantOk {
				t.Fatalf("signalContext.Deadline() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_signalContext_Value(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want interface{}
	}{
		{
			name: "Valid",
			args: "test",
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := &signalContext{}
			if got := ctx.Value(tt.args); !cmp.Equal(got, tt.want) {
				t.Fatalf("signalContext.Value(%v) should be return nil = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
