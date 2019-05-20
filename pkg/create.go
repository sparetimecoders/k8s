package pkg

import (
	"errors"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"gitlab.com/sparetimecoders/k8s-go/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
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
	price := awsSvc.InstancePrice(clusterConfig.Nodes.InstanceType, clusterConfig.Region)
	autoscaler := config.ClusterAutoscaler{}
	autoscale := clusterConfig.GetAddon(autoscaler) != nil

	setInstanceGroupToSpotPricesAndSize(cluster, "nodes", clusterConfig.Nodes.Min, clusterConfig.Nodes.Max, price, autoscale)
}

func setMasterInstanceGroupsToSpotPricesAndSize(awsSvc aws.Service, cluster kops.Cluster, config config.ClusterConfig) {
	for _, zone := range config.MasterZones {
		igName := fmt.Sprintf("master-%v%v", config.Region, zone)
		price := awsSvc.InstancePrice(config.MasterInstanceType, config.Region)
		setInstanceGroupToSpotPricesAndSize(cluster, igName, 1, 1, price, false)
	}
}

func setInstanceGroupToSpotPricesAndSize(cluster kops.Cluster, igName string, min int, max int, price float64, autoScale bool) {
	group, err := cluster.GetInstanceGroup(igName)
	if err != nil {
		log.Fatalf("Failed to get instancegroup %v, %v", igName, err)
	}

	group = group.MinSize(min).MaxSize(max).MaxPrice(price)
	if autoScale {
		group = group.AutoScale()
	}
	err = cluster.UpdateInstanceGroup(group)
	if err != nil {
		log.Fatalf("Failed to update instancegroup %v, %v", igName, err)
	}
}
