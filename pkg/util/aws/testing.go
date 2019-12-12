// +build !prod

package aws

import (
	"gitlab.com/sparetimecoders/k8s-go/pkg/config"
)

type MockService struct {
	ExistingCluster bool
	_               struct{}
}

func (awsSvc MockService) ClusterExist(config config.ClusterConfig) bool {
	return awsSvc.ExistingCluster
}

func (awsSvc MockService) InstancePrice(instanceType string, region string) float64 {
	return 0.7
}

func (awsSvc MockService) GetStateStore(config config.ClusterConfig) string {
	return "s3://dummy.state.store"
}
