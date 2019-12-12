package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"strings"
)

type RemoteFluent struct {
	Host      string `yaml:"host"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	SharedKey string `yaml:"sharedKey"`
}

type RemoteEs struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port" default:"9200"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Fluent struct {
	RemoteFluent *RemoteFluent `yaml:"remoteFluent" optional:"true"`
	RemoteEs     *RemoteEs     `yaml:"remoteElasticSearch" optional:"true"`
}

func (i Fluent) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/fluentd")

	if i.RemoteFluent == nil && i.RemoteEs == nil {
		return "", errors.New("neither 'remoteFluent' nor 'remoteElasticSearch' was specified")
	}

	replacementStrings := make(map[string]string)

	if i.RemoteFluent != nil {
		replacementStrings["$output"] = "fluent"
		replacementStrings["$fluent_host"] = i.RemoteFluent.Host
		replacementStrings["$fluent_username"] = base64.StdEncoding.EncodeToString([]byte(i.RemoteFluent.Username))
		replacementStrings["$fluent_password"] = base64.StdEncoding.EncodeToString([]byte(i.RemoteFluent.Password))
		replacementStrings["$fluent_shared_key"] = base64.StdEncoding.EncodeToString([]byte(i.RemoteFluent.SharedKey))
	} else {
		replacementStrings["$output"] = "es"
		replacementStrings["$es_host_and_port"] = fmt.Sprintf("%s:%d", i.RemoteEs.Host, i.RemoteEs.Port)
		replacementStrings["$es_user"] = fmt.Sprintf("user %s", i.RemoteEs.Username)
		replacementStrings["$es_pass"] = fmt.Sprintf("password %s", i.RemoteEs.Password)
	}
	replacementStrings["$cluster_name"] = clusterConfig.Name

	config, _ := replace(box.MustString("fluentd_config.yaml"), replacementStrings)
	daemonset, _ := replace(box.MustString("fluentd_daemonset.yaml"), replacementStrings)
	envConfig, _ := replace(box.MustString("fluentd_env_config.yaml"), replacementStrings)
	envSecret, _ := replace(box.MustString("fluentd_env_secret.yaml"), replacementStrings)

	return strings.Join([]string{config, envConfig, envSecret, daemonset}, "\n---\n"), nil
}

func (i Fluent) Name() string {
	return "Fluent"
}

func (i Fluent) Policies() Policies {
	return Policies{}
}
