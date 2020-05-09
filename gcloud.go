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

func (c *GCloudConfig) Setup(projectID string, port int) error {
	cmd := exec.Command("gcloud", "config", "configurations", "create", "--no-activate", c.Name)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "auth/disable_credentials", "true")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "project", projectID)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "api_endpoint_overrides/spanner", fmt.Sprintf("http://127.0.0.1:%d/", port))
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *GCloudConfig) CleanUp() error {
	cmd := exec.Command("gcloud", "-q", "config", "configurations", "delete", c.Name)
	return cmd.Run()
}

func generateConfigName() string {
	return fmt.Sprintf("spanemu-%d", time.Now().Unix())
}
