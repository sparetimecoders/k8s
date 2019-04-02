package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddons_List(t *testing.T) {
	c := config(t)

	assert.Equal(t, 2, len(c.AllAddons()))
	assert.Equal(t, 10, c.Addons.Ingress.Aws.Timeout)
	assert.Equal(t, "https", c.Addons.Ingress.Aws.SSLPort)
}

func TestAddons_GetConfiguredAddon(t *testing.T) {
	c := config(t)

	ingress := c.GetAddon(Ingress{}).(*Ingress)
	assert.Equal(t, 10, ingress.Aws.Timeout)
}

func TestAddons_GetNonConfiguredAddon(t *testing.T) {
	c := config(t)

	scaler := c.GetAddon(ClusterAutoscaler{})
	assert.Nil(t, scaler)
}

func config(t *testing.T) ClusterConfig {
	c, err := ParseConfigData([]byte(`
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
      certificateARN: "arn:...."
  externalDns: {}
`))
	assert.Nil(t, err)
	return c
}
