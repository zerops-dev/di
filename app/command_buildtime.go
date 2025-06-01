package app

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func CreateBuildTimeCmd(app *ApplicationSetup) *cobra.Command {

	var format string

	cmd := &cobra.Command{
		Use:   "buildtime",
		Short: "application build time",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(os.Stdout, app.BuildTime.Format(format))
			return nil

		},
	}

	cmd.Flags().StringVar(&format, "format", time.RFC3339, "timestampFormat to use for display when a full timestamp is printed ("+time.RFC3339+")")
	return cmd
}
