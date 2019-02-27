package ingress

import (
	"github.com/GeertJohan/go.rice"
	"log"
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
}

// TODO Handle Azure and other cloud providers?
func (ingress Ingress) Create() {
	log.Println("Creating ingress from configuration")
	b, _ := readManifestFile()

	log.Println(b)
}

func readManifestFile() (string, error) {
	box := rice.MustFindBox("./manifests")
	s, err := box.String("ingress.yaml")
	if err != nil {
		return "", err
	}
	return s, nil
}
