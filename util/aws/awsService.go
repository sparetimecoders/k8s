package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/sparetimecoders/k8s-go/config"
)

type Service interface {
	ClusterExist(config config.ClusterConfig) bool
	GetStateStore(config config.ClusterConfig) string
	InstancePrice(instanceType string, region string) float64
}

type defaultAwsService struct {
	_ struct{}
}

func New() Service {
	return defaultAwsService{}
}

func (awsSvc defaultAwsService) awsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{}))
}
