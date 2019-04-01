package kops

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var versions = []struct {
	input    string
	expected string
}{

	{"Version 1.12.0-alpha.3 (git-08aa9ac7f)", "1.12.0-alpha.3"},
	{"Version 1.11.1", "1.11.1"},
}

func TestKopsVersionString(t *testing.T) {
	r := make(chan string, 10)
	handler := MockHandler{Cmds: make(chan string, 10), Responses: r}
	k := kops{Handler: handler}

	for _, tt := range versions {
		r <- tt.input
		actual, _ := k.Version()
		assert.Equal(t, tt.expected, actual)
	}
}

var minVersions = []struct {
	actual  string
	minimum string
	ok      bool
}{

	{"Version 1.12.0-alpha.3 (git-08aa9ac7f)", "1.12.7", true},
	{"Version 1.11.1", "1.11.1", true},
	{"Version 1.10.10", "1.11.1", false},
}

func TestKopsMinimumKopsVersionInstalled(t *testing.T) {
	r := make(chan string, 10)
	handler := MockHandler{Cmds: make(chan string, 10), Responses: r}
	k :=  kops{Handler: handler}

	for _, tt := range minVersions {
		r <- tt.actual
		actual := k.MinimumKopsVersionInstalled(tt.minimum)
		assert.Equal(t, tt.ok, actual)
	}

}
