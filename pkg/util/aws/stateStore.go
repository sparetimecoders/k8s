package aws

import (
	"bufio"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/pkg/config"
	"log"
	"os"
)

func (awsSvc defaultAwsService) GetStateStore(config config.ClusterConfig) string {
	bucketName := awsSvc.stateStoreBucketName(config.Region)
	if awsSvc.stateStoreBucketExist(config.Region) {
		fmt.Printf("Using existing statestore: %v \n", bucketName)
	} else {
		fmt.Printf("No statestore S3 bucket found with name: %v \n", bucketName)
		fmt.Print("Continue and create statestore (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		if r, _, _ := reader.ReadRune(); r == 'y' || r == 'Y' {
			awsSvc.createStateStoreBucket(config.Region)
		} else {
			log.Fatalln("Aborting...")
		}
	}

	return bucketName
}
