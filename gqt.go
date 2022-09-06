package gqt

import (
	"fmt"
	"strconv"
)

type OperationType int8

const (
	_ OperationType = iota
	OperationTypeQuery
	OperationTypeMutation
	OperationTypeSubscription
)

func (t OperationType) String() string {
	switch t {
	case OperationTypeQuery:
		return "query"
	case OperationTypeMutation:
		return "mutation"
	case OperationTypeSubscription:
		return "subscription"
	}
	return ""
}

type (
	Location struct {
		Index, Line, Column int
	}

	Operation struct {
		Location
		Type       OperationType
		Selections []Selection
	}

	// Selection can be either of:
	//
	//	SelectionField
	//	SelectionMax
	//	SelectionInlineFrag
	Selection any

	SelectionField struct {
		Location
		Name       string
		Arguments  []*Argument
		Selections []Selection
	}

	Argument struct {
		Location
		Name                   string
		AssociatedVariableName string
		Constraint             Expression
	}

	// Expression can be either of:
	//
	//	ValueInt
	//	ValueFloat
	//	ValueString
	//	ValueTrue
	//	ValueFalse
	//	ValueNull
	//	ValueEnum
	//	ValueArray
	//	ValueObject
	//	ExprModulo
	//	ExprDivision
	//	ExprMultiplication
	//	ExprAddition
	//	ExprSubtraction
	//	ExprEqual
	//	ExprNotEqual
	//  ExprLessOrEqual
	//  ExprGreaterOrEqual
	//  ExprLess
	//  ExprGreater
	//	ExprLogicalNegation
	//	Variable
	//	ExprParentheses
	//  ExprConstrEquals
	//  ExprConstrNotEquals
	//  ExprConstrLess
	//  ExprConstrLessOrEqual
	//  ExprConstrGreater
	//  ExprConstrGreaterOrEqual
	//  ExprConstrLenEquals
	//  ExprConstrLenNotEquals
	//  ExprConstrLenLess
	//  ExprConstrLenLessOrEqual
	//  ExprConstrLenGreater
	//  ExprConstrLenGreaterOrEqual
	//	ExprLogicalAnd
	//	ExprLogicalOr
	Expression any

	ExprConstrEquals struct {
		Location
		Value Expression
	}

	ExprConstrNotEquals struct {
		Location
		Value Expression
	}

	ExprConstrLess struct {
		Location
		Value Expression
	}

	ExprConstrLessOrEqual struct {
		Location
		Value Expression
	}

	ExprConstrGreater struct {
		Location
		Value Expression
	}

	ExprConstrGreaterOrEqual struct {
		Location
		Value Expression
	}

	ExprConstrLenEquals struct {
		Location
		Value Expression
	}

	ExprConstrLenNotEquals struct {
		Location
		Value Expression
	}

	ExprConstrLenLess struct {
		Location
		Value Expression
	}

	ExprConstrLenLessOrEqual struct {
		Location
		Value Expression
	}

	ExprConstrLenGreater struct {
		Location
		Value Expression
	}

	ExprConstrLenGreaterOrEqual struct {
		Location
		Value Expression
	}

	ExprParentheses struct {
		Location
		Expression Expression
	}

	ExprLogicalNegation struct {
		Location
		Expression Expression
	}

	ExprModulo struct {
		Location
		Dividend Expression
		Divisor  Expression
	}

	ExprDivision struct {
		Location
		Dividend Expression
		Divisor  Expression
	}

	ExprMultiplication struct {
		Location
		Multiplicant  Expression
		Multiplicator Expression
	}

	ExprAddition struct {
		Location
		AddendLeft  Expression
		AddendRight Expression
	}

	ExprSubtraction struct {
		Location
		Minuend    Expression
		Subtrahend Expression
	}

	ExprEqual struct {
		Location
		Left  Expression
		Right Expression
	}

	ExprNotEqual struct {
		Location
		Left  Expression
		Right Expression
	}

	ExprLess struct {
		Location
		Left  Expression
		Right Expression
	}

	ExprLessOrEqual struct {
		Location
		Left  Expression
		Right Expression
	}

	ExprGreater struct {
		Location
		Left  Expression
		Right Expression
	}

	ExprGreaterOrEqual struct {
		Location
		Left  Expression
		Right Expression
	}

	ExprLogicalAnd struct {
		Location
		Expressions []Expression
	}

	ExprLogicalOr struct {
		Location
		Expressions []Expression
	}

	ValueInt struct {
		Location
		Value int64
	}

	ValueFloat struct {
		Location
		Value float64
	}

	ValueString struct {
		Location
		Value string
	}

	ValueTrue struct {
		Location
	}

	ValueFalse struct {
		Location
	}

	ValueNull struct {
		Location
	}

	ValueEnum struct {
		Location
		Value string
	}

	ValueArray struct {
		Location
		Items []Expression
	}

	ValueObject struct {
		Location
		Fields []ObjectField
	}

	ObjectField struct {
		Location
		Name  string
		Value Expression
	}

	Variable struct {
		Location
		Name string
	}

	SelectionMax struct {
		Location
		Limit      int
		Selections []Selection
	}

	SelectionInlineFrag struct {
		Location
		TypeCondition string
		Selections    []Selection
	}
)

