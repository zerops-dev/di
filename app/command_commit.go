package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CreateCommitCmd(app *ApplicationSetup) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "application build commit",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(os.Stdout, app.Commit)
			return nil
		},
	}
	return cmd
}
