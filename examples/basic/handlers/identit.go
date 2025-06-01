package handlers

import (
	"context"

	"github.com/zerops-dev/di/examples/basic/services/server"
)

func NewIdentity() *Identity {
	return &Identity{}
}

type Identity struct{}

func (h *Identity) GetTool() (string, server.HandlerFunc) {
	return "identity", h.Handle
}

func (h *Identity) Handle(_ context.Context, in string) string {
	return in
}
