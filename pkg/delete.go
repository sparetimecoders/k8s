package pkg

import (
	"errors"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/util"
)

func Delete(file string, f util.Factory) error {
	if clusterConfig, err := config.Load(file); err != nil {
		return err
	} else {
		awsSvc := f.Aws()

		if !awsSvc.ClusterExist(clusterConfig) {
			return errors.New(fmt.Sprintf("Cluster %v does not exist", clusterConfig.ClusterName()))
		}
		stateStore := awsSvc.GetStateStore(clusterConfig)
		k := f.Kops(clusterConfig.ClusterName(), stateStore)
		return k.DeleteCluster(clusterConfig)
	}
}
