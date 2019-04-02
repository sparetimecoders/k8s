package pkg

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"testing"
)

var baseConfig = func() config.ClusterConfig {
	c, _ := config.ParseConfigData([]byte(`
name: gotest
dnsZone: example.com
domain: example.com
kubernetesVersion: 1.12.2
masterZones:
  - a
cloudLabels: {}
`))
	return c
}

var baseConfig2, _ = config.ParseConfigData([]byte(`
name: gotest
dnsZone: example.com
domain: example.com
kubernetesVersion: 1.12.2
masterZones:
  - a
cloudLabels: {}
`))

func TestPolicies(t *testing.T) {
	// No policies
	p := policies(baseConfig())
	assert.Equal(t, config.Policies{Node: []config.Policy{}, Master: []config.Policy{}}, p)
}

func TestPolicies2(t *testing.T) {
	policy := config.Policy{Actions: []string{"Action"}, Effect: "Allow", Resources: []string{"Resource"}}
	clusterConfig := baseConfig()
	clusterConfig.Nodes.Policies = []config.Policy{policy}

	p := policies(clusterConfig)
	assert.Equal(t, 1, len(p.Node))
	assert.Equal(t, policy, p.Node[0])
	assert.Equal(t, []config.Policy{}, p.Master)
}
func TestPolicies3(t *testing.T) {
	clusterConfig := baseConfig()
	clusterConfig.Addons = &config.Addons{
		ClusterAutoscaler: &config.ClusterAutoscaler{},
	}

	p := policies(clusterConfig)
	assert.Equal(t, 1, len(p.Node))
	assert.Equal(t, config.ClusterAutoscaler{}.Policies().Node, p.Node)
	assert.Equal(t, []config.Policy{}, p.Master)
}

func TestPolicies4(t *testing.T) {
	policy := config.Policy{Actions: []string{"Action"}, Effect: "Allow", Resources: []string{"Resource1", "Resource2"}}
	clusterConfig := baseConfig()
	clusterConfig.Nodes.Policies = []config.Policy{policy}

	clusterConfig.Addons = &config.Addons{
		ClusterAutoscaler: &config.ClusterAutoscaler{},
	}
	// No policies
	p := policies(clusterConfig)
	assert.Equal(t, 2, len(p.Node))
	assert.Equal(t, config.ClusterAutoscaler{}.Policies().Node[0], p.Node[1])
	assert.Equal(t, policy, p.Node[0])
	assert.Equal(t, []config.Policy{}, p.Master)
}
