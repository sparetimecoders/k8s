package config

import (
	"reflect"
	"strings"
)

type Addons struct {
	Ingress           *Ingress           `yaml:"ingress" optional:"true"`
	ExternalDNS       *ExternalDNS       `yaml:"externalDns" optional:"true"`
	ClusterAutoscaler *ClusterAutoscaler `yaml:"clusterAutoscaler" optional:"true"`
}

func (addons Addons) AllAddons() []Addon {
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

func (addons Addons) GetAddon(t Addon) Addon {
	for _, addon := range addons.AllAddons() {
		x := t.Name()
		y := addon.Name()
		if x == y {
			return addon
		}
	}
	return nil
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
