package logger

import (
	"log/slog"
	"os"
)

func New(config Config) *slog.Logger {
	var output slog.Handler
	switch config.Output {
	case OutputJson:
		output = slog.NewJSONHandler(os.Stdout, nil)
	case OutputDiscard:
		output = slog.DiscardHandler
	default:
		output = slog.NewTextHandler(os.Stdout, nil)
	}
	return slog.New(output)
}
