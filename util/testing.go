// +build !prod

package util

import (
	"gitlab.com/sparetimecoders/k8s-go/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
)

type MockFactory struct {
	ClusterExists bool
	Handler       kops.MockHandler
}

func NewMockFactory() *MockFactory {
	factory := &MockFactory{}
	factory.Handler = kops.MockHandler{
		Cmds: make(chan string, 100),
	}
	return factory
}

func (c *MockFactory) Aws() aws.Service {
	return aws.MockService{
		ExistingCluster: c.ClusterExists,
	}
}

func (c *MockFactory) Kops(stateStore string) kops.Kops {
	k := kops.New(stateStore)
	k.Handler = c.Handler
	return k
}
