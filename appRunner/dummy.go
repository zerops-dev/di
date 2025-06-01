package appRunner

import "context"

type DummyRegister struct {
	list []interface{}
}

func NewDummyRegister() *DummyRegister {
	return &DummyRegister{}
}

func (d *DummyRegister) Add(a interface{}) {
	d.list = append(d.list, a)
}

func (d *DummyRegister) RunBefore(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	return RunBefore(ctx, cancel, d.list...)
}

func (d *DummyRegister) RunAfter() error {
	return RunAfter(d.list...)
}
