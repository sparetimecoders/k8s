package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"gitlab.com/sparetimecoders/k8s-go/pkg/config"
	"log"
)

func (awsSvc defaultAwsService) stateStoreBucketName(region string) string {
	sess := awsSvc.awsSession(region)
	identity, _ := sts.New(sess).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	return fmt.Sprintf("%v-%v-kops-storage", *identity.Account, region)
}

func (awsSvc defaultAwsService) createStateStoreBucket(region string) {
	sess := awsSvc.awsSession(region)
	s3Svc := s3.New(sess)
	bucketName := awsSvc.stateStoreBucketName(region)
	log.Printf("Creating statestore %v", bucketName)

	if _, err := s3Svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName, CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: &region}}); err == nil {
		if _, err := s3Svc.PutBucketVersioning(&s3.PutBucketVersioningInput{
			Bucket: aws.String(bucketName),
			VersioningConfiguration: &s3.VersioningConfiguration{
				Status: aws.String("Enabled"),
			},
		}); err != nil {
			log.Fatalf("Failed to set versioning for statestore %v, %v", bucketName, err)
		}
	} else {
		log.Fatalf("Failed to create statestore %v, %v", bucketName, err)
	}
}

func (awsSvc defaultAwsService) stateStoreBucketExist(region string) bool {
	sess := awsSvc.awsSession(region)
	s3Svc := s3.New(sess)
	result, err := s3Svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println(err)
	}
	bucketName := awsSvc.stateStoreBucketName(region)
	for _, b := range result.Buckets {
		if *b.Name == bucketName {
			return true
		}
	}

	return false
}

func (awsSvc defaultAwsService) ClusterExist(config config.ClusterConfig) bool {
	sess := awsSvc.awsSession(config.Region)

	s3Svc := s3.New(sess)

	stateStoreBucketName := awsSvc.stateStoreBucketName(config.Region)
	if list, err := s3Svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(stateStoreBucketName),
		Prefix: aws.String(config.ClusterName()),
	}); err != nil {
		log.Fatalf("Could not list statestore %v, %v", stateStoreBucketName, err)
	} else {
		if len(list.Contents) > 0 {
			return true
		}
	}
	return false
}
