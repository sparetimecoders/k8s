package config

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidConfig(t *testing.T) {
	c, err := ParseConfigData([]byte(`
name: test-cluster
dnsZone: example.com
vpc: vpc-123
networkCIDR: 172.21.0.0/16
masters:
  zones:
  - a
  - b
  - c
  spot: true
  type: m3.large
cloudLabels:
  environment: prod
  organisation: org
`))

	assert.Nil(t, err)

	assert.Equal(t, ClusterConfig{
		Name:              "test-cluster",
		KubernetesVersion: "1.15.5",
		DnsZone:           "example.com",
		Region:            "eu-west-1",
		Vpc:"vpc-123",
		NetworkCIDR: "172.21.0.0/16",
		Masters: MasterNodes{
			Zones:        []string{"a", "b", "c"},
			Spot:         true,
			InstanceType: "m3.large",
		},
		Nodes: Nodes{
			Min:          1,
			Max:          2,
			InstanceType: "t3.medium",
			Zones:        []string{"a","b","c"},
		},
		CloudLabels: map[string]string{
			"environment":  "prod",
			"organisation": "org",
		},
		SshKeyPath: "~/.ssh/id_rsa.pub",
	}, c)
}

func TestInvalidConfig(t *testing.T) {
	_, err := ParseConfigData([]byte(`
name: es
`))

	assert.Equal(t, "Missing required value for field(s): '[CloudLabels]'\n", err.Error())
}

func TestDefaultValuesConfig(t *testing.T) {

	c, err := ParseConfigData([]byte(`
name: test-cluster
dnsZone: example.com
cloudLabels:
  environment: prod
  organisation: org
  
`))

	assert.Nil(t, err)

	assert.Equal(t, ClusterConfig{
		Name:              "test-cluster",
		KubernetesVersion: "1.15.5",
		DnsZone:           "example.com",
		Region:            "eu-west-1",
		Masters: MasterNodes{
			Zones:        []string{"a"},
			Spot:         false,
			InstanceType: "t3.small",
		}, NetworkCIDR: "172.21.0.0/22",
		Nodes: Nodes{
			Min:          1,
			Max:          2,
			InstanceType: "t3.medium",
			Zones:        []string{"a","b","c"},
		},
		CloudLabels: map[string]string{
			"environment":  "prod",
			"organisation": "org",
		},
		SshKeyPath: "~/.ssh/id_rsa.pub",
	}, c)

}

func TestDefaultValuesWithSomeGiven(t *testing.T) {
	c, err := ParseConfigData([]byte(`
name: test-cluster
dnsZone: example.com
cloudLabels:
  environment: prod
  organisation: org
nodes:
  max: 10
`))
	assert.Nil(t, err)
	assert.Equal(t, 10, c.Nodes.Max)
	assert.Equal(t, 1, c.Nodes.Min)
}

func TestIllegalYaml(t *testing.T) {
	c, err := parseConfig(strings.NewReader(`as:a`))
	assert.NotNil(t, c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors")
}

func TestReaderError(t *testing.T) {
	_, err := parseConfig(MockReader{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected failure")
}

type MockReader struct{}

func (m MockReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("expected failure")
}
