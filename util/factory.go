package util

import (
	"gitlab.com/sparetimecoders/k8s-go/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
)

type Factory interface {
	Aws() aws.Service
	Kops(stateStore string) kops.Kops
}

type DefaultFactory struct{}

func NewFactory() Factory {
	return &DefaultFactory{}
}

func (c *DefaultFactory) Aws() aws.Service {
	return aws.New()
}

func (c *DefaultFactory) Kops(stateStore string) kops.Kops {
	return kops.New(stateStore)
}
