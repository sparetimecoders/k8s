package main

import (
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/kops"
	"log"
)

func main() {
	config, err := config.ParseConfigFile("./config.yaml")

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", config)
	// TODO Build statestore from config if not supplied?
	// TODO statestore part of config.Cluster ?

	kops := kops.New("s3://k8s.sparetimecoders.com-kops-storage")
	_ = kops.CreateCluster(config)
}
