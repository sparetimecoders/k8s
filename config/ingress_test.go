package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIngress_readManifestFile(t *testing.T) {
	s, _ := Ingress{}.Manifests(ClusterConfig{})
	assert.Contains(t, s, "prometheus.io/port: \"10254\"")
	assert.Contains(t, s, "server-names-hash-bucket-size: \"128\"")
	assert.NotContains(t, s, "annotations_placeholder")
}

func TestIngress_readManifestFileWithAwsSettings(t *testing.T) {
	s, _ := Ingress{Aws: &Aws{Protocol: "protocol", Port: "port", SecurityPolicy: "policy", CertificateARN: "ARN", Timeout: 10}}.Manifests(ClusterConfig{})
	assert.NotContains(t, s, "annotations_placeholder")
	assert.Contains(t, s, "    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: ARN")
	assert.Contains(t, s, "    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: protocol")
	assert.Contains(t, s, "    service.beta.kubernetes.io/aws-load-balancer-ssl-negotiation-policy: policy")
	assert.Contains(t, s, "    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: port")
	assert.Contains(t, s, "    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: 10")
}
