// +build !prod

package aws

import (
	"gitlab.com/sparetimecoders/k8s-go/config"
)

type MockService struct {
	_ struct{}
}

func (awsSvc MockService) ClusterExist(config config.ClusterConfig) bool {
	return false
}

func (awsSvc MockService) OnDemandPrice(instanceType string, region string) (float64, error) {
	return 0.03, nil
}

func (awsSvc MockService) GetStateStore(config config.ClusterConfig) string {
	return "s3://dummy.state.store"
}
