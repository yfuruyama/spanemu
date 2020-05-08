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
	Stdout io.Writer
	Stderr io.Writer
	Cmd    *exec.Cmd
}

func (e *Emulator) Start() error {
	if _, err := exec.LookPath("docker"); err != nil {
		if err == exec.ErrNotFound {
			return errors.New("spanemu uses docker, but docker is not found")
		}
		return fmt.Errorf("finding docker failed: %w", err)
	}

	cmd := exec.Command("docker", "run", "-p", "127.0.0.1:9010:9010", "-p", "127.0.0.1:9020:9020", "gcr.io/cloud-spanner-emulator/emulator:0.7.28")
	// https://stackoverflow.com/questions/33165530/prevent-ctrlc-from-interrupting-exec-command-in-golang
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr
	e.Cmd = cmd

	return cmd.Start()
}

// WaitForReady waits until Cloud Spanner Emulator is up and ready.
func (e *Emulator) WaitForReady() error {
	timeout := time.NewTimer(time.Second * 60)
	for {
		select {
		case <-timeout.C:
			return errors.New("waited for emulator to be up, but timeout")
		default:
			// only check REST server since Cloud Spanner Emulator starts REST server after gRPC server is ready
			resp, err := http.Get("http://localhost:9020/v1/projects/test-project/instanceConfigs")
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(time.Second)
	}
}

func (e * Emulator) WaitForFinish() error {
	if e.Cmd == nil {
		return errors.New("emulator not started")
	}
	return e.Cmd.Wait()
}

func (e *Emulator) Shutdown() error {
	if e.Cmd == nil {
		return errors.New("emulator not started")
	}
	return e.Cmd.Process.Signal(os.Interrupt)
}
