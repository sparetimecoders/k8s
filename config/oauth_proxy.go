package config

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
)

type OauthProxy struct {
	Provider string `yaml:"provider" default:"azure"`
	AzureTenantId string `yaml:"azureTenantId" default:"common"`
	ClientId string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	CookieSecret string `yaml:"cookieSecret"`
}

func (o OauthProxy) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/oauth_proxy")
	manifest := box.MustString("oauth_proxy.yaml")
	additionalArgs:=""
	if o.Provider == "azure" {
		additionalArgs = fmt.Sprintf("- --azure-tenant=%v", o.AzureTenantId)
	}

	return replace(manifest, map[string]string{
		"$additional_args": additionalArgs,
		"$oauth2_proxy_client_id": o.ClientId,
		"$oauth2_proxy_client_secret": o.ClientSecret,
		"$oauth2_proxy_cookie_secret": o.CookieSecret,})
}

func (o OauthProxy) Name() string {
	return "OauthProxy"
}

func (o OauthProxy) Policies() Policies {
	return Policies{}
}
