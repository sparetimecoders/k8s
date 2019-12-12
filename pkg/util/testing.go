// +build !prod

package util

import (
	aws2 "gitlab.com/sparetimecoders/k8s-go/pkg/util/aws"
	kops2 "gitlab.com/sparetimecoders/k8s-go/pkg/util/kops"
)

type MockFactory struct {
	ClusterExists bool
	Handler       kops2.MockHandler
}

func NewMockFactory() *MockFactory {
	factory := &MockFactory{}
	factory.Handler = kops2.MockHandler{
		Cmds:      make(chan string, 100),
		Responses: make(chan string, 100),
	}
	return factory
}

func (c *MockFactory) Aws() aws2.Service {
	return aws2.MockService{
		ExistingCluster: c.ClusterExists,
	}
}

func (c *MockFactory) Kops(clusterName string, stateStore string) kops2.Kops {
	return kops2.NewMock(clusterName, c.Handler)
}
