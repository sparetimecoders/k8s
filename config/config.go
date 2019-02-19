package config

type Nodes struct {
	Min          int    `yaml:"min" default:"1"`
	Max          int    `yaml:"max" default:"2"`
	InstanceType string `yaml:"instanceType" default:"t3.medium"`
}

type Cluster struct {
	Name               string            `yaml:"name"`
	KubernetesVersion  string            `yaml:"kubernetesVersion" default:"1.11.7"`
	DnsZone            string            `yaml:"dnsZone"`
	Region             string            `yaml:"region" default:"eu-west-1"`
	MasterZones        []string          `yaml:"masterZones" default:"a"`
	NetworkCIDR        string            `yaml:"networkCIDR" default:"172.21.0.0/22"`
	Nodes              Nodes             `yaml:"nodes"`
	MasterInstanceType string            `yaml:"masterInstanceType" default:"t3.small"`
	CloudLabels        map[string]string `yaml:"cloudLabels" default:""`
	SshKeyPath         string            `yaml:"sshKeyPath" default:"~/.ssh/id_rsa.pub"`
}
