package gqt

import (
	"bytes"
	"fmt"
	"sort"
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
	//  *Variable
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

	Variable struct {
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

func (e *True) GetParent() Expression     { return e.Parent }
func (e *False) GetParent() Expression    { return e.Parent }
func (e *Int) GetParent() Expression      { return e.Parent }
func (e *Float) GetParent() Expression    { return e.Parent }
func (e *String) GetParent() Expression   { return e.Parent }
func (e *Null) GetParent() Expression     { return e.Parent }
func (e *Enum) GetParent() Expression     { return e.Parent }
func (e *Array) GetParent() Expression    { return e.Parent }
func (e *Object) GetParent() Expression   { return e.Parent }
func (e *Variable) GetParent() Expression { return e.Parent }

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

func (e *True) IsFloat() bool     { return false }
func (e *False) IsFloat() bool    { return false }
func (e *Int) IsFloat() bool      { return false }
func (e *Float) IsFloat() bool    { return true }
func (e *String) IsFloat() bool   { return false }
func (e *Null) IsFloat() bool     { return false }
func (e *Enum) IsFloat() bool     { return false }
func (e *Array) IsFloat() bool    { return false }
func (e *Object) IsFloat() bool   { return false }
func (e *Variable) IsFloat() bool { return false }

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
	return typeDesignationRelational(e)
}
func (e *ConstrLessOrEqual) TypeDesignation() string {
	return typeDesignationRelational(e)
}
func (e *ConstrGreater) TypeDesignation() string {
	return typeDesignationRelational(e)
}
func (e *ConstrGreaterOrEqual) TypeDesignation() string {
	return typeDesignationRelational(e)
}
func (e *ConstrLenEquals) TypeDesignation() string {
	return typeDesignationLen(e)
}
func (e *ConstrLenNotEquals) TypeDesignation() string {
	return typeDesignationLen(e)
}
func (e *ConstrLenLess) TypeDesignation() string {
	return typeDesignationLen(e)
}
func (e *ConstrLenLessOrEqual) TypeDesignation() string {
	return typeDesignationLen(e)
}
func (e *ConstrLenGreater) TypeDesignation() string {
	return typeDesignationLen(e)
}
func (e *ConstrLenGreaterOrEqual) TypeDesignation() string {
	return typeDesignationLen(e)
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
		// Determine type designation based on schema
		return "[" + e.ItemTypeDef.Name + "]"
	} else if len(e.Items) > 0 {
		// Determine type designation based on
		// constraint expression of the first item
		return "[" + e.Items[0].TypeDesignation() + "]"
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
func (e *Variable) TypeDesignation() string {
	switch p := e.Declaration.Parent.(type) {
	case *Argument:
		if p.Def != nil {
			return p.Def.Type.String()
		}
		return p.Constraint.TypeDesignation()
	case *ObjectField:
		if p.Def != nil {
			return p.Def.Type.String()
		}
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
	References []Expression

	// Parent can be any of:
	//
	//  *Argument
	//  *ObjectField
	Parent Expression
}

func (v *VariableDeclaration) GetInfo() (schemaType *ast.Type, constr Expression) {
	switch p := v.Parent.(type) {
	case *Argument:
		constr = p.Constraint
		if p.Def != nil {
			schemaType = p.Def.Type
		}
	case *ObjectField:
		constr = p.Constraint
		if p.Def != nil {
			schemaType = p.Def.Type
		}
	}
	return schemaType, constr
}

type Parser struct {
	schema   *ast.Schema
	enumVal  map[string]*ast.Definition
	varDecls map[string]*VariableDeclaration
	varRefs  []*Variable
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
	p.varRefs = make([]*Variable, 0)

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
		sort.Slice(p.errors, func(i, j int) bool {
			return p.errors[i].Index < p.errors[j].Index
		})
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
		case *ConstrAny, *Int, *Float, *True, *False, *String, *Variable:
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
				e.ItemTypeDef = p.schema.Types[exp.Elem.NamedType]
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

	fields := map[string]struct{}{}
	typeConds := map[string]struct{}{}
	maxBlock := false

	for _, s := range s.Selections {
		switch s := s.(type) {
		case *SelectionField:
			if _, decl := fields[s.Name]; decl {
				ok = false
				p.errRedeclField(s)
				continue
			}
			fields[s.Name] = struct{}{}

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
			if _, decl := typeConds[s.TypeCondition.TypeName]; decl {
				ok = false
				p.errRedeclTypeCond(s)
				continue
			}
			typeConds[s.TypeCondition.TypeName] = struct{}{}

			p.validateInlineFrag(s, expect)
		case *SelectionMax:
			if maxBlock {
				p.errRedeclMax(s)
				continue
			}

			maxBlock = true
			for _, s := range s.Options.Selections {
				switch s := s.(type) {
				case *SelectionField:
					if _, decl := fields[s.Name]; decl {
						ok = false
						p.errRedeclField(s)
						continue
					}
					fields[s.Name] = struct{}{}

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
					if _, decl := typeConds[s.TypeCondition.TypeName]; decl {
						ok = false
						p.errRedeclTypeCond(s)
						continue
					}
					typeConds[s.TypeCondition.TypeName] = struct{}{}

					if !p.validateInlineFrag(s, expect) {
						ok = false
						continue
					}
				case *SelectionMax:
					p.errNestedMaxSet(s)
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

		p.validateExpr([]Expression{a}, a.Constraint, exp)
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
	pathToOriginArg []Expression,
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
		expect := top.Expect

	TYPESWITCH:
		switch e := top.Expr.(type) {
		case *ConstrAny:
		case *ConstrEquals:
			push(e.Value, expect)
		case *ConstrNotEquals:
			push(e.Value, expect)
		case *ConstrGreater, *ConstrLess,
			*ConstrGreaterOrEqual, *ConstrLessOrEqual:
			var val Expression
			switch e := e.(type) {
			case *ConstrGreater:
				val = e.Value
			case *ConstrLess:
				val = e.Value
			case *ConstrGreaterOrEqual:
				val = e.Value
			case *ConstrLessOrEqual:
				val = e.Value
			}
			if expect == nil {
				if !p.isNumeric(val) {
					ok = false
					p.errExpectedNum(val)
					break TYPESWITCH
				}
				push(val, expect)
				break TYPESWITCH
			} else if !expectationIsNum(expect) {
				p.errCantApplyConstrRel(
					val,
					mustMakeExpectNumNonNull(expect),
				)
				ok = false
				break TYPESWITCH
			}
			if containsNullConstant(val) {
				p.errUnexpectedNull(
					val.GetLocation(),
					mustMakeExpectNumNonNull(expect),
				)
				ok = false
				break TYPESWITCH
			} else if !p.isNumeric(val) {
				p.errUnexpType(
					mustMakeExpectNumNonNull(expect),
					val,
				)
				ok = false
				break TYPESWITCH
			}
			push(val, nil)
		case *ConstrLenEquals, *ConstrLenNotEquals,
			*ConstrLenGreater, *ConstrLenLess,
			*ConstrLenGreaterOrEqual, *ConstrLenLessOrEqual:
			var val Expression
			switch e := e.(type) {
			case *ConstrLenEquals:
				val = e.Value
			case *ConstrLenNotEquals:
				val = e.Value
			case *ConstrLenGreater:
				val = e.Value
			case *ConstrLenLess:
				val = e.Value
			case *ConstrLenGreaterOrEqual:
				val = e.Value
			case *ConstrLenLessOrEqual:
				val = e.Value
			}
			if expect != nil && !expectationHasLength(expect) {
				p.errCantApplyConstrLen(e, mustMakeExpectNumNonNull(expect))
				ok = false
				break TYPESWITCH
			}
			if containsNullConstant(val) {
				p.errExpectedNumGotNull(val.GetLocation())
				ok = false
				break TYPESWITCH
			} else if !p.isNumeric(val) {
				p.errExpectedNum(val)
				ok = false
				break TYPESWITCH
			}
			push(val, nil)
		case *ConstrMap:
			var exp *ast.Type
			if expect != nil {
				exp = expect.Elem
			}
			push(e.Constraint, exp)
		case *ExprParentheses:
			push(e.Expression, expect)
		case *ExprEqual, *ExprNotEqual:
			var left, right Expression
			var loc Location
			switch e := e.(type) {
			case *ExprEqual:
				loc, left, right = e.Location, e.Left, e.Right
			case *ExprNotEqual:
				loc, left, right = e.Location, e.Left, e.Right
			}
			if !p.assumeValue(left) {
				ok = false
			}
			if !p.assumeValue(right) {
				ok = false
			}
			if !ok {
				break TYPESWITCH
			}
			if !p.assumeComparableValues(loc, left, right) {
				ok = false
				break TYPESWITCH
			}
			push(right, nil)
			push(left, nil)
		case *ExprLogicalNegation:
			if expect != nil && !expectationIsTypeBoolean(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			} else if expect == nil {
				expect = typeBoolean
			}
			push(e.Expression, expect)
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
			}
			for i := len(e.Expressions) - 1; i >= 0; i-- {
				push(e.Expressions[i], expect)
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
			}
			for i := len(e.Expressions) - 1; i >= 0; i-- {
				push(e.Expressions[i], expect)
			}
		case *ExprAddition, *ExprSubtraction,
			*ExprMultiplication, *ExprDivision, *ExprModulo:
			var left, right Expression
			switch e := e.(type) {
			case *ExprAddition:
				left, right = e.AddendLeft, e.AddendRight
			case *ExprSubtraction:
				left, right = e.Minuend, e.Subtrahend
			case *ExprMultiplication:
				left, right = e.Multiplicant, e.Multiplicator
			case *ExprDivision:
				left, right = e.Dividend, e.Divisor
			case *ExprModulo:
				left, right = e.Dividend, e.Divisor
			}

			if expect == nil {
				br := false
				if !p.isNumeric(left) {
					p.errExpectedNum(left)
					br = true
				}
				if !p.isNumeric(right) {
					p.errExpectedNum(right)
					br = true
				}
				if br {
					ok = false
					break TYPESWITCH
				}
			} else if !expectationIsNum(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			}
			if containsNullConstant(left) {
				p.errUnexpectedNull(
					left.GetLocation(),
					mustMakeExpectNumNonNull(expect),
				)
				ok = false
				break TYPESWITCH
			}
			if containsNullConstant(right) {
				p.errUnexpectedNull(
					right.GetLocation(),
					mustMakeExpectNumNonNull(expect),
				)
				ok = false
				break TYPESWITCH
			}
			push(right, mustMakeExpectNumNonNull(expect))
			push(left, mustMakeExpectNumNonNull(expect))
		case *ExprGreater, *ExprLess,
			*ExprGreaterOrEqual, *ExprLessOrEqual:
			var left, right Expression
			switch e := e.(type) {
			case *ExprGreater:
				left, right = e.Left, e.Right
			case *ExprLess:
				left, right = e.Left, e.Right
			case *ExprGreaterOrEqual:
				left, right = e.Left, e.Right
			case *ExprLessOrEqual:
				left, right = e.Left, e.Right
			}

			if expect == nil {
				br := false
				if !p.isNumeric(left) {
					p.errExpectedNum(left)
					br = true
				}
				if !p.isNumeric(right) {
					p.errExpectedNum(right)
					br = true
				}
				if br {
					ok = false
					break TYPESWITCH
				}
			} else if !expectationIsNum(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			}
			push(right, mustMakeExpectNumNonNull(expect))
			push(left, mustMakeExpectNumNonNull(expect))
		case *Int, *Float:
			if expect != nil && !expectationIsNum(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
			}
		case *True:
			if expect != nil && !expectationIsTypeBoolean(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
			}
		case *False:
			if expect != nil && !expectationIsTypeBoolean(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
			}
		case *Null:
			if expect != nil && expect.NonNull {
				ok = false
				p.errUnexpType(expect, top.Expr)
			}
		case *Array:
			if expect != nil && !expectationIsArray(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			}
			for i := len(e.Items) - 1; i >= 0; i-- {
				var exp *ast.Type
				if expect != nil {
					exp = expect.Elem
				}
				push(e.Items[i], exp)
			}
		case *Variable:
			for _, o := range pathToOriginArg {
				if o == e.Declaration.Parent {
					ok = false
					p.errConstrSelf(e)
					break TYPESWITCH
				}
			}

			tp, constr := e.Declaration.GetInfo()
			if expect != nil {
				if tp == nil {
					// Check type based on expression of the variable origin
					switch expect.NamedType {
					case "Int":
						if !p.isInt(constr) {
							p.errUnexpType(expect, e)
							ok = false
							break TYPESWITCH
						}
					case "Float":
						if !p.isFloat(constr) {
							p.errUnexpType(expect, e)
							ok = false
							break TYPESWITCH
						}
					default:
						panic(fmt.Errorf(
							"unhandled type name: %q", expect.NamedType,
						))
					}
				} else if !areTypesCompatible(expect, tp) {
					p.errUnexpType(expect, e)
				}
			}
		case *Enum:
			if p.schema == nil && expect != nil {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			} else if p.schema != nil && e.TypeDef == nil {
				ok = false
				p.errUndefEnumVal(e)
				break TYPESWITCH
			} else if e.TypeDef != nil &&
				expect != nil &&
				expect.NamedType != e.TypeDef.Name {
				ok = false
				p.errUnexpType(expect, top.Expr)
			}
		case *String:
			if expect != nil && !expectationIsTypeString(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
			}
		case *Object:
			p.validateObject(pathToOriginArg, e, expect)
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
	pathToOriginArg []Expression,
	o *Object,
	exp *ast.Type,
) (ok bool) {
	ok = true

	if exp != nil {
		if p.schema != nil {
			def := p.schema.Types[exp.NamedType]
			if def != nil && def.Kind != ast.InputObject {
				ok = false
				// Make sure the object doesn't contain any
				// recursive variable references, otherwise
				// printing the object type will result in an endless loop.
				if p.checkObjectVarRefs(o) {
					p.errUnexpType(exp, o)
				}
				return
			}
		} else if exp.Elem != nil ||
			exp.NamedType == "ID" ||
			exp.NamedType == "String" ||
			exp.NamedType == "Int" ||
			exp.NamedType == "Float" ||
			exp.NamedType == "Boolean" {
			ok = false
			// Make sure the object doesn't contain any
			// recursive variable references, otherwise
			// printing the object type will result in an endless loop.
			if p.checkObjectVarRefs(o) {
				p.errUnexpType(exp, o)
			}
		}
	}

	fieldNames := make(map[string]struct{}, len(o.Fields))
	validFields := o.Fields
	if p.schema != nil && exp != nil {
		validFields = make([]*ObjectField, 0, len(o.Fields))
		if td := p.schema.Types[exp.NamedType]; td != nil {
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
	}

	// Check constraints
	for _, f := range validFields {
		var exp *ast.Type
		if f.Def != nil {
			exp = f.Def.Type
		}
		if !p.validateExpr(
			makeAppend[Expression](pathToOriginArg, f),
			f.Constraint,
			exp,
		) {
			ok = false
		}
	}

	return ok
}

func (p *Parser) errConstrSelf(v *Variable) {
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

func (p *Parser) errNestedMaxSet(s *SelectionMax) {
	p.errors = append(p.errors, Error{
		Location: s.Location,
		Msg:      "nested max set",
	})
}

func (p *Parser) errRedeclMax(s *SelectionMax) {
	p.errors = append(p.errors, Error{
		Location: s.Location,
		Msg:      "redeclared max set",
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

func (p *Parser) errUnexpectedNull(l Location, expected *ast.Type) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg:      "expected type " + expected.String() + " but received null",
	})
}

func (p *Parser) errCompareWithNull(l Location, e Expression) {
	d := e.TypeDesignation()
	p.errors = append(p.errors, Error{
		Location: l,
		Msg:      fmt.Sprintf("mismatching types %s and null", d),
	})
}

func (p *Parser) errRedeclField(f *SelectionField) {
	p.newErr(f.Location, fmt.Sprintf("redeclared field %q", f.Name))
}

func (p *Parser) errRedeclVar(v *VariableDeclaration) {
	p.newErr(v.Location, fmt.Sprintf("redeclared variable %q", v.Name))
}

func (p *Parser) errRedeclTypeCond(f *SelectionInlineFrag) {
	p.newErr(
		f.TypeCondition.Location,
		"redeclared condition for type "+f.TypeCondition.TypeName,
	)
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
			if s, fragInline = p.ParseInlineFrag(s); s.stop() {
				return stop(), SelectionSet{}
			}
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
						"limit of options must be an unsigned integer greater 0",
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
						"max set must have at least 2 selection options",
					)
				} else if maxNum > int64(len(options.Selections)-1) {
					p.newErr(
						lBeforeMaxNum,
						"max limit exceeds number of options-1",
					)
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
				p.errRedeclVar(def)
			} else {
				p.varDecls[def.Name] = def
			}
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
	var num Expression
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
		v := &Variable{Location: l}

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
						"declaration of variables inside "+
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
					p.errRedeclVar(def)
				} else {
					p.varDecls[def.Name] = def
				}

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

func (p *Parser) ParseNumber(s source) (_ source, number Expression) {
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

func (p *Parser) isNumeric(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
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
	case *ConstrAny, *Float, *Int, *ExprAddition, *ExprSubtraction,
		*ExprMultiplication, *ExprDivision, *ExprModulo:
		return true
	case *ConstrEquals:
		return p.isNumeric(e.Value)
	case *ConstrNotEquals:
		return p.isNumeric(e.Value)
	case *ExprParentheses:
		return p.isNumeric(e.Expression)
	}
	return false
}

func (p *Parser) isInt(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
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
	case *ConstrAny, *Int:
		return true
	case *ConstrEquals:
		return p.isInt(e.Value)
	case *ConstrNotEquals:
		return p.isInt(e.Value)
	case *ExprParentheses:
		return p.isInt(e.Expression)
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

func (p *Parser) isFloat(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			// Check type based on schema
			return t.Elem == nil && t.NamedType == "Float"
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isFloat(v.Constraint)
		case *ObjectField:
			return p.isFloat(v.Constraint)
		}
	case *ConstrAny, *Float:
		return true
	case *ExprParentheses:
		return p.isFloat(e.Expression)
	case *ExprAddition:
		return e.Float
	case *ExprSubtraction:
		return e.Float
	case *ExprMultiplication:
		return e.Float
	case *ExprDivision:
		return e.Float
	case *ExprModulo:
		return e.Float
	}
	return false
}

func (p *Parser) isBoolean(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
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
	case *ConstrAny, *True, *False, *ConstrGreater, *ConstrGreaterOrEqual,
		*ConstrLess, *ConstrLessOrEqual, *ExprEqual, *ExprNotEqual,
		*ExprLess, *ExprGreater, *ExprLessOrEqual, *ExprGreaterOrEqual,
		*ExprLogicalNegation, *ExprLogicalAnd, *ExprLogicalOr:
		return true
	case *ConstrEquals:
		return p.isBoolean(e.Value)
	case *ConstrNotEquals:
		return p.isBoolean(e.Value)
	case *ExprParentheses:
		return p.isBoolean(e.Expression)
	}
	return false
}

func (p *Parser) isString(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
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
	case *ConstrAny,
		*String,
		*ConstrLenGreater, *ConstrLenLess,
		*ConstrLenGreaterOrEqual, *ConstrLenLessOrEqual,
		*ConstrLenEquals, *ConstrLenNotEquals:
		return true
	case *ConstrEquals:
		return p.isString(e.Value)
	case *ConstrNotEquals:
		return p.isString(e.Value)
	case *ExprParentheses:
		return p.isString(e.Expression)
	}
	return false
}

func (p *Parser) isEnum(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			// Check type based on schema
			tp := p.schema.Types[t.NamedType]
			return t.Elem == nil && tp.Kind == ast.Enum
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isEnum(v.Constraint)
		case *ObjectField:
			return p.isEnum(v.Constraint)
		}
	case *ConstrAny, *Enum:
		return true
	case *ConstrEquals:
		return p.isEnum(e.Value)
	case *ConstrNotEquals:
		return p.isEnum(e.Value)
	case *ExprParentheses:
		return p.isEnum(e.Expression)
	}
	return false
}

func (p *Parser) isNull(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			// Check type based on schema
			return !t.NonNull
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isNull(v.Constraint)
		case *ObjectField:
			return p.isNull(v.Constraint)
		}
	case *ConstrAny, *Null:
		return true
	case *ConstrEquals:
		return p.isNull(e.Value)
	case *ConstrNotEquals:
		return p.isNull(e.Value)
	case *ExprParentheses:
		return p.isNull(e.Expression)
	}
	return false
}