func Parse(
	src []byte,
) (*Operation, Error) {
	s := source{
		Location: Location{
			Line:   1,
			Column: 1,
		},
		s: src,
	}

	s = s.consumeIgnored()
	o := &Operation{Location: s.Location}

	var tok []byte
	si := s
	s, tok = s.consumeToken()

	if string(tok) == "query" {
		o.Type = OperationTypeQuery
	} else if string(tok) == "mutation" {
		o.Type = OperationTypeMutation
	} else if string(tok) == "subscription" {
		o.Type = OperationTypeSubscription
	} else {
		return nil, newError(
			si, "expected query, mutation, or subscription "+
				"operation definition",
		)
	}

	s = s.consumeIgnored()
	var err Error
	s, o.Selections, err = ParseSelectionSet(s)
	if err.IsErr() {
		return nil, err
	}

	return o, Error{}
}

func newError(s source, msg string) Error {
	prefix := "unexpected token, "
	if s.Index >= len(s.s) {
		prefix = "unexpected end of file, "
	}
	return Error{
		Location: s.Location,
		Msg:      prefix + msg,
	}
}

func newErrorNoPrefix(s source, msg string) Error {
	return Error{
		Location: s.Location,
		Msg:      msg,
	}
}

func ParseSelectionSet(s source) (source, []Selection, Error) {
	var ok bool
	si := s
	if s, ok = s.consume("{"); !ok {
		return s, nil, newError(s, "expected selection set")
	}

	var selections []Selection
	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			return s, nil, newError(s, "expected selection")
		}

		if s, ok = s.consume("}"); ok {
			break
		}

		s = s.consumeIgnored()

		if _, ok := s.consume("..."); ok {
			var err Error
			var fragInline SelectionInlineFrag
			if s, fragInline, err = ParseInlineFrag(s); err.IsErr() {
				return s, nil, err
			}
			selections = append(selections, fragInline)
			continue
		}

		sel := &SelectionField{Location: s.Location}
		if s, name = s.consumeName(); name == nil {
			return s, nil, newError(s, "expected selection")
		}
		sel.Name = string(name)

		s = s.consumeIgnored()

		if s.peek1('(') {
			var err Error
			s, sel.Arguments, err = ParseArguments(s)
			if err.IsErr() {
				return s, nil, err
			}
		}

		s = s.consumeIgnored()

		if s.peek1('{') {
			var err Error
			s, sel.Selections, err = ParseSelectionSet(s)
			if err.IsErr() {
				return s, nil, err
			}
		}

		selections = append(selections, sel)
	}

	if len(selections) < 1 {
		return si, nil, newErrorNoPrefix(si, "empty selection set")
	}

	return s, selections, Error{}
}

