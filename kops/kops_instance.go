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

func (k kops) UpdateInstanceGroup(group InstanceGroup) error {

	params := strings.TrimSpace(fmt.Sprintf(`replace ig %v --name %v --state %v -f -`, group.Metadata.Name, group.Metadata.Labels["kops.k8s.io/cluster"], k.stateStore))

	cmd := exec.Command(k.cmd, strings.Split(params, " ")...)
	data, err := yaml.Marshal(group)
	if err != nil {
		log.Println("Failed to convert to yaml")
		return err
	}
	cmd.Stdin = bytes.NewBuffer(data)
	err = cmd.Start()
	if err != nil {
		log.Printf("Failed to update instancegroup %v\n %v\n", group.Metadata.Name, err)
		return err
	}

	_ = cmd.Wait()
	return nil
}

func (k kops) GetInstanceGroup(name string, clusterName string) (InstanceGroup, error) {
	params := strings.TrimSpace(fmt.Sprintf(`get ig %v --name %v --state %v -o yaml`, name, clusterName, k.stateStore))

	cmd := exec.Command(k.cmd, strings.Split(params, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return InstanceGroup{}, err
	}
	_ = cmd.Start()
	_ = cmd.Wait()

	return parseInstanceGroup(out)
}

func parseInstanceGroup(data []byte) (InstanceGroup, error) {
	ig := InstanceGroup{}
	if err := yaml.UnmarshalStrict(data, &ig); err != nil {
		return InstanceGroup{}, err
	}
	return ig, nil
}
