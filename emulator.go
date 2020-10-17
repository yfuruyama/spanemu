package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Emulator struct {
	Stdout   io.Writer
	Stderr   io.Writer
	GRPCPort int
	RESTPort int
	Cmd      *exec.Cmd
}

const emulatorVersion = "1.1.1"

func (e *Emulator) Start() error {
	if _, err := exec.LookPath("docker"); err != nil {
		if err == exec.ErrNotFound {
			return errors.New("spanemu uses docker, but docker is not found")
		}
		return fmt.Errorf("finding docker failed: %w", err)
	}

	cmd := exec.Command("docker", "run",
		"-p", fmt.Sprintf("127.0.0.1:%d:%d", e.GRPCPort, e.GRPCPort),
		"-p", fmt.Sprintf("127.0.0.1:%d:%d", e.RESTPort, e.RESTPort),
		fmt.Sprintf("gcr.io/cloud-spanner-emulator/emulator:%s", emulatorVersion))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// https://stackoverflow.com/questions/33165530/prevent-ctrlc-from-interrupting-exec-command-in-golang
		Setpgid: true,
	}
	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr
	e.Cmd = cmd

	return cmd.Start()
}

func (e *Emulator) WaitForReady() error {
	timeout := time.NewTimer(time.Second * 60)
	for {
		select {
		case <-timeout.C:
			return errors.New("waited for emulator to be up, but timeout")
		default:
			// Wait only REST server since Cloud Spanner Emulator itself waits for gRPC server to be up before starting REST server.
			url := fmt.Sprintf("http://127.0.0.1:%d/v1/projects/test-project/instanceConfigs", e.RESTPort)
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(time.Second)
	}
}

func (e *Emulator) WaitForFinish() error {
	if e.Cmd == nil {
		return errors.New("emulator not started")
	}
	return e.Cmd.Wait()
}

func (e *Emulator) Shutdown() error {
	if e.Cmd == nil {
		return errors.New("emulator not started")
	}
	if e.Cmd.ProcessState != nil && e.Cmd.ProcessState.Exited() {
		return nil
	}

	return e.Cmd.Process.Signal(os.Interrupt)
}
