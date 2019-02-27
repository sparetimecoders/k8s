package ingress

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIngress_readManifestFile(t *testing.T) {
	s, _ := readManifestFile()
	assert.Contains(t, s, "prometheus.io/port: \"10254\"")
}
