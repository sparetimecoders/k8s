package pkg

import (
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
	"log"
)

func Delete(file string) {
	if clusterConfig, err := config.Load(file); err != nil {
		log.Fatal(err)
	} else {
		stateStore := aws.GetStateStore(clusterConfig)
		awsSvc := aws.New()

		if !awsSvc.ClusterExist(clusterConfig) {
			log.Fatalf("Cluster %v does not exist", clusterConfig.ClusterName())
		}
		k := kops.New(stateStore)
		err := k.DeleteCluster(clusterConfig)
		if err != nil {
			log.Fatal(err)
		}
	}
}
