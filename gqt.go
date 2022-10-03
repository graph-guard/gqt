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
		return "Query"
	case OperationTypeMutation:
		return "Mutation"
	case OperationTypeSubscription:
		return "Subscription"
	}
	return ""
}

type (
	Location struct{ Index, Line, Column int }

	Operation struct {
		Location
		Type OperationType
		SelectionSet
		Def *ast.Definition
	}

	// Selection can be either of:
	//
	//	*SelectionField
	//	*SelectionInlineFrag
	//	*SelectionMax
	Selection Expression

	SelectionField struct {
		Location
		Parent Expression
		Name   string
		ArgumentList
		SelectionSet
		Def *ast.FieldDefinition
	}

	SelectionSet struct {
		Location
		Selections []Selection
	}

	ArgumentList struct {
		Location
		Arguments []*Argument
	}

	Argument struct {
		Location
		Parent             Expression
		Name               string
		AssociatedVariable *VariableDeclaration
		Constraint         Expression
		Def                *ast.ArgumentDefinition
	}

	// Expression can be either of:
	//
	//  *Operation
	//  *ConstrAny
	//  *ConstrEquals
	//  *ConstrNotEquals
	//  *ConstrLess
	//  *ConstrLessOrEqual
	//  *ConstrGreater
	//  *ConstrGreaterOrEqual
	//  *ConstrLenEquals
	//  *ConstrLenNotEquals
	//  *ConstrLenLess
	//  *ConstrLenLessOrEqual
	//  *ConstrLenGreater
	//  *ConstrLenGreaterOrEqual
	//  *ConstrMap
	//  *ExprParentheses
	//  *ExprModulo
	//  *ExprDivision
	//  *ExprMultiplication
	//  *ExprAddition
	//  *ExprSubtraction
	//  *ExprLogicalNegation
	//  *ExprEqual
	//  *ExprNotEqual
	//  *ExprLess
	//  *ExprLessOrEqual
	//  *ExprGreater
	//  *ExprGreaterOrEqual
	//  *ExprLogicalAnd
	//  *ExprLogicalOr
	//  *True
	//  *False
	//  *Int
	//  *Float
	//  *String
	//  *Null
	//  *Enum
	//  *Array
	//  *Object
	//  *VariableRef
	//  *SelectionInlineFrag
	//  *ObjectField
	//  *SelectionField
	//  *SelectionMax
	//  *Argument
	Expression interface {
		GetParent() Expression
		GetLocation() Location
		TypeDesignation() string
		IsFloat() bool
	}

	ConstrAny struct {
		Location
		Parent Expression
	}

	ConstrEquals struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrNotEquals struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLess struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLessOrEqual struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrGreater struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrGreaterOrEqual struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLenEquals struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLenNotEquals struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLenLess struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLenLessOrEqual struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLenGreater struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrLenGreaterOrEqual struct {
		Location
		Parent Expression
		Value  Expression
	}

	ConstrMap struct {
		Location
		Parent     Expression
		Constraint Expression
	}

	ExprParentheses struct {
		Location
		Parent     Expression
		Expression Expression
	}

	ExprLogicalNegation struct {
		Location
		Parent     Expression
		Expression Expression
	}

	ExprModulo struct {
		Location
		Parent   Expression
		Dividend Expression
		Divisor  Expression
		Float    bool
	}

	ExprDivision struct {
		Location
		Parent   Expression
		Dividend Expression
		Divisor  Expression
		Float    bool
	}

	ExprMultiplication struct {
		Location
		Parent        Expression
		Multiplicant  Expression
		Multiplicator Expression
		Float         bool
	}

	ExprAddition struct {
		Location
		Parent      Expression
		AddendLeft  Expression
		AddendRight Expression
		Float       bool
	}

	ExprSubtraction struct {
		Location
		Parent     Expression
		Minuend    Expression
		Subtrahend Expression
		Float      bool
	}

	ExprEqual struct {
		Location
		Parent Expression
		Left   Expression
		Right  Expression
	}

	ExprNotEqual struct {
		Location
		Parent Expression
		Left   Expression
		Right  Expression
	}

	ExprLess struct {
		Location
		Parent Expression
		Left   Expression
		Right  Expression
	}

	ExprLessOrEqual struct {
		Location
		Parent Expression
		Left   Expression
		Right  Expression
	}

	ExprGreater struct {
		Location
		Parent Expression
		Left   Expression
		Right  Expression
	}

	ExprGreaterOrEqual struct {
		Location
		Parent Expression
		Left   Expression
		Right  Expression
	}

	ExprLogicalAnd struct {
		Location
		Parent      Expression
		Expressions []Expression
	}

	ExprLogicalOr struct {
		Location
		Parent      Expression
		Expressions []Expression
	}

	Int struct {
		Location
		Parent Expression
		Value  int64
	}

	Float struct {
		Location
		Parent Expression
		Value  float64
	}

	String struct {
		Location
		Parent Expression
		Value  string
	}

	True struct {
		Location
		Parent Expression
	}

	False struct {
		Location
		Parent Expression
	}

	Null struct {
		Location
		Parent  Expression
		TypeDef *ast.Definition
	}

	Enum struct {
		Location
		Parent  Expression
		Value   string
		TypeDef *ast.Definition
	}

	Array struct {
		Location
		Parent      Expression
		Items       []Expression
		ItemTypeDef *ast.Definition
	}

	Object struct {
		Location
		Parent  Expression
		Fields  []*ObjectField
		TypeDef *ast.Definition
	}

	ObjectField struct {
		Location
		Parent             Expression
		Name               string
		AssociatedVariable *VariableDeclaration
		Constraint         Expression
		Def                *ast.FieldDefinition
	}

	VariableRef struct {
		Location
		Parent      Expression
		Name        string
		Declaration *VariableDeclaration
	}

	SelectionMax struct {
		Location
		Parent  Expression
		Limit   int
		Options SelectionSet
	}

	SelectionInlineFrag struct {
		Location
		Parent        Expression
		TypeCondition TypeCondition
		SelectionSet
	}

	TypeCondition struct {
		Location
		TypeName string
		TypeDef  *ast.Definition
	}
)

func (e *Operation) GetParent() Expression               { return nil }
func (e *ConstrAny) GetParent() Expression               { return e.Parent }
func (e *ConstrEquals) GetParent() Expression            { return e.Parent }
func (e *ConstrNotEquals) GetParent() Expression         { return e.Parent }
func (e *ConstrLess) GetParent() Expression              { return e.Parent }
func (e *ConstrLessOrEqual) GetParent() Expression       { return e.Parent }
func (e *ConstrGreater) GetParent() Expression           { return e.Parent }
func (e *ConstrGreaterOrEqual) GetParent() Expression    { return e.Parent }
func (e *ConstrLenEquals) GetParent() Expression         { return e.Parent }
func (e *ConstrLenNotEquals) GetParent() Expression      { return e.Parent }
func (e *ConstrLenLess) GetParent() Expression           { return e.Parent }
func (e *ConstrLenLessOrEqual) GetParent() Expression    { return e.Parent }
func (e *ConstrLenGreater) GetParent() Expression        { return e.Parent }
func (e *ConstrLenGreaterOrEqual) GetParent() Expression { return e.Parent }
func (e *ConstrMap) GetParent() Expression               { return e.Parent }

func (e *ExprParentheses) GetParent() Expression    { return e.Parent }
func (e *ExprModulo) GetParent() Expression         { return e.Parent }
func (e *ExprDivision) GetParent() Expression       { return e.Parent }
func (e *ExprMultiplication) GetParent() Expression { return e.Parent }
func (e *ExprAddition) GetParent() Expression       { return e.Parent }
func (e *ExprSubtraction) GetParent() Expression    { return e.Parent }

func (e *ExprLogicalNegation) GetParent() Expression { return e.Parent }
func (e *ExprEqual) GetParent() Expression           { return e.Parent }
func (e *ExprNotEqual) GetParent() Expression        { return e.Parent }
func (e *ExprLess) GetParent() Expression            { return e.Parent }
func (e *ExprLessOrEqual) GetParent() Expression     { return e.Parent }
func (e *ExprGreater) GetParent() Expression         { return e.Parent }
func (e *ExprGreaterOrEqual) GetParent() Expression  { return e.Parent }
func (e *ExprLogicalAnd) GetParent() Expression      { return e.Parent }
func (e *ExprLogicalOr) GetParent() Expression       { return e.Parent }

func (e *True) GetParent() Expression        { return e.Parent }
func (e *False) GetParent() Expression       { return e.Parent }
func (e *Int) GetParent() Expression         { return e.Parent }
func (e *Float) GetParent() Expression       { return e.Parent }
func (e *String) GetParent() Expression      { return e.Parent }
func (e *Null) GetParent() Expression        { return e.Parent }
func (e *Enum) GetParent() Expression        { return e.Parent }
func (e *Array) GetParent() Expression       { return e.Parent }
func (e *Object) GetParent() Expression      { return e.Parent }
func (e *VariableRef) GetParent() Expression { return e.Parent }

func (e *SelectionInlineFrag) GetParent() Expression { return e.Parent }
func (e *ObjectField) GetParent() Expression         { return e.Parent }
func (e *SelectionField) GetParent() Expression      { return e.Parent }
func (e *SelectionMax) GetParent() Expression        { return e.Parent }
func (e *Argument) GetParent() Expression            { return e.Parent }

func (e *Operation) GetLocation() Location               { return e.Location }
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

