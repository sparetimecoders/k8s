package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContent(t *testing.T) {
	s, _ := ExternalDNS{Domain: "replaced_domain"}.Manifests(ClusterConfig{Name: "replaced_cluster_name", DnsZone: "dns"})
	assert.Contains(t, s, "registry.opensource.zalan.do/teapot/external-dns:v0.5.14")
	assert.Contains(t, s, "domain-filter=replaced_domain")
	assert.Contains(t, s, "txt-owner-id=replaced_cluster_name.dns")
}
