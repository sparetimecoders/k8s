package addons

import (
	"gitlab.com/sparetimecoders/k8s-go/addons/ingress"
	"reflect"
)

type Addon interface {
	Content() (string, error)
	Name() string
}

type Addons struct {
	Ingress *ingress.Ingress `yaml:"ingress"`
	_       struct{}
}

func (addons Addons) List() []Addon {
	var result []Addon
	a := reflect.TypeOf(addons)
	value := reflect.ValueOf(&addons).Elem()
	for i := 0; i < a.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() != reflect.Struct && !field.IsNil() {
			result = append(result, field.Interface().(Addon))
		}
	}
	return result
}