func (e *True) GetLocation() Location        { return e.Location }
func (e *False) GetLocation() Location       { return e.Location }
func (e *Int) GetLocation() Location         { return e.Location }
func (e *Float) GetLocation() Location       { return e.Location }
func (e *String) GetLocation() Location      { return e.Location }
func (e *Null) GetLocation() Location        { return e.Location }
func (e *Enum) GetLocation() Location        { return e.Location }
func (e *Array) GetLocation() Location       { return e.Location }
func (e *ConstrMap) GetLocation() Location   { return e.Location }
func (e *Object) GetLocation() Location      { return e.Location }
func (e *VariableRef) GetLocation() Location { return e.Location }

func (e *SelectionInlineFrag) GetLocation() Location { return e.Location }
func (e *ObjectField) GetLocation() Location         { return e.Location }
func (e *SelectionField) GetLocation() Location      { return e.Location }
func (e *SelectionMax) GetLocation() Location        { return e.Location }
func (e *Argument) GetLocation() Location            { return e.Location }

func (e *Operation) IsFloat() bool               { return false }
func (e *ConstrAny) IsFloat() bool               { return false }
func (e *ConstrEquals) IsFloat() bool            { return false }
func (e *ConstrNotEquals) IsFloat() bool         { return false }
func (e *ConstrLess) IsFloat() bool              { return false }
func (e *ConstrLessOrEqual) IsFloat() bool       { return false }
func (e *ConstrGreater) IsFloat() bool           { return false }
func (e *ConstrGreaterOrEqual) IsFloat() bool    { return false }
func (e *ConstrLenEquals) IsFloat() bool         { return false }
func (e *ConstrLenNotEquals) IsFloat() bool      { return false }
func (e *ConstrLenLess) IsFloat() bool           { return false }
func (e *ConstrLenLessOrEqual) IsFloat() bool    { return false }
func (e *ConstrLenGreater) IsFloat() bool        { return false }
func (e *ConstrLenGreaterOrEqual) IsFloat() bool { return false }
func (e *ConstrMap) IsFloat() bool               { return false }

func (e *ExprParentheses) IsFloat() bool    { return e.Expression.IsFloat() }
func (e *ExprModulo) IsFloat() bool         { return e.Float }
func (e *ExprDivision) IsFloat() bool       { return e.Float }
func (e *ExprMultiplication) IsFloat() bool { return e.Float }
func (e *ExprAddition) IsFloat() bool       { return e.Float }
func (e *ExprSubtraction) IsFloat() bool    { return e.Float }

func (e *ExprLogicalNegation) IsFloat() bool { return false }
func (e *ExprEqual) IsFloat() bool           { return false }
func (e *ExprNotEqual) IsFloat() bool        { return false }
func (e *ExprLess) IsFloat() bool            { return false }
func (e *ExprLessOrEqual) IsFloat() bool     { return false }
func (e *ExprGreater) IsFloat() bool         { return false }
func (e *ExprGreaterOrEqual) IsFloat() bool  { return false }
func (e *ExprLogicalAnd) IsFloat() bool      { return false }
func (e *ExprLogicalOr) IsFloat() bool       { return false }

func (e *True) IsFloat() bool        { return false }
func (e *False) IsFloat() bool       { return false }
func (e *Int) IsFloat() bool         { return false }
func (e *Float) IsFloat() bool       { return true }
func (e *String) IsFloat() bool      { return false }
func (e *Null) IsFloat() bool        { return false }
func (e *Enum) IsFloat() bool        { return false }
func (e *Array) IsFloat() bool       { return false }
func (e *Object) IsFloat() bool      { return false }
func (e *VariableRef) IsFloat() bool { return false }

func (e *Argument) IsFloat() bool            { return e.Constraint.IsFloat() }
func (e *SelectionInlineFrag) IsFloat() bool { return false }
func (e *ObjectField) IsFloat() bool         { return false }
func (e *SelectionField) IsFloat() bool      { return false }
func (e *SelectionMax) IsFloat() bool        { return false }

func (e *Operation) TypeDesignation() string {
	return e.Type.String()
}
func (e *ConstrAny) TypeDesignation() string {
	switch p := e.Parent.(type) {
	case *Argument:
		if p.Def != nil {
			return p.Def.Type.String()
		}
	case *ObjectField:
		if p.Def != nil {
			return p.Def.Type.String()
		}
	}
	return "*"
}
func (e *ConstrEquals) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrNotEquals) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLess) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLessOrEqual) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrGreater) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrGreaterOrEqual) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLenEquals) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLenNotEquals) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLenLess) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLenLessOrEqual) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLenGreater) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrLenGreaterOrEqual) TypeDesignation() string {
	return e.Value.TypeDesignation()
}
func (e *ConstrMap) TypeDesignation() string {
	return e.Constraint.TypeDesignation()
}

func (e *ExprParentheses) TypeDesignation() string {
	return e.Expression.TypeDesignation()
}
func (e *ExprModulo) TypeDesignation() string {
	if e.Float {
		return "Float"
	}
	return "Int"
}
func (e *ExprDivision) TypeDesignation() string {
	if e.Float {
		return "Float"
	}
	return "Int"
}
func (e *ExprMultiplication) TypeDesignation() string {
	if e.Float {
		return "Float"
	}
	return "Int"
}
func (e *ExprAddition) TypeDesignation() string {
	if e.Float {
		return "Float"
	}
	return "Int"
}
func (e *ExprSubtraction) TypeDesignation() string {
	if e.Float {
		return "Float"
	}
	return "Int"
}

func (e *ExprLogicalNegation) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprEqual) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprNotEqual) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprLess) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprLessOrEqual) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprGreater) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprGreaterOrEqual) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprLogicalAnd) TypeDesignation() string {
	return "Boolean"
}
func (e *ExprLogicalOr) TypeDesignation() string {
	return "Boolean"
}

func (e *True) TypeDesignation() string {
	return "Boolean"
}
func (e *False) TypeDesignation() string {
	return "Boolean"
}
func (e *Int) TypeDesignation() string {
	return "Int"
}
func (e *Float) TypeDesignation() string {
	return "Float"
}
func (e *String) TypeDesignation() string {
	return "String"
}
func (e *Null) TypeDesignation() string {
	return "null"
}
func (e *Enum) TypeDesignation() string {
	if e.TypeDef != nil {
		return e.TypeDef.Name
	}
	return "enum"
}
func (e *Array) TypeDesignation() string {
	if e.ItemTypeDef != nil {
		return "[" + e.ItemTypeDef.Name + "]"
	}
	return "array"
}
func (e *Object) TypeDesignation() string {
	if e.TypeDef != nil {
		return e.TypeDef.Name
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, f := range e.Fields {
		b.WriteString(f.Name)
		b.WriteByte(':')
		b.WriteString(f.Constraint.TypeDesignation())
		if i+1 < len(e.Fields) {
			b.WriteByte(',')
		}
	}
	b.WriteByte('}')
	return b.String()
}
func (e *VariableRef) TypeDesignation() string {
	switch p := e.Declaration.Parent.(type) {
	case *Argument:
		return p.Constraint.TypeDesignation()
	case *ObjectField:
		return p.Constraint.TypeDesignation()
	}
	panic(fmt.Errorf(
		"unhandled variable declaration parent: %T",
		e.Declaration.Parent,
	))
}

func (e *Argument) TypeDesignation() string {
	return e.Constraint.TypeDesignation()
}
func (e *SelectionInlineFrag) TypeDesignation() string { return "" }
func (e *ObjectField) TypeDesignation() string         { return "" }
func (e *SelectionField) TypeDesignation() string      { return "" }
func (e *SelectionMax) TypeDesignation() string        { return "" }

type VariableDeclaration struct {
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
	Parent Expression
}

func (v *VariableDeclaration) Type() *ast.Type {
	switch p := v.Parent.(type) {
	case *Argument:
		if p.Def != nil {
			return p.Def.Type
		}
	case *ObjectField:
		if p.Def != nil {
			return p.Def.Type
		}
	}
	return nil
}

type Parser struct {
	schema   *ast.Schema
	enumVal  map[string]*ast.Definition
	varDecls map[string]*VariableDeclaration
	varRefs  []*VariableRef
	errors   []Error
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
	variables map[string]*VariableDeclaration,
	errors []Error,
) {
	return (&Parser{}).Parse(src)
}

func (p *Parser) Parse(src []byte) (
	operation *Operation,
	variables map[string]*VariableDeclaration,
	errors []Error,
) {
	p.errors = p.errors[:0]
	p.varDecls = make(map[string]*VariableDeclaration)
	p.varRefs = make([]*VariableRef, 0)

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
		p.errUnexpTok(
			si, "expected query, mutation, or subscription "+
				"operation definition",
		)
		return nil, nil, p.errors
	}

	s = s.consumeIgnored()
	if s, o.SelectionSet = p.ParseSelectionSet(s); s.stop() {
		return nil, nil, p.errors
	}
	for _, sel := range o.Selections {
		setParent(sel, o)
	}
	if !p.validateSelectionSet(
		s.s, o.Selections,
	) {
		return nil, nil, p.errors
	}

	// Link variable declarations and references
	for _, r := range p.varRefs {
		if v, ok := p.varDecls[r.Name]; !ok {
			sc := s
			sc.Location = r.Location
			p.newErr(sc.Location, "undefined variable")
			return nil, nil, p.errors
		} else {
			r.Declaration = v
			v.References = append(v.References, r)
		}
	}

	p.setTypes(o)
	p.validate(o)

	if len(p.errors) > 0 {
		o = nil
		p.varDecls = nil
	}

	return o, p.varDecls, p.errors
}

