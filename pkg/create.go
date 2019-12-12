package pkg

import (
	"errors"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/pkg/config"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util/kops"
	"io"
	"log"
	"time"
)

func Create(file string, f util.Factory, out io.Writer) error {
	if clusterConfig, err := config.Load(file); err != nil {
		return err
	} else {
		awsSvc := f.Aws()
		stateStore := awsSvc.GetStateStore(clusterConfig)

		if awsSvc.ClusterExist(clusterConfig) {
			return errors.New(fmt.Sprintf("Cluster %v already exists", clusterConfig.ClusterName()))
		}
		k := f.Kops(clusterConfig.ClusterName(), stateStore)
		cluster, err := k.CreateCluster(clusterConfig)
		if err != nil {
			return err
		}

		policies := policies(clusterConfig)

		if err := cluster.SetIamPolicies(policies); err != nil {
			return err
		}

		setNodeInstanceGroupToSpotPricesAndSize(awsSvc, cluster, clusterConfig)
		setMasterInstanceGroupsToSpotPricesAndSize(awsSvc, cluster, clusterConfig)
		_ = cluster.CreateClusterResources()
		// Wait for completion/valid cluster...
		waitTime := 5 * time.Minute
		if clusterConfig.DnsZone == "k8s.local" {
			waitTime = 10 * time.Minute
		}
		if cluster.WaitForValidState(int(waitTime.Seconds())) {
			addons(clusterConfig)
			return nil
		} else {
			return fmt.Errorf("failed to validate cluster in time, state unknown")
		}
	}
}

func setNodeInstanceGroupToSpotPricesAndSize(awsSvc aws.Service, cluster kops.Cluster, clusterConfig config.ClusterConfig) {
	var price float64
	if clusterConfig.Nodes.Spot {
		price = awsSvc.InstancePrice(clusterConfig.Nodes.InstanceType, clusterConfig.Region)
	}
	autoscaler := config.ClusterAutoscaler{}
	autoscale := clusterConfig.GetAddon(autoscaler) != nil

	setInstanceGroupToSpotPricesAndSize(cluster, "nodes", clusterConfig.Nodes.Min, clusterConfig.Nodes.Max, price, autoscale)
}

func setMasterInstanceGroupsToSpotPricesAndSize(awsSvc aws.Service, cluster kops.Cluster, clusterConfig config.ClusterConfig) {
	var price float64
	if clusterConfig.Masters.Spot {
		price = awsSvc.InstancePrice(clusterConfig.Masters.InstanceType, clusterConfig.Region)
	}
	for _, zone := range clusterConfig.Masters.Zones {
		igName := fmt.Sprintf("master-%v%v", clusterConfig.Region, zone)
		setInstanceGroupToSpotPricesAndSize(cluster, igName, 1, 1, price, false)
	}
}

func setInstanceGroupToSpotPricesAndSize(cluster kops.Cluster, igName string, min int, max int, price float64, autoScale bool) {
	group, err := cluster.GetInstanceGroup(igName)
	if err != nil {
		log.Fatalf("Failed to get instancegroup %v, %v", igName, err)
	}

	group = group.MinSize(min).MaxSize(max)
	if price > 0 {
		group = group.MaxPrice(price)
	}
	if autoScale {
		group = group.AutoScale()
	}
	err = cluster.UpdateInstanceGroup(group)
	if err != nil {
		log.Fatalf("Failed to update instancegroup %v, %v", igName, err)
	}
}
