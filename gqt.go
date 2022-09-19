package gqt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	gqlparser "github.com/vektah/gqlparser/v2"
	ast "github.com/vektah/gqlparser/v2/ast"
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
	Location struct{ Index, Line, Column int }

	Operation struct {
		Location
		Type       OperationType
		Selections []Selection
	}

	// Selection can be either of:
	//
	//	*SelectionField
	//	*SelectionMax
	//	*SelectionInlineFrag
	Selection any

	SelectionField struct {
		Location
		Parent     any
		Name       string
		Arguments  []*Argument
		Selections []Selection
	}

	Argument struct {
		Location
		Parent             any
		Name               string
		AssociatedVariable *VariableDefinition
		Constraint         Expression
	}

	// Expression can be either of:
	//
	//	*Variable
	//	*Int
	//	*Float
	//	*String
	//	*True
	//	*False
	//	*Null
	//	*Enum
	//	*Array
	//	*Object
	//	*ExprModulo
	//	*ExprDivision
	//	*ExprMultiplication
	//	*ExprAddition
	//	*ExprSubtraction
	//	*ExprEqual
	//	*ExprNotEqual
	//  *ExprLess
	//  *ExprGreater
	//  *ExprLessOrEqual
	//  *ExprGreaterOrEqual
	//	*ExprParentheses
	//	*ExprLogicalAnd
	//	*ExprLogicalOr
	//	*ExprLogicalNegation
	//  *ConstrEquals
	//  *ConstrNotEquals
	//  *ConstrLess
	//  *ConstrGreater
	//  *ConstrLessOrEqual
	//  *ConstrGreaterOrEqual
	//  *ConstrLenEquals
	//  *ConstrLenNotEquals
	//  *ConstrLenLess
	//  *ConstrLenGreater
	//  *ConstrLenLessOrEqual
	//  *ConstrLenGreaterOrEqual
	Expression interface {
		GetLocation() Location
		ValueType(p *Parser) string
	}

	ConstrAny struct {
		Location
		Parent any
	}

	ConstrEquals struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrNotEquals struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLess struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLessOrEqual struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrGreater struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrGreaterOrEqual struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLenEquals struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLenNotEquals struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLenLess struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLenLessOrEqual struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLenGreater struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrLenGreaterOrEqual struct {
		Location
		Parent any
		Value  Expression
	}

	ConstrMap struct {
		Location
		Parent     any
		Constraint Expression
	}

	ExprParentheses struct {
		Location
		Parent     any
		Expression Expression
	}

	ExprLogicalNegation struct {
		Location
		Parent     any
		Expression Expression
	}

	ExprModulo struct {
		Location
		Parent   any
		Dividend Expression
		Divisor  Expression
	}

	ExprDivision struct {
		Location
		Parent   any
		Dividend Expression
		Divisor  Expression
	}

	ExprMultiplication struct {
		Location
		Parent        any
		Multiplicant  Expression
		Multiplicator Expression
	}

	ExprAddition struct {
		Location
		Parent      any
		AddendLeft  Expression
		AddendRight Expression
	}

	ExprSubtraction struct {
		Location
		Parent     any
		Minuend    Expression
		Subtrahend Expression
	}

	ExprEqual struct {
		Location
		Parent any
		Left   Expression
		Right  Expression
	}

	ExprNotEqual struct {
		Location
		Parent any
		Left   Expression
		Right  Expression
	}

	ExprLess struct {
		Location
		Parent any
		Left   Expression
		Right  Expression
	}

	ExprLessOrEqual struct {
		Location
		Parent any
		Left   Expression
		Right  Expression
	}

	ExprGreater struct {
		Location
		Parent any
		Left   Expression
		Right  Expression
	}

	ExprGreaterOrEqual struct {
		Location
		Parent any
		Left   Expression
		Right  Expression
	}

	ExprLogicalAnd struct {
		Location
		Parent      any
		Expressions []Expression
	}

	ExprLogicalOr struct {
		Location
		Parent      any
		Expressions []Expression
	}

	Int struct {
		Location
		Parent any
		Value  int64
	}

	Float struct {
		Location
		Parent any
		Value  float64
	}

	String struct {
		Location
		Parent any
		Value  string
	}

	True struct {
		Location
		Parent any
	}

	False struct {
		Location
		Parent any
	}

	Null struct {
		Location
		Parent any
	}

	Enum struct {
		Location
		Parent any
		Value  string
	}

	Array struct {
		Location
		Parent any
		Items  []Expression
	}

	Object struct {
		Location
		Parent any
		Fields []*ObjectField
	}

	ObjectField struct {
		Location
		Parent             any
		Name               string
		AssociatedVariable *VariableDefinition
		Constraint         Expression
	}

	Variable struct {
		Location
		Parent any
		Name   string
	}

	SelectionMax struct {
		Location
		Parent  any
		Limit   int
		Options []Selection
	}

	SelectionInlineFrag struct {
		Location
		Parent        any
		TypeCondition string
		Selections    []Selection
	}
)

func (e *ExprParentheses) GetLocation() Location         { return e.Location }
func (e *ConstrAny) GetLocation() Location               { return e.Location }
func (e *ConstrEquals) GetLocation() Location            { return e.Location }
func (e *ConstrNotEquals) GetLocation() Location         { return e.Location }
func (e *ConstrLess) GetLocation() Location              { return e.Location }
func (e *ConstrLessOrEqual) GetLocation() Location       { return e.Location }
func (e *ConstrGreater) GetLocation() Location           { return e.Location }
func (e *ConstrGreaterOrEqual) GetLocation() Location    { return e.Location }
func (e *ConstrLenEquals) GetLocation() Location         { return e.Location }
func (e *ConstrLenNotEquals) GetLocation() Location      { return e.Location }
func (e *ConstrLenLess) GetLocation() Location           { return e.Location }
func (e *ConstrLenLessOrEqual) GetLocation() Location    { return e.Location }
func (e *ConstrLenGreater) GetLocation() Location        { return e.Location }
func (e *ConstrLenGreaterOrEqual) GetLocation() Location { return e.Location }