func (p *Parser) setTypes(o *Operation) {
	if p.schema != nil {
		switch o.Type {
		case OperationTypeQuery:
			o.Def = p.schema.Query
		case OperationTypeMutation:
			o.Def = p.schema.Mutation
		case OperationTypeSubscription:
			o.Def = p.schema.Subscription
		}
	}
	var fields ast.FieldList
	if o.Def != nil {
		fields = o.Def.Fields
	}
	p.setTypesSelSet(o.SelectionSet, fields)
}

func (p *Parser) setTypesSelSet(s SelectionSet, defs []*ast.FieldDefinition) {
	find := func(n string) *ast.FieldDefinition {
		for _, s := range defs {
			if s.Name == n {
				return s
			}
		}
		return nil
	}
	for _, s := range s.Selections {
		switch s := s.(type) {
		case *SelectionField:
			s.Def = find(s.Name)
			if s.Def != nil {
				if t := p.schema.Types[s.Def.Type.NamedType]; t != nil {
					p.setTypesSelSet(s.SelectionSet, t.Fields)
				}
				for _, a := range s.Arguments {
					a.Def = s.Def.Arguments.ForName(a.Name)
					if a.Def == nil {
						continue
					}
					var exp *ast.Type
					if a.Def.Type != nil {
						exp = a.Def.Type
					}
					p.setTypesExpr(a.Constraint, exp)
				}
			}
		case *SelectionInlineFrag:
			if p.schema != nil {
				s.TypeCondition.TypeDef = p.schema.Types[s.TypeCondition.TypeName]
				var fields ast.FieldList
				if s.TypeCondition.TypeDef != nil {
					fields = s.TypeCondition.TypeDef.Fields
				}
				p.setTypesSelSet(s.SelectionSet, fields)
			} else {
				p.setTypesSelSet(s.SelectionSet, nil)
			}
		case *SelectionMax:
			p.setTypesSelSet(s.Options, defs)
		}
	}
}

func (p *Parser) setTypesExpr(e Expression, exp *ast.Type) {
	type S struct {
		Expr   Expression
		Expect *ast.Type
	}

	stack := make([]S, 0, 64)
	push := func(expression Expression, expect *ast.Type) {
		stack = append(stack, S{Expr: expression, Expect: expect})
	}

	push(e, exp)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		exp := top.Expect

		switch e := top.Expr.(type) {
		case *ConstrAny, *Int, *Float, *True, *False, *String, *VariableRef:
		case *Null:
			if exp != nil {
				t := p.schema.Types[exp.NamedType]
				e.TypeDef = t
			}
		case *ConstrEquals:
			push(e.Value, exp)
		case *ConstrNotEquals:
			push(e.Value, exp)
		case *ConstrGreater:
			push(e.Value, exp)
		case *ConstrGreaterOrEqual:
			push(e.Value, exp)
		case *ConstrLess:
			push(e.Value, exp)
		case *ConstrLessOrEqual:
			push(e.Value, exp)
		case *ConstrLenEquals:
			push(e.Value, exp)
		case *ConstrLenNotEquals:
			push(e.Value, exp)
		case *ConstrLenGreater:
			push(e.Value, exp)
		case *ConstrLenGreaterOrEqual:
			push(e.Value, exp)
		case *ConstrLenLess:
			push(e.Value, exp)
		case *ConstrLenLessOrEqual:
			push(e.Value, exp)
		case *ConstrMap:
			push(e.Constraint, exp.Elem)
		case *ExprParentheses:
			push(e.Expression, exp)
		case *ExprEqual:
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprNotEqual:
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprLogicalNegation:
			push(e.Expression, exp)
		case *ExprLogicalOr:
			for _, e := range e.Expressions {
				push(e, exp)
			}
		case *ExprLogicalAnd:
			for _, e := range e.Expressions {
				push(e, exp)
			}
		case *ExprAddition:
			push(e.AddendLeft, exp)
			push(e.AddendRight, exp)
		case *ExprSubtraction:
			push(e.Minuend, exp)
			push(e.Subtrahend, exp)
		case *ExprMultiplication:
			push(e.Multiplicant, exp)
			push(e.Multiplicator, exp)
		case *ExprDivision:
			push(e.Dividend, exp)
			push(e.Divisor, exp)
		case *ExprModulo:
			push(e.Dividend, exp)
			push(e.Divisor, exp)
		case *Array:
			if exp != nil && exp.Elem != nil {
				for _, i := range e.Items {
					push(i, exp.Elem)
				}
			}
		case *Enum:
			e.TypeDef = p.enumVal[e.Value]
		case *Object:
			if exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.InputObject {
					e.TypeDef = t
					for _, f := range e.Fields {
						var tp *ast.Type
						f.Def = t.Fields.ForName(f.Name)
						if f.Def != nil {
							tp = f.Def.Type
						}
						p.setTypesExpr(f.Constraint, tp)
					}
				}
			}
		default:
			panic(fmt.Errorf("unhandled type: %T", top.Expr))
		}
	}
}

func (p *Parser) validate(o *Operation) (ok bool) {
	var def *ast.Definition
	switch o.Type {
	case OperationTypeQuery:
		if p.schema != nil {
			if p.schema.Query == nil {
				p.newErr(o.Location, "type Query is undefined")
				return false
			}
			def = p.schema.Query
		}
		return p.validateSelSet(o, def)
	case OperationTypeMutation:
		if p.schema != nil {
			if p.schema.Mutation == nil {
				p.newErr(o.Location, "type Mutation is undefined")
				return false
			}
			def = p.schema.Mutation
		}
		return p.validateSelSet(o, def)
	case OperationTypeSubscription:
		if p.schema != nil {
			if p.schema.Subscription == nil {
				p.newErr(o.Location, "type Subscription is undefined")
				return false
			}
			def = p.schema.Subscription
		}
		return p.validateSelSet(o, def)
	}
	panic(fmt.Sprintf("unhandled operation type: %v", o.Type))
}

func (p *Parser) validateSelSet(
	host Expression,
	expect *ast.Definition,
) (ok bool) {
	ok = true
	var l Location
	var s SelectionSet
	var hostTypeName string
	switch h := host.(type) {
	case *Operation:
		s, l = h.SelectionSet, h.Location
		hostTypeName = h.Type.String()
	case *SelectionField:
		s, l = h.SelectionSet, h.Location
		if h.Def != nil {
			hostTypeName = h.Def.Type.String()
		}
	case *SelectionInlineFrag:
		s, l = h.SelectionSet, h.Location
		hostTypeName = h.TypeCondition.TypeName
	default:
		panic(fmt.Errorf("unsupported type: %T", host))
	}

	parent, _ := host.(*SelectionField)
	if parent != nil &&
		expect != nil &&
		len(s.Selections) < 1 &&
		len(expect.Fields) > 0 {
		p.errMissingSelSet(l, parent.Name, hostTypeName)
		return false
	}

	for _, s := range s.Selections {
		switch s := s.(type) {
		case *SelectionField:
			var def *ast.FieldDefinition
			if expect != nil {
				def = expect.Fields.ForName(s.Name)
				if def == nil {
					if s.Name == "__typename" {
						continue
					}
					p.errUndefFieldInType(s.Location, s.Name, expect.Name)
					ok = false
					continue
				}
			}
			if !p.validateField(expect, s, def) {
				ok = false
				continue
			}
		case *SelectionInlineFrag:
			p.validateInlineFrag(s, expect)
		case *SelectionMax:
			for _, s := range s.Options.Selections {
				switch s := s.(type) {
				case *SelectionField:
					if s.Name == "__typename" {
						ok = false
						p.errTypenameInMax(s)
						continue
					}
					var def *ast.FieldDefinition
					if expect != nil {
						def = expect.Fields.ForName(s.Name)
						if def == nil {
							ok = false
							p.errUndefFieldInType(
								s.Location, s.Name, expect.Name,
							)
							continue
						}
					}
					if !p.validateField(expect, s, def) {
						ok = false
						continue
					}
				case *SelectionInlineFrag:
					if !p.validateInlineFrag(s, expect) {
						ok = false
						continue
					}
				case *SelectionMax:
					p.errNestedMaxBlock(s)
					ok = false
					continue
				}
			}
		}
	}
	return ok
}

func (p *Parser) validateField(
	host *ast.Definition,
	f *SelectionField,
	expect *ast.FieldDefinition,
) (ok bool) {
	if expect != nil {
		args := make(map[string]*Argument, len(f.Arguments))
		for _, a := range f.Arguments {
			// Check undefined arguments
			if ad := expect.Arguments.ForName(a.Name); ad == nil {
				p.errUndefArg(a, f, host.Name)
			}
			args[a.Name] = a
		}

		// Check required arguments
		for _, a := range expect.Arguments {
			if a.Type.NonNull && a.DefaultValue == nil {
				if _, ok := args[a.Name]; !ok {
					l := f.ArgumentList.Location
					if len(f.Arguments) < 1 {
						l = f.Location
					}
					p.errMissingArg(l, a)
					return false
				}
			}
		}
	}

	// Check constraints
	for _, a := range f.Arguments {
		var exp *ast.Type
		if a.Def != nil {
			exp = a.Def.Type
		}

		p.validateExpr([]any{a}, a.Constraint, exp)
	}

	var def *ast.Definition
	if p.schema != nil {
		def = p.schema.Types[expect.Type.NamedType]
	}
	return p.validateSelSet(f, def)
}

