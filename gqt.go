package gqt

import (
	"bytes"
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
		Name               string
		AssociatedVariable *VariableDefinition
		Constraint         Expression
	}

	// Expression can be either of:
	//
	//	Variable
	//	Int
	//	Float
	//	String
	//	True
	//	False
	//	Null
	//	Enum
	//	Array
	//	Object
	//	ExprModulo
	//	ExprDivision
	//	ExprMultiplication
	//	ExprAddition
	//	ExprSubtraction
	//	ExprEqual
	//	ExprNotEqual
	//  ExprLess
	//  ExprGreater
	//  ExprLessOrEqual
	//  ExprGreaterOrEqual
	//	ExprParentheses
	//  ExprConstrEquals
	//  ExprConstrNotEquals
	//  ExprConstrLess
	//  ExprConstrGreater
	//  ExprConstrLessOrEqual
	//  ExprConstrGreaterOrEqual
	//  ExprConstrLenEquals
	//  ExprConstrLenNotEquals
	//  ExprConstrLenLess
	//  ExprConstrLenGreater
	//  ExprConstrLenLessOrEqual
	//  ExprConstrLenGreaterOrEqual
	//	ExprLogicalAnd
	//	ExprLogicalOr
	//	ExprLogicalNegation
	Expression interface {
		Type() string
	}

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

	Int struct {
		Location
		Value int64
	}

	Float struct {
		Location
		Value float64
	}

	String struct {
		Location
		Value string
	}

	True struct {
		Location
	}

	False struct {
		Location
	}

	Null struct {
		Location
	}

	Enum struct {
		Location
		Value string
	}

	Array struct {
		Location
		Items []Expression
	}

	Object struct {
		Location
		Fields []*ObjectField
	}

	ObjectField struct {
		Location
		Name               string
		AssociatedVariable *VariableDefinition
		Constraint         Expression
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

func (e *ExprParentheses) Type() string {
	if e.Expression != nil {
		return e.Expression.Type()
	}
	return ""
}

func (e *ExprConstrEquals) Type() string            { return "constraint" }
func (e *ExprConstrNotEquals) Type() string         { return "constraint" }
func (e *ExprConstrLess) Type() string              { return "constraint" }
func (e *ExprConstrLessOrEqual) Type() string       { return "constraint" }
func (e *ExprConstrGreater) Type() string           { return "constraint" }
func (e *ExprConstrGreaterOrEqual) Type() string    { return "constraint" }
func (e *ExprConstrLenEquals) Type() string         { return "constraint" }
func (e *ExprConstrLenNotEquals) Type() string      { return "constraint" }
func (e *ExprConstrLenLess) Type() string           { return "constraint" }
func (e *ExprConstrLenLessOrEqual) Type() string    { return "constraint" }
func (e *ExprConstrLenGreater) Type() string        { return "constraint" }
func (e *ExprConstrLenGreaterOrEqual) Type() string { return "constraint" }

func (*ExprLogicalNegation) Type() string { return "boolean" }
func (*ExprModulo) Type() string          { return "number" }
func (*ExprDivision) Type() string        { return "number" }
func (*ExprMultiplication) Type() string  { return "number" }
func (*ExprAddition) Type() string        { return "number" }
func (*ExprSubtraction) Type() string     { return "number" }
func (*Int) Type() string                 { return "number" }
func (*Float) Type() string               { return "number" }

func (*True) Type() string               { return "boolean" }
func (*False) Type() string              { return "boolean" }
func (*ExprEqual) Type() string          { return "boolean" }
func (*ExprNotEqual) Type() string       { return "boolean" }
func (*ExprLess) Type() string           { return "boolean" }
func (*ExprLessOrEqual) Type() string    { return "boolean" }
func (*ExprGreater) Type() string        { return "boolean" }
func (*ExprGreaterOrEqual) Type() string { return "boolean" }
func (*ExprLogicalAnd) Type() string     { return "boolean" }
func (*ExprLogicalOr) Type() string      { return "boolean" }

func (*String) Type() string { return "string" }

func (*Null) Type() string { return "null" }

func (*Enum) Type() string { return "enum" }

func (*Array) Type() string { return "array" }

func (*Object) Type() string { return "object" }

func (*Variable) Type() string { return "variable" }

type VariableDefinition struct {
	Location Location
	Name     string

	// References can be any of:
	//
	//  *Argument
	//  *ObjectField
	References []any

	// Parent can be any of:
	//
	//  *Argument
	//  *ObjectField
	Parent any
}

func Parse(
	src []byte,
) (
	operation *Operation,
	variables map[string]*VariableDefinition,
	err Error,
) {
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
		return nil, nil, errUnexp(
			si, "expected query, mutation, or subscription "+
				"operation definition",
		)
	}

	variables = make(map[string]*VariableDefinition)
	varRefs := []*Variable{}

	s = s.consumeIgnored()
	if s, o.Selections, err = ParseSelectionSet(
		s, variables, &varRefs,
	); err.IsErr() {
		return nil, nil, err
	}

	for _, r := range varRefs {
		if v, ok := variables[r.Name]; !ok {
			sc := s
			sc.Location = r.Location
			return nil, nil, errMsg(sc, "undefined variable")
		} else {
			v.References = append(v.References, r)
		}
	}

	return o, variables, Error{}
}

