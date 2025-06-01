package app

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func CreateDocCmd() *cobra.Command {

	var directory string

	cmd := &cobra.Command{
		Use:   "doc",
		Short: "generate doc",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			rootCmd := cmd.Root()
			if err := os.MkdirAll(directory, os.ModePerm); err != nil {
				return err
			}
			return doc.GenMarkdownTree(rootCmd, directory)

		},
	}
	cmd.Flags().StringVar(&directory, "directory", "doc", "output directory")
	return cmd
}