func (p *Parser) validateExpr(
	// hostPath defines the path to the origin argument
	// can contain any of: *Argument, *ObjectField
	pathToOriginArg []any,
	e Expression,
	expect *ast.Type,
) (ok bool) {
	ok = true
	type S struct {
		Expr   Expression
		Expect *ast.Type
	}

	stack := make([]S, 0, 64)
	push := func(expression Expression, expect *ast.Type) {
		stack = append(stack, S{Expr: expression, Expect: expect})
	}

	push(e, expect)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch e := top.Expr.(type) {
		case *ConstrAny:
		case *ConstrEquals:
			push(e.Value, top.Expect)
		case *ConstrNotEquals:
			push(e.Value, top.Expect)
		case *ConstrGreater:
			if top.Expect == nil {
				if !p.isNumeric(e.Value) {
					ok = false
					p.errExpectedNum(e.Value)
					break
				}
				push(e.Value, top.Expect)
				break
			}
			if top.Expect.NamedType != "Float" &&
				top.Expect.NamedType != "Int" {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
				break
			}
			push(e.Value, top.Expect)
		case *ConstrGreaterOrEqual:
			if top.Expect == nil {
				if !p.isNumeric(e.Value) {
					ok = false
					p.errExpectedNum(e.Value)
					break
				}
				push(e.Value, top.Expect)
				break
			}
			if top.Expect.NamedType != "Float" &&
				top.Expect.NamedType != "Int" {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
				break
			}
			push(e.Value, top.Expect)
		case *ConstrLess:
			if top.Expect == nil {
				if !p.isNumeric(e.Value) {
					ok = false
					p.errExpectedNum(e.Value)
					break
				}
				push(e.Value, top.Expect)
				break
			}
			if top.Expect.NamedType != "Float" &&
				top.Expect.NamedType != "Int" {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
				break
			}
			push(e.Value, top.Expect)
		case *ConstrLessOrEqual:
			if top.Expect == nil {
				if !p.isNumeric(e.Value) {
					ok = false
					p.errExpectedNum(e.Value)
					break
				}
				push(e.Value, top.Expect)
				break
			}
			if top.Expect.NamedType != "Float" &&
				top.Expect.NamedType != "Int" {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
				break
			}
			push(e.Value, top.Expect)
		case *ConstrLenEquals:
			push(e.Value, &ast.Type{
				NamedType: "Int",
				NonNull:   true,
			})
		case *ConstrLenNotEquals:
			push(e.Value, &ast.Type{
				NamedType: "Int",
				NonNull:   true,
			})
		case *ConstrLenGreater:
			push(e.Value, &ast.Type{
				NamedType: "Int",
				NonNull:   true,
			})
		case *ConstrLenGreaterOrEqual:
			push(e.Value, &ast.Type{
				NamedType: "Int",
				NonNull:   true,
			})
		case *ConstrLenLess:
			push(e.Value, &ast.Type{
				NamedType: "Int",
				NonNull:   true,
			})
		case *ConstrLenLessOrEqual:
			push(e.Value, &ast.Type{
				NamedType: "Int",
				NonNull:   true,
			})
		case *ConstrMap:
			var exp *ast.Type
			if top.Expect != nil {
				exp = top.Expect.Elem
			}
			push(e.Constraint, exp)
		case *ExprParentheses:
			push(e.Expression, top.Expect)
		case *ExprEqual:
			if !p.assumeComparableType(e.Location, e.Left, e.Right) {
				ok = false
				break
			}
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprNotEqual:
			if !p.assumeComparableType(e.Location, e.Left, e.Right) {
				ok = false
				break
			}
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprLogicalNegation:
			if top.Expect != nil && !isTypeBool(top.Expect) {
				p.errUnexpType(top.Expect, e)
				ok = false
				break
			} else if top.Expect == nil {
				top.Expect = &ast.Type{NamedType: "Boolean"}
			}
			push(e.Expression, top.Expect)
		case *ExprLogicalOr:
			objEncountered := false
			for _, e := range e.Expressions {
				o := getConstrEqValue[*Object](e)
				if o != nil && objEncountered {
					ok = false
					p.newErr(o.Location, "use single object "+
						"with multiple field constraints instead of "+
						"multiple object variants in an OR statement")
					continue
				} else if o != nil {
					objEncountered = true
				}
				push(o, top.Expect)
			}
		case *ExprLogicalAnd:
			objEncountered := false
			for _, e := range e.Expressions {
				o := getConstrEqValue[*Object](e)
				if o != nil && objEncountered {
					ok = false
					p.newErr(o.Location, "use single object "+
						"with multiple field constraints instead of "+
						"multiple object variants in an AND statement")
					continue
				} else if o != nil {
					objEncountered = true
				}
				push(o, top.Expect)
			}
		case *ExprAddition:
			if top.Expect == nil {
				br := false
				if !p.isNumeric(e.AddendLeft) {
					p.errExpectedNum(e.AddendLeft)
					br = true
				}
				if !p.isNumeric(e.AddendRight) {
					p.errExpectedNum(e.AddendRight)
					br = true
				}
				if br {
					ok = false
					break
				}
			} else if !isTypeNum(top.Expect) {
				p.errUnexpType(top.Expect, e)
				ok = false
				break
			}
			push(e.AddendLeft, top.Expect)
			push(e.AddendRight, top.Expect)
		case *ExprSubtraction:
			if top.Expect == nil {
				br := false
				if !p.isNumeric(e.Minuend) {
					p.errExpectedNum(e.Minuend)
					br = true
				}
				if !p.isNumeric(e.Subtrahend) {
					p.errExpectedNum(e.Subtrahend)
					br = true
				}
				if br {
					ok = false
					break
				}
			} else if !isTypeNum(top.Expect) {
				p.errUnexpType(top.Expect, e)
				ok = false
				break
			}
			push(e.Minuend, top.Expect)
			push(e.Subtrahend, top.Expect)
		case *ExprMultiplication:
			if top.Expect == nil {
				br := false
				if !p.isNumeric(e.Multiplicant) {
					p.errExpectedNum(e.Multiplicant)
					br = true
				}
				if !p.isNumeric(e.Multiplicator) {
					p.errExpectedNum(e.Multiplicator)
					br = true
				}
				if br {
					ok = false
					break
				}
			} else if !isTypeNum(top.Expect) {
				p.errUnexpType(top.Expect, e)
				ok = false
				break
			}
			push(e.Multiplicant, top.Expect)
			push(e.Multiplicator, top.Expect)
		case *ExprDivision:
			if top.Expect == nil {
				br := false
				if !p.isNumeric(e.Dividend) {
					p.errExpectedNum(e.Dividend)
					br = true
				}
				if !p.isNumeric(e.Divisor) {
					p.errExpectedNum(e.Divisor)
					br = true
				}
				if br {
					ok = false
					break
				}
			} else if !isTypeNum(top.Expect) {
				p.errUnexpType(top.Expect, e)
				ok = false
				break
			}
			push(e.Dividend, top.Expect)
			push(e.Divisor, top.Expect)
		case *ExprModulo:
			if top.Expect == nil {
				br := false
				if !p.isNumeric(e.Dividend) {
					p.errExpectedNum(e.Dividend)
					br = true
				}
				if !p.isNumeric(e.Divisor) {
					p.errExpectedNum(e.Divisor)
					br = true
				}
				if br {
					ok = false
					break
				}
			} else if !isTypeNum(top.Expect) {
				p.errUnexpType(top.Expect, e)
				ok = false
				break
			}
			push(e.Dividend, top.Expect)
			push(e.Divisor, top.Expect)
		case *Int:
			if top.Expect != nil && !isTypeNum(top.Expect) {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *Float:
			if top.Expect != nil && !isTypeFloat(top.Expect) {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *True:
			if top.Expect != nil && !isTypeBool(top.Expect) {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *False:
			if top.Expect != nil && !isTypeBool(top.Expect) {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *Null:
			if top.Expect != nil && top.Expect.NonNull {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *Array:
			if top.Expect != nil && !typeIsArray(top.Expect) {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
				break
			}
			for _, i := range e.Items {
				var exp *ast.Type
				if top.Expect != nil {
					exp = top.Expect.Elem
				}
				push(i, exp)
			}
		case *VariableRef:
			for _, o := range pathToOriginArg {
				if o == e.Declaration.Parent {
					ok = false
					p.errConstrSelf(e)
					break
				}
			}

			if top.Expect != nil {
				if e.Declaration.Type().NamedType != top.Expect.NamedType {
					p.errUnexpType(top.Expect, top.Expr)
				}
			}
		case *Enum:
			if p.schema != nil && e.TypeDef == nil {
				ok = false
				p.errUndefEnumVal(e)
				break
			} else if e.TypeDef != nil &&
				top.Expect != nil &&
				top.Expect.NamedType != e.TypeDef.Name {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *String:
			if top.Expect != nil && top.Expect.NamedType != "String" {
				ok = false
				p.errUnexpType(top.Expect, top.Expr)
			}
		case *Object:
			if top.Expect != nil && p.schema != nil {
				def := p.schema.Types[top.Expect.NamedType]
				if def != nil && def.Kind != ast.InputObject {
					ok = false
					p.errUnexpType(top.Expect, top.Expr)
					break
				}
			}
			p.validateObject(pathToOriginArg, e, top.Expect)
		default:
			panic(fmt.Errorf("unhandled type: %T", top.Expr))
		}
	}
	return ok
}

// getConstrEqValue returns E if e is ConstrEquals or ConstrNotEquals
// and the value is of type E, otherwise returns the zero value of E.
func getConstrEqValue[E Expression](e Expression) (x E) {
	switch e := e.(type) {
	case *ConstrEquals:
		if y, ok := e.Value.(E); ok {
			return y
		}
	case *ConstrNotEquals:
		if y, ok := e.Value.(E); ok {
			return y
		}
	}
	return
}

func (p *Parser) validateObject(
	pathToOriginArg []any,
	o *Object,
	exp *ast.Type,
) (ok bool) {
	ok = true
	if exp != nil && (exp.Elem != nil ||
		exp.NamedType == "ID" ||
		exp.NamedType == "String" ||
		exp.NamedType == "Int" ||
		exp.NamedType == "Float" ||
		exp.NamedType == "Boolean") {
		ok = false
		p.errUnexpType(exp, o)
	}

	fieldNames := make(map[string]struct{}, len(o.Fields))
	validFields := o.Fields
	if p.schema != nil {
		validFields = make([]*ObjectField, 0, len(o.Fields))
		td := p.schema.Types[exp.NamedType]

		for _, f := range o.Fields {
			// Check undefined fields
			fd := td.Fields.ForName(f.Name)
			if fd == nil {
				ok = false
				p.errUndefField(f, exp.NamedType)
				continue
			}
			fieldNames[f.Name] = struct{}{}
			validFields = append(validFields, f)
		}

		// Check required fields
		for _, f := range td.Fields {
			if f.Type.NonNull && f.DefaultValue == nil {
				if _, ok := fieldNames[f.Name]; !ok {
					ok = false
					p.errMissingInputField(o.Location, f)
				}
			}
		}
	}

	// Check constraints
	for _, f := range validFields {
		var exp *ast.Type
		if f.Def != nil {
			exp = f.Def.Type
		}
		if !p.validateExpr(
			makeAppend[any](pathToOriginArg, f),
			f.Constraint,
			exp,
		) {
			ok = false
		}
	}

	return ok
}

func (p *Parser) errConstrSelf(v *VariableRef) {
	var on string
	switch v.Declaration.Parent.(type) {
	case *Argument:
		on = "argument"
	case *ObjectField:
		on = "object field"
	}
	p.errors = append(p.errors, Error{
		Location: v.Location,
		Msg:      on + " self reference in constraint",
	})
}

func (p *Parser) errUnexpTok(s source, msg string) {
	prefix := "unexpected token, "
	if s.Index >= len(s.s) {
		prefix = "unexpected end of file, "
	}
	p.errors = append(p.errors, Error{
		Location: s.Location,
		Msg:      prefix + msg,
	})
}

func (p *Parser) errTypenameInMax(f *SelectionField) {
	p.errors = append(p.errors, Error{
		Location: f.Location,
		Msg:      "avoid __typename in max selectors",
	})
}

func (p *Parser) errNestedMaxBlock(s *SelectionMax) {
	p.errors = append(p.errors, Error{
		Location: s.Location,
		Msg:      "nested max block",
	})
}

func (p *Parser) errUndefType(l Location, name string) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg:      "type " + name + " is undefined in schema",
	})
}

func errFieldTypenameCantHaveArgs() string {
	return "built-in field \"__typename\" can never have arguments"
}

func errFieldTypenameCantHaveSels() string {
	return "built-in field \"__typename\" can never have subselections"
}

func (p *Parser) errUndefField(f *ObjectField, hostTypeName string) {
	p.errors = append(p.errors, Error{
		Location: f.Location,
		Msg: fmt.Sprintf(
			"field %q is undefined in type %s",
			f.Name, hostTypeName,
		),
	})
}

func (p *Parser) errUndefArg(
	a *Argument,
	f *SelectionField,
	hostTypeName string,
) {
	p.errors = append(p.errors, Error{
		Location: a.Location,
		Msg: fmt.Sprintf(
			"argument %q is undefined on field %q in type %s",
			a.Name, f.Name, hostTypeName,
		),
	})
}

func (p *Parser) errMismatchingTypes(l Location, left, right Expression) {
	ld := left.TypeDesignation()
	rd := right.TypeDesignation()
	p.errors = append(p.errors, Error{
		Location: l,
		Msg:      fmt.Sprintf("mismatching types %s and %s", ld, rd),
	})
}

func (p *Parser) newErr(l Location, msg string) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg:      msg,
	})
}

