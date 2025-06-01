package env

import "testing"

func TestExpand(t *testing.T) {

	exists := func(in string) (string, bool) {
		value, exists := map[string]string{
			"x":     "xxx",
			"y":     "yy",
			"z":     "z",
			"a":     "$q",
			"b":     "q",
			"aetor": "xxx",
		}[in]
		return value, exists
	}
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Expand simple",
			s:    "x: $x y: $y",
			want: "x: xxx y: yy",
		},
		{
			name: "Expand simple with {}",
			s:    "x: ${x} y: ${y}",
			want: "x: xxx y: yy",
		},
		{
			name: "Expand unknown",
			s:    "w: $w zetor: $zetor",
			want: "w: $w zetor: $zetor",
		},
		{
			name: "Expand unknown 2",
			s:    "w: $w zetor: $aetor",
			want: "w: $w zetor: xxx",
		},
		{
			name: "Expand unknown 3",
			s:    "w: $w zetor: ${a}etor ${b}tor",
			want: "w: $w zetor: $qetor qtor",
		},
		{
			name: "Expand unknown with {}",
			s:    "w: ${w} t: ${t}",
			want: "w: ${w} t: ${t}",
		},
		{
			name: "Expand multiple",
			s:    "x: ${x} y: $y",
			want: "x: xxx y: yy",
		},
		{
			name: "Expand bad syntax",
			s:    "x: ${} y: ${y bla",
			want: "x: ${} y: ${y bla",
		},
		{
			name: "Expand special",
			s:    "x: $@ y: ${y bla",
			want: "x: $@ y: ${y bla",
		},
		{
			name: "Expand special2",
			s:    "x: ${@aa y: ${y bla",
			want: "x: ${@aa y: ${y bla",
		},
		{
			name: "Expand special2",
			s:    "x: ${@} y: ${y bla",
			want: "x: ${@} y: ${y bla",
		},
		{
			name: "Expand complex",
			s:    "xxx $x ${x} $c ${y} ${ahoj} ${aaaa ${} neco",
			want: "xxx xxx xxx $c yy ${ahoj} ${aaaa ${} neco",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Expand(tt.s, exists); got != tt.want {
				t.Errorf("Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}
