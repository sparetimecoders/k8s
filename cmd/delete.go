package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/pkg"
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
		pkg.Delete(deleteFile)
	},
}
