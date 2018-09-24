package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	resource "github.com/concourse/mock-resource"
	"github.com/sirupsen/logrus"
)

type InRequest struct {
	Source  resource.Source    `json:"source"`
	Version resource.Version   `json:"version"`
	Params  resource.GetParams `json:"params"`
}

type InResponse struct {
	Version  resource.Version         `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}

type ImageMetadata struct {
	Env  []string `json:"env"`
	User string   `json:"user"`
}

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	var req InRequest
	err := json.NewDecoder(os.Stdin).Decode(&req)
	if err != nil {
		logrus.Errorf("invalid payload: %s", err)
		os.Exit(1)
		return
	}

	if len(os.Args) < 2 {
		logrus.Errorf("destination path not specified")
		os.Exit(1)
		return
	}

	dest := os.Args[1]

	logrus.Printf("fetching version: %s", req.Version.Version)

	versionFile := filepath.Join(dest, "version")
	err = ioutil.WriteFile(versionFile, []byte(req.Version.Version+"\n"), os.ModePerm)
	if err != nil {
		logrus.Errorf("failed to write version file: %s", err)
		os.Exit(1)
		return
	}

	if req.Source.MirrorSelf || req.Params.MirrorSelfViaParams {
		logrus.WithFields(logrus.Fields{
			"via_params": req.Params.MirrorSelfViaParams,
		}).Printf("mirroring self image")

		replicateTo(filepath.Join(dest, "rootfs"))

		encTo(filepath.Join(dest, "metadata.json"), ImageMetadata{
			Env:  []string{"MIRRORED_VERSION=" + req.Version.Version},
			User: "root",
		})
	}

	json.NewEncoder(os.Stdout).Encode(InResponse{
		Version:  req.Version,
		Metadata: []resource.MetadataField{},
	})
}

func replicateTo(rootfs string) {
	err := os.MkdirAll(rootfs, os.ModePerm)
	if err != nil {
		logrus.Errorf("failed to create rootfs dir: %s", err)
		os.Exit(1)
		return
	}

	dirs, err := ioutil.ReadDir("/")
	if err != nil {
		logrus.Errorf("failed to read /: %s", err)
		os.Exit(1)
		return
	}

	for _, d := range dirs {
		rootfsDst := filepath.Join(rootfs, d.Name())

		switch d.Name() {
		case "tmp", "dev", "proc", "sys":
			// prevent recursing and copying wacky stuff
			err := os.MkdirAll(rootfsDst, d.Mode())
			if err != nil {
				logrus.Errorf("failed to create %s: %s", rootfsDst, err)
				os.Exit(1)
				return
			}

			continue
		}

		src := filepath.Join("/", d.Name())
		cp := exec.Command("cp", "-a", src, rootfsDst)
		cp.Stdout = os.Stderr
		cp.Stderr = os.Stderr
		err := cp.Run()
		if err != nil {
			logrus.Errorf("failed to copy from %s to %s: %s", src, rootfsDst, err)
			os.Exit(1)
			return
		}
	}
}

func encTo(path string, js interface{}) {
	meta, err := os.Create(path)
	if err != nil {
		logrus.Errorf("failed to create %s: %s", path, err)
		os.Exit(1)
		return
	}

	err = json.NewEncoder(meta).Encode(js)
	if err != nil {
		logrus.Errorf("failed to write %s: %s", path, err)
		os.Exit(1)
		return
	}

	err = meta.Close()
	if err != nil {
		logrus.Errorf("failed to close %s: %s", path, err)
		os.Exit(1)
		return
	}
}