func errUnexp(s source, msg string) Error {
	prefix := "unexpected token, "
	if s.Index >= len(s.s) {
		prefix = "unexpected end of file, "
	}
	return Error{
		Location: s.Location,
		Msg:      prefix + msg,
	}
}

func errTypef(s source, format string, v ...any) Error {
	return Error{
		Location: s.Location,
		Msg:      fmt.Sprintf("mismatching types: "+format, v...),
	}
}

func errCantCompare(s source, left, right Expression) Error {
	return Error{
		Location: s.Location,
		Msg: fmt.Sprintf(
			"can't compare %s and %s",
			left.Type(), right.Type(),
		),
	}
}

func errMsg(s source, msg string) Error {
	return Error{
		Location: s.Location,
		Msg:      msg,
	}
}

func newErrorNoPrefix(s source, msg string) Error {
	return Error{
		Location: s.Location,
		Msg:      msg,
	}
}

type expect int8

const (
	expectConstraint expect = 0
	expectValue      expect = 1
)

func ParseSelectionSet(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
) (source, []Selection, Error) {
	var ok bool
	si := s
	if s, ok = s.consume("{"); !ok {
		return s, nil, errUnexp(s, "expected selection set")
	}

	fields := map[string]struct{}{}
	typeConds := map[string]struct{}{}

	var selections []Selection
	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			return s, nil, errUnexp(s, "expected selection")
		}

		if s, ok = s.consume("}"); ok {
			break
		}

		s = s.consumeIgnored()

		if _, ok := s.consume("..."); ok {
			var err Error
			var fragInline SelectionInlineFrag
			sBefore := s
			if s, fragInline, err = ParseInlineFrag(
				s, variables, varRefs,
			); err.IsErr() {
				return s, nil, err
			}

			if _, ok := typeConds[fragInline.TypeCondition]; ok {
				return sBefore, nil, errMsg(sBefore, "redeclared type condition")
			}
			typeConds[fragInline.TypeCondition] = struct{}{}

			selections = append(selections, fragInline)
			continue
		}

		sBeforeName := s
		sel := &SelectionField{Location: s.Location}
		if s, name = s.consumeName(); name == nil {
			return s, nil, errUnexp(s, "expected selection")
		}
		sel.Name = string(name)

		if _, ok := fields[sel.Name]; ok {
			return sBeforeName, nil, errMsg(sBeforeName, "redeclared field")
		}
		fields[sel.Name] = struct{}{}

		s = s.consumeIgnored()

		if s.peek1('(') {
			var err Error
			s, sel.Arguments, err = ParseArguments(s, variables, varRefs)
			if err.IsErr() {
				return s, nil, err
			}
		}

		s = s.consumeIgnored()

		if s.peek1('{') {
			var err Error
			s, sel.Selections, err = ParseSelectionSet(s, variables, varRefs)
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

func ParseInlineFrag(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
) (source, SelectionInlineFrag, Error) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("..."); !ok {
		return s, SelectionInlineFrag{}, errUnexp(s, "expected '...'")
	}

	inlineFrag := SelectionInlineFrag{Location: l}
	s = s.consumeIgnored()

	var tok []byte
	sp := s
	s, tok = s.consumeToken()
	if string(tok) != "on" {
		return sp, SelectionInlineFrag{}, errUnexp(sp, "expected keyword 'on'")
	}

	s = s.consumeIgnored()

	var name []byte
	s, name = s.consumeName()
	if len(name) < 1 {
		return s, SelectionInlineFrag{}, errUnexp(s, "expected type condition")
	}
	inlineFrag.TypeCondition = string(name)

	s = s.consumeIgnored()
	var sels []Selection
	var err Error
	if s, sels, err = ParseSelectionSet(s, variables, varRefs); err.IsErr() {
		return s, SelectionInlineFrag{}, err
	}

	inlineFrag.Selections = sels
	return s, inlineFrag, Error{}
}

