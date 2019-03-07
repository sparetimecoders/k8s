package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsService struct {
	_      struct{}
}

func New() awsService {
	return awsService{}
}

func (awsSvc awsService) awsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{}))
}
