package config

import (
	"github.com/GeertJohan/go.rice"
)

type ClusterAutoscaler struct{}

func (i ClusterAutoscaler) Manifests(clusterConfig ClusterConfig) (string, error) {
	box := rice.MustFindBox("manifests/cluster_autoscaler")
	manifest := box.MustString("cluster_autoscaler.yaml")

	return replace(manifest, map[string]string{"$region": clusterConfig.Region, "$cluster_name": clusterConfig.ClusterName()})
}

func (i ClusterAutoscaler) Name() string {
	return "Cluster ClusterAutoscaler"
}

func (i ClusterAutoscaler) Policies() Policies {
	return Policies{Node: []Policy{
		{Actions: []string{"autoscaling:DescribeAutoScalingGroups",
			"autoscaling:DescribeAutoScalingInstances",
			"autoscaling:SetDesiredCapacity",
			"autoscaling:TerminateInstanceInAutoScalingGroup",
			"autoscaling:DescribeTags",
		}, Effect: "Allow", Resources: []string{"*"}},
	}}
}
