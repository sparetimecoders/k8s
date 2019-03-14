package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFluent_readManifestFile(t *testing.T) {
	_, err := Fluent{}.Manifests(ClusterConfig{})
	assert.Error(t, err, "")
}

func TestFluent_readManifestFileWithRemoteFluentSettings(t *testing.T) {
	s, _ := Fluent{RemoteFluent: &RemoteFluent{Host: "host", Username: "user", Password: "pass", SharedKey: "key"}}.Manifests(ClusterConfig{Name: "cluster"})
	assert.NotContains(t, s, "@include $output-output.conf")
	assert.Contains(t, s, "@include fluent-output.conf")
	assert.Contains(t, s, "name host")
	assert.Contains(t, s, "host host")
	assert.Contains(t, s, `CLUSTER_NAME: "cluster"`)
	assert.Contains(t, s, `FLUENT_USER: "user"`)
	assert.Contains(t, s, `FLUENT_PASSWORD: "pass"`)
	assert.Contains(t, s, `FLUENT_SHAREDKEY: "key"`)
}

func TestFluent_readManifestFileWithRemoteEsSettings(t *testing.T) {
	s, _ := Fluent{RemoteEs: &RemoteEs{Host: "host", Port: 1234}}.Manifests(ClusterConfig{})
	assert.NotContains(t, s, "@include $output-output.conf")
	assert.Contains(t, s, "@include es-output.conf")
	assert.Contains(t, s, "hosts host:1234")
}
