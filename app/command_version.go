package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CreateVersionCmd(app *ApplicationSetup) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "application version",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(os.Stdout, app.Version)
			return nil

		},
	}
	return cmd
}
