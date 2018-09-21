package main

import (
	"encoding/json"
	"os"

	resource "github.com/concourse/mirror-resource"
	"github.com/sirupsen/logrus"
)

type CheckRequest struct {
	Source  resource.Source   `json:"source"`
	Version *resource.Version `json:"version"`
}

type CheckResponse []resource.Version

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	var req CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		logrus.Errorf("invalid payload: %s", err)
		os.Exit(1)
		return
	}

	response := CheckResponse{}

	if req.Version != nil {
		response = append(response, *req.Version)
	} else if !req.Source.NoInitialVersion {
		response = append(response, resource.Version{
			Version: req.Source.InitialVersion(),
		})
	}

	json.NewEncoder(os.Stdout).Encode(response)
}
