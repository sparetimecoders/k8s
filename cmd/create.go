package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/pkg"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"io"
)

func NewCmdCreate(f util.Factory, out io.Writer) *cobra.Command {
	var file string

	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a K8S-cluster",
		Long:  `Create a new K8S-cluster based on the provided config-file`,
		Run: func(cmd *cobra.Command, args []string) {
			pkg.Create(file, f, out)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "config-file, use - for stdin (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
