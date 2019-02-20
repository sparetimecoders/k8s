package kops

import (
	"fmt"
	"log"
)

type Cluster struct {
	name string
	kops kops
}

func GetCluster(name string,stateStore string) Cluster {
	return Cluster{name, New(stateStore)}
}

func (c Cluster) CreateClusterResources() error {
	log.Printf("Creating cloud resources for %v", c.name)
	return c.kops.RunCmd(fmt.Sprintf("update cluster %v --yes", c.name), nil)
}
