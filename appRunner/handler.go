package appRunner

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/wire"
)

var RegisterSet = wire.NewSet(
	wire.Bind(new(Register), new(Handler)),
	New,
)

type runSignal chan struct{}

func (r runSignal) Run(_ context.Context) error {
	close(r)
	return nil
}

type Register interface {
	Add(interface{})
}

type Handler struct {
	list    []interface{}
	running runSignal
	log     *slog.Logger
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		log:     log,
		running: make(runSignal),
	}
}

func (s *Handler) Running() chan struct{} {
	return s.running
}

func (s *Handler) Add(item interface{}) {
	if s != nil {
		s.list = append(s.list, item)
	}
}

func (s *Handler) RunWithSigTerm(ctx context.Context, cancel context.CancelFunc) error {
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		cancel()
	}()
	return s.Run(ctx, cancel)
}

func (s *Handler) RunOnce(
	ctx context.Context,
	cancel context.CancelFunc,
	runFunc func() error,
) error {
	return s.RunOnceWithContext(
		ctx,
		cancel,
		func(_ context.Context) error {
			return runFunc()
		},
	)
}

func (s *Handler) RunOnceWithContext(
	ctx context.Context,
	cancel context.CancelFunc,
	runFunc func(context.Context) error,
) error {
	if err := RunBefore(ctx, cancel, s.list...); err != nil {
		return err
	}
	if err := runFunc(ctx); err != nil {
		return err
	}
	if err := RunAfter(s.list...); err != nil {
		return err
	}
	return nil
}

type RunnerHandler func(ctx context.Context) error

func (h RunnerHandler) Run(ctx context.Context) error {
	return h(ctx)
}

func (s Handler) RunFunc(
	ctx context.Context,
	cancel context.CancelFunc,
	runFunc func() error,
) (returnErr error) {
	if err := Runner(ctx, s.log, cancel, append(append(s.list, RunnerHandler(func(ctx context.Context) error {
		defer cancel()
		if err := runFunc(); err != nil {
			return fmt.Errorf("type: %T, %v", runFunc, err)
		}
		return nil
	})), s.running)...); err != nil {
		return err
	}
	return returnErr
}

func (s Handler) RunBefore(ctx context.Context, cancel context.CancelFunc) error {
	return RunBefore(ctx, cancel, s.list...)
}

func (s Handler) RunAfter() error {
	return RunAfter(s.list...)
}

func (s Handler) Run(ctx context.Context, cancel context.CancelFunc) error {
	return Runner(ctx, s.log, cancel, append(s.list, s.running)...)
}
