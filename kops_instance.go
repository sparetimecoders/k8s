package main

import (
	"gopkg.in/yaml.v2"
	"time"
)

type InstanceGroup struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name              string            `yaml:"name"`
		CreationTimestamp time.Time         `yaml:"creationTimestamp"`
		Labels            map[string]string `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Role        string            `yaml:"role"`
		Image       string            `yaml:"image"`
		MinSize     int               `yaml:"minSize"`
		MaxSize     int               `yaml:"maxSize"`
		MachineType string            `yaml:"machineType"`
		Subnets     []string          `yaml:"subnets"`
		MaxPrice    string            `yaml:"maxPrice"`
		CloudLabels map[string]string `yaml:"cloudLabels"`
		NodeLabels  map[string]string `yaml:"nodeLabels"`
	} `yaml:spec`
}

func ParseInstanceGroup(data []byte) (InstanceGroup, error) {
	ig := InstanceGroup{}
	if err := yaml.UnmarshalStrict(data, &ig); err != nil {
		return InstanceGroup{}, err
	}
	return ig, nil
}
