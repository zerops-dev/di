package vflags

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

type mapValuer struct {
	value         reflect.Value
	defaultValues reflect.Value
}

func (s mapValuer) Set(in interface{}) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = fmt.Errorf("%v", recoverErr)
		}
	}()

	switch s.value.Type() {
	case reflect.TypeOf(map[string]string{}):
		value, err := cast.ToStringMapStringE(in)
		if err != nil {
			return err
		}
		s.value.Set(reflect.ValueOf(value))
	case reflect.TypeOf(map[string]bool{}):
		value, err := cast.ToStringMapBoolE(in)
		if err != nil {
			return err
		}
		s.value.Set(reflect.ValueOf(value))
	case reflect.TypeOf(map[string]int{}):
		value, err := cast.ToStringMapIntE(in)
		if err != nil {
			return err
		}
		s.value.Set(reflect.ValueOf(value))
	case reflect.TypeOf(map[string][]string{}):
		value, err := cast.ToStringMapStringSliceE(in)
		if err != nil {
			return err
		}
		s.value.Set(reflect.ValueOf(value))
	default:
		inR := reflect.ValueOf(in)
		if inR.Kind() != reflect.Map {
			return nil
		}
		valueMap := reflect.MakeMap(s.value.Type())
		r := inR.MapRange()
		e := s.value.Type().Elem()
		for r.Next() {
			vv := reflect.New(e)
			v := defaultValuer{value: vv}
			v.Set(s.defaultValues.Interface())
			v.Set(r.Value().Interface())
			valueMap.SetMapIndex(r.Key(), vv.Elem())
		}
		s.value.Set(valueMap)
	}
	return nil
}

func (s mapValuer) Value() interface{} {
	return s.value.Interface()
}

func (s mapValuer) ValueAddr() interface{} {
	if s.value.CanAddr() {
		return s.value.Addr().Interface()
	}
	return nil
}

func (s mapValuer) Type() string {
	return ".<key>"
}

func (s mapValuer) String() string {
	return cast.ToString(s.value.Interface())
}
