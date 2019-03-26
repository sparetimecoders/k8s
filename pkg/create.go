package pkg

import (
	"errors"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"gitlab.com/sparetimecoders/k8s-go/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/util/creator"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
	"io"
	"log"
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
		k := f.Kops(stateStore)
		cluster, err := k.CreateCluster(clusterConfig)
		if err != nil {
			return err
		}

		policies := config.Policies{Node: clusterConfig.Nodes.Policies}
		for _, p := range clusterConfig.AllAddons() {
			policies = config.Policies{Node: append(policies.Node, p.Policies().Node...), Master: append(policies.Master, p.Policies().Master...)}
		}
		// TODO Move code out of this main method...
		if len(policies.Master) > 0 || len(policies.Node) > 0 {
			if err := cluster.SetIamPolicies(policies); err != nil {
				return err
			}
		}
		setNodeInstanceGroupToSpotPricesAndSize(awsSvc, cluster, clusterConfig)
		setMasterInstanceGroupsToSpotPricesAndSize(awsSvc, cluster, clusterConfig)
		_ = cluster.CreateClusterResources()
		// Wait for completion/valid cluster...
		cluster.WaitForValidState(500)
		addons(clusterConfig)
		return nil
	}
}

func addons(clusterConfig config.ClusterConfig) {
	addons := clusterConfig.AllAddons()
	if len(addons) == 0 {
		return
	}
	creator, err := creator.ForContext(clusterConfig.ClusterName())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating %d addon(s)\n", len(addons))

	for _, addon := range addons {
		log.Printf("Creating %s\n", addon.Name())
		s, err := addon.Manifests(clusterConfig)
		if err != nil {
			log.Fatal(err)
		}
		creator.Create(s)

		log.Printf("%s created\n", addon.Name())
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
