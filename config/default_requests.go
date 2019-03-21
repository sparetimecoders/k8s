package config

import (
	"github.com/GeertJohan/go.rice"
	"strings"
)

type DefaultRequests struct {
	ExcludedNs []string `yaml:"excludedNs" default:"kube-system,ingress-nginx"`
	Memory     string   `yaml:"memory" default:"1Pi"`
}

func (d DefaultRequests) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/default_requests")
	manifest := box.MustString("default_requests.yaml")

	return replace(manifest, map[string]string{
		"$excluded-ns": strings.Join(d.ExcludedNs, ","),
		"$memory":      d.Memory,
	})
}

func (d DefaultRequests) Name() string {
	return "DefaultRequests"
}

func (d DefaultRequests) Policies() Policies {
	return Policies{}
}
