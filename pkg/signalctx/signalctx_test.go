// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package signalctx

import (
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"
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
