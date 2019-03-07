package config

type ExternalDNS struct {
	ManifestLoader
}

func (e ExternalDNS) Name() string {
	return "ExternalDNS"
}

func (e ExternalDNS) Content(config ClusterConfig) (string, error) {
	return e.Load("external_dns", "external_dns.yaml")
}
