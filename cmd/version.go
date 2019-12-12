package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/pkg/util"
	"io"
)

var GitCommit, GitBranch, BuildDate, Version string = "sha", "branch", "now", "version"

func NewCmdVersion(f util.Factory, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of k8s-go",
		Long:  `All software has versions. This is k8s-go's'`,
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintf(out, "Version: %s, GitCommit: %s, GitBranch: %s, BuildDate: %s\n", Version, GitCommit, GitBranch, BuildDate)
		},
	}
}
