package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func generateConfigName() string {
	return fmt.Sprintf("spanemu-%d", time.Now().Unix())
}

func setupCloudSDK(projectID string) (string, error) {
	fmt.Println("Creating an ephemeral gcloud configuration...")
	configName := generateConfigName()

	cmd := exec.Command("gcloud", "config", "configurations", "create", "--no-activate", configName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", configName, "set", "auth/disable_credentials", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", configName, "set", "project", projectID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command("gcloud", "config", "--configuration", configName, "set", "api_endpoint_overrides/spanner", "http://localhost:9020/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return configName, nil
}

func cleanupCloudSDK() error {
	return nil
}