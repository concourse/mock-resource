package main

import (
	"encoding/json"
	"os"

	resource "github.com/concourse/mock-resource"
	"github.com/sirupsen/logrus"
)

type OutRequest struct {
	Source  resource.Source    `json:"source"`
	Version resource.Version   `json:"version"`
	Params  resource.PutParams `json:"params"`
}

type OutResponse struct {
	Version  resource.Version         `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	var req OutRequest
	err := decoder.Decode(&req)
	if err != nil {
		logrus.Errorf("invalid payload: %s", err)
		os.Exit(1)
		return
	}

	if len(os.Args) < 2 {
		logrus.Errorf("source path not specified")
		os.Exit(1)
		return
	}

	privileged, err := resource.IsPrivileged()
	if err != nil {
		logrus.Errorf("could not check privilege: %s", err)
		os.Exit(1)
		return
	}

	version := resource.Version{Version: req.Params.Version}

	if privileged {
		logrus.Printf("pushing in a privileged container")
		version.Privileged = "true"
	}

	logrus.Printf("pushing version: %s", req.Params.Version)

	json.NewEncoder(os.Stdout).Encode(OutResponse{
		Version:  version,
		Metadata: req.Source.Metadata,
	})
}