func ParseArguments(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
) (source, []*Argument, Error) {
	si := s
	var ok bool
	if s, ok = s.consume("("); !ok {
		return s, nil, errUnexp(s, "expected opening parenthesis")
	}

	names := map[string]struct{}{}

	var arguments []*Argument
	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			return s, nil, errUnexp(s, "expected argument")
		}

		if s, ok = s.consume(")"); ok {
			break
		}

		s = s.consumeIgnored()

		sBeforeName := s
		arg := &Argument{Location: s.Location}
		if s, name = s.consumeName(); name == nil {
			return s, nil, errUnexp(s, "expected argument name")
		}
		arg.Name = string(name)

		if _, ok := names[arg.Name]; ok {
			return sBeforeName, nil, errMsg(sBeforeName, "redeclared argument")
		}

		names[arg.Name] = struct{}{}

		s = s.consumeIgnored()

		if s, ok = s.consume("="); ok {
			// Has an associated variable name
			s = s.consumeIgnored()

			sBeforeDollar := s
			if s, ok = s.consume("$"); !ok {
				return s, nil, errUnexp(s, "expected variable name")
			}

			var name []byte
			if s, name = s.consumeName(); name == nil {
				return s, nil, errUnexp(s, "expected variable name")
			}

			def := &VariableDefinition{
				Location: sBeforeDollar.Location,
				Parent:   arg,
				Name:     string(name),
			}
			arg.AssociatedVariable = def
			if _, ok := variables[def.Name]; ok {
				return s, nil, errMsg(sBeforeDollar, "redeclared variable")
			}
			variables[def.Name] = def
			s = s.consumeIgnored()
		}

		if s, ok = s.consume(":"); ok {
			s = s.consumeIgnored()

			if s.isEOF() {
				return s, nil, errUnexp(s, "expected constraint")
			}

			// Has a constraint
			var expr Expression
			var err Error
			s, expr, err = ParseExprLogicalOr(
				s, variables, varRefs, expectConstraint,
			)
			if err.IsErr() {
				return s, nil, err
			}
			arg.Constraint = expr
		}

		arguments = append(arguments, arg)
		s = s.consumeIgnored()

		if s, ok = s.consume(","); !ok {
			if s, ok = s.consume(")"); !ok {
				return s, nil, errUnexp(s, "expected comma or end of argument list")
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

func ParseValue(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	l := s.Location

	if s.isEOF() {
		return s, nil, errUnexp(s, "expected value")
	}

	var str []byte
	var num any
	var ok bool
	var err Error
	if s, ok = s.consume("("); ok {
		e := &ExprParentheses{Location: l}

		s = s.consumeIgnored()

		var err Error
		s, e.Expression, err = ParseExprLogicalOr(s, variables, varRefs, expect)
		if err.IsErr() {
			return s, nil, err
		}

		if s, ok = s.consume(")"); !ok {
			return s, nil, errUnexp(s, "missing closing parenthesis")
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("$"); ok {
		v := &Variable{Location: l}

		var name []byte
		if s, name = s.consumeName(); name == nil {
			return s, nil, errUnexp(s, "expected variable name")
		}

		v.Name = string(name)

		*varRefs = append(*varRefs, v)

		s = s.consumeIgnored()

		return s, v, Error{}
	} else if s, str, ok = s.consumeString(); ok {
		return s, &String{Location: l, Value: string(str)}, Error{}
	}

	if s, num, err = s.ParseNumber(); err.IsErr() {
		return s, nil, err
	} else if num != nil {
		switch v := num.(type) {
		case *Int:
			return s, v, Error{}
		case *Float:
			return s, v, Error{}
		default:
			panic(fmt.Errorf("unexpected number type: %#v", num))
		}
	}

	if s, ok = s.consume("["); ok {
		e := &Array{Location: l}

		for {
			var err Error
			if s, ok = s.consume("]"); ok {
				break
			}

			var expr Expression
			s, expr, err = ParseExprLogicalOr(
				s, variables, varRefs, expectConstraint,
			)
			if err.IsErr() {
				return s, nil, err
			}

			e.Items = append(e.Items, expr)

			s = s.consumeIgnored()

			if s, ok = s.consume(","); !ok {
				if s, ok = s.consume("]"); !ok {
					return s, nil, errUnexp(s, "expected comma or end of array")
				}
				break
			}
			s = s.consumeIgnored()
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("{"); ok {
		// Object
		o := &Object{Location: l}
		si := s
		var ok bool
		fieldNames := map[string]struct{}{}

		for {
			s = s.consumeIgnored()
			var name []byte

			if s.isEOF() {
				return s, nil, errUnexp(s, "expected object field")
			}

			if s, ok = s.consume("}"); ok {
				break
			}

			s = s.consumeIgnored()

			sBeforeName := s
			fld := &ObjectField{Location: s.Location}
			if s, name = s.consumeName(); name == nil {
				return s, nil, errUnexp(s, "expected object field name")
			}
			fld.Name = string(name)

			if _, ok := fieldNames[fld.Name]; ok {
				return sBeforeName, nil, errMsg(
					sBeforeName, "redeclared object field",
				)
			}

			fieldNames[fld.Name] = struct{}{}

			s = s.consumeIgnored()

			if s, ok = s.consume("="); ok {
				// Has an associated variable name
				s = s.consumeIgnored()

				sBeforeDollar := s
				if s, ok = s.consume("$"); !ok {
					return s, nil, errUnexp(s, "expected variable name")
				}

				var name []byte
				if s, name = s.consumeName(); name == nil {
					return s, nil, errUnexp(s, "expected variable name")
				}

				def := &VariableDefinition{
					Location: sBeforeDollar.Location,
					Parent:   fld,
					Name:     string(name),
				}
				fld.AssociatedVariable = def
				if _, ok := variables[def.Name]; ok {
					return s, nil, errMsg(sBeforeDollar, "redeclared variable")
				}
				variables[def.Name] = def

				s = s.consumeIgnored()
			}

			if s, ok = s.consume(":"); ok {
				s = s.consumeIgnored()

				if s.isEOF() {
					return s, nil, errUnexp(s, "expected constraint")
				}

				// Has a constraint
				var expr Expression
				var err Error
				s, expr, err = ParseExprLogicalOr(
					s, variables, varRefs, expectConstraint,
				)
				if err.IsErr() {
					return s, nil, err
				}
				fld.Constraint = expr
			}

			o.Fields = append(o.Fields, fld)
			s = s.consumeIgnored()

			if s, ok = s.consume(","); !ok {
				if s, ok = s.consume("}"); !ok {
					return s, nil, errUnexp(s, "expected comma or end of object")
				}
				break
			}
			s = s.consumeIgnored()
		}

		if len(o.Fields) < 1 {
			return si, nil, newErrorNoPrefix(si, "empty input object")
		}

		return s, o, Error{}
	}

	s, str = s.consumeName()
	switch string(str) {
	case "true":
		return s, &True{Location: l}, Error{}
	case "false":
		return s, &False{Location: l}, Error{}
	case "null":
		return s, &Null{Location: l}, Error{}
	case "":
		return s, nil, errUnexp(s, "invalid value")
	default:
		return s, &Enum{Location: l, Value: string(str)}, Error{}
	}
}

func ParseExprUnary(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	l := s.Location
	var ok bool
	var err Error
	if s, ok = s.consume("!"); ok {
		e := &ExprLogicalNegation{Location: l}

		s = s.consumeIgnored()

		si := s
		s, e.Expression, err = ParseValue(s, variables, varRefs, expect)
		if err.IsErr() {
			return s, nil, err
		}

		if !isBoolean(e.Expression) {
			return s, nil, errTypef(
				si, "can't use %s as boolean", e.Expression.Type(),
			)
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	}
	return ParseValue(s, variables, varRefs, expect)
}

func ParseExprMultiplicative(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s

	var result Expression
	var err Error
	s, result, err = ParseExprUnary(s, variables, varRefs, expect)
	if err.IsErr() {
		return s, nil, err
	}

	for {
		s = s.consumeIgnored()
		var selected int
		s, selected = s.consumeEitherOf3("*", "/", "%")
		switch selected {
		case 0:
			e := &ExprMultiplication{Location: si.Location, Multiplicant: result}

			if !isNumeric(e.Multiplicant) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Multiplicant.Type(),
				)
			}

			s = s.consumeIgnored()
			si := s
			s, e.Multiplicator, err = ParseExprUnary(
				s, variables, varRefs, expect,
			)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(e.Multiplicator) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Multiplicator.Type(),
				)
			}

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprDivision{Location: si.Location, Dividend: result}

			if !isNumeric(e.Dividend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Dividend.Type(),
				)
			}

			s = s.consumeIgnored()
			si := s
			s, e.Divisor, err = ParseExprUnary(
				s, variables, varRefs, expect,
			)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(e.Divisor) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Divisor.Type(),
				)
			}

			s = s.consumeIgnored()
			result = e
		case 2:
			e := &ExprModulo{Location: si.Location, Dividend: result}

			if !isNumeric(e.Dividend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Dividend.Type(),
				)
			}

			s = s.consumeIgnored()
			si := s
			s, e.Divisor, err = ParseExprUnary(s, variables, varRefs, expect)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(e.Divisor) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Divisor.Type(),
				)
			}

			s = s.consumeIgnored()
			result = e
		default:
			s = s.consumeIgnored()
			return s, result, Error{}
		}
	}
}

