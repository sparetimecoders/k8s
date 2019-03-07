package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContent(t *testing.T) {
	s, _ := ExternalDNS{}.Content(ClusterConfig{Domain: "replaced_domain", Name: "replaced_cluster_name", DnsZone: "dns"})
	assert.Contains(t, s, "registry.opensource.zalan.do/teapot/external-dns:latest")
	assert.Contains(t, s, "domain-filter=replaced_domain")
	assert.Contains(t, s, "txt-owner-id=replaced_cluster_name.dns")
}
