// +build !prod

package util

import "gitlab.com/sparetimecoders/k8s-go/util/aws"

type MockFactory struct{}

func NewMockFactory() Factory {
	return &MockFactory{}
}

func (c *MockFactory) Aws() aws.Service {
	return aws.MockService{}
}
