package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoader(t *testing.T) {
	s, _ := ExternalDNS{}.Content(ClusterConfig{})
	assert.Contains(t, s, "registry.opensource.zalan.do/teapot/external-dns:latest")
}
