package main

import (
  "fmt"
  "github.com/stretchr/testify/assert"
  "gitlab.com/sparetimecoders/k8s-go/config"
  "testing"
)

func TestInvalidConfig(t *testing.T) {
  if testing.Short() {
    fmt.Println("Skipping test in short-mode")
  }

  _, err := config.ParseConfig([]byte(`
name: es
`))

  assert.Equal(t, "Missing required value for field(s): '[DnsZone CloudLabels]'\n", err.Error())
}

func TestValidConfig(t *testing.T) {
  if testing.Short() {
    fmt.Println("Skipping test in short-mode")
  }

  c, err := config.ParseConfig([]byte(`
name: es
dnsZone: example.com
masterZones:
  - a
  - b
  - c
cloudLabels:
  environment: prod
  organisation: dSPA
`))

  assert.Nil(t, err)

  assert.Equal(t, config.Cluster {
    Name: "es",
    KubernetesVersion: "1.11.7",
    DnsZone: "example.com",
    Region: "eu-west-1",
    MasterZones: []string {"a", "b", "c"},
    NetworkCIDR: "172.21.0.0/22",
    Nodes: config.Nodes{
      Min: 1,
      Max: 2,
      InstanceType: "t3.medium",
    },
    MasterInstanceType: "t3.small",
    CloudLabels: map[string]string {
      "environment": "prod",
      "organisation": "dSPA",
    },
    SshKeyPath: "~/.ssh/id_rsa.pub",
  }, c)
}
