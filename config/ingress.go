package config

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"strings"
)

type Aws struct {
	SecurityPolicy string `yaml:"securityPolicy" default:"ELBSecurityPolicy-TLS-1-2-2017-01"`
	CertificateARN string `yaml:"certificateARN"`
	Protocol       string `yaml:"protocol" default:"http"`
	SSLPort        string `yaml:"sslPort" default:"https"`
	Timeout        int    `yaml:"timeout" default:"60"`
}

type Ingress struct {
	Aws *Aws `yaml:"aws"`
	_   struct{}
}

func (i Ingress) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/ingress")
	manifest := box.MustString("ingress.yaml")

	var replacementString []string
	if i.Aws != nil {
		if i.Aws.CertificateARN != "" {
			replacementString = append(replacementString, fmt.Sprintf("    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: %v", i.Aws.CertificateARN))
		}
		replacementString = append(replacementString, fmt.Sprintf("    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: \"%d\"", i.Aws.Timeout))
		replacementString = append(replacementString, fmt.Sprintf("    service.beta.kubernetes.io/aws-load-balancer-ssl-negotiation-policy: %v", i.Aws.SecurityPolicy))
		replacementString = append(replacementString, fmt.Sprintf("    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: %v", i.Aws.SSLPort))
		replacementString = append(replacementString, fmt.Sprintf("    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: %v", i.Aws.Protocol))
	}
	manifest = strings.Replace(manifest, `    annotations_placeholder: ""`, strings.Join(replacementString, "\n"), 1)
	return strings.Join([]string{manifest, box.MustString("nginx-config.yaml")}, "\n---\n"), nil
}

func (i Ingress) Name() string {
	return "Ingress"
}

func (i Ingress) Policies() Policies {
	return Policies{}
}
