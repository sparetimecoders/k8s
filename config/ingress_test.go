package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIngress_readManifestFile(t *testing.T) {
	s, _ := Ingress{}.Manifests(ClusterConfig{})
	assert.Contains(t, s, "prometheus.io/port: \"10254\"")
	assert.Contains(t, s, "server-names-hash-bucket-size: \"128\"")
}
