package create

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate_readGetParts(t *testing.T) {
	assert.Equal(t, 0, len(getParts("")))

	assert.Equal(t, 2, len(getParts(`
test1: 1
---
test2: 2
`)))

	assert.Equal(t, 1, len(getParts(`
test:1
`)))

	assert.Equal(t, 2, len(getParts(`
test:1
---
test:1
---
`)))

}

func TestCreate_BuildUrl(t *testing.T) {
	namespaceUrl, _ := FromHost("localhost").buildUrl(`
apiVersion: v1
kind: Namespace
metadata:
  name: peter`)
	assert.Equal(t, "https://localhost:443/api/v1/namespaces", namespaceUrl)

	deploymentUrl, _ := FromHostPort("local", 8443).buildUrl(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
  namespace: ingress-nginx
`)
	assert.Equal(t, "https://local:8443/apis/apps/v1/namespaces/ingress-nginx/deployments", deploymentUrl)

	serviceUrl, _ := New("http", "host",1000).buildUrl(`
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
`)
	assert.Equal(t, "http://host:1000/api/v1/namespaces/default/services", serviceUrl)
}
