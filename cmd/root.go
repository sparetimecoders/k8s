package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/sparetimecoders/k8s-go/util"
	"gitlab.com/sparetimecoders/k8s-go/util/aws"
	"gitlab.com/sparetimecoders/k8s-go/util/kops"
	"io"
	"os"
)

type RootCmd struct {
	factory util.Factory

	cobraCommand *cobra.Command
}

var _ util.Factory = &RootCmd{}

var rootCommand = RootCmd{}

func NewCmdRoot(f util.Factory, out io.Writer) *cobra.Command {
	rootCommand.cobraCommand = &cobra.Command{
		Use:   "k8s-go",
		Short: "k8s-go yadda",
		Long:  `k8s-go yadda yadda`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(0)
		},
	}
	cmd := rootCommand.cobraCommand

	cmd.AddCommand(NewCmdVersion(f, out))
	cmd.AddCommand(NewCmdCreate(f, out))
	cmd.AddCommand(NewCmdDelete(f, out))
	cmd.AddCommand(NewCmdAddons(f, out))

	return cmd
}

func Execute() {
	factory := util.NewFactory()
	rootCommand.factory = factory

	NewCmdRoot(factory, os.Stdout)

	if err := rootCommand.cobraCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *RootCmd) Aws() aws.Service {
	return c.factory.Aws()
}

func (c *RootCmd) Kops(clusterName string, stateStore string) kops.Kops {
	return c.factory.Kops(clusterName, stateStore)
}
