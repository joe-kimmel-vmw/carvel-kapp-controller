// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"os/exec"
)

type CmdRunner interface {
	Run(*exec.Cmd) error
	RunWithCancel(cmd *exec.Cmd, cancelCh chan struct{}) error
}

type PlainCmdRunner struct{}

var _ CmdRunner = PlainCmdRunner{}

func NewPlainCmdRunner() PlainCmdRunner {
	return PlainCmdRunner{}
}

func (PlainCmdRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (PlainCmdRunner) RunWithCancel(cmd *exec.Cmd, cancelCh chan struct{}) error {
	select {
	case <-cancelCh:
		return fmt.Errorf("Already canceled")
	default:
		// continue with execution
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	doneCh := make(chan error, 1)
	go func() {
		doneCh <- cmd.Wait()
	}()

	select {
	case <-cancelCh:
		err := cmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("Killing process: %s", err)
		}
		return fmt.Errorf("Process was canceled")

	case err := <-doneCh:
		return err
	}
}
