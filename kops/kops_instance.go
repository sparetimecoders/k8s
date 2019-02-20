package kops

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os/exec"
	"strings"
	"time"
)

type InstanceGroup struct {
	ig instanceGroup
}

type instanceGroup struct {
	Kind       string `yaml:"kind"`
	ApiVersion string `yaml:"apiVersion"`
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
	} `yaml:"spec"`
}

func (ig InstanceGroup) MaxSize(n int) InstanceGroup {
	ig.ig.Spec.MaxSize = n
	return ig
}

func (ig InstanceGroup) MinSize(n int) InstanceGroup {
	ig.ig.Spec.MinSize = n
	return ig
}

func (ig InstanceGroup) MaxPrice(price float64) InstanceGroup {
	ig.ig.Spec.MaxPrice = fmt.Sprintf("%.4f", price)
	return ig
}

func (ig InstanceGroup) AutoScale() InstanceGroup {
	if ig.ig.Spec.CloudLabels == nil {
		ig.ig.Spec.CloudLabels = make(map[string]string)
	}
	ig.ig.Spec.CloudLabels["k8s.io/cluster-autoscaler/enabled"] = "true"
	ig.ig.Spec.CloudLabels[fmt.Sprintf("k8s.io/cluster-autoscaler/%v", ig.ig.Metadata.Labels["kops.k8s.io/cluster"])] = "true"
	return ig
}

func (c Cluster) UpdateInstanceGroup(group InstanceGroup) error {
	log.Printf("Updating instance group %v\n", group.ig.Metadata.Name)
	params := strings.TrimSpace(fmt.Sprintf(`replace ig %v --name %v --state %v -f -`, group.ig.Metadata.Name, c.name, c.kops.stateStore))

	cmd := exec.Command(c.kops.cmd, strings.Split(params, " ")...)

	data, err := yaml.Marshal(group.ig)

	if err != nil {
		log.Println("Failed to convert to yaml")
		return err
	}
	cmd.Stdin = bytes.NewBuffer(data)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("Failed to update instancegroup %v\n %v\n", group.ig.Metadata.Name, string(out))
		return err
	}
	log.Printf("Updated instance group %v\n", group.ig.Metadata.Name)
	return nil
}

func (c Cluster) GetInstanceGroup(name string) (InstanceGroup, error) {
	params := strings.TrimSpace(fmt.Sprintf(`get ig %v --name %v --state %v -o yaml`, name, c.name, c.kops.stateStore))

	cmd := exec.Command(c.kops.cmd, strings.Split(params, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return InstanceGroup{}, err
	}
	return parseInstanceGroup(out)
}

func parseInstanceGroup(data []byte) (InstanceGroup, error) {
	ig := instanceGroup{}
	if err := yaml.UnmarshalStrict(data, &ig); err != nil {
		return InstanceGroup{}, err
	}
	return InstanceGroup{ig}, nil
}
