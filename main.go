package main

import (
  "flag"
  "fmt"
  "gitlab.com/sparetimecoders/k8s-go/aws"
  "gitlab.com/sparetimecoders/k8s-go/config"
  "gitlab.com/sparetimecoders/k8s-go/kops"
  "log"
)

type args struct {
  Filename *string
  UseStdin bool
}

func main() {
  args := parseArgs()

  // TODO Build statestore from config if not supplied?
	// TODO statestore part of config.ClusterConfig ?

	if c, err := loadConfig(args); err != nil {
      log.Fatal(err)
  } else {
    kops := kops.New("s3://k8s.sparetimecoders.com-kops-storage")
    cluster, err := kops.CreateCluster(c)
    if err != nil {
      log.Fatal(err)
    }

    //cluster := kops.GetCluster("peter.sparetimecoders.com", "s3://k8s.sparetimecoders.com-kops-storage")
    setNodeInstanceGroupToSpotPricesAndSize(cluster, c)
    setMasterInstanceGroupsToSpotPricesAndSize(cluster, c)
    cluster.CreateClusterResources()
  }
}

func setNodeInstanceGroupToSpotPricesAndSize(cluster kops.Cluster, config config.ClusterConfig) {
	price := instancePrice(config.Nodes.InstanceType, config.Region)
	setInstanceGroupToSpotPricesAndSize(cluster, "nodes", config.Nodes.Min, config.Nodes.Max, price, true)
}

func setMasterInstanceGroupsToSpotPricesAndSize(cluster kops.Cluster, config config.ClusterConfig) {
	for _, zone := range config.MasterZones {
		igName := fmt.Sprintf("master-%v%v", config.Region, zone)
		price := instancePrice(config.MasterInstanceType, config.Region)
		setInstanceGroupToSpotPricesAndSize(cluster, igName, 1, 1, price, false)
	}
}

func instancePrice(instanceType string, region string) float64 {
	awsSvc := aws.New()
	price, err := awsSvc.OnDemandPrice(instanceType, region)
	if err != nil {
		log.Fatalf("Failed to get price for instancetype, %v, %v", instanceType, err)
	}
	log.Printf("Got price %v for instancetype %v", price, instanceType)
	return price
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

func loadConfig(a args) (config.ClusterConfig, error) {
  if a.UseStdin {
    return config.ParseConfigStdin()
  } else {
    return config.ParseConfigFile(*a.Filename)
  }
}

func parseArgs() args {
  args := args{
    Filename: flag.String("f", "-", "-f <filename>"),
  }

  flag.Parse()

  args.UseStdin = *args.Filename == "-"

  return args
}
