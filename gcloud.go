package main

import (
	"fmt"
	"os/exec"
	"time"
)

type GCloudConfig struct {
	Name string
}

func NewGCloudConfig() *GCloudConfig {
	return &GCloudConfig{generateConfigName()}
}

func (c *GCloudConfig) Setup(projectID string) error {
	cmd := exec.Command("gcloud", "config", "configurations", "create", "--no-activate", c.Name)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "auth/disable_credentials", "true")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "project", projectID)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "api_endpoint_overrides/spanner", "http://localhost:9020/")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *GCloudConfig) CleanUp() error {
	cmd := exec.Command("gcloud", "-q", "config", "configurations", "delete", c.Name)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func generateConfigName() string {
	return fmt.Sprintf("spanemu-%d", time.Now().Unix())
}
