package handlers

import (
	"context"
	"github.com/zerops-dev/di/examples/basic/services/server"
)

func NewReverse() *Reverse {
	return &Reverse{}
}

type Reverse struct{}

func (h *Reverse) GetTool() (string, server.HandlerFunc) {
	return "reverse", h.Handle
}

func (h *Reverse) Handle(_ context.Context, in string) string {
	return in
}
