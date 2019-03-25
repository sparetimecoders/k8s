package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/pkg"
)

var createFile string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "config-file, use - for stdin (required)")
	_ = createCmd.MarkFlagRequired("file")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a K8S-cluster",
	Long:  `Create a new K8S-cluster based on the provided config-file`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Create(createFile)
	},
}
