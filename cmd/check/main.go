package main

import (
	"encoding/json"
	"os"

	resource "github.com/concourse/mock-resource"
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

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	var req CheckRequest
	err := decoder.Decode(&req)
	if err != nil {
		logrus.Errorf("invalid payload: %s", err)
		os.Exit(1)
		return
	}

	response := CheckResponse{}

	if req.Source.ForceVersion != "" {
		response = append(response, resource.Version{
			Version: req.Source.ForceVersion,
		})
	} else if req.Version != nil {
		response = append(response, *req.Version)
	} else if !req.Source.NoInitialVersion {
		response = append(response, resource.Version{
			Version: req.Source.InitialVersion(),
		})
	}

	json.NewEncoder(os.Stdout).Encode(response)
}
