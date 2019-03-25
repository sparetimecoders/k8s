package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/aws"
	"gitlab.com/sparetimecoders/k8s-go/config"
	"gitlab.com/sparetimecoders/k8s-go/kops"
	"log"
)

var deleteFile string

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteFile, "file", "f", "", "config-file, use - for stdin (required)")
	_ = deleteCmd.MarkFlagRequired("file")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a K8S-cluster",
	Long:  `Delete an existing K8S-cluster based on the provided config-file`,
	Run: func(cmd *cobra.Command, args []string) {
		if clusterConfig, err := config.Load(deleteFile); err != nil {
			log.Fatal(err)
		} else {
			stateStore := getStateStore(clusterConfig)
			awsSvc := aws.New()

			if !awsSvc.ClusterExist(clusterConfig) {
				log.Fatalf("Cluster %v does not exist", clusterConfig.ClusterName())
			}
			k := kops.New(stateStore)
			err := k.DeleteCluster(clusterConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
