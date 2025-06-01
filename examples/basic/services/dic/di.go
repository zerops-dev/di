package dic

import (
	"github.com/zerops-dev/di/appRunner"
	"github.com/zerops-dev/di/di/s"
	"github.com/zerops-dev/di/examples/basic/handlers"
	"github.com/zerops-dev/di/examples/basic/services/server"
	"github.com/zerops-dev/di/logger"
	_ "golang.org/x/tools/imports"
)

//go:generate templater -tags "di" -templateTags "!templater,!di"
var _ = func() *s.Di {
	di := s.NewDi("example-app",
		s.WithCommand(
			"run",
			"run command",
			``,
		),
	)

	di.Add(
		s.Scope(
			logger.DiScope(),
			appRunner.DiScope(),
			s.Service(server.New, s.WithSetter(server.RegisterTools)),

			s.Service(handlers.NewIdentity),
			s.Service(handlers.NewReverse),

			s.Config(server.NewConfig, "server"),
		),
	)

	return di
}
