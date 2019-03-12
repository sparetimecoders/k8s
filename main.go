package main

import (
	"bufio"
	"flag"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/aws"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/creator"
	"gitlab.com/sparetimecoders/k8s-go/kops"
	"log"
	"os"
)

var GitCommit, GitBranch, BuildDate, Version string

type args struct {
	Filename *string
	UseStdin bool
	Version  *bool
	_        struct{}
}

func main() {
	args := parseArgs()

	if *args.Version {
		fmt.Printf("Version: %s, GitCommit: %s, GitBranch: %s, BuildDate: %s\n", Version, GitCommit, GitBranch, BuildDate)
		os.Exit(0)
	} else if args.UseStdin == false && *args.Filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if clusterConfig, err := loadConfig(args); err != nil {
		log.Fatal(err)
	} else {
		stateStore := getStateStore(clusterConfig)
		awsSvc := aws.New()

		if awsSvc.ClusterExist(clusterConfig) {
			log.Fatalf("Cluster %v already exists", clusterConfig.ClusterName())
		}
		k := kops.New(stateStore)
		cluster, err := k.CreateCluster(clusterConfig)
		if err != nil {
			log.Fatal(err)
		}

		policies := config.Policies{Node: clusterConfig.Nodes.Policies}
		for _, p := range clusterConfig.Addons.List() {
			policies = config.Policies{Node: append(policies.Node, p.Policies().Node...), Master: append(policies.Master, p.Policies().Master...)}
		}
		// TODO Move code out of this main method...
		if len(policies.Master) > 0 || len(policies.Node) > 0 {
			if err := cluster.SetIamPolicies(policies); err != nil {
				log.Fatal(err)
			}
		}
		setNodeInstanceGroupToSpotPricesAndSize(cluster, clusterConfig)
		setMasterInstanceGroupsToSpotPricesAndSize(cluster, clusterConfig)
		cluster.CreateClusterResources()
		// Wait for completion/valid cluster...
		cluster.WaitForValidState(500)
		addons(clusterConfig)
	}
}

func addons(clusterConfig config.ClusterConfig) {
	addons := clusterConfig.Addons.List()
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

func getStateStore(c config.ClusterConfig) string {
	awsSvc := aws.New()
	bucketName := awsSvc.StateStoreBucketName(c.DnsZone)
	if awsSvc.StateStoreBucketExist(c.DnsZone) {
		fmt.Printf("Using existing statestore: %v \n", bucketName)
	} else {
		fmt.Printf("No statestore S3 bucket found with name: %v \n", bucketName)
		fmt.Print("Continue and create statestore (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		if r, _, _ := reader.ReadRune(); r == 'y' || r == 'Y' {
			awsSvc.CreateStateStoreBucket(c.DnsZone, c.Region)
		} else {
			log.Fatalln("Aborting...")
		}
	}

	return bucketName
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
		Filename: flag.String("f", "", "filename to load, use - for stdin"),
		Version:  flag.Bool("v", false, "print version info"),
	}
	flag.Parse()

	args.UseStdin = *args.Filename == "-"

	return args
}
