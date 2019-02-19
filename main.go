package main

import (
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/kops"
	"gopkg.in/yaml.v2"
	"log"
)

func main() {
	// TODO Build statestore from config if not supplied?
	// TODO statestore part of config.ClusterConfig ?

	kops := kops.New("s3://k8s.sparetimecoders.com-kops-storage")
	if parsed, err := yaml.Marshal(config.ClusterConfig{Name: "peter", DnsZone: "sparetimecoders.com"}); err == nil {
		c, err := config.ParseConfig(parsed)
		if err != nil {
			log.Fatal(err)
		}
		cluster, err := kops.CreateCluster(c)
		if err != nil {
			log.Fatal(err)
		}

		nodes, err := cluster.GetInstanceGroup("nodes")
		if err != nil {
			log.Fatal(err)
		}

		if e := cluster.UpdateInstanceGroup(nodes.MaxPrice(10.1).MaxSize(3).MinSize(1)); e != nil {
			log.Fatal(e)
		}
	}

}
