package gqt

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

type Doc interface{}

type DocQuery struct {
	Selections []Selection
}

type DocMutation struct {
	Selections []Selection
}

type Selection struct {
	Name             FieldName
	InputConstraints []InputConstraint
	Selections       []Selection
}

type InputConstraint struct {
	Name       ParameterName
	Constraint Constraint
}

type Constraint interface{}

type (
	ConstraintIsEqual          struct{ Value Value }
	ConstraintIsNotEqual       struct{ Value Value }
	ConstraintIsGreater        struct{ Value float64 }
	ConstraintIsLess           struct{ Value float64 }
	ConstraintIsGreaterOrEqual struct{ Value float64 }
	ConstraintIsLessOrEqual    struct{ Value float64 }
	ConstraintIsType           struct{ TypeName string }

	ConstraintBytelenEqual          struct{ Value uint }
	ConstraintBytelenNotEqual       struct{ Value uint }
	ConstraintBytelenGreater        struct{ Value uint }
	ConstraintBytelenLess           struct{ Value uint }
	ConstraintBytelenGreaterOrEqual struct{ Value uint }
	ConstraintBytelenLessOrEqual    struct{ Value uint }

	ConstraintLenEqual          struct{ Value uint }
	ConstraintLenNotEqual       struct{ Value uint }
	ConstraintLenGreater        struct{ Value uint }
	ConstraintLenLess           struct{ Value uint }
	ConstraintLenGreaterOrEqual struct{ Value uint }
	ConstraintLenLessOrEqual    struct{ Value uint }
)

type Value interface{}

type FieldName = string
type ParameterName = string

func Parse(s []byte) (Doc, Error) {
	return parse(source{s, 0})
}

func parse(s source) (Doc, Error) {
	s = s.skipIrrelevant()

	var ok bool
	var err Error
	s, ok = s.consume(keywordQuery)
	if ok {
		s = s.skipIrrelevant()
		var selections []Selection
		if s, selections, err = parseSelectionSet(s); err.IsErr() {
			return nil, err
		}
		s = s.skipIrrelevant()
		if err = s.expectEOF(); err.IsErr() {
			return nil, err
		}
		return DocQuery{
			Selections: selections,
		}, Error{}
	}

	s, ok = s.consume(keywordMutation)
	if ok {
		s = s.skipIrrelevant()
		var selections []Selection
		if s, selections, err = parseSelectionSet(s); err.IsErr() {
			return nil, err
		}
		s = s.skipIrrelevant()
		if err = s.expectEOF(); err.IsErr() {
			return nil, err
		}
		return DocMutation{
			Selections: selections,
		}, Error{}
	}

	return nil, s.err("unexpected definition")
}

func parseSelectionSet(s source) (source, []Selection, Error) {
	var ok bool
	s, ok = s.consume(curlyBracketLeft)
	if !ok {
		return s, nil, s.err("expected selection set")
	}

	var selections []Selection

	s = s.skipIrrelevant()
	for {
		var err Error
		var selection Selection
		if s, selection, err = parseSelection(s); err.IsErr() {
			return s, nil, err
		}
		selections = append(selections, selection)
		s = s.skipIrrelevant()
		if s, ok = s.consume(curlyBracketRight); ok {
			s = s.skipIrrelevant()
			break
		}
	}

	return s, selections, Error{}
}

func parseSelection(s source) (_ source, sel Selection, err Error) {
	var ok bool
	var name []byte
	if s, name = s.consumeName(); name == nil {
		return s, Selection{}, s.err("expected field name")
	}
	sel.Name = FieldName(name)

	s = s.skipIrrelevant()

	if s, ok = s.consume(parenthesisLeft); ok {
		for {
			s = s.skipIrrelevant()
			if s, name = s.consumeName(); name == nil {
				return s, Selection{}, s.err("expected parameter name")
			}
			s = s.skipIrrelevant()
			if s, ok = s.consume(column); !ok {
				return s, Selection{}, s.err("expected column after parameter name")
			}
			s = s.skipIrrelevant()
			var inputConstraint Constraint
			if s, inputConstraint, err = parseConstraint(s); err.IsErr() {
				return s, Selection{}, err
			}
			sel.InputConstraints = append(sel.InputConstraints, InputConstraint{
				Name:       ParameterName(name),
				Constraint: inputConstraint,
			})
			s = s.skipIrrelevant()
			if s, ok = s.consume(parenthesisRight); ok {
				s = s.skipIrrelevant()
				break
			}
		}
	}

	if bytes.HasPrefix(s.s, curlyBracketLeft) {
		s = s.skipIrrelevant()
		if s, sel.Selections, err = parseSelectionSet(s); err.IsErr() {
			return s, Selection{}, err
		}
	}

	return s, sel, Error{}
}

