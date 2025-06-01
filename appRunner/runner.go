package appRunner

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sync"
)

type Executor struct {
	BeforeDecorator
	RunDecorator
	AfterDecorator
}

type BeforeDecorator func(context.Context) error

func (d BeforeDecorator) Before(ctx context.Context) error {
	if d == nil {
		return nil
	}
	return d(ctx)
}

type RunDecorator func(context.Context, chan error)

func (d RunDecorator) Run(ctx context.Context, done chan error) {
	if d == nil {
		close(done)
		return
	}
	d(ctx, done)
}

type AfterDecorator func() error

func (d AfterDecorator) After() error {
	if d == nil {
		return nil
	}
	return d()
}

type Before interface {
	Before(context.Context) error
}

type Run interface {
	Run(context.Context) error
}

type After interface {
	After() error
}

type wrapper struct {
	*runBeforeWrapper
	*runWrapper
}

type runBeforeWrapper struct {
	beforeRun func(context.Context) error
}

type runWrapper struct {
	run func(context.Context) error
}

func BeforeWrapper(beforeRun func(context.Context) error) *runBeforeWrapper {
	return &runBeforeWrapper{
		beforeRun: beforeRun,
	}
}

func RunWrapper(run func(context.Context) error) *runWrapper {
	return &runWrapper{
		run: run,
	}
}

func Wrapper(beforeRun func(context.Context) error, run func(context.Context) error) *wrapper {
	return &wrapper{
		runBeforeWrapper: BeforeWrapper(beforeRun),
		runWrapper:       RunWrapper(run),
	}
}

func (s *runBeforeWrapper) RunBefore(ctx context.Context) error {
	return s.beforeRun(ctx)
}

func (s *runWrapper) Run(ctx context.Context) error {
	return s.run(ctx)
}

func before(ctx context.Context, cancel context.CancelFunc, runners ...interface{}) (resultRunners []interface{}, err error) {
	for _, r := range runners {
		if before, isBefore := r.(Before); isBefore {
			if err = before.Before(ctx); err != nil {
				beforeType := reflect.TypeOf(before)
				x := beforeType.String()
				err = fmt.Errorf("beforeRunner %s, %v", x, err)
				cancel()
				return
			}
		}
		resultRunners = append(resultRunners, r)
	}
	return
}

func RunBefore(ctx context.Context, cancel context.CancelFunc, runners ...interface{}) error {
	beforeRunners, err := before(ctx, cancel, runners...)
	if err != nil {
		afterRunError := RunAfter(beforeRunners...)
		if afterRunError != nil {
			return fmt.Errorf("afterRunErr: %v, %v", afterRunError, err)
		}
		return err
	}
	return nil
}

func RunAfter(runners ...interface{}) (err error) {
	for i := len(runners) - 1; i >= 0; i-- {
		if after, isAfter := runners[i].(After); isAfter {
			if afterErr := after.After(); afterErr != nil {
				afterType := reflect.TypeOf(after)
				afterErr = fmt.Errorf("afterRunner %s, %v", afterType.String(), afterErr)
				if err == nil {
					err = afterErr
				}
			}
		}
	}
	return
}

func Runner(ctx context.Context, log *slog.Logger, cancel context.CancelFunc, runners ...interface{}) error {

	log.Info("run before - start")
	if err := RunBefore(ctx, cancel, runners...); err != nil {
		log.With(err).Error("runBefore error")
		return err
	}
	log.Info("run before - end")

	var wg sync.WaitGroup

	once := sync.Once{}
	var runError error

	log.Info("run - start")
	// call lastCtxBreaker when ctx is canceled and all previous dependencies ended
	lastCtxBreaker := func() {}
	for _, r := range runners {
		if run, isRunner := r.(Run); isRunner {
			runType := reflect.TypeOf(run)
			runTypeLog := log.With("running", runType.String())
			runTypeLog.Info("run")
			depCtx, nextCtxCancel := context.WithCancel(context.Background())
			nextCtxBreaker := func() {
				runTypeLog.Info("calling cancel")
				nextCtxCancel()
			}
			wg.Add(1)
			done := make(chan error)

			go func() {
				defer close(done)
				if err := run.Run(depCtx); err != nil {
					done <- err
				}
			}()

			go func(depCtx context.Context, breakDepCtx context.CancelFunc) {
				defer func() { // when Run is done
					<-ctx.Done()    // and ctx is canceled
					<-depCtx.Done() // and all dependencies ended
					breakDepCtx()   // cascade ending to previously called function
				}()
				err, opened := <-done
				runTypeLog.Info("finished")
				if err != nil {
					runTypeLog.With(err).Error("finishedWithError")
				}
				if opened && err != nil {
					typedRunError := fmt.Errorf("runner %s %v", runType.String(), err)
					fmt.Fprint(os.Stderr, typedRunError.Error())
					once.Do(func() {
						runError = typedRunError
					})
					cancel()
				}
				wg.Done()
			}(depCtx, lastCtxBreaker)

			lastCtxBreaker = nextCtxBreaker
		}
	}

	// start context done cascade
	go func() {
		<-ctx.Done()
		lastCtxBreaker()
	}()

	<-ctx.Done()
	wg.Wait()

	log.Info("run - end")

	log.Info("run after - start")
	err := RunAfter(runners...)
	log.Info("run after - end")
	if err != nil {
		log.With(err).Error("runAfter error")
		return err
	}

	return runError
}
