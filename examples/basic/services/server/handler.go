package server

import (
	"context"
	"log/slog"

	"github.com/zerops-dev/di/appRunner"
)

type HandlerFunc func(context.Context, string) string

type Tools interface {
	GetTool() (string, HandlerFunc)
}

func RegisterTools(r *Handler, p Tools) {
	name, handler := p.GetTool()
	r.handlers[name] = handler
}

func New(
	config Config,
	log *slog.Logger,
	register appRunner.Register,
) (h *Handler, _ error) {
	defer func() { register.Add(h) }()
	h = &Handler{
		config:   config,
		log:      log,
		handlers: make(map[string]HandlerFunc),
	}
	return h, nil
}

type Handler struct {
	config   Config
	log      *slog.Logger
	handlers map[string]HandlerFunc
}

func (h *Handler) Run(ctx context.Context) error {
	for name, handler := range h.handlers {
		h.log.Debug("Running handler", "name", name, handler(ctx, "test"))
	}
	return nil
}
