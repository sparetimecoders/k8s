package pkg

import (
	"errors"
	"fmt"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"gitlab.com/sparetimecoders/k8s-go/util/creator"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
	"io"
	"log"
)

func Addons(file string, f util.Factory, out io.Writer) error {
	if clusterConfig, err := config.Load(file); err != nil {
		return err
	} else {
		awsSvc := f.Aws()
		stateStore := awsSvc.GetStateStore(clusterConfig)

		if !awsSvc.ClusterExist(clusterConfig) {
			return errors.New(fmt.Sprintf("Cluster %v does not exist", clusterConfig.ClusterName()))
		}

		cluster := kops.GetCluster(clusterConfig.ClusterName(), stateStore)
		cluster.WaitForValidState(500)
		if err := cluster.SetIamPolicies(policies(clusterConfig)); err != nil {
			return err
		}

		addons(clusterConfig)
	}
	return nil
}

func addons(clusterConfig config.ClusterConfig) {
	log.Printf("Addons for cluster %s", clusterConfig.ClusterName())

	addons := clusterConfig.AllAddons()
	if len(addons) == 0 {
		return
	}
	creator:=creator.New(clusterConfig.ClusterName())
/*	creator, err := creator.ForContext(clusterConfig.ClusterName())
	if err != nil {
		log.Fatal(err)
	}
*/
	log.Printf("Creating %d addon(s)\n", len(addons))

	for _, addon := range addons {
		log.Printf("Creating %s\n", addon.Name())
		s, err := addon.Manifests(clusterConfig)
		if err != nil {
			log.Fatal(err)
		}
		creator.Apply(s)
		//creator.Create(s)

		log.Printf("%s created\n", addon.Name())
	}
}
