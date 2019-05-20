package config

import (
	"github.com/GeertJohan/go.rice"
)

type ExternalDNS struct {
	Domain string `yaml:"domain" optional:"true"`
}

func (e ExternalDNS) Name() string {
	return "ExternalDNS"
}

func (e ExternalDNS) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/external_dns")
	s := box.MustString("external_dns.yaml")

	return replace(s, map[string]string{"$domain": e.Domain, "$cluster_name": clusterConfig.ClusterName()})
}

func (e ExternalDNS) Policies() Policies {
	return Policies{Node: []Policy{
		{Actions: []string{"route53:ChangeResourceRecordSets"}, Effect: "Allow", Resources: []string{"arn:aws:route53:::hostedzone/*"}},
		{Actions: []string{"route53:ListHostedZones", "route53:ListResourceRecordSets"}, Effect: "Allow", Resources: []string{"*"}},
	}}
}
