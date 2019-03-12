package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidConfig(t *testing.T) {
	if testing.Short() {
		fmt.Println("Skipping test in short-mode")
	}

	c, err := ParseConfig([]byte(`
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
	if testing.Short() {
		fmt.Println("Skipping test in short-mode")
	}

	_, err := ParseConfig([]byte(`
name: es
`))

	assert.Equal(t, "Missing required value for field(s): '[DnsZone Domain CloudLabels]'\n", err.Error())
}

func TestDefaultValuesConfig(t *testing.T) {
	if testing.Short() {
		fmt.Println("Skipping test in short-mode")
	}

	c, err := ParseConfig([]byte(`
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

func TestAddons(t *testing.T) {
	c, err := ParseConfig([]byte(`
name: es
dnsZone: example.com
domain: example.com
cloudLabels:
  environment: prod
  organisation: dSPA
addons:
  ingress: 
    aws:
      timeout: 10
  externalDns: {}
`))
	assert.Nil(t, err)

	assert.Equal(t, 2, len(c.Addons.List()))
	assert.Equal(t, 10, c.Addons.Ingress.Aws.Timeout)
	assert.Equal(t, "https", c.Addons.Ingress.Aws.SSLPort)

}
