package vflags

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Value struct {
	Name         string
	Description  string
	Configurable bool
	Valuer       Valuer
	Virtual      bool
}

func (v Value) String() string {
	return v.Valuer.String()
}

func (v Value) Set(in string) error {
	if v.Virtual {
		return nil
	}
	if err := v.Valuer.Set(in); err != nil {
		return fmt.Errorf("%s: %v", v.Name, err)
	}
	return nil
}

func (v Value) SetValue(in interface{}) (err error) {
	if v.Virtual {
		return nil
	}
	if err := v.Valuer.Set(in); err != nil {
		return fmt.Errorf("%s: %v", v.Name, err)
	}
	return nil
}

func (v Value) Type() string {
	return v.Valuer.Type()
}

func (v Value) AppendToFlags(flags *pflag.FlagSet) {
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Sprintf("%s => %v", v.Name, e))
		}
	}()
	bInterface := v.Valuer.ValueAddr()
	switch b := (bInterface).(type) {
	case *bool:
		flags.BoolVar(b, v.Name, v.Valuer.Value().(bool), v.Description)
	case *string:
		flags.StringVar(b, v.Name, v.Valuer.Value().(string), v.Description)
	case *int:
		flags.IntVar(b, v.Name, v.Valuer.Value().(int), v.Description)
	default:
		flags.Var(v, v.Name, v.Description)
	}
}