func (p *Parser) newErrorNoPrefix(s source, msg string) {
	p.errors = append(p.errors, Error{
		Location: s.Location,
		Msg:      msg,
	})
}

type expect int8

const (
	expectConstraint      expect = 0
	expectConstraintInMap expect = 1
	expectValue           expect = 2
	expectValueInMap      expect = 3
)

func (p *Parser) ParseSelectionSet(s source) (source, SelectionSet) {
	selset := SelectionSet{Location: s.Location}
	si := s
	var ok bool
	if s, ok = s.consume("{"); !ok {
		p.errUnexpTok(s, "expected selection set")
		return stop(), SelectionSet{}
	}

	fields := map[string]struct{}{}
	typeConds := map[string]struct{}{}

	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			p.errUnexpTok(s, "expected selection")
			return stop(), SelectionSet{}
		}

		if s, ok = s.consume("}"); ok {
			break
		}

		s = s.consumeIgnored()

		if _, ok = s.consume("..."); ok {
			var fragInline *SelectionInlineFrag
			sBefore := s
			if s, fragInline = p.ParseInlineFrag(s); s.stop() {
				return stop(), SelectionSet{}
			}

			if _, ok = typeConds[fragInline.TypeCondition.TypeName]; ok {
				p.newErr(sBefore.Location, "redeclared type condition")
				return stop(), SelectionSet{}
			}
			typeConds[fragInline.TypeCondition.TypeName] = struct{}{}

			selset.Selections = append(selset.Selections, fragInline)
			continue
		}

		lBeforeName := s.Location
		sel := &SelectionField{Location: s.Location}
		if s, name = s.consumeName(); name == nil {
			p.errUnexpTok(s, "expected selection")
			return stop(), SelectionSet{}
		}
		sel.Name = string(name)

		s = s.consumeIgnored()

		if sel.Name == "max" {
			lBeforeMaxNum := s.Location
			var maxNum int64
			if s, maxNum, ok = s.consumeUnsignedInt(); ok {
				// Max
				s = s.consumeIgnored()

				if maxNum < 1 {
					p.newErr(
						lBeforeMaxNum,
						"maximum number of options"+
							" must be an unsigned integer greater 0",
					)
					return stop(), SelectionSet{}
				}

				var options SelectionSet
				sBeforeOptionsBlock := s
				if s, options = p.ParseSelectionSet(s); s.stop() {
					return stop(), SelectionSet{}
				}

				if len(options.Selections) < 2 {
					p.newErr(
						sBeforeOptionsBlock.Location,
						"max block must have at least 2 selection options",
					)
					return stop(), SelectionSet{}
				} else if maxNum > int64(len(options.Selections)-1) {
					p.newErr(
						lBeforeMaxNum,
						"max selections number exceeds number of options-1",
					)
					return stop(), SelectionSet{}
				}

				e := &SelectionMax{
					Location: lBeforeName,
					Limit:    int(maxNum),
					Options:  options,
				}
				selset.Selections = append(selset.Selections, e)
				continue
			}
		}

		if _, ok := fields[sel.Name]; ok {
			p.newErr(lBeforeName, "redeclared field")
			return stop(), SelectionSet{}
		}
		fields[sel.Name] = struct{}{}

		if s.peek1('(') {
			if sel.Name == "__typename" {
				p.newErr(s.Location, errFieldTypenameCantHaveArgs())
				return stop(), SelectionSet{}
			}
			if s, sel.ArgumentList = p.ParseArguments(s); s.stop() {
				return stop(), SelectionSet{}
			}
			for _, arg := range sel.Arguments {
				setParent(arg, sel)
			}
		}

		s = s.consumeIgnored()

		if s.peek1('{') {
			if sel.Name == "__typename" {
				p.newErr(s.Location, errFieldTypenameCantHaveSels())
				return stop(), SelectionSet{}
			}

			if s, sel.SelectionSet = p.ParseSelectionSet(s); s.stop() {
				return stop(), SelectionSet{}
			}
			for _, sub := range sel.Selections {
				setParent(sub, sel)
			}
			if !p.validateSelectionSet(s.s, sel.Selections) {
				return stop(), SelectionSet{}
			}
		}

		selset.Selections = append(selset.Selections, sel)
	}

	if len(selset.Selections) < 1 {
		p.newErrorNoPrefix(si, "empty selection set")
		return stop(), SelectionSet{}
	}

	return s, selset
}

func (p *Parser) ParseInlineFrag(s source) (source, *SelectionInlineFrag) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("..."); !ok {
		p.errUnexpTok(s, "expected '...'")
		return stop(), nil
	}

	inlineFrag := &SelectionInlineFrag{Location: l}
	s = s.consumeIgnored()

	var tok []byte
	sp := s
	s, tok = s.consumeToken()
	if string(tok) != "on" {
		p.errUnexpTok(sp, "expected keyword 'on'")
		return stop(), nil
	}

	s = s.consumeIgnored()

	sBeforeName := s
	var name []byte
	s, name = s.consumeName()
	if len(name) < 1 {
		p.errUnexpTok(s, "expected type condition")
		return stop(), nil
	}
	inlineFrag.TypeCondition = TypeCondition{
		Location: sBeforeName.Location,
		TypeName: string(name),
	}

	s = s.consumeIgnored()
	var selset SelectionSet
	if s, selset = p.ParseSelectionSet(s); s.stop() {
		return stop(), nil
	}
	for _, sel := range selset.Selections {
		setParent(sel, inlineFrag)
	}

	inlineFrag.SelectionSet = selset
	return s, inlineFrag
}

