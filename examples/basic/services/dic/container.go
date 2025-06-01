//go:build wireinject
// +build wireinject

package dic

import (
	"github.com/google/wire"

	"github.com/zerops-dev/di/app"
)

//go:generate wire
func CreateApp(*app.ApplicationSetup, Config) (*App, error) {
	wire.Build(
		DepSet,
	)

	return &App{}, nil
}
