package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pricing"
	"strconv"
)

func (awsSvc awsService) OnDemandPrice(instanceType string, region string) (float64, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		// API Endpoint for price must be eu-east-1
		Region: aws.String("us-east-1"),
	}))

	pricingSvc := pricing.New(sess)

	input := priceForInstanceRequest(instanceType, regionLocation(region))
	out, err := pricingSvc.GetProducts(input)
	if err != nil {
		return -1, err
	}
	return findPriceInUSD(out.PriceList[0]["terms"].(map[string]interface{})["OnDemand"].(map[string]interface{}))
}

func regionLocation(region string) string {
	var regions = map[string]string{
		"eu-west-1": "EU (Ireland)",
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
	return strconv.ParseFloat(price.(map[string]interface{})["USD"].(string), 64)
}

func findKey(key string, m map[string]interface{}) (interface{}, error) {
	for k, v := range m {
		if k == key {
			return m[k], nil
		} else {
			switch v.(type) {
			case map[string]interface{}:
				return findKey(key, m[k].(map[string]interface{}))
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Failed to findKey key: %v", key))
}