package config

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidConfig(t *testing.T) {
	c, err := ParseConfigData([]byte(`
name: es
dnsZone: example.com
domain: example.com
masterZones:
  - a
  - b
  - c
cloudLabels:
  environment: prod
  organisation: dSPA
`))

	assert.Nil(t, err)

	assert.Equal(t, ClusterConfig{
		Name:              "es",
		KubernetesVersion: "1.11.7",
		DnsZone:           "example.com",
		Domain:            "example.com",
		Region:            "eu-west-1",
		MasterZones:       []string{"a", "b", "c"},
		NetworkCIDR:       "172.21.0.0/22",
		Nodes: Nodes{
			Min:          1,
			Max:          2,
			InstanceType: "t3.medium",
		},
		MasterInstanceType: "t3.small",
		CloudLabels: map[string]string{
			"environment":  "prod",
			"organisation": "dSPA",
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
name: es
dnsZone: example.com
domain: example.com
cloudLabels:
  environment: prod
  organisation: dSPA
  
`))

	assert.Nil(t, err)

	assert.Equal(t, ClusterConfig{
		Name:              "es",
		KubernetesVersion: "1.11.7",
		DnsZone:           "example.com",
		Domain:            "example.com",
		Region:            "eu-west-1",
		MasterZones:       []string{"a"},
		NetworkCIDR:       "172.21.0.0/22",
		Nodes: Nodes{
			Min:          1,
			Max:          2,
			InstanceType: "t3.medium",
		},
		MasterInstanceType: "t3.small",
		CloudLabels: map[string]string{
			"environment":  "prod",
			"organisation": "dSPA",
		},
		SshKeyPath: "~/.ssh/id_rsa.pub",
	}, c)

}

func TestDefaultValuesWithSomeGiven(t *testing.T) {
	c, err := ParseConfigData([]byte(`
name: es
dnsZone: example.com
domain: example.com
cloudLabels:
  environment: prod
  organisation: dSPA
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
