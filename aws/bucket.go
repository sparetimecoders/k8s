package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"log"
)

func (awsSvc awsService) StateStoreBucketName(dns string) string {
	sess := awsSvc.awsSession()
	identity, _ := sts.New(sess).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	return fmt.Sprintf("%v-%v-kops-storage", *identity.Account, dns)
}

func (awsSvc awsService) CreateStateStoreBucket(dns string, region string) {
	sess := awsSvc.awsSession()
	s3Svc := s3.New(sess)
	bucketName := awsSvc.StateStoreBucketName(dns)
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

func (awsSvc awsService) StateStoreBucketExist(dns string) bool {
	sess := awsSvc.awsSession()

	s3Svc := s3.New(sess)
	result, err := s3Svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println(err)
	}
	bucketName := awsSvc.StateStoreBucketName(dns)
	for _, b := range result.Buckets {
		if *b.Name == bucketName {
			return true
		}
	}

	return false
}

func (awsSvc awsService) ClusterExist(config config.ClusterConfig) bool {
	sess := awsSvc.awsSession()

	s3Svc := s3.New(sess)
	buckets, _ := s3Svc.ListBuckets(&s3.ListBucketsInput{})

	for _, b := range buckets.Buckets {
		fmt.Println(b)
	}
	stateStoreBucketName := awsSvc.StateStoreBucketName(config.DnsZone)
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