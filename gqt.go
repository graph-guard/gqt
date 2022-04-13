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
	Name             string
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

	ConstraintAreEqual          struct{ Value Value }
	ConstraintAreNotEqual       struct{ Value Value }
	ConstraintAreGreater        struct{ Value float64 }
	ConstraintAreLess           struct{ Value float64 }
	ConstraintAreGreaterOrEqual struct{ Value float64 }
	ConstraintAreLessOrEqual    struct{ Value float64 }
	ConstraintAreType           struct{ TypeName string }
)

type Value interface{}

type ValueArray struct {
	Items []Constraint
}

type ValueObject struct {
	Fields []ObjectField
}

type ObjectField struct {
	Name  string
	Value Constraint
}

type ParameterName = string

func Parse(s []byte) (Doc, Error) {
	return parse(source{s, s})
}

func parse(s source) (Doc, Error) {
	s = s.consumeIrrelevant()
	if s.isEOF() {
		return nil, s.err("expected definition")
	}

	var ok bool
	var err Error

	if s, ok = s.consume(keywordQuery); ok {
		s = s.consumeIrrelevant()

		var selections []Selection
		if s, selections, err = parseSelectionSet(s); err.IsErr() {
			return nil, err
		}
		s = s.consumeIrrelevant()
		if err = s.expectEOF(); err.IsErr() {
			return nil, err
		}
		return DocQuery{
			Selections: selections,
		}, Error{}
	} else if s, ok = s.consume(keywordMutation); ok {
		s = s.consumeIrrelevant()
		var selections []Selection
		if s, selections, err = parseSelectionSet(s); err.IsErr() {
			return nil, err
		}
		s = s.consumeIrrelevant()
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

	s = s.consumeIrrelevant()
	for {
		var err Error
		var selection Selection
		sb := s
		if s, selection, err = parseSelection(s); err.IsErr() {
			return sb, nil, err
		}

		// Check for redundancy
		for _, x := range selections {
			if x.Name == selection.Name {
				return sb, nil, sb.err("redundant selection")
			}
		}

		selections = append(selections, selection)
		s = s.consumeIrrelevant()
		if s, ok = s.consume(curlyBracketRight); ok {
			s = s.consumeIrrelevant()
			break
		}
	}

	return s, selections, Error{}
}

func parseSelection(s source) (n source, sel Selection, err Error) {
	var ok bool
	var name []byte
	if n, name = s.consumeName(); name == nil {
		return s, Selection{}, s.err("expected field name")
	}
	sel.Name = string(name)

	n = n.consumeIrrelevant()

	if n, ok = n.consume(parenthesisLeft); ok {
		for {
			s = n
			n = n.consumeIrrelevant()
			if n, name = n.consumeName(); name == nil {
				return s, Selection{}, n.err("expected parameter name")
			}
			n = n.consumeIrrelevant()
			if n, ok = n.consume(column); !ok {
				return s, Selection{}, n.err(
					"expected column after parameter name",
				)
			}
			n = n.consumeIrrelevant()
			var inputConstraint Constraint
			if n, inputConstraint, err = parseConstraint(n); err.IsErr() {
				return s, Selection{}, err
			}

			// Check for redundancy
			for _, x := range sel.InputConstraints {
				if x.Name == string(name) {
					return s, Selection{}, s.err("redundant constraint")
				}
			}

			sel.InputConstraints = append(
				sel.InputConstraints,
				InputConstraint{
					Name:       ParameterName(name),
					Constraint: inputConstraint,
				},
			)
			n = n.consumeIrrelevant()
			if n, ok = n.consume(parenthesisRight); ok {
				n = n.consumeIrrelevant()
				break
			}
		}
	}

	if bytes.HasPrefix(n.s, curlyBracketLeft) {
		n = n.consumeIrrelevant()
		if n, sel.Selections, err = parseSelectionSet(n); err.IsErr() {
			return s, Selection{}, err
		}
	}

	return n, sel, Error{}
}

func parseConstraint(s source) (_ source, c Constraint, err Error) {
	var name []byte
	if s, name = s.consumeName(); name == nil {
		return s, nil, s.err("expected constraint name")
	}

	s = s.consumeIrrelevant()
	if s.isEOF() {
		return s, nil, s.err("expected constraint operator")
	}

	si := s
	var ok bool

	switch string(name) {
	case "is":
		if s, ok = s.consume(operatorGreaterEqual); ok {
			// is >= x
			c = ConstraintIsGreaterOrEqual{}
		} else if s, ok = s.consume(operatorLesserEqual); ok {
			// is <= x
			c = ConstraintIsLessOrEqual{}
		} else if s, ok = s.consume(operatorEqual); ok {
			// is = x
			c = ConstraintIsEqual{}
		} else if s, ok = s.consume(operatorNotEqual); ok {
			// is != x
			c = ConstraintIsNotEqual{}
		} else if s, ok = s.consume(operatorGreater); ok {
			// is > x
			c = ConstraintIsGreater{}
		} else if s, ok = s.consume(operatorLesser); ok {
			// is < x
			c = ConstraintIsLess{}
		} else {
			return si, nil, s.err(
				"unsupported operator for 'is' constraint",
			)
		}

	case "len":
		if s, ok = s.consume(operatorGreaterEqual); ok {
			// len >= x
			c = ConstraintLenGreaterOrEqual{}
		} else if s, ok = s.consume(operatorLesserEqual); ok {
			// len <= x
			c = ConstraintLenLessOrEqual{}
		} else if s, ok = s.consume(operatorEqual); ok {
			// len = x
			c = ConstraintLenEqual{}
		} else if s, ok = s.consume(operatorNotEqual); ok {
			// len != x
			c = ConstraintLenNotEqual{}
		} else if s, ok = s.consume(operatorGreater); ok {
			// len > x
			c = ConstraintLenGreater{}
		} else if s, ok = s.consume(operatorLesser); ok {
			// len < x
			c = ConstraintLenLess{}
		} else {
			return si, nil, s.err(
				"unsupported operator for 'len' constraint",
			)
		}

	case "bytelen":
		if s, ok = s.consume(operatorGreaterEqual); ok {
			// bytelen >= x
			c = ConstraintBytelenGreaterOrEqual{}
		} else if s, ok = s.consume(operatorLesserEqual); ok {
			// bytelen <= x
			c = ConstraintBytelenLessOrEqual{}
		} else if s, ok = s.consume(operatorEqual); ok {
			// bytelen = x
			c = ConstraintBytelenEqual{}
		} else if s, ok = s.consume(operatorNotEqual); ok {
			// bytelen != x
			c = ConstraintBytelenNotEqual{}
		} else if s, ok = s.consume(operatorGreater); ok {
			// bytelen > x
			c = ConstraintBytelenGreater{}
		} else if s, ok = s.consume(operatorLesser); ok {
			// bytelen < x
			c = ConstraintBytelenLess{}
		} else {
			return si, nil, s.err(
				"unsupported operator for 'bytelen' constraint",
			)
		}

	case "are":
		if s, ok = s.consume(operatorGreaterEqual); ok {
			// are >= x
			c = ConstraintAreGreaterOrEqual{}
		} else if s, ok = s.consume(operatorLesserEqual); ok {
			// are <= x
			c = ConstraintAreLessOrEqual{}
		} else if s, ok = s.consume(operatorEqual); ok {
			// are = x
			c = ConstraintAreEqual{}
		} else if s, ok = s.consume(operatorNotEqual); ok {
			// are != x
			c = ConstraintAreNotEqual{}
		} else if s, ok = s.consume(operatorGreater); ok {
			// are > x
			c = ConstraintAreGreater{}
		} else if s, ok = s.consume(operatorLesser); ok {
			// are < x
			c = ConstraintAreLess{}
		} else {
			return si, nil, s.err(
				"unsupported operator for 'are' constraint",
			)
		}
	default:
		return s, nil, s.err("unsupported constraint function")
	}

	s = s.consumeIrrelevant()

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
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintLenNotEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenNotEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintLenGreater:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenGreater{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintLenLess:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenLess{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintLenGreaterOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenGreaterOrEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintLenLessOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintLenLessOrEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}

	// Bytelen
	case ConstraintBytelenEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintBytelenNotEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenNotEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintBytelenGreater:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenGreater{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintBytelenLess:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenLess{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintBytelenGreaterOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenGreaterOrEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}
	case ConstraintBytelenLessOrEqual:
		if v, ok := f64ToUint(v); ok {
			c = ConstraintBytelenLessOrEqual{Value: v}
		} else {
			return s, nil, s.err(
				"unexpected value type, expected unsigned integer",
			)
		}

	// Are
	case ConstraintAreEqual:
		c = ConstraintAreEqual{Value: v}
	case ConstraintAreNotEqual:
		c = ConstraintAreNotEqual{Value: v}
	case ConstraintAreGreater:
		if v, ok := v.(float64); ok {
			c = ConstraintAreGreater{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	case ConstraintAreLess:
		if v, ok := v.(float64); ok {
			c = ConstraintAreLess{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	case ConstraintAreGreaterOrEqual:
		if v, ok := v.(float64); ok {
			c = ConstraintAreGreaterOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	case ConstraintAreLessOrEqual:
		if v, ok := v.(float64); ok {
			c = ConstraintAreLessOrEqual{Value: v}
		} else {
			return s, nil, s.err("unexpected value type, expected number")
		}
	}

	return s, c, Error{}
}

func parseValue(s source) (_ source, v Value, err Error) {
	if s.isEOF() {
		return s, nil, s.err("expected value")
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
	} else if b == '[' {
		var a ValueArray
		s, a, err = parseValueArray(s)
		if err.IsErr() {
			return
		}
		return s, a, Error{}
	} else if b == '{' {
		var o ValueObject
		s, o, err = parseValueObject(s)
		if err.IsErr() {
			return
		}
		s = s.consumeIrrelevant()
		return s, o, Error{}
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

func parseValueArray(s source) (_ source, a ValueArray, err Error) {
	var ok bool
	if s, ok = s.consume(squareBracketLeft); !ok {
		return s, ValueArray{}, s.err("expected array")
	}

	for {
		s = s.consumeIrrelevant()
		if s, ok = s.consume(squareBracketRight); ok {
			break
		}
		var c Constraint
		if s, c, err = parseConstraint(s); err.IsErr() {
			return s, ValueArray{}, err
		}
		a.Items = append(a.Items, c)
	}

	for _, item := range a.Items {
		switch item.(type) {
		case ConstraintAreEqual:
		case ConstraintAreNotEqual:
		case ConstraintAreGreater:
		case ConstraintAreLess:
		case ConstraintAreGreaterOrEqual:
		case ConstraintAreLessOrEqual:
		case ConstraintAreType:
		default:
			continue
		}
		if len(a.Items) > 1 {
			return s, a, s.err("contains exclusive are constraint")
		}
		break
	}

	return s, a, Error{}
}

func parseValueObject(s source) (_ source, o ValueObject, err Error) {
	var ok bool
	if s, ok = s.consume(curlyBracketLeft); !ok {
		return s, ValueObject{}, s.err("expected object")
	}

	for {
		s = s.consumeIrrelevant()
		if s, ok = s.consume(curlyBracketRight); ok {
			break
		}

		var name []byte
		if s, name = s.consumeName(); name == nil {
			return s, ValueObject{}, s.err("expected field name")
		}

		s = s.consumeIrrelevant()

		if s, ok = s.consume(column); !ok {
			return s, ValueObject{}, s.err(
				"expected column after object field name",
			)
		}

		s = s.consumeIrrelevant()

		var c Constraint
		if s, c, err = parseConstraint(s); err.IsErr() {
			return s, ValueObject{}, err
		}
		o.Fields = append(o.Fields, ObjectField{
			Name:  string(name),
			Value: c,
		})
	}

	if len(o.Fields) < 1 {
		return s, ValueObject{}, s.err("empty object")
	}

	return s, o, Error{}
}

type source struct {
	original []byte
	s        []byte
}

func (s source) isEOF() bool {
	return len(s.s) < 1
}

func (s source) index() int {
	return len(s.original) - len(s.s)
}

func (s source) err(msg string) Error {
	return Error{
		Index: len(s.original) - len(s.s),
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
	keywordQuery         = []byte("query")
	keywordMutation      = []byte("mutation")
	squareBracketLeft    = []byte("[")
	squareBracketRight   = []byte("]")
	curlyBracketLeft     = []byte("{")
	curlyBracketRight    = []byte("}")
	parenthesisLeft      = []byte("(")
	parenthesisRight     = []byte(")")
	column               = []byte(":")
	valueNull            = []byte("null")
	valueFalse           = []byte("false")
	valueTrue            = []byte("true")
	operatorGreater      = []byte(">")
	operatorLesser       = []byte("<")
	operatorGreaterEqual = []byte(">=")
	operatorLesserEqual  = []byte("<=")
	operatorEqual        = []byte("=")
	operatorNotEqual     = []byte("!=")
)

func (s source) expectEOF() (e Error) {
	if len(s.s) > 0 {
		e.Index = s.index()
		e.Msg = "expected EOF"
	}
	return
}

// consumeIrrelevant skips spaces, tabs, line-feeds
// carriage-returns and comment sequences
func (s source) consumeIrrelevant() source {
MAIN:
	for len(s.s) > 0 {
		if s.s[0] == '#' {
			s.s = s.s[1:]
			for len(s.s) > 0 {
				if s.s[0] == '\n' {
					continue MAIN
				}
				s.s = s.s[1:]
			}
		} else if s.s[0] == ' ' ||
			s.s[0] == '\n' ||
			s.s[0] == '\t' ||
			s.s[0] == '\r' ||
			s.s[0] == ',' {
			s.s = s.s[1:]
			continue
		}
		break
	}
	return s
}

func (s source) consume(x []byte) (_ source, ok bool) {
	if bytes.HasPrefix(s.s, x) {
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
	i, ii := s.s, s.index()
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
			b == ')' ||
			b == '}' ||
			b == '#' {
			break
		}
		s.s = s.s[1:]
	}
	return s, i[:s.index()-ii]
}

func (s source) consumeName() (_ source, name []byte) {
	i, ii := s.s, s.index()
	if len(s.s) < 1 {
		return s, nil
	}
	if b := s.s[0]; b == '_' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') {
		s.s = s.s[1:]
	} else {
		return s, nil
	}
	for len(s.s) > 0 && (s.s[0] == '_' ||
		(s.s[0] >= 'a' && s.s[0] <= 'z') ||
		(s.s[0] >= 'A' && s.s[0] <= 'Z') ||
		(s.s[0] >= '0' && s.s[0] <= '9')) {
		s.s = s.s[1:]
	}
	return s, i[:s.index()-ii]
}

func (s source) consumeString() (n source, str []byte, ok bool) {
	escaped := false
	n = s

	if len(n.s) < 1 || n.s[0] != '"' {
		return s, nil, false
	}
	n.s = n.s[1:]
	ii := n

	for ; len(n.s) > 0; n.s = n.s[1:] {
		if n.s[0] < 0x20 {
			return s, nil, false
		} else if n.s[0] == '"' {
			if escaped {
				escaped = false
			} else {
				str = ii.s[:n.index()-ii.index()]
				n.s = n.s[1:]
				return n, str, true
			}
		} else if n.s[0] == '\\' {
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
