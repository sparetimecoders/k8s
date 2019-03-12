package config

import (
	"github.com/GeertJohan/go.rice"
)

type ExternalDNS struct {
	NodePolicies []Policy `yaml:"nodePolicies"`
}

func (e ExternalDNS) Name() string {
	return "ExternalDNS"
}

func (e ExternalDNS) Manifests(config ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/external_dns")
	s := box.MustString("external_dns.yaml")

	return replace(s, map[string]string{"$domain": config.Domain, "$cluster_name": config.ClusterName()})
}

func (e ExternalDNS) Policies() Policies {
	return Policies{Node: e.NodePolicies}
}