func (e *ExprModulo) GetLocation() Location         { return e.Location }
func (e *ExprDivision) GetLocation() Location       { return e.Location }
func (e *ExprMultiplication) GetLocation() Location { return e.Location }
func (e *ExprAddition) GetLocation() Location       { return e.Location }
func (e *ExprSubtraction) GetLocation() Location    { return e.Location }

func (e *ExprLogicalNegation) GetLocation() Location { return e.Location }
func (e *ExprEqual) GetLocation() Location           { return e.Location }
func (e *ExprNotEqual) GetLocation() Location        { return e.Location }
func (e *ExprLess) GetLocation() Location            { return e.Location }
func (e *ExprLessOrEqual) GetLocation() Location     { return e.Location }
func (e *ExprGreater) GetLocation() Location         { return e.Location }
func (e *ExprGreaterOrEqual) GetLocation() Location  { return e.Location }
func (e *ExprLogicalAnd) GetLocation() Location      { return e.Location }
func (e *ExprLogicalOr) GetLocation() Location       { return e.Location }

func (e *True) GetLocation() Location      { return e.Location }
func (e *False) GetLocation() Location     { return e.Location }
func (e *Int) GetLocation() Location       { return e.Location }
func (e *Float) GetLocation() Location     { return e.Location }
func (e *String) GetLocation() Location    { return e.Location }
func (e *Null) GetLocation() Location      { return e.Location }
func (e *Enum) GetLocation() Location      { return e.Location }
func (e *Array) GetLocation() Location     { return e.Location }
func (e *ConstrMap) GetLocation() Location { return e.Location }
func (e *Object) GetLocation() Location    { return e.Location }
func (e *Variable) GetLocation() Location  { return e.Location }

func (e *ExprParentheses) ValueType(p *Parser) string {
	if e.Expression != nil {
		return e.Expression.ValueType(p)
	}
	return ""
}
func (e *ConstrAny) ValueType(p *Parser) string { return "*" }
func (e *ConstrEquals) ValueType(p *Parser) string {
	return e.Value.ValueType(p)
}
func (e *ConstrNotEquals) ValueType(p *Parser) string {
	return e.Value.ValueType(p)
}
func (e *ConstrLess) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLessOrEqual) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrGreater) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrGreaterOrEqual) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLenEquals) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLenNotEquals) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLenLess) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLenLessOrEqual) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLenGreater) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ConstrLenGreaterOrEqual) ValueType(p *Parser) string {
	return "Int|Float"
}

func (e *ExprModulo) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ExprDivision) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ExprMultiplication) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ExprAddition) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ExprSubtraction) ValueType(p *Parser) string {
	return "Int|Float"
}
func (e *ExprLogicalNegation) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprEqual) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprNotEqual) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprLess) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprLessOrEqual) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprGreater) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprGreaterOrEqual) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprLogicalAnd) ValueType(p *Parser) string {
	return "Boolean"
}
func (e *ExprLogicalOr) ValueType(p *Parser) string {
	return "Boolean"
}

func (e *True) ValueType(p *Parser) string   { return "Boolean" }
func (e *False) ValueType(p *Parser) string  { return "Boolean" }
func (e *Int) ValueType(p *Parser) string    { return "Int" }
func (e *Float) ValueType(p *Parser) string  { return "Float" }
func (e *String) ValueType(p *Parser) string { return "String" }
func (e *Null) ValueType(p *Parser) string   { return "null" }
func (e *Enum) ValueType(p *Parser) string {
	if p.schema == nil {
		return "enum"
	}

	t := p.enumVal[e.Value]
	if t == nil {
		return "enum"
	}
	return t.Name
}
func (e *Array) ValueType(p *Parser) string {
	if p.schema == nil {
		return "array"
	}

	var b strings.Builder
	notNull := true
	foundType := false

	b.WriteRune('[')
	for _, i := range e.Items {
		vt := i.ValueType(p)
		if vt == "null" {
			notNull = false
			continue
		}
		if !foundType {
			b.WriteString(vt)
			foundType = true
		}
	}
	if notNull {
		b.WriteRune('!')
	}
	b.WriteRune(']')
	return b.String()
}
func (e *ConstrMap) ValueType(p *Parser) string {
	return "Array"
}
func (e *Object) ValueType(p *Parser) string {
	if p.schema == nil {
		return "object"
	}
	return "object"
}
func (e *Variable) ValueType(p *Parser) string {
	return "variable"
}

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

type Parser struct {
	schema  *ast.Schema
	enumVal map[string]*ast.Definition
}

type Source struct {
	Name    string
	Content string
}

func NewParser(schema []Source) (*Parser, error) {
	if len(schema) < 1 {
		return &Parser{}, nil
	}

	in := make([]*ast.Source, len(schema))
	for i, s := range schema {
		in[i] = &ast.Source{
			Input: s.Content,
			Name:  s.Name,
		}
	}

	s, err := gqlparser.LoadSchema(in...)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		schema:  s,
		enumVal: map[string]*ast.Definition{},
	}
	for _, t := range s.Types {
		for _, v := range t.EnumValues {
			p.enumVal[v.Name] = t
		}
	}

	return p, nil
}

func Parse(src []byte) (
	operation *Operation,
	variables map[string]*VariableDefinition,
	err Error,
) {
	return (&Parser{}).Parse(src)
}

