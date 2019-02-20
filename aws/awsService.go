package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
)

type awsService struct {
	region string
}

func New(region string) awsService {
	// TODO Credentials...
	return awsService{region}
}
func (awsSvc awsService )awsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsSvc.region),
	}))
}


func (awsSvc awsService) StateStoreBucketName(dns string) string {
	sess := awsSvc.awsSession()
	identity,_ :=sts.New(sess).GetCallerIdentity(&sts.GetCallerIdentityInput{})
	return fmt.Sprintf("%v-%v-kops-storage", *identity.Account,dns)
}

func (awsSvc awsService) CreateStateStoreBucket(dns string) {
	sess := awsSvc.awsSession()

	s3Svc := s3.New(sess)
	bucketName := awsSvc.StateStoreBucketName(dns)

	s3Svc.CreateBucket(&s3.CreateBucketInput{Bucket: &bucketName, CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: &awsSvc.region}})

}

func (awsSvc awsService) StateStoreBucketExist(dns string) bool {
	sess := awsSvc.awsSession()

	s3Svc := s3.New(sess)
	result, err :=s3Svc.ListBuckets(&s3.ListBucketsInput{})
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


	/*
	cluster:s3_bucket(){
  local bucket_name="${1}"
  local region=${2}
  exists=$(aws s3api list-buckets \
      --query "Buckets[].Name" \
      | jq ". | index(\"${bucket_name}\") != null"
  )
  if [[ "${exists}" == "false" ]]; then
    printf "${RED}Creating S3 bucket ${bucket_name}${NC}\n"
    aws s3api create-bucket \
        --bucket ${bucket_name} \
        --region ${region} \
        --create-bucket-configuration LocationConstraint=${region}

    aws s3api put-bucket-versioning \
      --bucket ${bucket_name} \
      --versioning-configuration Status=Enabled \
      --region ${region}
  else
    printf "Using existing S3 bucket ${GREEN}${bucket_name}${NC}\n"
  fi
}
	 */

}
