package cmd

import (
	"github.com/spf13/cobra"
	runApp "github.com/zerops-dev/di/app"
	"github.com/zerops-dev/di/examples/basic/services/dic"
)

func AddCommands(applicationSetup *runApp.ApplicationSetup, rootCommand *cobra.Command) {
	rootCommand.AddCommand(dic.CreateCommand(applicationSetup))
}
