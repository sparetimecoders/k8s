package main

import (
	"errors"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/kops"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	data, err := ioutil.ReadFile("./config.yaml")

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	config, err2 := ParseConfig(data)

	if err2 != nil {
		log.Fatalf("error: %v", err2)
	}

	fmt.Printf("%v\n", config)
	kops := kops.New("#")
	_ = kops.CreateCluster(config)
}

func ParseConfig(content []byte) (config.Cluster, error) {
	config := config.Cluster{}
	if err := yaml.UnmarshalStrict(content, &config); err != nil {
		return config, err
	} else {
		var missingFields []string

		if err = handleDefaultValues(reflect.ValueOf(&config).Elem(), &missingFields, ""); err != nil {
			panic(err)
		}

		if len(missingFields) != 0 {
			return config, errors.New(fmt.Sprintf("Missing required value for field(s): '%v'\n", missingFields))
		}
		return config, nil
	}
}

func handleDefaultValues(t reflect.Value, missingFields *[]string, prefix string) error {
	refType := t.Type()
	for i := 0; i < refType.NumField(); i++ {
		name := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, refType.Field(i).Name), ".")
		value := t.Field(i)
		defaultValue := refType.Field(i).Tag.Get("default")

		if isZeroOfUnderlyingType(value) {
			if value.Kind() == reflect.Struct {
				if err := set(value, name, defaultValue, missingFields); err != nil {
					return err
				}
			} else if defaultValue == "" {
				*missingFields = append(*missingFields, name)
			} else {
				fmt.Printf("Setting default value for field '%s' = '%s'\n", name, defaultValue)
				if err := set(value, name, defaultValue, missingFields); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func isZeroOfUnderlyingType(x reflect.Value) bool {
	return reflect.DeepEqual(x.Interface(), reflect.Zero(x.Type()).Interface())
}

func set(field reflect.Value, name string, value string, missingFields *[]string) error {
	switch field.Kind() {
	case reflect.Slice:
		arr := strings.Split(value, ",")
		field.Set(reflect.ValueOf(arr))
	case reflect.Struct:
		s := reflect.New(field.Type())
		if err := handleDefaultValues(s.Elem(), missingFields, name); err != nil {
			return err
		}
		field.Set(s.Elem())
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	}
	return nil
}