func parseConstraint(s source) (_ source, c Constraint, err Error) {
	var name []byte
	if s, name = s.consumeName(); name == nil {
		return s, nil, s.err("expected constraint function name")
	}
	s = s.skipIrrelevant()

	si := s
	var t []byte
	s, t = s.consumeToken()

	switch string(name) {
	case "is":
		switch string(t) {
		case "=":
			// is = x
			c = ConstraintIsEqual{}
		case "!=":
			// is != x
			c = ConstraintIsNotEqual{}
		case ">":
			// is > x
			c = ConstraintIsGreater{}
		case "<":
			// is < x
			c = ConstraintIsLess{}
		case ">=":
			// is >= x
			c = ConstraintIsGreaterOrEqual{}
		case "<=":
			// is <= x
			c = ConstraintIsLessOrEqual{}
		default:
			return si, nil, s.err(
				"unsupported operator for 'is' constraint",
			)
		}
	case "len":
		switch string(t) {
		case "=":
			// len = x
			c = ConstraintLenEqual{}
		case "!=":
			// len != x
			c = ConstraintLenNotEqual{}
		case ">":
			// len > x
			c = ConstraintLenGreater{}
		case "<":
			// len < x
			c = ConstraintLenLess{}
		case ">=":
			// len >= x
			c = ConstraintLenGreaterOrEqual{}
		case "<=":
			// len <= x
			c = ConstraintLenLessOrEqual{}
		default:
			return si, Selection{}, s.err(
				"unsupported operator for 'len' constraint function",
			)
		}

	case "bytelen":
		switch string(t) {
		case "=":
			// bytelen = x
			c = ConstraintBytelenEqual{}
		case "!=":
			// bytelen != x
			c = ConstraintBytelenNotEqual{}
		case ">":
			// bytelen > x
			c = ConstraintBytelenGreater{}
		case "<":
			// bytelen < x
			c = ConstraintBytelenLess{}
		case ">=":
			// bytelen >= x
			c = ConstraintBytelenGreaterOrEqual{}
		case "<=":
			// bytelen <= x
			c = ConstraintBytelenLessOrEqual{}
		default:
			return s, nil, s.err(
				"unsupported operator for 'bytelen' constraint function",
			)
		}
	default:
		return s, nil, s.err("unsupported constraint function")
	}

	s = s.skipIrrelevant()

	var v Value
	if s, v, err = parseValue(s); err.IsErr() {
		return s, nil, err
	}

	switch c.(type) {
	// Is
	case ConstraintIsEqual:
		c = ConstraintIsEqual{Value: v}
	case ConstraintIsNotEqual:
		c = ConstraintIsNotEqual{Value: v}
	case ConstraintIsGreater:
		if v, ok := v.(float64); ok {
			c = ConstraintIsGreater{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	case ConstraintIsLess:
		if v, ok := v.(float64); ok {
			c = ConstraintIsLess{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	case ConstraintIsGreaterOrEqual:
		if v, ok := v.(float64); ok {
			c = ConstraintIsGreaterOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	case ConstraintIsLessOrEqual:
		if v, ok := v.(float64); ok {
			c = ConstraintIsLessOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}

	// Len
	case ConstraintLenEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintLenNotEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenNotEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintLenGreater:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenGreater{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintLenLess:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenLess{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintLenGreaterOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenGreaterOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintLenLessOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenLessOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}

	// Bytelen
	case ConstraintBytelenEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintBytelenNotEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenNotEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintBytelenGreater:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenGreater{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintBytelenLess:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenLess{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintBytelenGreaterOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenGreaterOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	case ConstraintBytelenLessOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenLessOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected unsigned integer")
		}
	}

	return s, c, Error{}
}

func parseValue(s source) (sn source, v Value, err Error) {
	if len(s.s) < 1 {
		return s, nil, s.err("unexpected end of file")
	}

	if b := s.s[0]; b == '-' ||
		b == '+' ||
		b == '0' ||
		b == '1' ||
		b == '2' ||
		b == '3' ||
		b == '4' ||
		b == '5' ||
		b == '6' ||
		b == '7' ||
		b == '8' ||
		b == '9' {
		var num float64
		var ok bool
		if s, num, ok = s.consumeNumber(); ok {
			return s, num, Error{}
		}
	} else if b == '"' {
		var str []byte
		var ok bool
		if s, str, ok = s.consumeString(); ok {
			return s, string(str), Error{}
		}
	}

	var t []byte
	si := s
	s, t = s.consumeToken()
	switch string(t) {
	case "null":
		return s, nil, Error{}
	case "true":
		return s, true, Error{}
	case "false":
		return s, false, Error{}
	}

	return si, nil, s.err("invalid value")
}

type source struct {
	s     []byte
	index int
}

func (s source) err(msg string) Error {
	return Error{
		Index: s.index,
		Msg:   msg,
	}
}

type Error struct {
	Index int
	Msg   string
}

func (e Error) IsErr() bool {
	return e.Msg != ""
}

func (e Error) Error() string {
	if e.Msg == "" {
		return ""
	}
	return fmt.Sprintf("error at %d: %s", e.Index, e.Msg)
}

var (
	keywordQuery      = []byte("query")
	keywordMutation   = []byte("mutation")
	curlyBracketLeft  = []byte("{")
	curlyBracketRight = []byte("}")
	parenthesisLeft   = []byte("(")
	parenthesisRight  = []byte(")")
	column            = []byte(":")
	valueNull         = []byte("null")
	valueFalse        = []byte("false")
	valueTrue         = []byte("true")
)

func (s source) expectEOF() (e Error) {
	if s.index < len(s.s) {
		e.Index = s.index
		e.Msg = "expected EOF"
	}
	return
}

// skipIrrelevant skips spaces, tabs, line-feeds
// carriage-returns and comment sequences
func (s source) skipIrrelevant() source {
MAIN:
	for i := range s.s {
		if s.s[i] == '#' {
			c := s.s[i:]
			for i := range c {
				if c[i] == '\n' {
					s.index += i
					s.s = s.s[i:]
					continue MAIN
				}
			}
		} else if s.s[i] == ' ' ||
			s.s[i] == '\n' ||
			s.s[i] == '\t' ||
			s.s[i] == '\r' ||
			s.s[i] == ',' {
			continue
		}
		s.index += i
		s.s = s.s[i:]
		break
	}
	return s
}

func (s source) consume(x []byte) (_ source, ok bool) {
	if bytes.HasPrefix(s.s, x) {
		s.index += len(x)
		s.s = s.s[len(x):]
		return s, true
	}
	return s, false
}

func (s source) consumeNumber() (_ source, number float64, ok bool) {
	if s, token := s.consumeToken(); token != nil {
		if f, err := strconv.ParseFloat(string(token), 64); err == nil {
			return s, f, true
		}
	}
	return s, 0, false
}

func (s source) consumeToken() (_ source, token []byte) {
	i, ii := s.s, s.index
	if len(s.s) < 1 {
		return s, nil
	}
	for len(s.s) > 0 {
		if b := s.s[0]; b == ' ' ||
			b == '\n' ||
			b == '\t' ||
			b == '\r' ||
			b == ',' ||
			b == '{' ||
			b == '(' ||
			b == '#' {
			break
		}
		s.index++
		s.s = s.s[1:]
	}
	return s, i[:s.index-ii]
}

func (s source) consumeName() (_ source, name []byte) {
	i, ii := s.s, s.index
	if len(s.s) < 1 {
		return s, nil
	}
	if b := s.s[0]; b == '_' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') {
		s.index++
		s.s = s.s[1:]
	} else {
		return s, nil
	}
	for len(s.s) > 0 && (s.s[0] == '_' ||
		(s.s[0] >= 'a' && s.s[0] <= 'z') ||
		(s.s[0] >= 'A' && s.s[0] <= 'Z') ||
		(s.s[0] >= '0' && s.s[0] <= '9')) {
		s.index++
		s.s = s.s[1:]
	}
	return s, i[:s.index-ii]
}

func (s source) consumeString() (_ source, str []byte, ok bool) {
	escaped := false

	if len(s.s) < 1 || s.s[0] != '"' {
		return s, nil, false
	}
	s.index++
	s.s = s.s[1:]
	ii := s

	for ; len(s.s) > 0; s.s, s.index = s.s[1:], s.index+1 {
		if s.s[0] < 0x20 {
			return s, nil, false
		} else if s.s[0] == '"' {
			if escaped {
				escaped = false
			} else {
				str = ii.s[:s.index-ii.index]
				s.index++
				s.s = s.s[1:]
				return s, str, true
			}
		} else if s.s[0] == '\\' {
			escaped = !escaped
		}
	}
	return s, nil, false
}

func f64ToUint(v interface{}) (uint, bool) {
	if f, ok := v.(float64); ok && !math.Signbit(f) && f == float64(uint(f)) {
		return uint(f), true
	}
	return 0, false
}
