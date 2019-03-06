package addons

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/sparetimecoders/k8s-go/addons/ingress"
	"testing"
)

func TestAddons_List(t *testing.T) {
	assert.Equal(t, 0, len(Addons{}.List()))
	assert.Equal(t, 1, len(Addons{Ingress: &ingress.Ingress{}}.List()))
	addons := Addons{Ingress: &ingress.Ingress{}}
	for _, addon := range addons.List() {
		res, _ := addon.Content()
		assert.Equal(t, "PETER", res)
	}
}
