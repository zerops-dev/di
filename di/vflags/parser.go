package vflags

import (
	"fmt"
	"reflect"
	"strings"
)

type Valuer interface {
	Set(interface{}) error
	Type() string
	String() string
	Value() interface{}
	ValueAddr() interface{}
}

func Parse(name string, in interface{}) []Value {
	inValue := reflect.ValueOf(in)
	alreadyDefined := make(map[string]struct{})
	return parse(name, inValue, false, alreadyDefined)
}

func parse(name string, in reflect.Value, virtual bool, alreadyDefined map[string]struct{}) []Value {

	if in.Kind() == reflect.Ptr {
		return parse(name, in.Elem(), virtual, alreadyDefined)
	}
	if in.Kind() != reflect.Struct {
		return nil
	}
	var result []Value
	inType := in.Type()
	for i := 0; i < inType.NumField(); i++ {

		f := inType.Field(i)
		configName := f.Tag.Get("flag")
		subName := lcFirst(f.Name)
		if configName == "-" {
			continue
		}
		if configName != "" {
			subName = configName
		}
		newName := subName
		if name != "" {
			newName = name + "." + subName
		}
		if subName == "<" {
			newName = name
			delete(alreadyDefined, name) // means we are flattening the config
		}
		if _, ok := alreadyDefined[newName]; ok && subName != "<" {
			panic(fmt.Sprint("field ", newName, " already defined in config"))
		}
		alreadyDefined[newName] = struct{}{}
		switch k := in.Field(i).Kind(); k {
		case reflect.Slice:
			valuer := defaultValuer{value: in.FieldByName(f.Name)}
			result = append(result, Value{
				Name:         newName,
				Description:  f.Tag.Get("desc"),
				Configurable: f.Tag.Get("configurable") != "false",
				Valuer:       valuer,
				Virtual:      virtual,
			})
		case reflect.Map:
			elType := in.Field(i).Type().Elem()
			elDefault := reflect.New(elType)
			{
				r := in.Field(i).MapRange()
				if r.Next() {
					elDefault = r.Value()
				}
				result = append(result, parse(newName+".<key>", elDefault, true, alreadyDefined)...)
			}

			valuer := mapValuer{
				value:         in.FieldByName(f.Name),
				defaultValues: elDefault,
			}
			result = append(result, Value{
				Name:         newName,
				Description:  f.Tag.Get("desc"),
				Configurable: f.Tag.Get("configurable") != "false",
				Valuer:       valuer,
				Virtual:      virtual,
			})
		case reflect.Struct:
			result = append(result, parse(newName, in.Field(i), virtual, alreadyDefined)...)
		default:
			valuer := defaultValuer{value: in.FieldByName(f.Name)}
			result = append(result, Value{
				Name:         newName,
				Description:  desc(f.Tag, newName),
				Configurable: f.Tag.Get("configurable") != "false",
				Valuer:       valuer,
				Virtual:      virtual,
			})
		}
	}
	return result
}

func desc(tag reflect.StructTag, name string) string {
	name = strings.NewReplacer(".", "_").Replace(name)
	name = strings.ToUpper(name)
	if desc, ok := tag.Lookup("desc"); ok {
		return fmt.Sprintf("%s (env %s)", desc, name)
	}
	return fmt.Sprintf("(env %s)", name)
}