func ParseExprAdditive(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s

	var result Expression
	var err Error
	s, result, err = ParseExprMultiplicative(s, variables, varRefs, expect)
	if err.IsErr() {
		return s, nil, err
	}

	for {
		s = s.consumeIgnored()
		var selected int
		s, selected = s.consumeEitherOf3("+", "-", "")
		switch selected {
		case 0:
			e := &ExprAddition{Location: si.Location, AddendLeft: result}

			if !isNumeric(e.AddendLeft) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.AddendLeft.Type(),
				)
			}

			s = s.consumeIgnored()
			si := s
			s, e.AddendRight, err = ParseExprMultiplicative(
				s, variables, varRefs, expect,
			)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(e.AddendRight) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.AddendRight.Type(),
				)
			}

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprSubtraction{Location: si.Location, Minuend: result}

			if !isNumeric(e.Minuend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Minuend.Type(),
				)
			}

			s = s.consumeIgnored()
			si := s
			s, e.Subtrahend, err = ParseExprMultiplicative(
				s, variables, varRefs, expect,
			)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(e.Subtrahend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Subtrahend.Type(),
				)
			}

			s = s.consumeIgnored()
			result = e
		default:
			s = s.consumeIgnored()
			return s, result, Error{}
		}
	}
}

func ParseExprRelational(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	l := s.Location

	var left Expression
	var err Error
	s, left, err = ParseExprAdditive(
		s, variables, varRefs, expect,
	)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("<="); ok {
		e := &ExprLessOrEqual{Location: l, Left: left}

		s = s.consumeIgnored()

		s, e.Right, err = ParseExprAdditive(
			s, variables, varRefs, expect,
		)
		if err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume(">="); ok {
		e := &ExprGreaterOrEqual{Location: l, Left: left}

		s = s.consumeIgnored()

		s, e.Right, err = ParseExprAdditive(
			s, variables, varRefs, expect,
		)
		if err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("<"); ok {
		e := &ExprLess{Location: l, Left: left}

		s = s.consumeIgnored()

		s, e.Right, err = ParseExprAdditive(
			s, variables, varRefs, expect,
		)
		if err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume(">"); ok {
		e := &ExprGreater{Location: l, Left: left}

		s = s.consumeIgnored()

		si := s
		s, e.Right, err = ParseExprAdditive(
			s, variables, varRefs, expect,
		)
		if err.IsErr() {
			return s, nil, err
		}

		if !isNumeric(e.Right) {
			return s, nil, errTypef(si, "can't use %s as number", e.Type())
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	}
	s = s.consumeIgnored()
	return s, left, Error{}
}

func ParseExprEquality(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s

	var err Error
	var exprLeft Expression
	s, exprLeft, err = ParseExprRelational(s, variables, varRefs, expect)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("=="); ok {
		e := &ExprEqual{Location: si.Location, Left: exprLeft}

		s = s.consumeIgnored()

		s, e.Right, err = ParseExprRelational(s, variables, varRefs, expect)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()

		if err = assumeSameType(si, e.Left, e.Right); err.IsErr() {
			return si, nil, err
		}

		return s, e, Error{}
	}
	if s, ok = s.consume("!="); ok {
		e := &ExprNotEqual{Location: si.Location, Left: exprLeft}

		s = s.consumeIgnored()

		s, e.Right, err = ParseExprRelational(s, variables, varRefs, expect)
		if err.IsErr() {
			return s, nil, err
		}
		s = s.consumeIgnored()

		if err = assumeSameType(si, e.Left, e.Right); err.IsErr() {
			return si, nil, err
		}

		return s, e, Error{}
	}

	s = s.consumeIgnored()
	return s, exprLeft, Error{}
}

func ParseExprLogicalOr(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	e := &ExprLogicalOr{Location: s.Location}

	var err Error
	for {
		var expr Expression
		s, expr, err = ParseExprLogicalAnd(s, variables, varRefs, expect)
		if err.IsErr() {
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

func ParseExprLogicalAnd(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	e := &ExprLogicalAnd{Location: s.Location}

	var err Error
	for {
		var expr Expression
		s, expr, err = ParseExprConstr(s, variables, varRefs, expect)
		if err.IsErr() {
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

func ParseExprConstr(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	l := s.Location
	var ok bool
	var expr Expression
	var err Error

	if expect == expectValue {
		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, expr, Error{}
	}

	if s, ok = s.consume("!="); ok {
		s = s.consumeIgnored()

		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, &ExprConstrNotEquals{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume("<="); ok {
		s = s.consumeIgnored()

		si := s
		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.Type(),
			)
		}

		s = s.consumeIgnored()

		return s, &ExprConstrLessOrEqual{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume(">="); ok {
		s = s.consumeIgnored()

		si := s
		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.Type(),
			)
		}

		s = s.consumeIgnored()

		return s, &ExprConstrGreaterOrEqual{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume("<"); ok {
		s = s.consumeIgnored()

		si := s
		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.Type(),
			)
		}

		s = s.consumeIgnored()

		return s, &ExprConstrLess{
			Location: l,
			Value:    expr,
		}, Error{}
	} else if s, ok = s.consume(">"); ok {
		s = s.consumeIgnored()

		si := s
		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.Type(),
			)
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

			si := s
			s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.Type(),
				)
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenNotEquals{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume("<="); ok {
			s = s.consumeIgnored()

			si := s
			s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.Type(),
				)
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenLessOrEqual{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume(">="); ok {
			s = s.consumeIgnored()

			si := s
			s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.Type(),
				)
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenGreaterOrEqual{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume("<"); ok {
			s = s.consumeIgnored()

			si := s
			s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.Type(),
				)
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenLess{
				Location: l,
				Value:    expr,
			}, Error{}
		} else if s, ok = s.consume(">"); ok {
			s = s.consumeIgnored()

			si := s
			s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
			if err.IsErr() {
				return s, nil, err
			}

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.Type(),
				)
			}

			s = s.consumeIgnored()

			return s, &ExprConstrLenGreater{
				Location: l,
				Value:    expr,
			}, Error{}
		}

		si := s
		s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
		if err.IsErr() {
			return s, nil, err
		}

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.Type(),
			)
		}

		s = s.consumeIgnored()

		return s, &ExprConstrGreater{
			Location: l,
			Value:    expr,
		}, Error{}
	}

	s, expr, err = ParseExprEquality(s, variables, varRefs, expectValue)
	if err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	return s, &ExprConstrEquals{
		Location: l,
		Value:    expr,
	}, Error{}
}

func (s source) ParseNumber() (_ source, number any, err Error) {
	si := s
	if s.Index >= len(s.s) {
		return si, nil, Error{}
	}

	if s.s[s.Index] == '-' || s.s[s.Index] == '+' {
		s.Index++
		s.Column++
	}

	for s.Index < len(s.s) {
		if (s.s[s.Index] < '0' || s.s[s.Index] > '9') &&
			s.s[s.Index] != '.' &&
			s.s[s.Index] != 'e' &&
			s.s[s.Index] != 'E' {
			break
		}
		if s.s[s.Index] == 'e' || s.s[s.Index] == 'E' {
			if s.Index == si.Index {
				return si, nil, Error{}
			}
			s.Index++
			s.Column++
			if s.Index >= len(s.s) {
				return si, nil, errMsg(s, "exponent has no digits")
			}
			s, _ = s.consumeEitherOf3("+", "-", "")
			if !s.isDigit() {
				return si, nil, errMsg(s, "exponent has no digits")
			}
		}

		s.Index++
		s.Column++
	}
	str := string(s.s[si.Index:s.Index])
	if str == "" {
		return si, nil, Error{}
	} else if v, err := strconv.ParseInt(str, 10, 64); err == nil {
		return s, &Int{
			Location: si.Location,
			Value:    v,
		}, Error{}
	} else if v, err := strconv.ParseFloat(str, 64); err == nil {
		return s, &Float{
			Location: si.Location,
			Value:    v,
		}, Error{}
	}
	return si, nil, Error{}
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

func isNumeric(expr any) bool {
	switch e := expr.(type) {
	case *Variable:
		// TODO: check schema
		return true
	case *ExprParentheses:
		return isNumeric(e.Expression)
	case *Float:
		return true
	case *Int:
		return true
	case *ExprAddition:
		return true
	case *ExprSubtraction:
		return true
	case *ExprMultiplication:
		return true
	case *ExprDivision:
		return true
	case *ExprModulo:
		return true
	}
	return false
}

func isBoolean(expr any) bool {
	switch e := expr.(type) {
	case *Variable:
		// TODO: check schema
		return true
	case *ExprParentheses:
		return isBoolean(e.Expression)
	case *True:
		return true
	case *False:
		return true
	case *ExprEqual:
		return true
	case *ExprNotEqual:
		return true
	case *ExprLess:
		return true
	case *ExprGreater:
		return true
	case *ExprLessOrEqual:
		return true
	case *ExprGreaterOrEqual:
		return true
	case *ExprLogicalNegation:
		return true
	case *ExprLogicalAnd:
		return true
	case *ExprLogicalOr:
		return true
	}
	return false
}

func assumeSameType(s source, left, right Expression) Error {
	_, lIsVar := left.(*Variable)
	_, rIsVar := right.(*Variable)
	if lIsVar || rIsVar {
		// TODO check variable types
		return Error{}
	}

	if isBoolean(left) && !isBoolean(right) {
		return errCantCompare(s, left, right)
	} else if isNumeric(left) && !isNumeric(right) {
		return errCantCompare(s, left, right)
	}
	switch left.(type) {
	case *String:
		if _, ok := right.(*String); !ok {
			return errCantCompare(s, left, right)
		}
	case *Enum:
		if _, ok := right.(*Enum); !ok {
			return errCantCompare(s, left, right)
		}
	case *Array:
		if _, ok := right.(*Array); !ok {
			return errCantCompare(s, left, right)
		}
	case *Object:
		if _, ok := right.(*Object); !ok {
			return errCantCompare(s, left, right)
		}
	}
	return Error{}
}
