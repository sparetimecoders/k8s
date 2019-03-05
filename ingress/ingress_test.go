package ingress

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIngress_readManifestFile(t *testing.T) {
	s, _ := readManifestFile()
	assert.Contains(t, s, "prometheus.io/port: \"10254\"")
}

func TestIngress_readGetParts(t *testing.T) {
	s, _ := readManifestFile()
	parts := getParts(s)
	assert.Equal(t, 3, len(parts))
}

func TestIngress_BuildUrl(t *testing.T) {
	namespaceUrl := buildUrl(`
apiVersion: v1
kind: Namespace
metadata:
  name: peter`)
	assert.Equal(t, "http://localhost:8080/api/v1/namespaces", namespaceUrl)

	deploymentUrl := buildUrl(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
  namespace: ingress-nginx
`)
	assert.Equal(t, "http://localhost:8080/apis/apps/v1/namespaces/ingress-nginx/deployments", deploymentUrl)

	serviceUrl := buildUrl(`
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
`)
	assert.Equal(t, "http://localhost:8080/api/v1/namespaces/default/services", serviceUrl)
}
