package creator

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// minimal k8s YAML descriptor
type k8sMinimalYaml struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

// Response from K8s API
type k8sResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type creator struct {
	client    http.Client
	serverUrl string
	_         struct{}
}

func ForContext(context string) (creator, error) {
	cfg, err := kubeConfigForContext(context)
	if err != nil {
		log.Fatalf("Failed to get kube config, %s", err)
	}
	return fromConfig(*cfg)
}

func fromConfig(c rest.Config) (creator, error) {
	config := rest.CopyConfig(&c)
	config.GroupVersion = &schema.GroupVersion{}
	config.AcceptContentTypes = "application/json"
	config.ContentType = "application/json"
	codec := runtime.NoopEncoder{Decoder: scheme.Codecs.UniversalDecoder()}
	config.NegotiatedSerializer = serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: codec})
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	serverUrl, _, err := defaultServerUrlFor(config)
	if err != nil {
		log.Fatalf("Failed to build server url, %s", err)
	}

	client, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatalf("Failed to build http client, %s", err)
	}
	return creator{client: *client.Client, serverUrl: serverUrl.String()}, nil

}

// Create tries to create the resources descripted in the passed YAML content
func (c creator) Create(yamlContent string) {
	parts := getParts(yamlContent)
	log.Printf("Found %d parts in yaml content", len(parts))
	for _, part := range parts {
		c.post(part)
	}
}

// Copied from client-go

// defaultServerUrlFor is shared between IsConfigTransportTLS and RESTClientFor. It
// requires Host and Version to be set prior to being called.
func defaultServerUrlFor(config *rest.Config) (*url.URL, string, error) {
	// TODO: move the default to secure when the apiserver supports TLS by default
	// config.Insecure is taken to mean "I want HTTPS but don't bother checking the certs against a CA."
	hasCA := len(config.CAFile) != 0 || len(config.CAData) != 0
	hasCert := len(config.CertFile) != 0 || len(config.CertData) != 0
	defaultTLS := hasCA || hasCert || config.Insecure
	host := config.Host
	if host == "" {
		host = "localhost"
	}

	return rest.DefaultServerURL(host, config.APIPath, schema.GroupVersion{}, defaultTLS)
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
	yamlData := k8sMinimalYaml{}

	if err := yaml.Unmarshal([]byte(yamlDescriptor), &yamlData); err != nil {
		return "", err
	}

	var resultingUrl, part string
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
	resultingUrl = fmt.Sprintf("%s%s", c.serverUrl, part)
	return resultingUrl, nil
}

// TODO Error and results
func (c creator) post(yamlContent string) {
	url, err := c.buildUrl(yamlContent)
	if err != nil {
		log.Fatalln(err)

	}
	log.Printf("Posting to url: %s\n", url)
	res, err := c.client.Post(url, "application/yaml", strings.NewReader(yamlContent))
	if err != nil {
		log.Fatalf("Failed to create resource, %s\n", err)
	}
	if res.StatusCode != 201 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Failed to read response, %s\n", err)
		}
		r := k8sResponse{}
		err = yaml.Unmarshal(b, &r)
		if err != nil {
			log.Fatalf("Failed to read response as yaml, %s\n%s\n", string(b), err)
		}
		log.Fatalf("Failed to crete resource, %s\n", r.Message)
	}
	log.Println("Resource created")

}

func kubeConfigForContext(context string) (*rest.Config, error) {
	var kubeconfig *string
	if e := os.Getenv("KUBECONFIG"); e != "" {
		kubeconfig = &e
	} else if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// use the current context in kubeconfig
	loadingRules := clientcmd.ClientConfigLoadingRules{
		Precedence: strings.Split(*kubeconfig, ":"),
	}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&loadingRules, &clientcmd.ConfigOverrides{CurrentContext: context})
	return kubeConfig.ClientConfig()

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