func ParseInlineFrag(s source) (source, SelectionInlineFrag, Error) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("..."); !ok {
		return s, SelectionInlineFrag{}, newError(s, "expected '...'")
	}

	inlineFrag := SelectionInlineFrag{Location: l}
	s = s.consumeIgnored()

	var tok []byte
	sp := s
	s, tok = s.consumeToken()
	if string(tok) != "on" {
		return sp, SelectionInlineFrag{}, newError(sp, "expected keyword 'on'")
	}

	s = s.consumeIgnored()

	var name []byte
	s, name = s.consumeName()
	if len(name) < 1 {
		return s, SelectionInlineFrag{}, newError(s, "expected type condition")
	}
	inlineFrag.TypeCondition = string(name)

	s = s.consumeIgnored()
	var sels []Selection
	var err Error
	if s, sels, err = ParseSelectionSet(s); err.IsErr() {
		return s, SelectionInlineFrag{}, err
	}

	inlineFrag.Selections = sels
	return s, inlineFrag, Error{}
}

func ParseArguments(s source) (source, []*Argument, Error) {
	si := s
	var ok bool
	if s, ok = s.consume("("); !ok {
		return s, nil, newError(s, "expected opening parenthesis")
	}

	var arguments []*Argument
	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			return s, nil, newError(s, "expected argument")
		}

		if s, ok = s.consume(")"); ok {
			break
		}

		s = s.consumeIgnored()

		arg := Argument{Location: s.Location}
		if s, name = s.consumeName(); name == nil {
			return s, nil, newError(s, "expected argument name")
		}
		arg.Name = string(name)

		s = s.consumeIgnored()

		if s, ok = s.consume("="); ok {
			// Has an associated variable name
			s = s.consumeIgnored()

			if s, ok = s.consume("$"); !ok {
				return s, nil, newError(s, "expected variable name")
			}

			var name []byte
			if s, name = s.consumeName(); name == nil {
				return s, nil, newError(s, "expected variable name")
			}

			arg.AssociatedVariableName = string(name)
			s = s.consumeIgnored()
		}

		if s, ok = s.consume(":"); ok {
			s = s.consumeIgnored()

			// Has a constraint
			var expr Expression
			var err Error
			if s, expr, err = ParseExprLogicalOr(s); err.IsErr() {
				return s, nil, newError(s, "expected constraint expression")
			}
			arg.Constraint = expr
		}

		arguments = append(arguments, &arg)
		s = s.consumeIgnored()

		if s, ok = s.consume(","); !ok {
			if s, ok = s.consume(")"); !ok {
				return s, nil, newError(s, "expected comma or end of argument list")
			}
			break
		}
		s = s.consumeIgnored()
	}

	if len(arguments) < 1 {
		return si, nil, newErrorNoPrefix(si, "empty argument list")
	}

	return s, arguments, Error{}
}

