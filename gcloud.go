package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type GCloudConfig struct {
	Name string
}

func NewGCloudConfig() *GCloudConfig {
	return &GCloudConfig{generateConfigName()}
}

func (c *GCloudConfig) setup(projectID string) error {
	fmt.Println("Creating an ephemeral gcloud configuration...")

	cmd := exec.Command("gcloud", "config", "configurations", "create", "--no-activate", c.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "auth/disable_credentials", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "project", projectID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", c.Name, "set", "api_endpoint_overrides/spanner", "http://localhost:9020/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil

}

func (c *GCloudConfig) cleanUp() error {
	fmt.Println("Deleting an ephemeral gcloud configuration...")
	cmd := exec.Command("gcloud", "-q", "config", "configurations", "delete", c.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func generateConfigName() string {
	return fmt.Sprintf("spanemu-%d", time.Now().Unix())
}
