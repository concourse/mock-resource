package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	var req InRequest
	err := decoder.Decode(&req)
	if err != nil {
		logrus.Fatalf("invalid payload: %s", err)
		return
	}

	if len(os.Args) < 2 {
		logrus.Fatal("destination path not specified")
		return
	}

	dest := os.Args[1]

	if req.Source.Log != "" {
		logrus.Info(req.Source.Log)
	}

	privileged, err := resource.IsPrivileged()
	if err != nil {
		logrus.Fatalf("could not check privilege: %s", err)
		return
	}

	if privileged {
		logrus.Printf("fetching in a privileged container")
	}

	encTo(filepath.Join(dest, "privileged"), privileged)

	logrus.Printf("fetching version: %s", req.Version.Version)

	versionFile := filepath.Join(dest, "version")
	err = ioutil.WriteFile(versionFile, []byte(req.Version.Version+"\n"), os.ModePerm)
	if err != nil {
		logrus.Fatalf("failed to write version file: %s", err)
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

	files := map[string]interface{}{}

	for path, content := range req.Source.CreateFiles {
		files[path] = content
	}

	for path, content := range req.Params.CreateFiles {
		files[path] = content
	}

	for path, content := range files {
		var bs []byte

		if str, ok := content.(string); ok {
			bs = []byte(str)
		} else {
			var err error
			bs, err = json.Marshal(content)
			if err != nil {
				logrus.Fatalf("failed to marshal content (%v): %s", content, err)
				return
			}
		}

		str := string(bs)
		str = strings.ReplaceAll(str, "START_VAR", "((")
		str = strings.ReplaceAll(str, "END_VAR", "))")
		bs = []byte(str)

		filePath := filepath.Join(dest, path)

		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			logrus.Fatalf("failed to create directory '%s': %s", filepath.Dir(filePath), err)
			return
		}

		err = ioutil.WriteFile(filePath, bs, 0644)
		if err != nil {
			logrus.Fatalf("failed to write to file '%s': %s", path, err)
			return
		}
	}

	json.NewEncoder(os.Stdout).Encode(InResponse{
		Version:  req.Version,
		Metadata: req.Source.Metadata,
	})
}

func replicateTo(rootfs string) {
	err := os.MkdirAll(rootfs, os.ModePerm)
	if err != nil {
		logrus.Fatalf("failed to create rootfs dir: %s", err)
		return
	}

	dirs, err := ioutil.ReadDir("/")
	if err != nil {
		logrus.Fatalf("failed to read /: %s", err)
		return
	}

	for _, d := range dirs {
		rootfsDst := filepath.Join(rootfs, d.Name())

		switch d.Name() {
		case "tmp", "dev", "proc", "sys":
			// prevent recursing and copying wacky stuff
			err := os.MkdirAll(rootfsDst, d.Mode())
			if err != nil {
				logrus.Fatalf("failed to create %s: %s", rootfsDst, err)
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
			logrus.Fatalf("failed to copy from %s to %s: %s", src, rootfsDst, err)
			return
		}
	}
}

func encTo(path string, js interface{}) {
	meta, err := os.Create(path)
	if err != nil {
		logrus.Fatalf("failed to create %s: %s", path, err)
		return
	}

	err = json.NewEncoder(meta).Encode(js)
	if err != nil {
		logrus.Fatalf("failed to write %s: %s", path, err)
		return
	}

	err = meta.Close()
	if err != nil {
		logrus.Fatalf("failed to close %s: %s", path, err)
		return
	}
}
