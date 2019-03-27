package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pricing"
	"log"
	"strconv"
)

func onDemandPrice(instanceType string, region string) (float64, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		// API Endpoint for price must be eu-east-1
		Region: aws.String("us-east-1"),
	}))

	pricingSvc := pricing.New(sess)
	location := regionLocation(region)
	input := priceForInstanceRequest(instanceType, location)
	if location == "" {
		log.Fatalf("No region matching %s\n", region)
	}
	out, err := pricingSvc.GetProducts(input)
	if err != nil {
		return -1, err
	}
	return findPriceInUSD(out.PriceList[0]["terms"].(map[string]interface{})["OnDemand"].(map[string]interface{}))
}

func (awsSvc defaultAwsService) InstancePrice(instanceType string, region string) float64 {
	price, err := onDemandPrice(instanceType, region)
	if err != nil {
		log.Fatalf("Failed to get price for instancetype, %v, %v", instanceType, err)
	}
	log.Printf("Got price %v for instancetype %v", price, instanceType)
	return price
}

func regionLocation(region string) string {
	var regions = map[string]string{
		"us-east-1":    "US East (N. Virginia)",
		"us-east-2":    "US East (Ohio)",
		"us-west-1":    "US West (N. California)",
		"us-west-2":    "US West (Oregon)",
		"ca-central-1": "Canada (Central)",
		"eu-central-1": "EU (Frankfurt)",
		"eu-west-1":    "EU (Ireland)",
		"eu-west-2":    "EU (London)",
		"eu-west-3":    "EU (Paris)",
		"eu-north-1":   "EU (Stockholm)",
	}

	return regions[region]
}

func priceForInstanceRequest(instanceType string, region string) *pricing.GetProductsInput {

	return &pricing.GetProductsInput{
		Filters: []*pricing.Filter{
			{
				Field: aws.String("instanceType"),
				Type:  aws.String(pricing.FilterTypeTermMatch),
				Value: aws.String(instanceType),
			}, {
				Field: aws.String("location"),
				Type:  aws.String(pricing.FilterTypeTermMatch),
				Value: aws.String(region),
			}, {
				Field: aws.String("operatingSystem"),
				Type:  aws.String(pricing.FilterTypeTermMatch),
				Value: aws.String("Linux"),
			}, {
				Field: aws.String("tenancy"),
				Type:  aws.String(pricing.FilterTypeTermMatch),
				Value: aws.String("Shared"),
			}, {
				Field: aws.String("capacitystatus"),
				Type:  aws.String(pricing.FilterTypeTermMatch),
				Value: aws.String("Used"),
			}, {
				Field: aws.String("preInstalledSw"),
				Type:  aws.String(pricing.FilterTypeTermMatch),
				Value: aws.String("NA"),
			},
		},
		ServiceCode: aws.String("AmazonEC2"),
	}
}

func findPriceInUSD(m map[string]interface{}) (float64, error) {
	price, err := findKey("pricePerUnit", m)
	if err != nil {
		return -1, err
	}
	usdPrice := price.(map[string]interface{})["USD"].(string)
	if usdPrice == "" {
		return -1, errors.New("empty string for price")
	}
	return strconv.ParseFloat(usdPrice, 64)
}

func findKey(key string, m map[string]interface{}) (interface{}, error) {
	found := recursiveFind(key, m)
	if found != nil {
		return found, nil
	}
	return nil, errors.New(fmt.Sprintf("Failed to findKey key: %v", key))
}

func recursiveFind(key string, m map[string]interface{}) interface{} {
	for k, v := range m {
		if k == key {
			return m[k]
		} else {
			switch v.(type) {
			case map[string]interface{}:
				found := recursiveFind(key, m[k].(map[string]interface{}))
				if found != nil {
					return found
				}
			}
		}
	}
	return nil
}