func ParseValue(s source) (source, Expression, Error) {
	l := s.Location

	var str []byte
	var num any
	var ok bool
	if s, ok = s.consume("("); ok {
		e := &ExprParentheses{Location: l}

		var err Error
		s, e.Expression, err = ParseExprLogicalOr(s)
		if err.IsErr() {
			return s, nil, err
		}

		if s, ok = s.consume(")"); !ok {
			return s, nil, newError(s, "missing closing parenthesis")
		}

		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume("$"); ok {
		v := &Variable{Location: l}

		var name []byte
		if s, name = s.consumeName(); name == nil {
			return s, nil, newError(s, "expected variable name")
		}

		v.Name = string(name)

		s = s.consumeIgnored()

		return s, v, Error{}
	} else if s, str, ok = s.consumeString(); ok {
		return s, &ValueString{Location: l, Value: string(str)}, Error{}
	} else if s, num, ok = s.consumeNumber(); ok {
		return s, num, Error{}
	} else if s, ok = s.consume("["); ok {
		e := &ValueArray{Location: l}

		for {
			var err Error
			if s, ok = s.consume("]"); ok {
				break
			}

			var expr Expression
			if s, expr, err = ParseExprLogicalOr(s); err.IsErr() {
				return s, nil, err
			}

			e.Items = append(e.Items, expr)

			s = s.consumeIgnored()

			if s, ok = s.consume(","); !ok {
				if s, ok = s.consume("]"); !ok {
					return s, nil, newError(s, "expected comma or end of array")
				}
				break
			}
			s = s.consumeIgnored()
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("{"); ok {
		//TODO
		// e := &ValueObject{}
		panic("object aren't supported yet")
	}

	s, str = s.consumeName()
	switch string(str) {
	case "true":
		return s, &ValueTrue{Location: l}, Error{}
	case "false":
		return s, &ValueFalse{Location: l}, Error{}
	case "null":
		return s, &ValueNull{Location: l}, Error{}
	case "":
		return s, nil, newError(s, "invalid value")
	default:
		return s, &ValueEnum{Location: l, Value: string(str)}, Error{}
	}
}

func ParseExprUnary(s source) (source, Expression, Error) {
	l := s.Location
	var ok bool
	var err Error
	if s, ok = s.consume("!"); ok {
		e := &ExprLogicalNegation{Location: l}
		s, e.Expression, err = ParseValue(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	}
	return ParseValue(s)
}

func ParseExprMultiplicative(s source) (source, Expression, Error) {
	l := s.Location

	var left Expression
	var err Error
	s, left, err = ParseExprUnary(s)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("*"); ok {
		e := &ExprMultiplication{Location: l, Multiplicant: left}
		s, e.Multiplicator, err = ParseExprUnary(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("/"); ok {
		e := &ExprDivision{Location: l, Dividend: left}
		s, e.Divisor, err = ParseExprUnary(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("%"); ok {
		e := &ExprModulo{Location: l, Dividend: left}
		s, e.Divisor, err = ParseExprUnary(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	}
	s = s.consumeIgnored()
	return s, left, Error{}
}

func ParseExprAdditive(s source) (source, Expression, Error) {
	l := s.Location

	var left Expression
	var err Error
	s, left, err = ParseExprMultiplicative(s)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("+"); ok {
		e := &ExprAddition{Location: l, AddendLeft: left}
		s, e.AddendRight, err = ParseExprAdditive(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("-"); ok {
		e := &ExprSubtraction{Location: l, Minuend: left}
		s, e.Subtrahend, err = ParseExprAdditive(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	}
	s = s.consumeIgnored()
	return s, left, Error{}
}

func ParseExprRelational(s source) (source, Expression, Error) {
	l := s.Location

	var left Expression
	var err Error
	s, left, err = ParseExprAdditive(s)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("<="); ok {
		e := &ExprLessOrEqual{Location: l, Left: left}
		s, e.Right, err = ParseExprAdditive(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume(">="); ok {
		e := &ExprGreaterOrEqual{Location: l, Left: left}
		s, e.Right, err = ParseExprAdditive(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("<"); ok {
		e := &ExprLess{Location: l, Left: left}
		s, e.Right, err = ParseExprAdditive(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume(">"); ok {
		e := &ExprGreater{Location: l, Left: left}
		s, e.Right, err = ParseExprAdditive(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	}
	s = s.consumeIgnored()
	return s, left, Error{}
}

func ParseExprEquality(s source) (source, Expression, Error) {
	e := &ExprEqual{Location: s.Location}

	var err Error
	s, e.Left, err = ParseExprRelational(s)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("=="); ok {
		s, e.Right, err = ParseExprRelational(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	}
	if s, ok = s.consume("!="); ok {
		s, e.Right, err = ParseExprRelational(s)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()
		return s, e, Error{}
	}
	s = s.consumeIgnored()
	return s, e.Left, Error{}
}

func ParseExprLogicalOr(s source) (source, Expression, Error) {
	e := &ExprLogicalOr{Location: s.Location}

	var err Error
	for {
		var expr Expression
		if s, expr, err = ParseExprLogicalAnd(s); err.IsErr() {
			return s, nil, err
		}
		e.Expressions = append(e.Expressions, expr)

		s = s.consumeIgnored()

		var ok bool
		if s, ok = s.consume("||"); !ok {
			if len(e.Expressions) < 2 {
				return s, e.Expressions[0], Error{}
			}
			return s, e, Error{}
		}

		s = s.consumeIgnored()
	}
}

func ParseExprLogicalAnd(s source) (source, Expression, Error) {
	e := &ExprLogicalAnd{Location: s.Location}

	var err Error
	for {
		var expr Expression
		if s, expr, err = ParseExprConstr(s); err.IsErr() {
			return s, nil, err
		}
		e.Expressions = append(e.Expressions, expr)

		s = s.consumeIgnored()

		var ok bool
		if s, ok = s.consume("&&"); !ok {
			if len(e.Expressions) < 2 {
				return s, e.Expressions[0], Error{}
			}
			return s, e, Error{}
		}

		s = s.consumeIgnored()
	}
}

func ParseExprConstr(s source) (source, Expression, Error) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("!="); ok {
		s = s.consumeIgnored()

		var expr Expression
		var err Error
		if s, expr, err = ParseExprEquality(s); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrNotEquals{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume("<="); ok {
		s = s.consumeIgnored()

		var expr Expression
		var err Error
		if s, expr, err = ParseExprEquality(s); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrLessOrEqual{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume(">="); ok {
		s = s.consumeIgnored()

		var expr Expression
		var err Error
		if s, expr, err = ParseExprEquality(s); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrGreaterOrEqual{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume("<"); ok {
		s = s.consumeIgnored()

		var expr Expression
		var err Error
		if s, expr, err = ParseExprEquality(s); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrLess{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume(">"); ok {
		s = s.consumeIgnored()

		var expr Expression
		var err Error
		if s, expr, err = ParseExprEquality(s); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrGreater{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume("len"); ok {
		s = s.consumeIgnored()

		if s, ok = s.consume("!="); ok {
			s = s.consumeIgnored()

			var expr Expression
			var err Error
			if s, expr, err = ParseExprEquality(s); err.IsErr() {
				return s, nil, err
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenNotEquals{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume("<="); ok {
			s = s.consumeIgnored()

			var expr Expression
			var err Error
			if s, expr, err = ParseExprEquality(s); err.IsErr() {
				return s, nil, err
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenLessOrEqual{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume(">="); ok {
			s = s.consumeIgnored()

			var expr Expression
			var err Error
			if s, expr, err = ParseExprEquality(s); err.IsErr() {
				return s, nil, err
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenGreaterOrEqual{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume("<"); ok {
			s = s.consumeIgnored()

			var expr Expression
			var err Error
			if s, expr, err = ParseExprEquality(s); err.IsErr() {
				return s, nil, err
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenLess{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume(">"); ok {
			s = s.consumeIgnored()

			var expr Expression
			var err Error
			if s, expr, err = ParseExprEquality(s); err.IsErr() {
				return s, nil, err
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenGreater{
				Location: l,
				Value:    expr,
			}, Error{}
		}

		var expr Expression
		var err Error
		if s, expr, err = ParseExprEquality(s); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrGreater{
			Location: l,
			Value:    expr,
		}, Error{}
	}

	var expr Expression
	var err Error
	if s, expr, err = ParseExprEquality(s); err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	return s, &ExprConstrEquals{
		Location: l,
		Value:    expr,
	}, Error{}
}

type Error struct {
	Location
	Msg string
}

func (e Error) IsErr() bool {
	return e.Msg != ""
}

func (e Error) Error() string {
	if !e.IsErr() {
		return ""
	}
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Column, e.Msg)
}

type source struct {
	Location
	s []byte
}

func (s source) String() string {
	c := s.s[s.Index:]
	if len(c) > 10 {
		c = c[:10]
	}
	return fmt.Sprintf("%d:%d: %q", s.Line, s.Column, c)
}

func (s source) isEOF() bool {
	return s.Index >= len(s.s)
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

func (s source) consumeNumber() (_ source, number any, ok bool) {
	loc := s.Location
	if s, token := s.consumeToken(); token != nil {
		if i, err := strconv.ParseInt(string(token), 10, 64); err == nil {
			return s, &ValueInt{
				Location: loc,
				Value:    i,
			}, true
		} else if f, err := strconv.ParseFloat(string(token), 64); err == nil {
			return s, &ValueFloat{
				Location: loc,
				Value:    f,
			}, true
		}
	}
	return s, number, false
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
