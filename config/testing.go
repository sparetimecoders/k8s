// +build !prod

package config

func ParseConfigData(content []byte) (ClusterConfig, error) {
	return parseConfigData(content)
}
