package env

// Expand replaces ${var} or $var in the string based on the mapping function.
// ⚠ Customized version that does not expand unknown variables into empty strings, but leaves them untouched.
func Expand(s string, mapping func(string) (string, bool)) string {
	var buf []byte
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, orig, w := getShellName(s[j+1:])
			switch {
			case name == "" && w > 0:
				// Encountered invalid syntax; eat the
				// characters.
			case name == "":
				// Valid syntax, but $ was not followed by a
				// name. Leave the dollar character untouched.
				buf = append(buf, s[j])
			default:
				if value, exists := mapping(name); exists {
					buf = append(buf, value...)
				} else {
					buf = append(buf, "$"...)
					buf = append(buf, orig...)
				}
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s
	}
	return string(buf) + s[i:]
}

// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore.
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// getShellName returns the name that begins the string and the number of bytes
// consumed to extract it. If the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name.
func getShellName(s string) (string, string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == '}' {
			return s[1:3], s[0:3], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return s[0:2], s[0:2], 2 // Bad syntax; eat "${}"
				}
				return s[1:i], s[0 : i+1], i + 1
			}
		}
		return s[0:1], s[0:1], 1 // Bad syntax; "${"
	case isShellSpecialVar(s[0]):
		return s[0:1], s[0:1], 1
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}
	return s[:i], s[:i], i
}
