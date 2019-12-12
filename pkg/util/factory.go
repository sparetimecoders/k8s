package util

import (
	aws2 "gitlab.com/sparetimecoders/k8s-go/pkg/util/aws"
	kops2 "gitlab.com/sparetimecoders/k8s-go/pkg/util/kops"
)

type Factory interface {
	Aws() aws2.Service
	Kops(clusterName string, stateStore string) kops2.Kops
}

type DefaultFactory struct{}

func NewFactory() Factory {
	return &DefaultFactory{}
}

func (c *DefaultFactory) Aws() aws2.Service {
	return aws2.New()
}

func (c *DefaultFactory) Kops(clusterName string, stateStore string) kops2.Kops {
	return kops2.New(clusterName, stateStore)
}