func (p *Parser) isObject(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			// Check type based on schema
			tp := p.schema.Types[t.NamedType]
			return t.Elem == nil && tp.Kind == ast.InputObject
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isObject(v.Constraint)
		case *ObjectField:
			return p.isObject(v.Constraint)
		}
	case *ConstrAny, *Object:
		return true
	case *ConstrEquals:
		return p.isObject(e.Value)
	case *ConstrNotEquals:
		return p.isObject(e.Value)
	case *ExprParentheses:
		return p.isObject(e.Expression)
	}
	return false
}

func (p *Parser) isArray(expr Expression) bool {
	switch e := expr.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			// Check type based on schema
			return t.Elem != nil
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isArray(v.Constraint)
		case *ObjectField:
			return p.isArray(v.Constraint)
		}
	case *ConstrAny,
		*Array,
		*ConstrLenGreater, *ConstrLenLess,
		*ConstrLenGreaterOrEqual, *ConstrLenLessOrEqual,
		*ConstrLenEquals, *ConstrLenNotEquals:
		return true
	case *ConstrEquals:
		return p.isArray(e.Value)
	case *ConstrNotEquals:
		return p.isArray(e.Value)
	case *ExprParentheses:
		return p.isArray(e.Expression)
	}
	return false
}

func (p *Parser) assumeComparableValues(
	l Location,
	left, right Expression,
) (ok bool) {
	switch {
	case p.isString(left):
		if !p.isString(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	case p.isBoolean(left):
		if !p.isBoolean(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	case p.isNumeric(left):
		if !p.isNumeric(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	// case p.isInt(left):
	// 	if !p.isNumeric(right) {
	// 		p.errMismatchingTypes(l, left, right)
	// 		return false
	// 	}
	// 	if containsNull(right) {
	// 		p.errCompareWithNull(l, left)
	// 		return false
	// 	}
	case p.isEnum(left):
		if !p.isEnum(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	case p.isNull(left):
		if !p.isNull(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	case p.isObject(left):
		if !p.isObject(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	case p.isArray(left):
		if !p.isArray(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		if containsNullConstant(right) {
			p.errCompareWithNull(l, left)
			return false
		}
	}
	return ok
}

// assumeValue return true if expression e doesn't contain
// any non-ConstrEquals constraints recursively. Otherwise,
// returns false and pushes errors onto the error stack.
func (p *Parser) assumeValue(e Expression) (ok bool) {
	ok = true
	switch v := e.(type) {
	case *ConstrAny,
		*ConstrNotEquals,
		*ConstrLess,
		*ConstrLessOrEqual,
		*ConstrGreater,
		*ConstrGreaterOrEqual,
		*ConstrLenEquals,
		*ConstrLenNotEquals,
		*ConstrLenLess,
		*ConstrLenLessOrEqual,
		*ConstrLenGreater,
		*ConstrLenGreaterOrEqual,
		*ConstrMap:
		p.newErr(e.GetLocation(), "unexpected constraint in value definition")
		return false
	case *ConstrEquals:
		return p.assumeValue(v.Value)
	case *ExprParentheses:
		return p.assumeValue(v.Expression)
	case *Array:
		for _, i := range v.Items {
			if !p.assumeValue(i) {
				ok = false
			}
		}
	case *Object:
		for _, f := range v.Fields {
			if !p.assumeValue(f.Constraint) {
				ok = false
			}
		}
	}
	return ok
}

func setParent(t, parent Expression) {
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

func (p *Parser) validateInlineFrag(
	frag *SelectionInlineFrag,
	hostDef *ast.Definition,
) (ok bool) {
	if p.schema == nil {
		if !p.validateCond(frag) {
			return false
		}
		return p.validateSelSet(frag, nil)
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
			return p.validateSelSet(frag, def)
		}
		for _, c := range def.Interfaces {
			if c == hostDef.Name {
				// condition Object implements the host interface
				return p.validateSelSet(frag, def)
			}
		}
		for _, c := range hostDef.Types {
			if c == def.Name {
				// condition Object is part of the host union
				return p.validateSelSet(frag, def)
			}
		}
	case ast.Interface:
		if hostDef.Name == def.Name {
			return p.validateSelSet(frag, def)
		}
		impls := p.schema.GetPossibleTypes(def)
		for _, t := range impls {
			if t == hostDef {
				return p.validateSelSet(frag, def)
			}
			possibleTypes := p.schema.GetPossibleTypes(hostDef)
			for _, ht := range possibleTypes {
				if ht == t {
					return p.validateSelSet(frag, def)
				}
			}
		}
	case ast.Union:
		if hostDef.Name == def.Name {
			return p.validateSelSet(frag, def)
		}
		impls := p.schema.GetPossibleTypes(def)
		for _, t := range impls {
			if t == hostDef {
				return p.validateSelSet(frag, def)
			}
			possibleTypes := p.schema.GetPossibleTypes(hostDef)
			for _, ht := range possibleTypes {
				if ht == t {
					return p.validateSelSet(frag, def)
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

func (p *Parser) errCantApplyConstrRel(c Expression, expect *ast.Type) {
	var designation, description string
	switch c.(type) {
	case *ConstrGreater:
		designation, description = "'>'", "greater than"
	case *ConstrGreaterOrEqual:
		designation, description = "'>='", "greater than or equal"
	case *ConstrLess:
		designation, description = "'<'", "less than"
	case *ConstrLessOrEqual:
		designation, description = "'<='", "less than or equal"
	default:
		panic(fmt.Errorf("unhandled constraint type: %T", c))
	}
	p.newErr(c.GetLocation(), "can't apply relational constraint "+
		designation+" ("+description+") to type "+
		expect.String())
}

func (p *Parser) errCantApplyConstrLen(c Expression, expect *ast.Type) {
	var designation, description string
	switch c.(type) {
	case *ConstrLenEquals:
		designation, description = "'len '", "length equal"
	case *ConstrLenNotEquals:
		designation, description = "'len !='", "length not equal"
	case *ConstrLenGreater:
		designation, description = "'len >'", "length greater than"
	case *ConstrLenGreaterOrEqual:
		designation, description = "'len >='", "length greater than or equal"
	case *ConstrLenLess:
		designation, description = "'len <'", "length less than"
	case *ConstrLenLessOrEqual:
		designation, description = "'len <='", "length less than or equal"
	default:
		panic(fmt.Errorf("unhandled constraint type: %T", c))
	}
	p.newErr(c.GetLocation(), "can't apply length constraint "+
		designation+" ("+description+") to type "+
		expect.String())
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

func (p *Parser) errExpectedNumGotNull(l Location) {
	p.errors = append(p.errors, Error{
		Location: l,
		Msg:      "expected number but received null",
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

func expectationIsNum(t *ast.Type) bool {
	return t.Elem == nil && (t.NamedType == "Int" || t.NamedType == "Float")
}

func expectationIsTypeBoolean(t *ast.Type) bool {
	return t.Elem == nil && t.NamedType == "Boolean"
}

func expectationIsTypeString(t *ast.Type) bool {
	return t.Elem == nil && t.NamedType == "String"
}

func expectationHasLength(t *ast.Type) bool {
	return t.NamedType == "String" || expectationIsArray(t)
}

func expectationIsArray(t *ast.Type) bool {
	return t.NamedType == "" && t.Elem != nil
}

func makeAppend[T any](a []T, b ...T) (new []T) {
	new = make([]T, len(a)+len(b))
	copy(new, a)
	copy(new[len(a):], b)
	return new
}

// checkObjectVarRefs checks for recursive variable references
// and returns false if none are found. Otherwise returns true
// and pushes an error onto the error stack.
func (p *Parser) checkObjectVarRefs(o *Object) (ok bool) {
	ok = true
	stack := make([]Expression, 0, 64)
	push := func(expression Expression) {
		stack = append(stack, expression)
	}

	push(o)
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch e := top.(type) {
		case *ConstrEquals:
			push(e.Value)
		case *ConstrNotEquals:
			push(e.Value)
		case *ConstrGreater:
			push(e.Value)
		case *ConstrGreaterOrEqual:
			push(e.Value)
		case *ConstrLess:
			push(e.Value)
		case *ConstrLessOrEqual:
			push(e.Value)
		case *ConstrLenEquals:
			push(e.Value)
		case *ConstrLenNotEquals:
			push(e.Value)
		case *ConstrLenGreater:
			push(e.Value)
		case *ConstrLenGreaterOrEqual:
			push(e.Value)
		case *ConstrLenLess:
			push(e.Value)
		case *ConstrLenLessOrEqual:
			push(e.Value)
		case *ConstrMap:
			push(e.Constraint)
		case *ExprParentheses:
			push(e.Expression)
		case *ExprEqual:
			push(e.Left)
			push(e.Right)
		case *ExprNotEqual:
			push(e.Left)
			push(e.Right)
		case *ExprLogicalNegation:
			push(e.Expression)
		case *ExprLogicalOr:
			for _, e := range e.Expressions {
				push(e)
			}
		case *ExprLogicalAnd:
			for _, e := range e.Expressions {
				push(e)
			}
		case *ExprAddition:
			push(e.AddendLeft)
			push(e.AddendRight)
		case *ExprSubtraction:
			push(e.Minuend)
			push(e.Subtrahend)
		case *ExprMultiplication:
			push(e.Multiplicant)
			push(e.Multiplicator)
		case *ExprDivision:
			push(e.Dividend)
			push(e.Divisor)
		case *ExprModulo:
			push(e.Dividend)
			push(e.Divisor)
		case *Array:
			for _, i := range e.Items {
				push(i)
			}
		case *Variable:
			for pr := e.Parent; pr != nil; pr = pr.GetParent() {
				switch pv := pr.(type) {
				case *Argument:
					if pv == e.Declaration.Parent {
						p.errConstrSelf(e)
						return false
					}
				case *ObjectField:
					if pv == e.Declaration.Parent {
						p.errConstrSelf(e)
						return false
					}
				}
			}
		case *Object:
			for _, f := range e.Fields {
				push(f)
			}
		case *ObjectField:
			push(e.Constraint)
		case *Int, *Float, *True, *False, *Null,
			*Enum, *String, *ConstrAny:
		default:
			panic(fmt.Errorf("unhandled type: %T", top))
		}
	}
	return ok
}

var typeIntNotNull = &ast.Type{
	NamedType: "Int",
	NonNull:   true,
}
var typeFloatNotNull = &ast.Type{
	NamedType: "Float",
	NonNull:   true,
}
var typeBoolean = &ast.Type{
	NamedType: "Boolean",
}

// mustMakeExpectNumNonNull returns either typeIntNotNull or
// typeFloatNotNull if t is Int or Float respectively.
func mustMakeExpectNumNonNull(t *ast.Type) *ast.Type {
	if t == nil {
		return nil
	}
	if t.Elem != nil || (t.NamedType != "Float" && t.NamedType != "Int") {
		panic(fmt.Errorf("expected numeric type, received: %#v", t))
	}
	if t.NamedType == "Int" {
		return typeIntNotNull
	}
	return typeFloatNotNull
}

func typeDesignationLen(e Expression) string {
	switch p := e.GetParent().(type) {
	case *Argument:
		if p.Def != nil {
			return p.Def.Type.String()
		}
	case *ObjectField:
		if p.Def != nil {
			return p.Def.Type.String()
		}
	}
	return "String|array"
}

func typeDesignationRelational(e Expression) string {
	switch p := e.GetParent().(type) {
	case *Argument:
		if p.Def != nil {
			return p.Def.Type.String()
		}
	case *ObjectField:
		if p.Def != nil {
			return p.Def.Type.String()
		}
	}
	return "number"
}

// containsNullConstant returns true if e contains a constraint expecting
// equality to null.
func containsNullConstant(e Expression) bool {
	switch e := e.(type) {
	case *ConstrAny:
		return false
	case *ConstrEquals:
		return containsNullConstant(e.Value)
	case *ExprParentheses:
		return containsNullConstant(e.Expression)
	case *ExprLogicalAnd:
		for _, i := range e.Expressions {
			if containsNullConstant(i) {
				return true
			}
		}
	case *ExprLogicalOr:
		for _, i := range e.Expressions {
			if containsNullConstant(i) {
				return true
			}
		}
	case *Null:
		return true
	case *Variable:
		switch p := e.Declaration.Parent.(type) {
		case *Argument:
			return containsNullConstant(p.Constraint)
		case *ObjectField:
			return containsNullConstant(p.Constraint)
		}
	}
	return false
}

// areTypesCompatible returns true if a and b are compatible
// not taking nullability into account.
func areTypesCompatible(a, b *ast.Type) bool {
	if a.Elem != nil && b.Elem != nil {
		return areTypesCompatible(a.Elem, b.Elem)
	} else if (a.Elem == nil && b.Elem != nil) ||
		(b.Elem == nil && a.Elem != nil) {
		return false
	} else if a.NamedType == "Float" &&
		(b.NamedType == "Int" || b.NamedType == "Float") {
		return true
	} else if a.NamedType == "Int" &&
		(b.NamedType == "Int" || b.NamedType == "Float") {
		return true
	} else if a.NamedType != b.NamedType {
		return false
	}
	return true
}
