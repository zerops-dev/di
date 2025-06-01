package s

import (
	_ "embed"
)

//go:embed di.tmpl
var _Di string

func (h Di) Template() string {
	return _Di
}
