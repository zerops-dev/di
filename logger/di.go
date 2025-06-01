package logger

import "github.com/zerops-dev/di/di/s"

func DiScope() *s.XScope {
	return s.Scope(
		s.Service(New, s.WithAppInject("logger")),
		s.Config(NewConfig, "log"),
	)
}
