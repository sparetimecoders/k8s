package kops

import (
	"fmt"
	"log"
	"time"
)

type Cluster struct {
	kops Kops
	_    struct{}
}

func GetCluster(name string, stateStore string) Cluster {
	return Cluster{kops: New(name, stateStore)}
}

func (c Cluster) CreateClusterResources() error {
	return c.kops.UpdateCluster()
}

func (c Cluster) WaitForValidState(maxWaitSeconds int) bool {
	log.Printf("Validating cluster, will wait max %v seconds\n", maxWaitSeconds)
	fmt.Printf("Validating cluster, will wait max %v seconds\n", maxWaitSeconds)
	endTime := time.Now().Add(time.Second * time.Duration(maxWaitSeconds))
	done := false
	out := ""
	for time.Now().Before(endTime) {
		out, done = c.checkValidState()
		if done {
			log.Println("Cluster up and running")
			return true
		} else {
			time.Sleep(5 * time.Second)
			fmt.Printf(".")
		}
	}
	log.Printf("Failed to validate cluster in time, %v\n", out)
	return false
}

func (c Cluster) checkValidState() (string, bool) {
	return c.kops.ValidateCluster()
}
