package config

import (
	"reflect"
	"strings"
)

type Addons struct {
	Ingress     *Ingress     `yaml:"ingress"`
	ExternalDNS *ExternalDNS `yaml:"externalDns"`
	_           struct{}
}

func (addons Addons) List() []Addon {
	var result []Addon
	a := reflect.TypeOf(addons)
	value := reflect.ValueOf(&addons).Elem()
	for i := 0; i < a.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() != reflect.Struct && !field.IsNil() {
			field.Interface()
			result = append(result, field.Interface().(Addon))
		}
	}
	return result
}

type Addon interface {
	Manifests(config ClusterConfig) (string, error)
	Name() string
	Policies() Policies
	//Validate(config config.ClusterConfig) bool

}

func replace(org string, a map[string]string) (string, error) {
	copy := org
	for k, v := range a {
		copy = strings.ReplaceAll(copy, k, v)
	}
	return copy, nil
}
