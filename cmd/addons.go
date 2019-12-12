package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/pkg"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util"
	"io"
)

func NewCmdAddons(f util.Factory, out io.Writer) *cobra.Command {
	var file string

	var cmd = &cobra.Command{
		Use:   "addons",
		Short: "(Re-)Applies addons to a K8S-cluster",
		Long:  `(Re-)Applies addons based on the provided config-file`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := pkg.Addons(file, f, out); err != nil {
				_, _ = out.Write([]byte(err.Error()))
			}
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "config-file, use - for stdin (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
