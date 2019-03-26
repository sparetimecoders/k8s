package util

import "gitlab.com/sparetimecoders/k8s-go/util/aws"

type Factory interface {
	Aws() aws.Service
}

type DefaultFactory struct{}

func NewFactory() Factory {
	return &DefaultFactory{}
}

func (c *DefaultFactory) Aws() aws.Service {
	return aws.New()
}
