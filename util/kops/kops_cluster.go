package kops

import (
	"encoding/json"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"log"
	"strings"
)

func (c Cluster) kopsClusterConfig() string {
	params := fmt.Sprintf(`get cluster %v -o yaml`, c.name)
	out, err := c.kops.Handler.QueryCmd(params, nil)
	if err != nil {
		log.Panicf("Failed to get clusterconfig %v", err)
		return ""
	}

	return string(out)
}

func policyString(instance string, policies []config.Policy) string {
	if len(policies) > 0 {
		jsonOut, err := json.Marshal(policies)
		if err != nil {
			log.Panicf("Failed to marshal policy for instance: %v, %v", instance, err)
		}
		return fmt.Sprintf("%v: '%v'", instance, string(jsonOut))
	} else {
		return ""
	}
}
func (c Cluster) SetIamPolicies(policies config.Policies) error {
	log.Println("Setting IAM policies for cluster")
	log.Printf("Master policies: %d, Node policies: %d\n ", len(policies.Master), len(policies.Node))
	kopsClusterConfig := c.kopsClusterConfig()
	node := policyString("node", policies.Node)
	master := policyString("master", policies.Master)

	replacement := fmt.Sprintf("spec:\n  additionalPolicies: \n    %v\n    %v", node, master)
	kopsClusterConfig = strings.Replace(kopsClusterConfig, "spec:", replacement, 1)

	params := fmt.Sprintf(`replace cluster %v -f -`, c.name)

	err := c.kops.Handler.RunCmd(params, []byte(kopsClusterConfig))

	if err != nil {
		log.Println("Failed to update cluster IAM policies")
		return err
	}
	log.Println("Updated IAM policies")

	return nil
}
