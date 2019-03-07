package config

import (
	"github.com/GeertJohan/go.rice"
)

type ExternalDNS struct{}

func (e ExternalDNS) Name() string {
	return "ExternalDNS"
}

func (e ExternalDNS) Content(config ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/external_dns")
	s := box.MustString("external_dns.yaml")

	return replace(s, map[string]string{"$domain": config.Domain, "$cluster_name": config.ClusterName()})
}
