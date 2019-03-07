package config

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"strings"
)

type ManifestLoader struct {
	_ struct{}
}

func (m ManifestLoader) Load(path string, files ...string) (string, error) {
	box := rice.MustFindBox(fmt.Sprintf("manifests/%s", path))
	var filesData []string
	for _, f := range files {
		fileData, err := box.String(f)
		if err != nil {
			return "", err
		}
		filesData = append(filesData, fileData)
	}
	return strings.Join(filesData, "\n---\n"), nil
}
