package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/pkg"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util"
	"io"
)

func NewCmdDelete(f util.Factory, out io.Writer) *cobra.Command {
	var file string

	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a K8S-cluster",
		Long:  `Delete an existing K8S-cluster based on the provided config-file`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := pkg.Delete(file, f); err != nil {
				_, _ = out.Write([]byte(err.Error()))
			}
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "config-file, use - for stdin (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
