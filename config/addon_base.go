package config

import (
	"strings"
)

func replace(org string, a map[string]string) (string, error) {
	copy := org
	for k, v := range a {
		copy = strings.ReplaceAll(copy, k, v)
	}
	return copy,nil
}
