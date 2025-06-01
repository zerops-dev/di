package runState

import (
	"sync"

	"errors"
)

const (
	stateRunning uint32 = iota
	stateStopped
)

var (
	stoppedErr = errors.New("state stopped")
)

type State struct {
	lock    sync.Mutex
	state   uint32
	counter sync.WaitGroup
}

func NewState() *State {
	return &State{}
}

func NewStoppedState() *State {
	return &State{
		state: stateStopped,
	}
}

func (r *State) Stop() {
	func() {
		r.lock.Lock()
		defer r.lock.Unlock()
		r.state = stateStopped
	}()
	r.counter.Wait()
}

func (r *State) Start() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.state = stateRunning
}

func (r *State) Process() (func(), error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	var once sync.Once
	if r.state == stateStopped {
		return nil, stoppedErr
	}
	r.counter.Add(1)
	return func() {
		once.Do(func() {
			r.counter.Done()
		})
	}, nil
}
