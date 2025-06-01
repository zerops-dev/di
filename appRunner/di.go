package appRunner

import (
	"github.com/zerops-dev/di/di/s"
)

func DiBind(i *Handler) Register {
	return i
}

func DiScope() *s.XScope {
	return s.Scope(
		s.Service(New, s.WithAppInject("runner")),
		s.Service(DiBind),
	)
}

func DiScopeNamed(appName string) *s.XScope {
	return s.Scope(
		s.Service(New, s.WithAppInject(appName)),
		s.Service(DiBind),
	)
}
