package main

import (
	"context"
	"path"
	"strings"

	runApp "github.com/zerops-dev/di/app"
	"github.com/zerops-dev/di/examples/basic/services/cmd"
)

var (
	Name            = "example-app"
	Service         = "example-app"
	Exec            = "/usr/local/bin/example-app"
	BuildTime       = "time"
	Commit          = "commit"
	Version         = "v0.0.0"
	Description     = "Example application"
	DescriptionLong = "Example application"
)

func main() {

	applicationSetup := runApp.New(context.Background(), Name)
	applicationSetup.Exec = strings.Join([]string{Exec, "run", "--config", path.Join("/etc", Service, "config.yml")}, " ")
	applicationSetup.Version = Version
	applicationSetup.Commit = Commit
	applicationSetup.Service = Service
	applicationSetup.Description = Description
	applicationSetup.DescriptionLong = DescriptionLong
	applicationSetup.SetBuildTime(BuildTime)

	applicationSetup.RegisterCommands()
	cmd.AddCommands(applicationSetup, applicationSetup.RootCommand())
	applicationSetup.Run()
}