func (p *Parser) ParseArguments(s source) (source, ArgumentList) {
	si := s
	var ok bool
	if s, ok = s.consume("("); !ok {
		p.errUnexpTok(s, "expected opening parenthesis")
		return stop(), ArgumentList{}
	}

	list := ArgumentList{Location: si.Location}
	names := map[string]struct{}{}
	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			p.errUnexpTok(s, "expected argument")
			return stop(), ArgumentList{}
		}

		if s, ok = s.consume(")"); ok {
			break
		}

		s = s.consumeIgnored()

		sBeforeName := s
		arg := &Argument{Location: s.Location}
		if s, name = s.consumeName(); name == nil {
			p.errUnexpTok(s, "expected argument name")
			return stop(), ArgumentList{}
		}
		arg.Name = string(name)

		if _, ok := names[arg.Name]; ok {
			p.newErr(sBeforeName.Location, "redeclared argument")
			return stop(), ArgumentList{}
		}

		names[arg.Name] = struct{}{}

		s = s.consumeIgnored()

		if s, ok = s.consume("="); ok {
			// Has an associated variable name
			s = s.consumeIgnored()

			sBeforeDollar := s
			if s, ok = s.consume("$"); !ok {
				p.errUnexpTok(s, "expected variable name")
				return stop(), ArgumentList{}
			}

			var name []byte
			if s, name = s.consumeName(); name == nil {
				p.errUnexpTok(s, "expected variable name")
				return stop(), ArgumentList{}
			}

			def := &VariableDeclaration{
				Location: sBeforeDollar.Location,
				Parent:   arg,
				Name:     string(name),
			}
			arg.AssociatedVariable = def
			if _, ok := p.varDecls[def.Name]; ok {
				p.newErr(sBeforeDollar.Location, "redeclared variable")
				return stop(), ArgumentList{}
			}
			p.varDecls[def.Name] = def
			s = s.consumeIgnored()
		}

		if s, ok = s.consume(":"); !ok {
			p.errUnexpTok(s, "expected colon")
			return stop(), ArgumentList{}
		}
		s = s.consumeIgnored()
		if s.isEOF() {
			p.errUnexpTok(s, "expected constraint")
			return stop(), ArgumentList{}
		}

		var expr Expression
		if s, expr = p.ParseExprLogicalOr(
			s, expectConstraint,
		); s.stop() {
			return s, ArgumentList{}
		}
		setParent(expr, arg)
		arg.Constraint = expr

		list.Arguments = append(list.Arguments, arg)
		s = s.consumeIgnored()

		if s, ok = s.consume(","); !ok {
			if s, ok = s.consume(")"); !ok {
				p.errUnexpTok(s, "expected comma or end of argument list")
				return stop(), ArgumentList{}
			}
			break
		}
		s = s.consumeIgnored()
	}

	if len(list.Arguments) < 1 {
		p.newErrorNoPrefix(si, "empty argument list")
		return stop(), ArgumentList{}
	}

	return s, list
}

func (p *Parser) ParseValue(
	s source,
	expect expect,
) (source, Expression) {
	l := s.Location

	if s.isEOF() {
		p.errUnexpTok(s, "expected value")
		return stop(), nil
	}

	var str []byte
	var num any
	var ok bool
	if s, ok = s.consume("("); ok {
		e := &ExprParentheses{Location: l}

		s = s.consumeIgnored()

		if s, e.Expression = p.ParseExprLogicalOr(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Expression, e)

		if s, ok = s.consume(")"); !ok {
			p.errUnexpTok(s, "missing closing parenthesis")
			return stop(), nil
		}

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume("$"); ok {
		v := &VariableRef{Location: l}

		var name []byte
		if s, name = s.consumeName(); name == nil {
			p.errUnexpTok(s, "expected variable name")
			return stop(), nil
		}

		v.Name = string(name)

		p.varRefs = append(p.varRefs, v)

		s = s.consumeIgnored()

		return s, v
	} else if s, str, ok = s.consumeString(); ok {
		return s, &String{Location: l, Value: string(str)}
	}

	if s, num = p.ParseNumber(s); s.stop() {
		return stop(), nil
	} else if num != nil {
		switch v := num.(type) {
		case *Int:
			return s, v
		case *Float:
			return s, v
		default:
			panic(fmt.Errorf("unexpected number type: %#v", num))
		}
	}

	if s, ok = s.consume("["); ok {
		e := &Array{Location: l}
		s = s.consumeIgnored()

		for {
			if s, ok = s.consume("]"); ok {
				break
			}

			var expr Expression
			if s, expr = p.ParseExprLogicalOr(s, expectConstraint); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)

			e.Items = append(e.Items, expr)

			s = s.consumeIgnored()

			if s, ok = s.consume(","); !ok {
				if s, ok = s.consume("]"); !ok {
					p.errUnexpTok(s, "expected comma or end of array")
					return stop(), nil
				}
				break
			}
			s = s.consumeIgnored()
		}

		s = s.consumeIgnored()
		return s, e
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
				p.errUnexpTok(s, "expected object field")
				return stop(), nil
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
				p.errUnexpTok(s, "expected object field name")
				return stop(), nil
			}
			fld.Name = string(name)

			if _, ok := fieldNames[fld.Name]; ok {
				p.newErr(sBeforeName.Location, "redeclared object field")
				return stop(), nil
			}

			fieldNames[fld.Name] = struct{}{}

			s = s.consumeIgnored()

			if s, ok = s.consume("="); ok {
				// Has an associated variable name
				s = s.consumeIgnored()

				sBeforeDollar := s
				if s, ok = s.consume("$"); !ok {
					p.errUnexpTok(s, "expected variable name")
					return stop(), nil
				}

				var name []byte
				if s, name = s.consumeName(); name == nil {
					p.errUnexpTok(s, "expected variable name")
					return stop(), nil
				}

				if expect == expectValueInMap {
					p.newErr(
						sBeforeDollar.Location,
						"the use of variables inside "+
							"map constraints is prohibited",
					)
					return stop(), nil
				}

				def := &VariableDeclaration{
					Location: sBeforeDollar.Location,
					Parent:   fld,
					Name:     string(name),
				}
				fld.AssociatedVariable = def
				if _, ok := p.varDecls[def.Name]; ok {
					p.newErr(sBeforeDollar.Location, "redeclared variable")
					return stop(), nil
				}
				p.varDecls[def.Name] = def

				s = s.consumeIgnored()
			}

			if s, ok = s.consume(":"); !ok {
				p.errUnexpTok(s, "expected colon")
				return stop(), nil
			}
			s = s.consumeIgnored()
			if s.isEOF() {
				p.errUnexpTok(s, "expected constraint")
				return stop(), nil
			}

			var expr Expression
			if s, expr = p.ParseExprLogicalOr(s, expectConstraint); s.stop() {
				return stop(), nil
			}
			setParent(expr, fld)
			fld.Constraint = expr

			o.Fields = append(o.Fields, fld)
			s = s.consumeIgnored()

			if s, ok = s.consume(","); !ok {
				if s, ok = s.consume("}"); !ok {
					p.errUnexpTok(s, "expected comma or end of object")
					return stop(), nil
				}
				break
			}
			s = s.consumeIgnored()
		}

		if len(o.Fields) < 1 {
			p.newErrorNoPrefix(si, "empty input object")
			return stop(), nil
		}

		return s, o
	}

	s, str = s.consumeName()
	switch string(str) {
	case "true":
		return s, &True{Location: l}
	case "false":
		return s, &False{Location: l}
	case "null":
		return s, &Null{Location: l}
	case "":
		p.errUnexpTok(s, "invalid value")
		return stop(), nil
	default:
		return s, &Enum{Location: l, Value: string(str)}
	}
}

func (p *Parser) ParseExprUnary(
	s source,
	expect expect,
) (source, Expression) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("!"); ok {
		e := &ExprLogicalNegation{Location: l}

		s = s.consumeIgnored()

		if s, e.Expression = p.ParseValue(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Expression, e)

		s = s.consumeIgnored()
		return s, e
	}
	return p.ParseValue(s, expect)
}

func (p *Parser) ParseExprMultiplicative(
	s source,
	expect expect,
) (source, Expression) {
	si := s

	var result Expression
	if s, result = p.ParseExprUnary(s, expect); s.stop() {
		return stop(), nil
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

			s = s.consumeIgnored()
			if s, e.Multiplicator = p.ParseExprUnary(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Multiplicator, e)

			e.Float = e.Multiplicant.IsFloat() || e.Multiplicator.IsFloat()

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprDivision{
				Location: si.Location,
				Dividend: result,
			}
			setParent(e.Dividend, e)

			s = s.consumeIgnored()
			if s, e.Divisor = p.ParseExprUnary(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Divisor, e)

			e.Float = e.Dividend.IsFloat() || e.Divisor.IsFloat()

			s = s.consumeIgnored()
			result = e
		case 2:
			e := &ExprModulo{
				Location: si.Location,
				Dividend: result,
			}
			setParent(e.Dividend, e)

			s = s.consumeIgnored()
			if s, e.Divisor = p.ParseExprUnary(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Divisor, e)

			e.Float = e.Dividend.IsFloat() || e.Divisor.IsFloat()

			s = s.consumeIgnored()
			result = e
		default:
			s = s.consumeIgnored()
			return s, result
		}
	}
}

