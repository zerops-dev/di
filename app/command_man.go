package app

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func CreateManCmd(app *ApplicationSetup) *cobra.Command {

	var remove bool

	cmd := &cobra.Command{
		Use:   "man",
		Short: "generate man pages",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {

			pathName := path.Join("/usr/local/man/man1")
			header := &doc.GenManHeader{
				Title:   app.Name,
				Section: "1",
			}

			if remove {
				files, err := os.ReadDir(pathName)
				if nil != err {
					return err
				}
				for _, file := range files {
					if !strings.HasPrefix(file.Name(), app.Name) {
						continue
					}
					os.Remove(path.Join(pathName, file.Name()))
				}
				return nil
			}
			if err := os.MkdirAll(pathName, 0755); err != nil {
				return err
			}
			return doc.GenManTree(cmd.Root(), header, pathName)

		},
	}
	cmd.Flags().BoolVar(&remove, "remove", false, "remove man pages")
	return cmd
}
