package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClusterAutoscaler_readManifestFile(t *testing.T) {
	s, _ := ClusterAutoscaler{}.Manifests(ClusterConfig{Name: "name", DnsZone: "zone", Region: "a region"})
	assert.Contains(t, s, "node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/name.zone")
	assert.Contains(t, s, `name: AWS_REGION
          value: a region`)
	assert.NotContains(t, s, "$")
}
