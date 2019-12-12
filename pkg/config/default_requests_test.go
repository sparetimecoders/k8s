package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultRequests_Manifests(t *testing.T) {
	s, _ := DefaultRequests{ExcludedNs: []string{"kube-system", "ingress-nginx"}, Memory: "1Pi"}.Manifests(ClusterConfig{})
	assert.Contains(t, s, "registry.gitlab.com/unboundsoftware/default-request-adder:1.0")
	assert.Contains(t, s, "-excluded-ns=kube-system,ingress-nginx")
	assert.Contains(t, s, "-memory=1Pi")
}
