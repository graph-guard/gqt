package gqt

import (
	"bytes"
	"strconv"
)

type source struct {
	Location
	s []byte
}

func (s source) stop() bool { return s.s == nil }

func stop() source { return source{} }

func (s source) isEOF() bool {
	return s.Index >= len(s.s)
}

func (s source) isDigit() bool {
	return s.Index < len(s.s) &&
		s.s[s.Index] >= '0' && s.s[s.Index] <= '9'
}

func (s source) peek1(b byte) bool {
	return s.Index < len(s.s) && s.s[s.Index] == b
}

// consumeIgnored skips spaces, tabs, line-feeds
// carriage-returns and comment sequences.
func (s source) consumeIgnored() source {
MAIN:
	for s.Index < len(s.s) {
		switch s.s[s.Index] {
		case '#':
			s.Index++
			s.Column++
			for s.Index < len(s.s) {
				if s.s[s.Index] == '\n' {
					s.Index++
					s.Line++
					s.Column = 1
					continue MAIN
				} else if s.s[s.Index] < 0x20 {
					return s
				}
				s.Index++
				s.Column++
			}
		case ' ', '\t', '\r':
			s.Index++
			s.Column++
			continue
		case '\n':
			s.Index++
			s.Line++
			s.Column = 1
			continue
		}
		break
	}
	return s
}

func (s source) consumeEitherOf3(a, b, c string) (_ source, selected int) {
	selected = -1
	if s.Index >= len(s.s) {
		return s, selected
	}
	x := s.s[s.Index:]
	if a != "" && bytes.HasPrefix(x, []byte(a)) {
		selected = 0
		for i := len(a); len(x) > 0 && i > 0; i-- {
			s.Index++
			if x[0] == '\n' {
				s.Line++
				s.Column = 1
			} else {
				s.Column++
			}
			x = x[1:]
		}
	} else if b != "" && bytes.HasPrefix(x, []byte(b)) {
		selected = 1
		for i := len(b); len(x) > 0 && i > 0; i-- {
			s.Index++
			if x[0] == '\n' {
				s.Line++
				s.Column = 1
			} else {
				s.Column++
			}
			x = x[1:]
		}
	} else if c != "" && bytes.HasPrefix(x, []byte(c)) {
		selected = 2
		for i := len(c); len(x) > 0 && i > 0; i-- {
			s.Index++
			if x[0] == '\n' {
				s.Line++
				s.Column = 1
			} else {
				s.Column++
			}
			x = x[1:]
		}
	}
	return s, selected
}

func (s source) consume(x string) (_ source, ok bool) {
	si, i := s, 0
	for ; s.Index < len(s.s); i++ {
		if i >= len(x) {
			return s, true
		}
		if s.s[s.Index] != x[i] {
			return si, false
		}
		if s.s[s.Index] == '\n' {
			s.Index++
			s.Line++
			s.Column = 1
		} else {
			s.Index++
			s.Column++
		}
	}
	return s, i >= len(x)
}

func (s source) consumeToken() (_ source, token []byte) {
	i, ii := s.s[s.Index:], s.Index
	for s.Index < len(s.s) {
		if b := s.s[s.Index]; b < 0x20 ||
			b == ' ' || b == '\n' || b == '\t' ||
			b == '\r' || b == ',' || b == '{' ||
			b == '}' || b == '(' || b == ')' ||
			b == ']' || b == '[' || b == '#' {
			break
		}
		s.Index++
		s.Column++
	}
	return s, i[:s.Index-ii]
}

func (s source) consumeName() (_ source, name []byte) {
	ii := s.Index
	if s.Index >= len(s.s) {
		return s, nil
	}
	if b := s.s[s.Index]; b == '_' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') {
		s.Index++
		s.Column++
	} else {
		return s, nil
	}
	for s.Index < len(s.s) && (s.s[s.Index] == '_' ||
		(s.s[s.Index] >= 'a' && s.s[s.Index] <= 'z') ||
		(s.s[s.Index] >= 'A' && s.s[s.Index] <= 'Z') ||
		(s.s[s.Index] >= '0' && s.s[s.Index] <= '9')) {
		s.Index++
		s.Column++
	}
	return s, s.s[ii:s.Index]
}

func (s source) consumeString() (n source, str []byte, ok bool) {
	escaped := false
	n = s

	if n.Index >= len(n.s) || n.s[n.Index] != '"' {
		return s, nil, false
	}
	n.Index++
	n.Column++
	ii := n.Index

	for n.Index < len(n.s) {
		switch {
		case n.s[n.Index] < 0x20:
			return s, nil, false
		case n.s[n.Index] == '"':
			if escaped {
				escaped = false
			} else {
				str = n.s[ii:n.Index]
				n.Index++
				n.Column++
				return n, str, true
			}
		case n.s[n.Index] == '\\':
			escaped = !escaped
		}
		n.Index++
		n.Column++
	}
	return s, nil, false
}

func (s source) consumeUnsignedInt() (next source, i int64, ok bool) {
	si := s
	for s.s[s.Index] >= '0' && s.s[s.Index] <= '9' {
		s.Index++
		s.Column++
	}
	i, err := strconv.ParseInt(string(s.s[si.Index:s.Index]), 10, 64)
	if err != nil {
		return si, 0, false
	}
	return s, i, true
}
