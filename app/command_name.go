package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CreateNameCmd(app *ApplicationSetup) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "name",
		Short: "application name",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(os.Stdout, app.Name)
			return nil

		},
	}
	return cmd
}