func (p *Parser) Parse(src []byte) (
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

	var typeDef *ast.Definition

	if string(tok) == "query" {
		o.Type = OperationTypeQuery
		if p.schema != nil {
			typeDef = p.schema.Query
		}
	} else if string(tok) == "mutation" {
		o.Type = OperationTypeMutation
		if p.schema != nil {
			typeDef = p.schema.Mutation
		}
	} else if string(tok) == "subscription" {
		o.Type = OperationTypeSubscription
		if p.schema != nil {
			typeDef = p.schema.Subscription
		}
	} else {
		return nil, nil, errUnexp(
			si, "expected query, mutation, or subscription "+
				"operation definition",
		)
	}

	variables = make(map[string]*VariableDefinition)
	varRefs := []*Variable{}

	s = s.consumeIgnored()
	if s, o.Selections, err = p.ParseSelectionSet(
		s, variables, &varRefs, typeDef,
	); err.IsErr() {
		return nil, nil, err
	}
	for _, sel := range o.Selections {
		setParent(sel, o)
	}
	if err := validateSelectionSet(s.s, o.Selections); err.IsErr() {
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

func errSchemaTypeUndef(s source) Error {
	return Error{
		Location: s.Location,
		Msg:      "type is undefined in schema",
	}
}

func errSchemaArgUndef(
	s source,
	argName string,
	fieldName string,
	hostType *ast.Definition,
) Error {
	return errMsg(s, fmt.Sprintf(
		"argument %q is undefined on field %q of type %q",
		argName, fieldName, hostType.Name,
	))
}

func errSchemaFieldUndef(
	s source,
	fieldName string,
	hostType *ast.Definition,
) Error {
	return Error{
		Location: s.Location,
		Msg: fmt.Sprintf(
			"field %q is undefined in type %q", fieldName, hostType.Name,
		),
	}
}

func (p *Parser) errCantCompare(s source, left, right Expression) Error {
	return Error{
		Location: s.Location,
		Msg: fmt.Sprintf(
			"can't compare %s and %s",
			left.ValueType(p), right.ValueType(p),
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

func (p *Parser) ParseSelectionSet(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	setDef *ast.Definition,
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
			var fragInline *SelectionInlineFrag
			sBefore := s
			if s, fragInline, err = p.ParseInlineFrag(
				s, variables, varRefs, setDef,
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

		s = s.consumeIgnored()

		if sel.Name == "max" {
			var maxNum any
			sBeforeMaxNum := s
			if s, maxNum, _ = s.ParseNumber(); maxNum != nil {
				// Max
				s = s.consumeIgnored()
				maxInt, ok := maxNum.(*Int)

				if !ok || maxInt.Value < 1 {
					return sBeforeMaxNum, nil, errMsg(
						sBeforeMaxNum, "maximum number of options"+
							" must be an unsigned integer greater 0",
					)
				}

				var options []Selection
				var err Error
				sBeforeOptionsBlock := s
				if s, options, err = p.ParseSelectionSet(
					s, variables, varRefs, setDef,
				); err.IsErr() {
					return s, nil, err
				}

				for _, option := range options {
					if v, ok := option.(*SelectionMax); ok {
						return s, nil, errMsg(source{
							s:        s.s,
							Location: v.Location,
						}, "nested max block")
					}
				}

				if len(options) < 2 {
					return sBeforeOptionsBlock, nil, errMsg(
						sBeforeOptionsBlock,
						"max block must have at least 2 selection options",
					)
				} else if maxInt.Value > int64(len(options)-1) {
					return sBeforeMaxNum, nil, errMsg(
						sBeforeMaxNum,
						"max selections number exceeds number of options-1",
					)
				}

				e := &SelectionMax{
					Location: sBeforeName.Location,
					Limit:    int(maxInt.Value),
					Options:  options,
				}
				selections = append(selections, e)
				continue
			}
		}

		if _, ok := fields[sel.Name]; ok {
			return sBeforeName, nil, errMsg(sBeforeName, "redeclared field")
		}
		fields[sel.Name] = struct{}{}

		var fieldDef *ast.FieldDefinition

		if p.schema != nil {
			if sel.Name != "__typename" {
				fieldDef = setDef.Fields.ForName(sel.Name)
				if fieldDef == nil {
					return sBeforeName, nil, errSchemaFieldUndef(
						sBeforeName, sel.Name, setDef,
					)
				}
			}
		}

		if s.peek1('(') {
			var err Error
			if s, sel.Arguments, err = p.ParseArguments(
				s, variables, varRefs, setDef, fieldDef,
			); err.IsErr() {
				return s, nil, err
			}
			for _, arg := range sel.Arguments {
				setParent(arg, sel)
			}
		} else if fieldDef != nil {
			for _, a := range fieldDef.Arguments {
				if a.Type.NonNull {
					return sBeforeName, nil, errMsg(
						sBeforeName, errMsgArgMissing(a),
					)
				}
			}
		}

		s = s.consumeIgnored()

		var typeDef *ast.Definition
		if p.schema != nil && fieldDef != nil {
			typeDef = p.schema.Types[fieldDef.Type.NamedType]
		}

		if s.peek1('{') {
			var err Error
			if s, sel.Selections, err = p.ParseSelectionSet(
				s, variables, varRefs, typeDef,
			); err.IsErr() {
				return s, nil, err
			}
			for _, sub := range sel.Selections {
				setParent(sub, sel)
			}
			if err := validateSelectionSet(s.s, sel.Selections); err.IsErr() {
				return s, nil, err
			}
		} else if p.schema != nil && typeDef != nil {
			if typeDef.Kind == ast.Object || typeDef.Kind == ast.Interface {
				return s, nil, errMsg(s, fmt.Sprintf(
					"missing selection set for selection %q of type %s",
					sel.Name, typeDef.Name,
				))
			}
		}

		selections = append(selections, sel)
	}

	if len(selections) < 1 {
		return si, nil, newErrorNoPrefix(si, "empty selection set")
	}

	return s, selections, Error{}
}

func (p *Parser) ParseInlineFrag(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	setType *ast.Definition,
) (source, *SelectionInlineFrag, Error) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("..."); !ok {
		return s, nil, errUnexp(s, "expected '...'")
	}

	inlineFrag := &SelectionInlineFrag{Location: l}
	s = s.consumeIgnored()

	var tok []byte
	sp := s
	s, tok = s.consumeToken()
	if string(tok) != "on" {
		return sp, nil, errUnexp(sp, "expected keyword 'on'")
	}

	s = s.consumeIgnored()

	sBeforeName := s
	var name []byte
	s, name = s.consumeName()
	if len(name) < 1 {
		return s, nil, errUnexp(s, "expected type condition")
	}
	inlineFrag.TypeCondition = string(name)

	var typeDef *ast.Definition
	if p.schema != nil {
		if typeDef, ok = p.schema.Types[inlineFrag.TypeCondition]; !ok {
			return sBeforeName, nil, errSchemaTypeUndef(sBeforeName)
		}
		if err := p.checkTypeCond(
			sBeforeName, setType, typeDef,
		); err.IsErr() {
			return sBeforeName, nil, err
		}
	}

	s = s.consumeIgnored()
	var sels []Selection
	var err Error
	s, sels, err = p.ParseSelectionSet(s, variables, varRefs, typeDef)
	if err.IsErr() {
		return s, nil, err
	}
	for _, sel := range sels {
		setParent(sel, inlineFrag)
	}

	inlineFrag.Selections = sels
	return s, inlineFrag, Error{}
}

func (p *Parser) ParseArguments(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	hostType *ast.Definition,
	fieldDef *ast.FieldDefinition,
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

		var argDef *ast.ArgumentDefinition
		if p.schema != nil {
			if argDef = fieldDef.Arguments.ForName(arg.Name); argDef == nil {
				return sBeforeName, nil, errSchemaArgUndef(
					sBeforeName, arg.Name, fieldDef.Name, hostType,
				)
			}
		}

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

		if s, ok = s.consume(":"); !ok {
			return s, nil, errUnexp(s, "expected colon")
		}
		s = s.consumeIgnored()
		if s.isEOF() {
			return s, nil, errUnexp(s, "expected constraint")
		}

		var expr Expression
		var err Error
		sBeforeConstr := s
		s, expr, err = p.ParseExprLogicalOr(
			s, variables, varRefs, expectConstraint,
		)
		if err.IsErr() {
			return s, nil, err
		}
		setParent(expr, arg)
		arg.Constraint = expr

		if argDef != nil {
			if err := p.checkArgConstraint(
				sBeforeConstr, arg.Constraint, argDef.Type,
			); err.IsErr() {
				return sBeforeConstr, nil, err
			}
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

	if p.schema != nil {
	ARG_SCAN:
		for _, a := range fieldDef.Arguments {
			if !a.Type.NonNull || a.DefaultValue != nil {
				continue
			}
			// Required argument with no default value must be included
			for _, arg := range arguments {
				if arg.Name == a.Name {
					continue ARG_SCAN
				}
			}
			return si, nil, errMsg(si, errMsgArgMissing(a))
		}
	}

	return s, arguments, Error{}
}

func (p *Parser) ParseValue(
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
		s, e.Expression, err = p.ParseExprLogicalOr(s, variables, varRefs, expect)
		if err.IsErr() {
			return s, nil, err
		}
		setParent(e.Expression, e)

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
		s = s.consumeIgnored()

		for {
			var err Error
			if s, ok = s.consume("]"); ok {
				break
			}

			var expr Expression
			s, expr, err = p.ParseExprLogicalOr(
				s, variables, varRefs, expectConstraint,
			)
			if err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)

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
			fld := &ObjectField{
				Location: s.Location,
				Parent:   o,
			}
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
				s, expr, err = p.ParseExprLogicalOr(
					s, variables, varRefs, expectConstraint,
				)
				if err.IsErr() {
					return s, nil, err
				}
				setParent(expr, fld)
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

func (p *Parser) ParseExprUnary(
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
		if s, e.Expression, err = p.ParseValue(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Expression, e)

		if !isBoolean(e.Expression) {
			return s, nil, errTypef(
				si, "can't use %s as Boolean", e.Expression.ValueType(p),
			)
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	}
	return p.ParseValue(s, variables, varRefs, expect)
}

func (p *Parser) ParseExprMultiplicative(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s

	var result Expression
	var err Error
	if s, result, err = p.ParseExprUnary(
		s, variables, varRefs, expect,
	); err.IsErr() {
		return s, nil, err
	}

	for {
		s = s.consumeIgnored()
		var selected int
		s, selected = s.consumeEitherOf3("*", "/", "%")
		switch selected {
		case 0:
			e := &ExprMultiplication{
				Location:     si.Location,
				Multiplicant: result,
			}
			setParent(e.Multiplicant, e)

			if !isNumeric(e.Multiplicant) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Multiplicant.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			si := s
			if s, e.Multiplicator, err = p.ParseExprUnary(
				s, variables, varRefs, expect,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(e.Multiplicator, e)

			if !isNumeric(e.Multiplicator) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Multiplicator.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprDivision{
				Location: si.Location,
				Dividend: result,
			}
			setParent(e.Dividend, e)

			if !isNumeric(e.Dividend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Dividend.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			si := s
			if s, e.Divisor, err = p.ParseExprUnary(
				s, variables, varRefs, expect,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(e.Divisor, e)

			if !isNumeric(e.Divisor) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Divisor.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			result = e
		case 2:
			e := &ExprModulo{
				Location: si.Location,
				Dividend: result,
			}
			setParent(e.Dividend, e)

			if !isNumeric(e.Dividend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Dividend.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			si := s
			if s, e.Divisor, err = p.ParseExprUnary(
				s, variables, varRefs, expect,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(e.Divisor, e)

			if !isNumeric(e.Divisor) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Divisor.ValueType(p),
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

func (p *Parser) ParseExprAdditive(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s

	var result Expression
	var err Error
	if s, result, err = p.ParseExprMultiplicative(
		s, variables, varRefs, expect,
	); err.IsErr() {
		return s, nil, err
	}

	for {
		s = s.consumeIgnored()
		var selected int
		s, selected = s.consumeEitherOf3("+", "-", "")
		switch selected {
		case 0:
			e := &ExprAddition{
				Location:   si.Location,
				AddendLeft: result,
			}
			setParent(e.AddendLeft, e)

			if !isNumeric(e.AddendLeft) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.AddendLeft.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			si := s
			if s, e.AddendRight, err = p.ParseExprMultiplicative(
				s, variables, varRefs, expect,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(e.AddendRight, e)

			if !isNumeric(e.AddendRight) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.AddendRight.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprSubtraction{
				Location: si.Location,
				Minuend:  result,
			}
			setParent(e.Minuend, e)

			if !isNumeric(e.Minuend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Minuend.ValueType(p),
				)
			}

			s = s.consumeIgnored()
			si := s
			if s, e.Subtrahend, err = p.ParseExprMultiplicative(
				s, variables, varRefs, expect,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(e.Subtrahend, e)

			if !isNumeric(e.Subtrahend) {
				return s, nil, errTypef(
					si, "can't use %s as number", e.Subtrahend.ValueType(p),
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

func (p *Parser) ParseExprRelational(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	l := s.Location

	var left Expression
	var err Error
	if s, left, err = p.ParseExprAdditive(
		s, variables, varRefs, expect,
	); err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("<="); ok {
		e := &ExprLessOrEqual{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right, err = p.ParseExprAdditive(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume(">="); ok {
		e := &ExprGreaterOrEqual{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right, err = p.ParseExprAdditive(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume("<"); ok {
		e := &ExprLess{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right, err = p.ParseExprAdditive(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e, Error{}
	} else if s, ok = s.consume(">"); ok {
		e := &ExprGreater{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		si := s
		if s, e.Right, err = p.ParseExprAdditive(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Right, e)

		if !isNumeric(e.Right) {
			return s, nil, errTypef(si, "can't use %s as number", e.ValueType(p))
		}

		s = s.consumeIgnored()
		return s, e, Error{}
	}
	s = s.consumeIgnored()
	return s, left, Error{}
}

func (p *Parser) ParseExprEquality(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s

	var err Error
	var exprLeft Expression
	if s, exprLeft, err = p.ParseExprRelational(
		s, variables, varRefs, expect,
	); err.IsErr() {
		return s, nil, err
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("=="); ok {
		e := &ExprEqual{
			Location: si.Location,
			Left:     exprLeft,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right, err = p.ParseExprRelational(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Right, e)
		s = s.consumeIgnored()

		if err = p.assumeSameType(si, e.Left, e.Right); err.IsErr() {
			return si, nil, err
		}

		return s, e, Error{}
	}
	if s, ok = s.consume("!="); ok {
		e := &ExprNotEqual{
			Location: si.Location,
			Left:     exprLeft,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right, err = p.ParseExprRelational(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(e.Right, e)
		s = s.consumeIgnored()

		if err = p.assumeSameType(si, e.Left, e.Right); err.IsErr() {
			return si, nil, err
		}

		return s, e, Error{}
	}

	s = s.consumeIgnored()
	return s, exprLeft, Error{}
}

func (p *Parser) ParseExprLogicalOr(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	e := &ExprLogicalOr{Location: s.Location}

	var err Error
	for {
		var expr Expression
		if s, expr, err = p.ParseExprLogicalAnd(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
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

func (p *Parser) ParseExprLogicalAnd(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	e := &ExprLogicalAnd{Location: s.Location}

	var err Error
	for {
		var expr Expression
		if s, expr, err = p.ParseConstr(
			s, variables, varRefs, expect,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
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

func (p *Parser) ParseConstr(
	s source,
	variables map[string]*VariableDefinition,
	varRefs *[]*Variable,
	expect expect,
) (source, Expression, Error) {
	si := s
	var ok bool
	var expr Expression
	var err Error

	if expect == expectValue {
		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}

		s = s.consumeIgnored()

		return s, expr, Error{}
	}

	if s, ok = s.consume("*"); ok {
		e := &ConstrAny{Location: si.Location}
		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume("!="); ok {
		e := &ConstrNotEquals{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume("<="); ok {
		e := &ConstrLessOrEqual{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		si := s
		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
		e.Value = expr

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.ValueType(p),
			)
		}

		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume(">="); ok {
		e := &ConstrGreaterOrEqual{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		si := s
		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
		e.Value = expr

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.ValueType(p),
			)
		}

		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume("<"); ok {
		e := &ConstrLess{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		si := s
		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
		e.Value = expr

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.ValueType(p),
			)
		}

		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume(">"); ok {
		e := &ConstrGreater{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		si := s
		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
		e.Value = expr

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.ValueType(p),
			)
		}

		s = s.consumeIgnored()

		return s, e, Error{}
	} else if s, ok = s.consume("len"); ok {
		s = s.consumeIgnored()

		if s, ok = s.consume("!="); ok {
			e := &ConstrLenNotEquals{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			var expr Expression
			var err Error

			si := s
			if s, expr, err = p.ParseExprEquality(
				s, variables, varRefs, expectValue,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)
			e.Value = expr

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.ValueType(p),
				)
			}

			s = s.consumeIgnored()

			return s, e, Error{}
		} else if s, ok = s.consume("<="); ok {
			e := &ConstrLenLessOrEqual{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			si := s
			if s, expr, err = p.ParseExprEquality(
				s, variables, varRefs, expectValue,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)
			e.Value = expr

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.ValueType(p),
				)
			}

			s = s.consumeIgnored()

			return s, e, Error{}
		} else if s, ok = s.consume(">="); ok {
			e := &ConstrLenGreaterOrEqual{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			si := s
			if s, expr, err = p.ParseExprEquality(
				s, variables, varRefs, expectValue,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)
			e.Value = expr

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.ValueType(p),
				)
			}

			s = s.consumeIgnored()

			return s, e, Error{}
		} else if s, ok = s.consume("<"); ok {
			e := &ConstrLenLess{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			si := s
			if s, expr, err = p.ParseExprEquality(
				s, variables, varRefs, expectValue,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)
			e.Value = expr

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.ValueType(p),
				)
			}

			s = s.consumeIgnored()

			return s, e, Error{}
		} else if s, ok = s.consume(">"); ok {
			e := &ConstrLenGreater{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			si := s
			if s, expr, err = p.ParseExprEquality(
				s, variables, varRefs, expectValue,
			); err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)
			e.Value = expr

			if !isNumeric(expr) {
				return s, nil, errTypef(
					si, "can't use %s as number", expr.ValueType(p),
				)
			}

			s = s.consumeIgnored()

			return s, e, Error{}
		}

		e := &ConstrGreater{
			Location: si.Location,
			Value:    expr,
		}

		si := s
		if s, expr, err = p.ParseExprEquality(
			s, variables, varRefs, expectValue,
		); err.IsErr() {
			return s, nil, err
		}
		setParent(expr, e)
		e.Value = expr

		if !isNumeric(expr) {
			return s, nil, errTypef(
				si, "can't use %s as number", expr.ValueType(p),
			)
		}

		s = s.consumeIgnored()

		return s, e, Error{}
	}

	if s, ok = s.consume("["); ok {
		s = s.consumeIgnored()
		if s, ok = s.consume("..."); ok {
			e := &ConstrMap{Location: si.Location}
			s = s.consumeIgnored()

			if s.isEOF() {
				return s, nil, errUnexp(s, "expected map constraint")
			}

			var expr Expression
			var err Error
			s, expr, err = p.ParseExprLogicalOr(
				s, variables, varRefs, expectConstraint,
			)
			if err.IsErr() {
				return s, nil, err
			}
			setParent(expr, e)
			e.Constraint = expr

			s = s.consumeIgnored()
			if s, ok = s.consume("]"); !ok {
				return s, nil, errUnexp(s, "expected end of map constraint ']'")
			}

			return s, e, Error{}
		} else {
			s = si
		}
	}

	e := &ConstrEquals{
		Location: si.Location,
		Value:    expr,
	}

	if s, expr, err = p.ParseExprEquality(
		s, variables, varRefs, expectValue,
	); err.IsErr() {
		return s, nil, err
	}
	setParent(expr, e)
	e.Value = expr

	s = s.consumeIgnored()

	return s, e, Error{}
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

func (p *Parser) assumeSameType(s source, left, right Expression) Error {
	_, lIsVar := left.(*Variable)
	_, rIsVar := right.(*Variable)
	if lIsVar || rIsVar {
		// TODO check variable types
		return Error{}
	}

	if isBoolean(left) && !isBoolean(right) {
		return p.errCantCompare(s, left, right)
	} else if isNumeric(left) && !isNumeric(right) {
		return p.errCantCompare(s, left, right)
	}
	switch left.(type) {
	case *String:
		if _, ok := right.(*String); !ok {
			return p.errCantCompare(s, left, right)
		}
	case *Enum:
		if _, ok := right.(*Enum); !ok {
			return p.errCantCompare(s, left, right)
		}
	case *Array:
		if _, ok := right.(*Array); !ok {
			return p.errCantCompare(s, left, right)
		}
	case *Object:
		if _, ok := right.(*Object); !ok {
			return p.errCantCompare(s, left, right)
		}
	}
	return Error{}
}

func setParent(t, parent any) {
	switch v := t.(type) {
	case *Variable:
		v.Parent = parent
	case *Int:
		v.Parent = parent
	case *Float:
		v.Parent = parent
	case *String:
		v.Parent = parent
	case *True:
		v.Parent = parent
	case *False:
		v.Parent = parent
	case *Null:
		v.Parent = parent
	case *Enum:
		v.Parent = parent
	case *Array:
		v.Parent = parent
	case *Object:
		v.Parent = parent
	case *ExprModulo:
		v.Parent = parent
	case *ExprDivision:
		v.Parent = parent
	case *ExprMultiplication:
		v.Parent = parent
	case *ExprAddition:
		v.Parent = parent
	case *ExprSubtraction:
		v.Parent = parent
	case *ExprEqual:
		v.Parent = parent
	case *ExprNotEqual:
		v.Parent = parent
	case *ExprLess:
		v.Parent = parent
	case *ExprGreater:
		v.Parent = parent
	case *ExprLessOrEqual:
		v.Parent = parent
	case *ExprGreaterOrEqual:
		v.Parent = parent
	case *ExprParentheses:
		v.Parent = parent
	case *ConstrEquals:
		v.Parent = parent
	case *ConstrNotEquals:
		v.Parent = parent
	case *ConstrLess:
		v.Parent = parent
	case *ConstrGreater:
		v.Parent = parent
	case *ConstrLessOrEqual:
		v.Parent = parent
	case *ConstrGreaterOrEqual:
		v.Parent = parent
	case *ConstrLenEquals:
		v.Parent = parent
	case *ConstrLenNotEquals:
		v.Parent = parent
	case *ConstrLenLess:
		v.Parent = parent
	case *ConstrLenGreater:
		v.Parent = parent
	case *ConstrLenLessOrEqual:
		v.Parent = parent
	case *ConstrLenGreaterOrEqual:
		v.Parent = parent
	case *ExprLogicalAnd:
		v.Parent = parent
	case *ExprLogicalOr:
		v.Parent = parent
	case *ExprLogicalNegation:
		v.Parent = parent
	case *SelectionInlineFrag:
		v.Parent = parent
	case *ObjectField:
		v.Parent = parent
	case *SelectionField:
		v.Parent = parent
	case *Argument:
		v.Parent = parent
	case *SelectionMax:
		v.Parent = parent
	case *ConstrMap:
		v.Parent = parent
	case *ConstrAny:
		v.Parent = parent
	default:
		panic(fmt.Errorf("unsupported type: %T", t))
	}
}

func validateSelectionSet(s []byte, set []Selection) Error {
	var maxs []*SelectionMax
	fields := map[string]*SelectionField{}
	inlineFrags := map[string]*SelectionInlineFrag{}
	for _, s := range set {
		switch v := s.(type) {
		case *SelectionMax:
			maxs = append(maxs, v)
		case *SelectionField:
			fields[v.Name] = v
		case *SelectionInlineFrag:
			inlineFrags[v.TypeCondition] = v
		}
	}

	if len(maxs) > 1 {
		return errMsg(source{
			s:        s,
			Location: maxs[len(maxs)-1].Location,
		}, "multiple max blocks in one selection set")
	}

	for _, m := range maxs {
		for _, sel := range m.Options {
			switch v := sel.(type) {
			case *SelectionField:
				if f, ok := fields[v.Name]; ok {
					return errMsg(source{
						s: s,
						Location: furthestLocation(
							f.Location, v.Location,
						),
					}, "redeclared field")
				}
			case *SelectionInlineFrag:
				if c, ok := inlineFrags[v.TypeCondition]; ok {
					return errMsg(source{
						s: s,
						Location: furthestLocation(
							c.Location, v.Location,
						),
					}, "redeclared type condition")
				}
			}
		}
	}

	return Error{}
}

func furthestLocation(a, b Location) Location {
	if a.Index < b.Index {
		return b
	}
	return a
}

// func Find[T any](document *Operation, path string) (t T) {
// 	var currentNode any = document
// 	i := strings.IndexAny(path, ".")
// 	if i < 0 {
// 		return
// 	}
// 	switch path[:i] {
// 	case "query":
// 		if document.Type != OperationTypeQuery {
// 			return
// 		}
// 	case "mutation":
// 		if document.Type != OperationTypeMutation {
// 			return
// 		}
// 	case "subscription":
// 		if document.Type != OperationTypeSubscription {
// 			return
// 		}
// 	default:
// 		return
// 	}
// 	path = path[i:]

// 	for path != "" {
// 		i := strings.IndexAny(path, ".|?")
// 		if i < 0 {
// 			break
// 		}
// 		next := path[:i]
// 		if next == "" {
// 			return
// 		}
// 		switch path[i] {
// 		case '.':
// 			// Field
// 		case '|':
// 			// Argument
// 		case '?':
// 			// Type condition
// 		}
// 	}
// 	return
// }

func (p *Parser) checkTypeCond(
	s source,
	hostDef, condTypeDef *ast.Definition,
) Error {
	switch condTypeDef.Kind {
	case ast.Scalar:
		return errMsg(s, fmt.Sprintf(
			"fragment can't condition on scalar type %q",
			condTypeDef.Name,
		))
	case ast.Enum:
		return errMsg(s, fmt.Sprintf(
			"fragment can't condition on enum type %q",
			condTypeDef.Name,
		))
	case ast.InputObject:
		return errMsg(s, fmt.Sprintf(
			"fragment can't condition on input type %q",
			condTypeDef.Name,
		))
	case ast.Object:
		if hostDef.Name == condTypeDef.Name {
			return Error{}
		}
		for _, c := range condTypeDef.Interfaces {
			if c == hostDef.Name {
				// condition Object implements the host interface
				return Error{}
			}
		}
		for _, c := range hostDef.Types {
			if c == condTypeDef.Name {
				// condition Object is part of the host union
				return Error{}
			}
		}
	case ast.Interface:
		if hostDef.Name == condTypeDef.Name {
			return Error{}
		}
		impls := p.schema.GetPossibleTypes(condTypeDef)
		for _, t := range impls {
			if t == hostDef {
				return Error{}
			}
			possibleTypes := p.schema.GetPossibleTypes(hostDef)
			for _, ht := range possibleTypes {
				if ht == t {
					return Error{}
				}
			}
		}
	case ast.Union:
		if hostDef.Name == condTypeDef.Name {
			return Error{}
		}
		impls := p.schema.GetPossibleTypes(condTypeDef)
		for _, t := range impls {
			if t == hostDef {
				return Error{}
			}
			possibleTypes := p.schema.GetPossibleTypes(hostDef)
			for _, ht := range possibleTypes {
				if ht == t {
					return Error{}
				}
			}
		}
	}
	return errMsg(s, fmt.Sprintf(
		"type %q can never be of type %q",
		hostDef.Name, condTypeDef.Name,
	))
}

func errMsgArgMissing(missingArgument *ast.ArgumentDefinition) string {
	return fmt.Sprintf(
		"argument %q of type %q is required but missing",
		missingArgument.Name, missingArgument.Type,
	)
}

func errMsgInputFieldMissing(missingField *ast.FieldDefinition) string {
	return fmt.Sprintf(
		"field %q of type %q is required but missing",
		missingField.Name, missingField.Type,
	)
}

func (p *Parser) errMsgUnexpType(
	s source,
	expected *ast.Type,
	actual Expression,
) Error {
	return errMsg(s, fmt.Sprintf(
		"expected type %q but received %q",
		expected, actual.ValueType(p),
	))
}

func errMsgUndefEnumVal(val string) string {
	return fmt.Sprintf("undefined enum value %q", val)
}

// checkArgConstraint returns false if the argument constraint type
// doesn't match the schema argument type.
func (p *Parser) checkArgConstraint(
	s source,
	e Expression,
	def *ast.Type,
) Error {
	switch v := e.(type) {
	case *ConstrAny:
	case *ConstrEquals:
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrNotEquals:
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLess:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrGreater:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLessOrEqual:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrGreaterOrEqual:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLenEquals:
		if !defHasLength(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLenNotEquals:
		if !defHasLength(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLenLess:
		if !defHasLength(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLenGreater:
		if !defHasLength(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLenLessOrEqual:
		if !defHasLength(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrLenGreaterOrEqual:
		if !defHasLength(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Value.GetLocation()}
		return p.checkArgConstraint(s, v.Value, def)
	case *ConstrMap:
		if !defIsArray(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Constraint.GetLocation()}
		return p.checkArgConstraint(s, v.Constraint, def.Elem)
	case *ExprParentheses:
		s = source{s: s.s, Location: v.Expression.GetLocation()}
		return p.checkArgConstraint(s, v.Expression, def)
	case *ExprLogicalNegation:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Expression.GetLocation()}
		return p.checkArgConstraint(s, v.Expression, def)
	case *ExprModulo:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Dividend.GetLocation()}
		if err := p.checkArgConstraint(s, v.Dividend, def); err.IsErr() {
			return err
		}
		s = source{s: s.s, Location: v.Divisor.GetLocation()}
		return p.checkArgConstraint(s, v.Divisor, def)
	case *ExprDivision:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Dividend.GetLocation()}
		if err := p.checkArgConstraint(s, v.Dividend, def); err.IsErr() {
			return err
		}
		s = source{s: s.s, Location: v.Divisor.GetLocation()}
		return p.checkArgConstraint(s, v.Divisor, def)
	case *ExprMultiplication:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Multiplicant.GetLocation()}
		if err := p.checkArgConstraint(s, v.Multiplicant, def); err.IsErr() {
			return err
		}
		s = source{s: s.s, Location: v.Multiplicator.GetLocation()}
		return p.checkArgConstraint(s, v.Multiplicator, def)
	case *ExprAddition:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.AddendLeft.GetLocation()}
		if err := p.checkArgConstraint(s, v.AddendLeft, def); err.IsErr() {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.AddendRight.GetLocation()}
		return p.checkArgConstraint(s, v.AddendRight, def)
	case *ExprSubtraction:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Minuend.GetLocation()}
		if err := p.checkArgConstraint(s, v.Minuend, def); err.IsErr() {
			return p.errMsgUnexpType(s, def, v)
		}
		s = source{s: s.s, Location: v.Subtrahend.GetLocation()}
		return p.checkArgConstraint(s, v.Subtrahend, def)
	case *Int:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *Float:
		if !defIsNum(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *True:
		if !defIsBool(def) {
			fmt.Println("OKAY", v.Location)
			return p.errMsgUnexpType(s, def, v)
		}
	case *False:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *String:
		if def.NamedType != "String" {
			return p.errMsgUnexpType(s, def, v)
		}
	case *Array:
		if !defIsArray(def) {
			return p.errMsgUnexpType(s, def, v)
		}
		for _, i := range v.Items {
			if err := p.checkArgConstraint(source{
				s:        s.s,
				Location: i.GetLocation(),
			}, i, def.Elem); err.IsErr() {
				return err
			}
		}
	case *Null:
		if def.NonNull {
			return p.errMsgUnexpType(s, def, v)
		}
	case *Enum:
		t, ok := p.enumVal[v.Value]
		if !ok {
			return errMsg(s, errMsgUndefEnumVal(v.Value))
		}
		if def.NamedType != t.Name {
			return p.errMsgUnexpType(s, def, v)
		}
	case *Object:
		tp := p.schema.Types[def.Name()]

		if tp.Kind != ast.InputObject {
			return p.errMsgUnexpType(s, def, v)
		}

		declaredFields := make(
			map[string]*ast.FieldDefinition,
			len(tp.Fields),
		)

		for _, vf := range v.Fields {
			field := tp.Fields.ForName(vf.Name)
			if field == nil {
				return errMsg(source{
					s:        s.s,
					Location: vf.Location,
				}, fmt.Sprintf(
					"undefined field %q in %s",
					field.Name, def.Name(),
				))
			}

			if err := p.checkArgConstraint(
				source{
					s:        s.s,
					Location: vf.Constraint.GetLocation(),
				}, vf.Constraint, field.Type,
			); err.IsErr() {
				return err
			}

			declaredFields[vf.Name] = field
		}

		for _, f := range tp.Fields {
			if f.Type.NonNull {
				if _, ok := declaredFields[f.Name]; !ok {
					return errMsg(s, errMsgInputFieldMissing(f))
				}
			}
		}

	case *ExprEqual:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *ExprNotEqual:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *ExprLess:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *ExprLessOrEqual:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *ExprGreater:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *ExprGreaterOrEqual:
		if !defIsBool(def) {
			return p.errMsgUnexpType(s, def, v)
		}
	case *ExprLogicalAnd:
		for _, e := range v.Expressions {
			s = source{s: s.s, Location: e.GetLocation()}
			if err := p.checkArgConstraint(s, e, def); err.IsErr() {
				return err
			}
		}
	case *ExprLogicalOr:
		for _, e := range v.Expressions {
			s = source{s: s.s, Location: e.GetLocation()}
			if err := p.checkArgConstraint(s, e, def); err.IsErr() {
				return err
			}
		}
	default:
		panic(fmt.Errorf("unsupported constraint type: %T", e))
	}
	return Error{}
}

func defIsNum(t *ast.Type) bool {
	return t.NamedType == "Int" || t.NamedType == "Float"
}

func defIsBool(t *ast.Type) bool {
	return t.NamedType == "Boolean"
}

func defHasLength(t *ast.Type) bool {
	return t.NamedType == "String" || defIsArray(t)
}

func defIsArray(t *ast.Type) bool {
	return t.NamedType == "" && t.Elem != nil
}
