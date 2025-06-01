package vflags

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ToStringStringSlice(in []string) string {
	if in == nil {
		return "[]"
	}
	if len(in) == 0 {
		return "[]"
	}
	return "[" + strings.Join(in, ", ") + "]"

}

func ToStringDurationSlice(in []time.Duration) string {
	if in == nil {
		return "[]"
	}
	if len(in) == 0 {
		return "[]"
	}
	inS := make([]string, 0, len(in))
	for _, v := range in {
		inS = append(inS, v.String())
	}
	return ToStringStringSlice(inS)

}

func ToStringMapStringStringSlice(in map[string][]string) string {
	s := &strings.Builder{}
	s.WriteString("[ ")
	first := true
	for k, v := range in {
		if !first {
			s.WriteString(", ")
		}
		s.WriteString(k)
		s.WriteString(": ")
		s.WriteString(ToStringStringSlice(v))
		first = false
	}
	s.WriteString(" ]")
	return s.String()
}

func ToStringIntSlice(in []int) string {
	if in == nil {
		return "[]"
	}
	if len(in) == 0 {
		return "[]"
	}
	inS := make([]string, 0, len(in))
	for _, v := range in {
		inS = append(inS, strconv.Itoa(v))
	}
	return ToStringStringSlice(inS)

}

func ToStringBoolSlice(in []bool) string {
	if in == nil {
		return "[]"
	}
	if len(in) == 0 {
		return "[]"
	}
	inS := make([]string, 0, len(in))
	for _, v := range in {
		if v {
			inS = append(inS, "true")
		} else {
			inS = append(inS, "false")
		}

	}
	return ToStringStringSlice(inS)
}

func ToStringMapStringInt(in map[string]int) string {
	s := &strings.Builder{}
	s.WriteString("[ ")
	first := true
	for k, v := range in {
		if !first {
			s.WriteString(", ")
		}
		s.WriteString(k)
		s.WriteString(": ")
		s.WriteString(strconv.Itoa(v))
		first = false
	}
	s.WriteString(" ]")
	return s.String()
}

func ToStringMapStringString(in map[string]string) string {
	s := &strings.Builder{}
	s.WriteString("[ ")
	first := true
	for k, v := range in {
		if !first {
			s.WriteString(", ")
		}
		s.WriteString(k)
		s.WriteString(": ")
		s.WriteString(v)
		first = false
	}
	s.WriteString(" ]")
	return s.String()
}

func lcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
