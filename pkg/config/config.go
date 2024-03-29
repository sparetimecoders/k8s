package config

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Policy struct {
	Actions   []string `yaml:"actions" json:"Action"`
	Effect    string   `yaml:"effect" json:"Effect"`
	Resources []string `yaml:"resources" json:"Resource"`
}

type Nodes struct {
	Min          int      `yaml:"min" default:"1"`
	Max          int      `yaml:"max" default:"2"`
	Zones        []string `yaml:"zones" default:"a,b,c"`
	Spot         bool     `yaml:"spot" default:"false"`
	InstanceType string   `yaml:"type" default:"t3.medium"`
	Policies     []Policy `yaml:"policies" optional:"true"`
}
type MasterNodes struct {
	Zones        []string `yaml:"zones" default:"a" type:"quorum"`
	Spot         bool     `yaml:"spot" default:"false"`
	InstanceType string   `yaml:"type" default:"t3.small"`
	Policies     []Policy `yaml:"policies" optional:"true"`
}

type Policies struct {
	Node   []Policy `yaml:"node"`
	Master []Policy `yaml:"master"`
	_      struct{}
}

const LocalCluster = "k8s.local"

type ClusterConfig struct {
	Name              string            `yaml:"name"`
	KubernetesVersion string            `yaml:"kubernetesVersion" default:"1.15.5"`
	DnsZone           string            `yaml:"dnsZone" default:"k8s.local"`
	Region            string            `yaml:"region" default:"eu-west-1"`
	Vpc               string            `yaml:"vpc" optional:"true"`
	NetworkCIDR       string            `yaml:"networkCIDR" default:"172.21.0.0/22"`
	Masters           MasterNodes       `yaml:"masters"`
	Nodes             Nodes             `yaml:"nodes"`
	CloudLabels       map[string]string `yaml:"cloudLabels" default:""`
	SshKeyPath        string            `yaml:"sshKeyPath" default:"~/.ssh/id_rsa.pub"`
	Addons            *Addons           `yaml:"addons" optional:"true"`
}

func (p Policies) Exists() bool {
	return len(p.Master) > 0 || len(p.Node) > 0
}

func (config ClusterConfig) ClusterName() string {
	return fmt.Sprintf("%s.%s", config.Name, config.DnsZone)
}

func (config ClusterConfig) AllAddons() []Addon {
	var result []Addon
	if config.Addons != nil {
		a := reflect.TypeOf(*config.Addons)
		value := reflect.ValueOf(config.Addons).Elem()
		for i := 0; i < a.NumField(); i++ {
			field := value.Field(i)
			if field.Kind() != reflect.Struct && !field.IsNil() {
				field.Interface()
				result = append(result, field.Interface().(Addon))
			}
		}
	}
	return result
}

func (config ClusterConfig) GetAddon(t Addon) Addon {
	for _, addon := range config.AllAddons() {
		x := t.Name()
		y := addon.Name()
		if x == y {
			return addon
		}
	}
	return nil
}

func Load(file string) (ClusterConfig, error) {
	if file == "-" {
		return ParseConfigStdin()
	} else {
		return ParseConfigFile(file)
	}
}

func ParseConfigFile(file string) (ClusterConfig, error) {
	if r, err := os.Open(file); err != nil {
		return ClusterConfig{}, err
	} else {
		return parseConfig(r)
	}
}

func ParseConfigStdin() (ClusterConfig, error) {
	return parseConfig(bufio.NewReader(os.Stdin))
}

func parseConfig(r io.Reader) (ClusterConfig, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return ClusterConfig{}, err
	}
	return parseConfigData(data)

}

func parseConfigData(content []byte) (ClusterConfig, error) {
	config := ClusterConfig{}
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
		log.Printf("Successfully parsed configuration for %s", config.ClusterName())
		return config, nil
	}
}

func handleDefaultValues(t reflect.Value, missingFields *[]string, prefix string) error {
	refType := t.Type()
	for i := 0; i < refType.NumField(); i++ {
		name := strings.TrimPrefix(fmt.Sprintf("%s.%s", prefix, refType.Field(i).Name), ".")
		value := t.Field(i)
		defaultValue := refType.Field(i).Tag.Get("default")
		mandatory := refType.Field(i).Tag.Get("optional") != "true"
		quorum := refType.Field(i).Tag.Get("type") == "quorum"
		if value.Kind() == reflect.Struct {
			if err := handleDefaultValues(value, missingFields, name); err != nil {
				return err
			}
		} else if value.Kind() == reflect.Ptr && !value.IsNil() {
			if err := handleDefaultValues(value.Elem(), missingFields, name); err != nil {
				return err
			}
		} else if isZeroOfUnderlyingType(value) && mandatory {
			if defaultValue == "" {
				*missingFields = append(*missingFields, name)
			} else {
				log.Printf("Setting default value for field '%s' = '%s'", name, defaultValue)
				if err := set(value, name, defaultValue, missingFields); err != nil {
					return err
				}
			}
		}
		if quorum {
			if !(value.Len() == 1 || value.Len() == 3 || value.Len() == 5) {
				return fmt.Errorf("%d is not a valid quorum size, expected 1,3,5 for %s ", value.Len(), name)
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
