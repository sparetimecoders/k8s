package config

import (
	"github.com/GeertJohan/go.rice"
	"strings"
)

/**
ingress:create() {
  local cert_arn="${1}"
  printf "Creating ${BLUE}nginx-ingress-controller${NC}\n"
  local ingress_service_name="ingress-nginx"
  ${KUBECTL_CMD} apply -f "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/manifests/" &>/dev/null
  until ${KUBECTL_CMD} get service ${ingress_service_name} --namespace ${ingress_service_name} &>/dev/null ; do date; sleep 1; echo ""; done

  if [[ -n "${cert_arn}" ]]; then
    ${KUBECTL_CMD} annotate service \
        --overwrite \
        --namespace ${ingress_service_name} \
        ${ingress_service_name} \
         "service.beta.kubernetes.io/aws-load-balancer-ssl-cert"="${cert_arn}" \
         "service.beta.kubernetes.io/aws-load-balancer-backend-protocol"="http" \
         "service.beta.kubernetes.io/aws-load-balancer-ssl-ports"="https" \
         "service.beta.kubernetes.io/aws-load-balancer-ssl-negotiation-policy"="ELBSecurityPolicy-TLS-1-2-2017-01" &>/dev/null
  fi
  printf "Created ${BLUE}nginx-ingress-controller${NC}\n"

}
*/

type Ingress struct {
	AwsCertificate struct {
		AwsSecurityPolicy string `yaml:"awsSecurityPolicy" default:"ELBSecurityPolicy-TLS-1-2-2017-01"`
		AwsCertificateARN string `yaml:"awsCertificateARN" default:""`
	} `yaml:"awsCertificate"`
	_ struct{}
}

func (i Ingress) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/ingress")
	return strings.Join([]string{box.MustString("ingress.yaml"), box.MustString("nginx-config.yaml")}, "\n---\n"), nil
}

func (i Ingress) Name() string {
	return "Ingress"
}

func (i Ingress) Policies() Policies {
	return Policies{}
}
