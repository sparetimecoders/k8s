package create

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"strings"
)

// minimal k8s YAML descriptor
type k8sYaml struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

type creator struct {
	urlRoot string
	_       struct{}
}

func FromHostPort(host string, port int) creator {
	return New("https", host, port)
}

func FromHost(host string) creator {
	return FromHostPort(host, 443)
}
func New(protocol string, host string, port int) creator {
	return creator{urlRoot: fmt.Sprintf("%s://%s:%d", protocol, host, port)}
}

// Create tries to create the resources descripted in the passed YAML content
func (c creator) Create(yamlContent string) {
	parts := getParts(yamlContent)
	log.Printf("Found %d parts in yaml content", len(parts))
	for _, part := range parts {
		c.post(part)
	}
}

// getParts returns the different YAML parts in a string
// The triple-dash '---' is used as the separator
func getParts(yamlContent string) []string {
	scanner := bufio.NewScanner(strings.NewReader(yamlContent))
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

// buildUrl creates a k8s api url from the YAML descriptor passed
func (c creator) buildUrl(yamlDescriptor string) (string, error) {
	yamlData := k8sYaml{}

	if err := yaml.UnmarshalStrict([]byte(yamlDescriptor), &yamlData); err != nil {
		return "", err
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
	url = fmt.Sprintf("%s%s", c.urlRoot, part)
	return url, nil
}

// TODO Error and results
func (c creator) post(yamlContent string) {
	url, err := c.buildUrl(yamlContent)
	if err != nil {
		fmt.Print(err)
	}
	log.Printf("Posting to url: %s\ns", url)
	res, err := http.Post(url, "application/yaml", strings.NewReader(yamlContent))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(res)
}
