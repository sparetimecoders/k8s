package config

import (
	"strings"
)

type Addons struct {
	Ingress           *Ingress           `yaml:"ingress" optional:"true"`
	ExternalDNS       *ExternalDNS       `yaml:"externalDns" optional:"true"`
	ClusterAutoscaler *ClusterAutoscaler `yaml:"clusterAutoscaler" optional:"true"`
	OauthProxy        *OauthProxy        `yaml:"oauthProxy" optional:"true"`
	Fluent            *Fluent            `yaml:"fluent" optional:"true"`
	DefaultRequests   *DefaultRequests   `yaml:"defaultRequests" optional:"true"`
}

type Addon interface {
	Manifests(config ClusterConfig) (string, error)
	Name() string
	Policies() Policies
	//Validate(config config.ClusterConfig) bool
}

func replace(org string, a map[string]string) (string, error) {
	result := org
	for k, v := range a {
		result = strings.ReplaceAll(result, k, v)
	}
	return result, nil
}
