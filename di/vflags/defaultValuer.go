package vflags

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type defaultValuer struct {
	value reflect.Value
}

func (s defaultValuer) Set(in interface{}) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = fmt.Errorf("%v", recoverErr)
		}
	}()

	switch {
	case
		s.value.Type().Kind() == reflect.Ptr && s.value.Type().Elem().Kind() == reflect.Struct,
		s.value.Type().Kind() == reflect.Struct:
		setValue := s.value
		if s.value.Type().Kind() == reflect.Ptr {
			setValue = s.value.Elem()
		}

		inR := reflect.ValueOf(in)
		switch inR.Kind() {
		case reflect.Struct:
			for i := 0; i < inR.NumField(); i++ {
				setValue.Field(i).Set(inR.Field(i))
			}

		case reflect.Map:
			inR, isIn := in.(map[string]interface{})
			if !isIn {
				return fmt.Errorf("invalid input format, expected map[string]interface{}")
			}
			for i := 0; i < setValue.Type().NumField(); i++ {
				f := setValue.Type().Field(i)
				d := defaultValuer{value: setValue.Field(i)}
				if v, exists := inR[strings.ToLower(f.Name)]; exists {
					if err := d.Set(v); err != nil {
						return err
					}
				}
			}
		}

	default:
		err = fmt.Errorf("unable to parse %v as %s", in, s.Type())
		var value interface{}
		value = in
		switch s.value.Type() {
		case reflect.TypeOf(time.Second):
			value, err = cast.ToDurationE(in)
		case reflect.TypeOf([]int{}):
			value, err = cast.ToIntSliceE(in)
		case reflect.TypeOf([]string{}):
			value, err = cast.ToStringSliceE(in)
		case reflect.TypeOf([]time.Duration{}):
			value, err = cast.ToDurationSliceE(in)
		case reflect.TypeOf(1):
			value, err = cast.ToIntE(in)
		case reflect.TypeOf(""):
			value, err = cast.ToStringE(in)
		case reflect.TypeOf(true):
			value, err = cast.ToBoolE(in)
		case reflect.TypeOf(float64(1)):
			value, err = cast.ToFloat64E(in)
		case reflect.TypeOf(map[string]string{}):
			value, err = cast.ToStringMapStringE(in)
		case reflect.TypeOf(map[string]bool{}):
			value, err = cast.ToStringMapBoolE(in)
		case reflect.TypeOf(map[string]int{}):
			value, err = cast.ToStringMapIntE(in)
		case reflect.TypeOf(map[string][]string{}):
			value, err = cast.ToStringMapStringSliceE(in)
		case reflect.TypeOf(uint64(0)):
			value, err = cast.ToUint64E(in)
		case reflect.TypeOf(uint32(0)):
			value, err = cast.ToUint32E(in)
		case reflect.TypeOf(uint16(0)):
			value, err = cast.ToUint16E(in)
		case reflect.TypeOf(int64(0)):
			value, err = cast.ToInt64E(in)
		case reflect.TypeOf(int32(0)):
			value, err = cast.ToInt32E(in)
		case reflect.TypeOf(int16(0)):
			value, err = cast.ToInt16E(in)
		default:
			switch s.value.Type().Kind() {
			case reflect.Map:
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
					if err := v.Set(r.Value().Interface()); err != nil {
						return err
					}
					valueMap.SetMapIndex(r.Key(), vv.Elem())
				}
				value = valueMap.Interface()
				err = nil
			default:
				err = nil
			}
		}
		if err != nil {
			return err
		}

		newValue := reflect.ValueOf(value)

		if !newValue.CanConvert(s.value.Type()) {
			msg := fmt.Sprintf("cannot convert %s to %s", newValue.Type(), s.value.Type())
			panic(msg)
		}

		s.value.Set(newValue.Convert(s.value.Type()))
	}
	return nil
}

func (s defaultValuer) Value() interface{} {
	return s.value.Interface()
}

func (s defaultValuer) ValueAddr() interface{} {
	if s.value.CanAddr() {
		return s.value.Addr().Interface()
	}
	return nil
}

func (s defaultValuer) Type() string {
	switch s.value.Type() {
	case reflect.TypeOf([]string{}):
		return "stringSlice"
	case reflect.TypeOf([]int{}):
		return "intSlice"
	case reflect.TypeOf(1):
		return "int"
	case reflect.TypeOf(true):
		return "bool"
	case reflect.TypeOf(float64(1)):
		return "float"
	case reflect.TypeOf(time.Second):
		return "duration"
	default:
		return s.value.Type().String()
	}
}

func (s defaultValuer) String() string {
	switch s.value.Type() {
	case reflect.TypeOf(map[string][]string{}):
		return ToStringMapStringStringSlice(cast.ToStringMapStringSlice(s.value.Interface()))
	case reflect.TypeOf(map[string]int{}):
		return ToStringMapStringInt(cast.ToStringMapInt(s.value.Interface()))
	case reflect.TypeOf(map[string]string{}):
		return ToStringMapStringString(cast.ToStringMapString(s.value.Interface()))
	case reflect.TypeOf([]string{}):
		return ToStringStringSlice(cast.ToStringSlice(s.value.Interface()))
	case reflect.TypeOf([]int{}):
		return ToStringIntSlice(cast.ToIntSlice(s.value.Interface()))
	case reflect.TypeOf([]bool{}):
		return ToStringBoolSlice(cast.ToBoolSlice(s.value.Interface()))
	case reflect.TypeOf([]time.Duration{}):
		return ToStringDurationSlice(cast.ToDurationSlice(s.value.Interface()))
	case reflect.TypeOf(time.Second):
		return s.value.Interface().(time.Duration).String()
	default:
		return fmt.Sprintf("%v", s.value.Interface())
	}
}
