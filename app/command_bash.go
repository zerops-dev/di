package app

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

func CreateBashCmd(app *ApplicationSetup) *cobra.Command {

	var save bool
	var remove bool

	cmd := &cobra.Command{
		Use:   "bash",
		Short: "generating bash completions",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			rootCmd := cmd.Root()
			if remove {
				return os.Remove(path.Join("/etc/bash_completion.d", app.Name+".sh"))
			}
			if save {
				return rootCmd.GenBashCompletionFile(path.Join("/etc/bash_completion.d", app.Name+".sh"))
			}
			return cmd.Root().GenBashCompletion(os.Stdout)
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "Save bash file")
	cmd.Flags().BoolVar(&remove, "remove", false, "Remove bash file")

	return cmd
}
