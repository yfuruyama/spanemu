package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type InstanceRequest struct {
	InstanceID string   `json:"instanceId"`
	Instance   Instance `json:"instance"`
}

type Instance struct {
	Name        string `json:"name"`
	Config      string `json:"config"`
	DisplayName string `json:"displayName"`
	NodeCount   int    `json:"nodeCount"`
}

type DatabaseRequest struct {
	CreateStatement string `json:"createStatement"`
}

func createInstance(port int, projectID, instanceID string) error {
	fmt.Println("Creating instance...")
	req := InstanceRequest{
		InstanceID: instanceID,
		Instance: Instance{
			Name:        fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
			Config:      "emulator-config",
			DisplayName: "Test Instance",
			NodeCount:   1,
		},
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/v1/projects/%s/instances", port, projectID)
	_, err := http.Post(url, "application/json", &body)
	return err
}

func createDatabase(port int, projectID, instanceID, databaseID string) error {
	fmt.Println("Creating database...")
	req := DatabaseRequest{
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseID),
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/v1/projects/%s/instances/%s/databases", port, projectID, instanceID)
	_, err := http.Post(url, "application/json", &body)
	return err
}
