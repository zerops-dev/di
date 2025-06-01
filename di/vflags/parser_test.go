package vflags

import (
	"testing"
	"time"
)

type Sub struct {
	In  string
	Out string
}

type Embed struct {
	Notherdam bool
}

type Config struct {
	Name           string          `desc:"description"`
	Value1         string          `desc:"description"`
	Value2         int             `desc:"description"`
	Values         []string        `desc:"description"`
	ValuesDur      []time.Duration `desc:"description"`
	T              int             `flag:"aaaa" desc:"description"`
	ValuesInt      []int           `desc:"description"`
	MapValues      map[string]string
	D              time.Duration
	B              bool `desc:"bool value"`
	SubStruct      Sub
	MapInt         map[string]int
	MapStringSlice map[string][]string
	BoolSlice      []bool `desc:"coze"`
	Embed          `flag:"<"`
}

func NewConfig() Config {
	return Config{
		Name:      "default458",
		Value1:    "42",
		Value2:    42,
		T:         123,
		Values:    []string{"aaaa", "rrrr"},
		MapInt:    map[string]int{"aaaa": 21, "rrr": 678},
		ValuesInt: []int{3, 2, 1},
		ValuesDur: []time.Duration{time.Second, time.Minute},
		D:         time.Second * 5,
		SubStruct: Sub{
			Out: "outaaaaa",
			In:  "inaaa",
		},
		B: true,
		MapValues: map[string]string{
			"m1": "v1",
			"m2": "v2",
			"m3": "v3",
		},
		MapStringSlice: map[string][]string{
			"aaa": {"jedna", "dva"},
			"bbb": {"tri", "ctyri"},
		},
		BoolSlice: []bool{true, false, true},
		Embed: Embed{
			Notherdam: true,
		},
	}
}

func TestParse(t *testing.T) {
	c := NewConfig()
	values := Parse("app", &c)

	if err := values[0].SetValue(time.Second); err != nil {
		t.Error(err)
	}
	if c.Name != "1s" {
		t.Errorf("expected 1s, have %#v", c.Name)
	}
}

type AlreadyDefined struct {
	Config `flag:"<"`
	Name   string `desc:"redefined field"`
}

func TestAlreadyDefined(t *testing.T) {
	c := AlreadyDefined{NewConfig(), "adios"}
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
		}
	}()
	_ = Parse("app", &c)
	t.Fatal("should panic")
}