func (p *Parser) ParseExprAdditive(
	s source,
	expect expect,
) (source, Expression) {
	si := s

	var result Expression
	if s, result = p.ParseExprMultiplicative(s, expect); s.stop() {
		return stop(), nil
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

			s = s.consumeIgnored()
			if s, e.AddendRight = p.ParseExprMultiplicative(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.AddendRight, e)

			e.Float = e.AddendLeft.IsFloat() || e.AddendRight.IsFloat()

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprSubtraction{
				Location: si.Location,
				Minuend:  result,
			}
			setParent(e.Minuend, e)

			s = s.consumeIgnored()
			if s, e.Subtrahend = p.ParseExprMultiplicative(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Subtrahend, e)

			e.Float = e.Minuend.IsFloat() || e.Subtrahend.IsFloat()

			s = s.consumeIgnored()
			result = e
		default:
			s = s.consumeIgnored()
			return s, result
		}
	}
}

func (p *Parser) ParseExprRelational(
	s source,
	expect expect,
) (source, Expression) {
	l := s.Location

	var left Expression
	if s, left = p.ParseExprAdditive(s, expect); s.stop() {
		return stop(), nil
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

		if s, e.Right = p.ParseExprAdditive(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume(">="); ok {
		e := &ExprGreaterOrEqual{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.ParseExprAdditive(s, expect); s.stop() {
			return s, nil
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume("<"); ok {
		e := &ExprLess{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.ParseExprAdditive(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume(">"); ok {
		e := &ExprGreater{
			Location: l,
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.ParseExprAdditive(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)

		s = s.consumeIgnored()
		return s, e
	}
	s = s.consumeIgnored()
	return s, left
}

func (p *Parser) ParseExprEquality(
	s source,
	expect expect,
) (source, Expression) {
	si := s

	var exprLeft Expression
	if s, exprLeft = p.ParseExprRelational(s, expect); s.stop() {
		return stop(), nil
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

		if s, e.Right = p.ParseExprRelational(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		s = s.consumeIgnored()

		return s, e
	}
	if s, ok = s.consume("!="); ok {
		e := &ExprNotEqual{
			Location: si.Location,
			Left:     exprLeft,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.ParseExprRelational(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		s = s.consumeIgnored()

		return s, e
	}

	s = s.consumeIgnored()
	return s, exprLeft
}

func (p *Parser) ParseExprLogicalOr(
	s source,
	expect expect,
) (source, Expression) {
	e := &ExprLogicalOr{Location: s.Location}

	for {
		var expr Expression
		if s, expr = p.ParseExprLogicalAnd(
			s, expect,
		); s.stop() {
			return s, nil
		}
		setParent(expr, e)
		e.Expressions = append(e.Expressions, expr)

		s = s.consumeIgnored()

		var ok bool
		if s, ok = s.consume("||"); !ok {
			if len(e.Expressions) < 2 {
				return s, e.Expressions[0]
			}
			return s, e
		}

		s = s.consumeIgnored()
	}
}

func (p *Parser) ParseExprLogicalAnd(
	s source,
	expect expect,
) (source, Expression) {
	e := &ExprLogicalAnd{Location: s.Location}

	for {
		var expr Expression
		if s, expr = p.ParseConstr(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Expressions = append(e.Expressions, expr)

		s = s.consumeIgnored()

		var ok bool
		if s, ok = s.consume("&&"); !ok {
			if len(e.Expressions) < 2 {
				return s, e.Expressions[0]
			}
			return s, e
		}

		s = s.consumeIgnored()
	}
}

func (p *Parser) ParseConstr(
	s source,
	expect expect,
) (source, Expression) {
	si := s
	var ok bool
	var expr Expression

	if expect == expectValue {
		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}

		s = s.consumeIgnored()

		return s, expr
	}

	if s, ok = s.consume("*"); ok {
		e := &ConstrAny{Location: si.Location}
		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("!="); ok {
		e := &ConstrNotEquals{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("<="); ok {
		e := &ConstrLessOrEqual{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume(">="); ok {
		e := &ConstrGreaterOrEqual{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("<"); ok {
		e := &ConstrLess{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume(">"); ok {
		e := &ConstrGreater{
			Location: si.Location,
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("len"); ok {
		s = s.consumeIgnored()

		if s, ok = s.consume("!="); ok {
			e := &ConstrLenNotEquals{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			var expr Expression

			if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume("<="); ok {
			e := &ConstrLenLessOrEqual{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume(">="); ok {
			e := &ConstrLenGreaterOrEqual{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume("<"); ok {
			e := &ConstrLenLess{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume(">"); ok {
			e := &ConstrLenGreater{
				Location: si.Location,
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr

			s = s.consumeIgnored()

			return s, e
		}

		e := &ConstrLenEquals{
			Location: si.Location,
			Value:    expr,
		}

		if s, expr = p.ParseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr

		s = s.consumeIgnored()

		return s, e
	}

	if s, ok = s.consume("["); ok {
		s = s.consumeIgnored()
		if s, ok = s.consume("..."); ok {
			e := &ConstrMap{Location: si.Location}
			s = s.consumeIgnored()

			if s.isEOF() {
				p.errUnexpTok(s, "expected map constraint")
				return stop(), nil
			}

			var expr Expression
			if s, expr = p.ParseExprLogicalOr(
				s, expectConstraintInMap,
			); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Constraint = expr

			s = s.consumeIgnored()
			if s, ok = s.consume("]"); !ok {
				p.errUnexpTok(s, "expected end of map constraint ']'")
				return stop(), nil
			}

			return s, e
		} else {
			s = si
		}
	}

	e := &ConstrEquals{
		Location: si.Location,
		Value:    expr,
	}

	if expect == expectConstraintInMap {
		expect = expectValueInMap
	} else {
		expect = expectValue
	}

	if s, expr = p.ParseExprEquality(s, expect); s.stop() {
		return stop(), nil
	}
	setParent(expr, e)
	e.Value = expr

	s = s.consumeIgnored()

	return s, e
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

func (p *Parser) ParseNumber(s source) (_ source, number any) {
	si := s
	if s.Index >= len(s.s) {
		return si, nil
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
				return si, nil
			}
			s.Index++
			s.Column++
			if s.Index >= len(s.s) {
				p.newErr(s.Location, "exponent has no digits")
				return stop(), nil
			}
			s, _ = s.consumeEitherOf3("+", "-", "")
			if !s.isDigit() {
				p.newErr(s.Location, "exponent has no digits")
				return stop(), nil
			}
		}

		s.Index++
		s.Column++
	}
	str := string(s.s[si.Index:s.Index])
	if str == "" {
		return si, nil
	} else if v, err := strconv.ParseInt(str, 10, 64); err == nil {
		return s, &Int{
			Location: si.Location,
			Value:    v,
		}
	} else if v, err := strconv.ParseFloat(str, 64); err == nil {
		return s, &Float{
			Location: si.Location,
			Value:    v,
		}
	}
	return si, nil
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

func (s source) stop() bool { return s.s == nil }

func stop() source { return source{} }

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

func (p *Parser) isNumeric(expr any) bool {
	switch e := expr.(type) {
	case *VariableRef:
		if t := e.Declaration.Type(); t != nil {
			// Check type based on schema
			return t.Elem == nil &&
				(t.NamedType == "Int" ||
					t.NamedType == "Float")
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isNumeric(v.Constraint)
		case *ObjectField:
			return p.isNumeric(v.Constraint)
		}
	case *ConstrEquals:
		return p.isNumeric(e.Value)
	case *ConstrNotEquals:
		return p.isNumeric(e.Value)
	case *ConstrAny:
		return true
	case *ExprParentheses:
		return p.isNumeric(e.Expression)
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

func (p *Parser) isInt(expr any) bool {
	switch e := expr.(type) {
	case *VariableRef:
		if t := e.Declaration.Type(); t != nil {
			// Check type based on schema
			return t.Elem == nil && t.NamedType == "Int"
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isInt(v.Constraint)
		case *ObjectField:
			return p.isInt(v.Constraint)
		}
	case *ExprParentheses:
		return p.isInt(e.Expression)
	case *Int:
		return true
	case *ExprAddition:
		return !e.Float
	case *ExprSubtraction:
		return !e.Float
	case *ExprMultiplication:
		return !e.Float
	case *ExprDivision:
		return !e.Float
	case *ExprModulo:
		return !e.Float
	}
	return false
}

func (p *Parser) isBoolean(expr any) bool {
	switch e := expr.(type) {
	case *VariableRef:
		if t := e.Declaration.Type(); t != nil {
			// Check type based on schema
			return t.Elem == nil && t.NamedType == "Boolean"
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isBoolean(v.Constraint)
		case *ObjectField:
			return p.isBoolean(v.Constraint)
		}
	case *ConstrAny:
		return true
	case *ConstrEquals:
		return p.isBoolean(e.Value)
	case *ConstrNotEquals:
		return p.isBoolean(e.Value)
	case *ConstrGreater:
		return true
	case *ConstrGreaterOrEqual:
		return true
	case *ConstrLess:
		return true
	case *ConstrLessOrEqual:
		return true
	case *ExprParentheses:
		return p.isBoolean(e.Expression)
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

func (p *Parser) isString(expr any) bool {
	switch e := expr.(type) {
	case *VariableRef:
		if t := e.Declaration.Type(); t != nil {
			// Check type based on schema
			return t.Elem == nil && t.NamedType == "String"
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isString(v.Constraint)
		case *ObjectField:
			return p.isString(v.Constraint)
		}
	case *ConstrAny:
		return true
	case *ConstrEquals:
		return p.isString(e.Value)
	case *ConstrNotEquals:
		return p.isString(e.Value)
	case *ExprParentheses:
		return p.isString(e.Expression)
	case *String:
		return true
	}
	return false
}

func (p *Parser) assumeComparableType(
	l Location,
	left, right Expression,
) (ok bool) {
	if p.isBoolean(left) && !p.isBoolean(right) {
		p.errMismatchingTypes(l, left, right)
		return false
	} else if p.isNumeric(left) && !p.isNumeric(right) {
		p.errMismatchingTypes(l, left, right)
		return false
	} else if left.IsFloat() && !right.IsFloat() ||
		!left.IsFloat() && right.IsFloat() {
		p.errMismatchingTypes(l, left, right)
		return false
	}
	switch left.(type) {
	case *String:
		if _, ok := right.(*String); !ok {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case *Enum:
		if _, ok := right.(*Enum); !ok {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case *Array:
		if _, ok := right.(*Array); !ok {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case *Object:
		if _, ok := right.(*Object); !ok {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	}
	return true
}

func setParent(t, parent Expression) {
	switch v := t.(type) {
	case *VariableRef:
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

func (p *Parser) validateSelectionSet(s []byte, set []Selection) (ok bool) {
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
			inlineFrags[v.TypeCondition.TypeName] = v
		}
	}

	if len(maxs) > 1 {
		p.newErr(
			maxs[len(maxs)-1].Location,
			"multiple max blocks in one selection set",
		)
		return false
	}

	for _, m := range maxs {
		for _, sel := range m.Options.Selections {
			switch v := sel.(type) {
			case *SelectionField:
				if f, ok := fields[v.Name]; ok {
					p.newErr(
						furthestLocation(f.Location, v.Location),
						"redeclared field",
					)
					return false
				}
			case *SelectionInlineFrag:
				if c, ok := inlineFrags[v.TypeCondition.TypeName]; ok {
					p.newErr(
						furthestLocation(c.Location, v.Location),
						"redeclared type condition",
					)
					return false
				}
			}
		}
	}
	return true
}

func furthestLocation(a, b Location) Location {
	if a.Index < b.Index {
		return b
	}
	return a
}

func (p *Parser) validateInlineFrag(
	frag *SelectionInlineFrag,
	hostDef *ast.Definition,
) (ok bool) {
	if p.schema == nil {
		if !p.validateCond(frag) {
			return false
		}
		p.validateSelSet(frag, nil)
		return true
	}

	def := p.schema.Types[frag.TypeCondition.TypeName]
	if def == nil {
		p.errUndefType(
			frag.TypeCondition.Location,
			frag.TypeCondition.TypeName,
		)
		return false
	}

	switch def.Kind {
	case ast.Scalar:
		p.errCondOnScalarType(frag)
		return false
	case ast.Enum:
		p.errCondOnEnumType(frag)
		return false
	case ast.InputObject:
		p.errCondOnInputType(frag)
		return false
	case ast.Object:
		if hostDef != nil && hostDef.Name == def.Name {
			return true
		}
		for _, c := range def.Interfaces {
			if c == hostDef.Name {
				// condition Object implements the host interface
				return true
			}
		}
		for _, c := range hostDef.Types {
			if c == def.Name {
				// condition Object is part of the host union
				return true
			}
		}
	case ast.Interface:
		if hostDef.Name == def.Name {
			return true
		}
		impls := p.schema.GetPossibleTypes(def)
		for _, t := range impls {
			if t == hostDef {
				return true
			}
			possibleTypes := p.schema.GetPossibleTypes(hostDef)
			for _, ht := range possibleTypes {
				if ht == t {
					return true
				}
			}
		}
	case ast.Union:
		if hostDef.Name == def.Name {
			return true
		}
		impls := p.schema.GetPossibleTypes(def)
		for _, t := range impls {
			if t == hostDef {
				return true
			}
			possibleTypes := p.schema.GetPossibleTypes(hostDef)
			for _, ht := range possibleTypes {
				if ht == t {
					return true
				}
			}
		}
	}
	p.errCanNeverBeOfType(
		frag.TypeCondition.Location, hostDef.Name, def.Name,
	)
	return false
}

func (p *Parser) validateCond(frag *SelectionInlineFrag) (ok bool) {
	ok = true
	switch frag.TypeCondition.TypeName {
	case "String", "ID", "Boolean", "Int", "Float":
		p.errCondOnScalarType(frag)
		return false
	case "Mutation", "Query", "Subscription":
		for par := frag.Parent; par != nil; par = par.GetParent() {
			switch par := par.(type) {
			case *Operation:
				if par.Type.String() != frag.TypeCondition.TypeName {
					p.errCanNeverBeOfType(
						frag.TypeCondition.Location,
						par.Type.String(),
						frag.TypeCondition.TypeName,
					)
					return false
				}
			case *SelectionInlineFrag:
				if par.TypeCondition.TypeName != frag.TypeCondition.TypeName {
					p.errCanNeverBeOfType(
						frag.TypeCondition.Location,
						par.TypeCondition.TypeName,
						frag.TypeCondition.TypeName,
					)
					return false
				}
			}
		}
	default:
		for par := frag.Parent; par != nil; par = par.GetParent() {
			switch par := par.(type) {
			case *Operation:
				p.errCanNeverBeOfType(
					frag.TypeCondition.Location,
					par.Type.String(),
					frag.TypeCondition.TypeName,
				)
				return false
			case *SelectionInlineFrag:
				switch par.TypeCondition.TypeName {
				case "Query", "Mutation", "Subscription":
					p.errCanNeverBeOfType(
						frag.TypeCondition.Location,
						frag.TypeCondition.TypeName,
						par.TypeCondition.TypeName,
					)
					return false
				default:
					return true
				}
			default:
				return true
			}
		}
	}
	return ok
}

func (p *Parser) errCondOnScalarType(s *SelectionInlineFrag) {
	p.errors = append(p.errors, Error{
		Location: s.TypeCondition.Location,
		Msg: "fragment can't condition on scalar type " +
			s.TypeCondition.TypeName,
	})
}

func (p *Parser) errCondOnEnumType(s *SelectionInlineFrag) {
	p.errors = append(p.errors, Error{
		Location: s.TypeCondition.Location,
		Msg: "fragment can't condition on enum type " +
			s.TypeCondition.TypeName,
	})
}

func (p *Parser) errCondOnInputType(s *SelectionInlineFrag) {
	p.errors = append(p.errors, Error{
		Location: s.TypeCondition.Location,
		Msg: "fragment can't condition on input type " +
			s.TypeCondition.TypeName,
	})
}

func (p *Parser) errCanNeverBeOfType(
	l Location,
	thisType,
	cantBeOf string,
) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg: "type " + thisType +
			" can never be of type " + cantBeOf,
	})
}

func (p *Parser) errMissingArg(
	l Location, missingArgument *ast.ArgumentDefinition,
) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg: fmt.Sprintf(
			"argument %q of type %s is required but missing",
			missingArgument.Name, missingArgument.Type,
		),
	})
}

func (p *Parser) errMissingInputField(
	l Location,
	missingField *ast.FieldDefinition,
) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg: fmt.Sprintf(
			"field %q of type %q is required but missing",
			missingField.Name, missingField.Type,
		),
	})
}

func (p *Parser) errMissingSelSet(
	l Location, fieldName, hostTypeName string,
) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg: fmt.Sprintf(
			"missing selection set for field %q of type %s",
			fieldName, hostTypeName,
		),
	})
}

func (p *Parser) errExpectedNum(actual Expression) {
	td := actual.TypeDesignation()
	p.errors = append(p.errors, Error{
		Location: actual.GetLocation(),
		Msg:      "expected number but received " + td,
	})
}

func (p *Parser) errExpectedInt(actual Expression) {
	td := actual.TypeDesignation()
	p.errors = append(p.errors, Error{
		Location: actual.GetLocation(),
		Msg:      "expected Int but received " + td,
	})
}

func (p *Parser) errUnexpType(
	expected *ast.Type,
	actual Expression,
) {
	td := actual.TypeDesignation()
	p.errors = append(p.errors, Error{
		Location: actual.GetLocation(),
		Msg: fmt.Sprintf(
			"expected type %s but received %s", expected, td,
		),
	})
}

func (p *Parser) errUndefEnumVal(e *Enum) {
	p.errors = append(p.errors, Error{
		Location: e.Location,
		Msg:      fmt.Sprintf("undefined enum value %q", e.Value),
	})
}

func (p *Parser) errUndefFieldInType(
	l Location, fieldName, typeName string,
) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg: fmt.Sprintf(
			"field %q is undefined in type %s",
			fieldName, typeName,
		),
	})
}

func isTypeNum(t *ast.Type) bool {
	return t.Elem == nil && (t.NamedType == "Int" || t.NamedType == "Float")
}

func isTypeFloat(t *ast.Type) bool {
	return t.Elem == nil && t.NamedType == "Float"
}

func isTypeBool(t *ast.Type) bool {
	return t.Elem == nil && t.NamedType == "Boolean"
}

func isTypeString(t *ast.Type) bool {
	return t.Elem == nil && t.NamedType == "String"
}

func defHasLength(t *ast.Type) bool {
	return t.NamedType == "String" || typeIsArray(t)
}

func typeIsArray(t *ast.Type) bool {
	return t.NamedType == "" && t.Elem != nil
}

func indexOf[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

func makeAppend[T any](a []T, b ...T) (new []T) {
	new = make([]T, len(a)+len(b))
	copy(new, a)
	copy(new[len(a):], b)
	return new
}
