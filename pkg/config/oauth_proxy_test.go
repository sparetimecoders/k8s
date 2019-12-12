package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOauthProxy_Manifests_Azure(t *testing.T) {
	s, _ := OauthProxy{Provider: "azure", AzureTenantId: "TenantId", CookieSecret: "Cookie", ClientSecret: "Client", ClientId: "ClientId", EmailDomain: "domain"}.Manifests(ClusterConfig{})
	assert.Contains(t, s, `
          env:
            - name: OAUTH2_PROXY_CLIENT_ID
              value: ClientId
            - name: OAUTH2_PROXY_CLIENT_SECRET
              value: Client
            - name: OAUTH2_PROXY_COOKIE_SECRET
              value: Cookie`)
	assert.Contains(t, s, `- --provider=azure`)
	assert.Contains(t, s, `- --azure-tenant=TenantId`)
	assert.Contains(t, s, `- --email-domain=domain`)
}

func TestOauthProxy_Manifests(t *testing.T) {
	s, _ := OauthProxy{Provider: "none", AzureTenantId: "TenantId", CookieSecret: "Cookie", ClientSecret: "Client", ClientId: "ClientId"}.Manifests(ClusterConfig{})
	assert.Contains(t, s, `
          env:
            - name: OAUTH2_PROXY_CLIENT_ID
              value: ClientId
            - name: OAUTH2_PROXY_CLIENT_SECRET
              value: Client
            - name: OAUTH2_PROXY_COOKIE_SECRET
              value: Cookie`)
	assert.NotContains(t, s, `- --azure-tenant=TenantId`)
}
