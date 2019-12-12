package pkg

import (
	"gitlab.com/sparetimecoders/k8s-go/pkg/config"
)

func policies(clusterConfig config.ClusterConfig) config.Policies {
	policies := defaults(clusterConfig)
	for _, p := range clusterConfig.AllAddons() {
		policies = config.Policies{Node: append(policies.Node, p.Policies().Node...), Master: append(policies.Master, p.Policies().Master...)}
	}
	return policies
}

func defaults(clusterConfig config.ClusterConfig) config.Policies {
	Policies := clusterConfig.Nodes.Policies
	if Policies == nil {
		Policies = []config.Policy{}
	}
	return config.Policies{Node: Policies, Master: []config.Policy{}}

}
