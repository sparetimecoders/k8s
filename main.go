package main

import (
	"bufio"
	"flag"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/aws"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/ingress"
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

	if c, err := loadConfig(args); err != nil {
		log.Fatal(err)
	} else {
		//kops.GetCluster("peter.sparetimecoders.com", "s3://k8s.sparetimecoders.com-kops-storage").WaitForValidState(180)
		//os.Exit(1)
		stateStore := getStateStore(c)
		awsSvc := aws.New(c.Region)

		if awsSvc.ClusterExist(c) {
			log.Fatalf("Cluster %v already exists", c.ClusterName())
		}
		k := kops.New(stateStore)
		cluster, err := k.CreateCluster(c)
		if err != nil {
			log.Fatal(err)
		}

		setNodeInstanceGroupToSpotPricesAndSize(cluster, c)
		setMasterInstanceGroupsToSpotPricesAndSize(cluster, c)
		cluster.CreateClusterResources()
		// Wait for completion/valid cluster...
		cluster.WaitForValidState(500)

		// Add-ons

		if (c.Ingress != ingress.Ingress{}) {
			fmt.Println("Ingress configured")
			c.Ingress.Create()
		}
	}
}

func getStateStore(c config.ClusterConfig) string {
	awsSvc := aws.New(c.Region)
	bucketName := awsSvc.StateStoreBucketName(c.DnsZone)
	if awsSvc.StateStoreBucketExist(c.DnsZone) {
		fmt.Printf("Using existing statestore: %v \n", bucketName)
	} else {
		fmt.Printf("No statestore S3 bucket found with name: %v \n", bucketName)
		fmt.Print("Continue and create statestore (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		if r, _, _ := reader.ReadRune(); r == 'y' || r == 'Y' {
			awsSvc.CreateStateStoreBucket(c.DnsZone)
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
	awsSvc := aws.New(region)
	price, err := awsSvc.OnDemandPrice(instanceType)
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
