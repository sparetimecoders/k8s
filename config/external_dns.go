package config

import (
	"strings"
)

type ExternalDNS struct {
	ManifestLoader
}

func (e ExternalDNS) Name() string {
	return "ExternalDNS"
}

func (e ExternalDNS) Content(config ClusterConfig) (string, error) {
	s, err := e.Load("external_dns", "external_dns.yaml")
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ReplaceAll(s, "$domain", config.Domain), "$cluster_name", config.ClusterName()), nil
}
