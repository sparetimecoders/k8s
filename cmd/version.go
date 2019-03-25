package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var GitCommit, GitBranch, BuildDate, Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of k8s-go",
	Long:  `All software has versions. This is k8s-go's'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s, GitCommit: %s, GitBranch: %s, BuildDate: %s\n", Version, GitCommit, GitBranch, BuildDate)
	},
}
