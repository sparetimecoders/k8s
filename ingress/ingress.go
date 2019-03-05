package ingress

import (
	"bufio"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
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

type basicYaml struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

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

func getParts(s string) []string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	var yamls, current []string
	for scanner.Scan() {
		if line := scanner.Text(); line == "---" {
			yamls = append(yamls, strings.Join(current, "\n"))
			current = nil
		} else {
			current = append(current, line)
		}
	}
	if len(current) > 0 {
		yamls = append(yamls, strings.Join(current, "\n"))
	}
	return yamls
}

func buildUrl(yamlContent string) string {
	host := "localhost"
	port := 8080
	urlString := "http://%s:%d%s"
	yamlData := basicYaml{}

	if err := yaml.UnmarshalStrict([]byte(yamlContent), &yamlData); err != nil {

	}
	var url, part string
	if yamlData.Kind == "Namespace" {
		part = "/api/v1/namespaces"
	} else {
		namespace := yamlData.Metadata.Namespace
		if namespace == "" {
			namespace = "default"
		}
		switch yamlData.ApiVersion {
		case "v1":
			// Core
			part = fmt.Sprintf("/api/%s/namespaces/%s/%ss", yamlData.ApiVersion, namespace, strings.ToLower(yamlData.Kind))

		default:
			part = fmt.Sprintf("/apis/%s/namespaces/%s/%ss", yamlData.ApiVersion, namespace, strings.ToLower(yamlData.Kind))
		}

	}
	url = fmt.Sprintf(urlString, host, port, part)
	return url
}

// TODO Error and results
func post(yamlContent string) {
	url := buildUrl(yamlContent)
	res, err := http.Post(url, "application/yaml", strings.NewReader(yamlContent))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(res)
}
