package gqt

import (
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
	// Location defines the start location of an expression.
	Location struct{ Index, Line, Column int }

	// LocationEnd defines the end location of an expression.
	LocationEnd struct{ IndexEnd, LineEnd, ColumnEnd int }

	// LocRange defines the start and end locations of an expression.
	LocRange struct {
		Location
		LocationEnd
	}

	// Expression can be either of:
	//
	//   • *Operation
	//   • *ConstrAny
	//   • *ConstrEquals
	//   • *ConstrNotEquals
	//   • *ConstrLess
	//   • *ConstrLessOrEqual
	//   • *ConstrGreater
	//   • *ConstrGreaterOrEqual
	//   • *ConstrLenEquals
	//   • *ConstrLenNotEquals
	//   • *ConstrLenLess
	//   • *ConstrLenLessOrEqual
	//   • *ConstrLenGreater
	//   • *ConstrLenGreaterOrEqual
	//   • *ConstrMap
	//   • *ExprParentheses
	//   • *ExprModulo
	//   • *ExprDivision
	//   • *ExprMultiplication
	//   • *ExprAddition
	//   • *ExprSubtraction
	//   • *ExprLogicalNegation
	//   • *ExprNumericNegation
	//   • *ExprEqual
	//   • *ExprNotEqual
	//   • *ExprLess
	//   • *ExprLessOrEqual
	//   • *ExprGreater
	//   • *ExprGreaterOrEqual
	//   • *ExprLogicalAnd
	//   • *ExprLogicalOr
	//   • *True
	//   • *False
	//   • *Int
	//   • *Float
	//   • *String
	//   • *Null
	//   • *Enum
	//   • *Array
	//   • *Object
	//   • *Variable
	//   • *SelectionInlineFrag
	//   • *ObjectField
	//   • *SelectionField
	//   • *SelectionMax
	//   • *Argument
	Expression interface {
		GetParent() Expression
		GetLocation() LocRange
		TypeDesignation() string
		IsFloat() bool
	}

	// Selection can be either of:
	//
	//	*SelectionField
	//	*SelectionInlineFrag
	//	*SelectionMax
	Selection Expression

	// Operation is the root of the abstract syntax tree of an operation.
	Operation struct {
		LocRange
		Type OperationType
		SelectionSet
		Def *ast.Definition
	}

	Name struct {
		LocRange
		Name string
	}

	// SelectionField is a field selection.
	SelectionField struct {
		LocRange
		Name
		Parent Expression
		ArgumentList
		SelectionSet
		Def *ast.FieldDefinition
	}

	// SelectionSet is a selection set.
	SelectionSet struct {
		LocRange
		Selections []Selection
	}

	// ArgumentList is an argument list.
	ArgumentList struct {
		LocRange
		Arguments []*Argument
	}

	// Argument is an input argument.
	Argument struct {
		LocRange
		Name
		Parent             Expression
		AssociatedVariable *VariableDeclaration
		Constraint         Expression
		Def                *ast.ArgumentDefinition
	}

	// ConstrAny is the any value constraint (*)
	ConstrAny struct {
		LocRange
		Parent Expression
	}

	// ConstrEquals is the equality constraint
	ConstrEquals struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrNotEquals is the inequality constraint
	ConstrNotEquals struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLess is the relational "less than" constraint (<)
	ConstrLess struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLess is the relational "less than or equal" constraint (<=)
	ConstrLessOrEqual struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrGreater is the relational "greater than" constraint (>)
	ConstrGreater struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrGreater is the relational "greater than or equal" constraint (>=)
	ConstrGreaterOrEqual struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLenEquals is the relational length equality constraint
	// for String or array values.
	ConstrLenEquals struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLenNotEquals is the relational length inequality constraint
	// for String or array values.
	ConstrLenNotEquals struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLenLess is the relational "length less than" constraint
	// for String or array values (len <).
	ConstrLenLess struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLenLessOrEqual is the relational "length less than or equal"
	// constraint for String or array values (len <=).
	ConstrLenLessOrEqual struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLenGreater is the relational "greater than"
	// constraint for String or array values (len >).
	ConstrLenGreater struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrLenGreaterOrEqual is the relational "greater than or equal"
	// constraint for String or array values (len >=).
	ConstrLenGreaterOrEqual struct {
		LocRange
		Parent Expression
		Value  Expression
	}

	// ConstrMap maps a constraint to all items of the array.
	ConstrMap struct {
		LocRange
		Parent     Expression
		Constraint Expression
	}

	// ExprParentheses is an expression enclosed by parentheses.
	ExprParentheses struct {
		LocRange
		Parent     Expression
		Expression Expression
	}

	// ExprLogicalNegation is a logical negation expression
	// using the prefix operator "!".
	ExprLogicalNegation struct {
		LocRange
		Parent     Expression
		Expression Expression
	}

	// ExprNumericNegation is a numeric negation expression
	// using the prefix operator "-".
	ExprNumericNegation struct {
		LocRange
		Parent     Expression
		Expression Expression
		Float      bool
	}

	// ExprModulo is a modulo expression using the operator "%".
	ExprModulo struct {
		LocRange
		Parent   Expression
		Dividend Expression
		Divisor  Expression
		Float    bool
	}

	// ExprDivision is an arithmetic division expression
	// using the operator "/".
	ExprDivision struct {
		LocRange
		Parent   Expression
		Dividend Expression
		Divisor  Expression
		Float    bool
	}

	// ExprMultiplication is an arithmetic multiplication expression
	// using the operator "*".
	ExprMultiplication struct {
		LocRange
		Parent        Expression
		Multiplicant  Expression
		Multiplicator Expression
		Float         bool
	}

	// ExprAddition is an arithmetic addition expression
	// using the operator "+".
	ExprAddition struct {
		LocRange
		Parent      Expression
		AddendLeft  Expression
		AddendRight Expression
		Float       bool
	}

	// ExprSubtraction is an arithmetic subtraction expression
	// using the operator "-".
	ExprSubtraction struct {
		LocRange
		Parent     Expression
		Minuend    Expression
		Subtrahend Expression
		Float      bool
	}

	// ExprEqual is a boolean equality expression
	// using the operator "==".
	ExprEqual struct {
		LocRange
		Parent Expression
		Left   Expression
		Right  Expression
	}

	// ExprNotEqual is a boolean inequality expression
	// using the operator "!=".
	ExprNotEqual struct {
		LocRange
		Parent Expression
		Left   Expression
		Right  Expression
	}

	// ExprLess is a boolean relational "less than" expression
	// using the operator "<".
	ExprLess struct {
		LocRange
		Parent Expression
		Left   Expression
		Right  Expression
	}

	// ExprLessOrEqual is a boolean relational "less than or equal"
	// expression using the operator "<=".
	ExprLessOrEqual struct {
		LocRange
		Parent Expression
		Left   Expression
		Right  Expression
	}

	// ExprGreater is a boolean relational "greater than" expression
	// using the operator ">".
	ExprGreater struct {
		LocRange
		Parent Expression
		Left   Expression
		Right  Expression
	}

	// ExprGreaterOrEqual is a boolean relational "greater than or equal"
	// expression using the operator ">=".
	ExprGreaterOrEqual struct {
		LocRange
		Parent Expression
		Left   Expression
		Right  Expression
	}

	// ExprLogicalAnd is a boolean logical AND expression
	// using the operator "&&".
	ExprLogicalAnd struct {
		LocRange
		Parent      Expression
		Expressions []Expression
	}

	// ExprLogicalOr is a boolean logical OR expression
	// using the operator "||".
	ExprLogicalOr struct {
		LocRange
		Parent      Expression
		Expressions []Expression
	}

	// Int is a signed 32-bit integer value constant.
	Int struct {
		LocRange
		Parent  Expression
		Value   int64
		TypeDef *ast.Definition
	}

	// Float is a signed 64-bit floating point value constant.
	Float struct {
		LocRange
		Parent  Expression
		Value   float64
		TypeDef *ast.Definition
	}

	// String is a UTF-8 string value constant.
	String struct {
		LocRange
		Parent  Expression
		Value   string
		TypeDef *ast.Definition
	}

	// True is a boolean value constant.
	True struct {
		LocRange
		Parent  Expression
		TypeDef *ast.Definition
	}

	// False is a boolean value constant.
	False struct {
		LocRange
		Parent  Expression
		TypeDef *ast.Definition
	}

	// Null is a null-value constant.
	Null struct {
		LocRange
		Parent Expression
		Type   *ast.Type
	}

	// Enum is an enumeration value constant.
	Enum struct {
		LocRange
		Parent  Expression
		Value   string
		TypeDef *ast.Definition
	}

	// Array is an array value constant or an array constraint.
	Array struct {
		LocRange
		Parent Expression
		Items  []Expression
		Type   *ast.Type
	}

	// Object is an input object constraint.
	Object struct {
		LocRange
		Parent  Expression
		Fields  []*ObjectField
		TypeDef *ast.Definition
	}

	// ObjectField is an input object field.
	ObjectField struct {
		LocRange
		Name
		Parent             Expression
		AssociatedVariable *VariableDeclaration
		Constraint         Expression
		Def                *ast.FieldDefinition
	}

	// Variable is a named value placeholder.
	Variable struct {
		LocRange
		Name
		Parent      Expression
		Declaration *VariableDeclaration
	}

	// SelectionMax is the max selection set.
	SelectionMax struct {
		LocRange
		Parent  Expression
		Limit   int
		Options SelectionSet
	}

	// SelectionInlineFrag is an inline fragment.
	SelectionInlineFrag struct {
		LocRange
		Parent        Expression
		TypeCondition TypeCondition
		SelectionSet
	}

	// TypeCondition is the type condition inside an inline fragment.
	TypeCondition struct {
		LocRange
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
func (e *ExprNumericNegation) GetParent() Expression { return e.Parent }
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

func (e *Operation) GetLocation() LocRange               { return e.LocRange }
func (e *ExprParentheses) GetLocation() LocRange         { return e.LocRange }
func (e *ConstrAny) GetLocation() LocRange               { return e.LocRange }
func (e *ConstrEquals) GetLocation() LocRange            { return e.LocRange }
func (e *ConstrNotEquals) GetLocation() LocRange         { return e.LocRange }
func (e *ConstrLess) GetLocation() LocRange              { return e.LocRange }
func (e *ConstrLessOrEqual) GetLocation() LocRange       { return e.LocRange }
func (e *ConstrGreater) GetLocation() LocRange           { return e.LocRange }
func (e *ConstrGreaterOrEqual) GetLocation() LocRange    { return e.LocRange }
func (e *ConstrLenEquals) GetLocation() LocRange         { return e.LocRange }
func (e *ConstrLenNotEquals) GetLocation() LocRange      { return e.LocRange }
func (e *ConstrLenLess) GetLocation() LocRange           { return e.LocRange }
func (e *ConstrLenLessOrEqual) GetLocation() LocRange    { return e.LocRange }
func (e *ConstrLenGreater) GetLocation() LocRange        { return e.LocRange }
func (e *ConstrLenGreaterOrEqual) GetLocation() LocRange { return e.LocRange }
func (e *ExprModulo) GetLocation() LocRange              { return e.LocRange }
func (e *ExprDivision) GetLocation() LocRange            { return e.LocRange }
func (e *ExprMultiplication) GetLocation() LocRange      { return e.LocRange }
func (e *ExprAddition) GetLocation() LocRange            { return e.LocRange }
func (e *ExprSubtraction) GetLocation() LocRange         { return e.LocRange }
func (e *ExprLogicalNegation) GetLocation() LocRange     { return e.LocRange }
func (e *ExprNumericNegation) GetLocation() LocRange     { return e.LocRange }
func (e *ExprEqual) GetLocation() LocRange               { return e.LocRange }
func (e *ExprNotEqual) GetLocation() LocRange            { return e.LocRange }
func (e *ExprLess) GetLocation() LocRange                { return e.LocRange }
func (e *ExprLessOrEqual) GetLocation() LocRange         { return e.LocRange }
func (e *ExprGreater) GetLocation() LocRange             { return e.LocRange }
func (e *ExprGreaterOrEqual) GetLocation() LocRange      { return e.LocRange }
func (e *ExprLogicalAnd) GetLocation() LocRange          { return e.LocRange }
func (e *ExprLogicalOr) GetLocation() LocRange           { return e.LocRange }
func (e *True) GetLocation() LocRange                    { return e.LocRange }
func (e *False) GetLocation() LocRange                   { return e.LocRange }
func (e *Int) GetLocation() LocRange                     { return e.LocRange }
func (e *Float) GetLocation() LocRange                   { return e.LocRange }
func (e *String) GetLocation() LocRange                  { return e.LocRange }
func (e *Null) GetLocation() LocRange                    { return e.LocRange }
func (e *Enum) GetLocation() LocRange                    { return e.LocRange }
func (e *Array) GetLocation() LocRange                   { return e.LocRange }
func (e *ConstrMap) GetLocation() LocRange               { return e.LocRange }
func (e *Object) GetLocation() LocRange                  { return e.LocRange }
func (e *Variable) GetLocation() LocRange                { return e.LocRange }
func (e *SelectionInlineFrag) GetLocation() LocRange     { return e.LocRange }
func (e *ObjectField) GetLocation() LocRange             { return e.LocRange }
func (e *SelectionField) GetLocation() LocRange          { return e.LocRange }
func (e *SelectionMax) GetLocation() LocRange            { return e.LocRange }
func (e *Argument) GetLocation() LocRange                { return e.LocRange }

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
func (e *ExprNumericNegation) IsFloat() bool { return e.Expression.IsFloat() }
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

func (e *ExprLogicalNegation) TypeDesignation() string { return "Boolean" }
func (e *ExprNumericNegation) TypeDesignation() string {
	if e.Float {
		return "Float"
	}
	return "Int"
}
func (e *ExprEqual) TypeDesignation() string          { return "Boolean" }
func (e *ExprNotEqual) TypeDesignation() string       { return "Boolean" }
func (e *ExprLess) TypeDesignation() string           { return "Boolean" }
func (e *ExprLessOrEqual) TypeDesignation() string    { return "Boolean" }
func (e *ExprGreater) TypeDesignation() string        { return "Boolean" }
func (e *ExprGreaterOrEqual) TypeDesignation() string { return "Boolean" }
func (e *ExprLogicalAnd) TypeDesignation() string     { return "Boolean" }
func (e *ExprLogicalOr) TypeDesignation() string      { return "Boolean" }

func (e *True) TypeDesignation() string {
	if e.TypeDef != nil {
		return e.TypeDef.Name
	}
	return "Boolean"
}

func (e *False) TypeDesignation() string {
	if e.TypeDef != nil {
		return e.TypeDef.Name
	}
	return "Boolean"
}

func (e *Int) TypeDesignation() string {
	if e.TypeDef != nil && e.TypeDef.Name != "Boolean" &&
		e.TypeDef.Name != "String" &&
		e.TypeDef.Name != "ID" &&
		e.TypeDef.Name != "Float" {
		return e.TypeDef.Name
	}
	return "Int"
}

func (e *Float) TypeDesignation() string {
	if e.TypeDef != nil && e.TypeDef.Name != "Boolean" &&
		e.TypeDef.Name != "String" &&
		e.TypeDef.Name != "ID" &&
		e.TypeDef.Name != "Int" {
		return e.TypeDef.Name
	}
	return "Float"
}

func (e *String) TypeDesignation() string {
	if e.TypeDef != nil && e.TypeDef.Name != "Boolean" &&
		e.TypeDef.Name != "Float" &&
		e.TypeDef.Name != "Int" {
		return e.TypeDef.Name
	}
	return "String"
}

func (e *Null) TypeDesignation() string {
	if e.Type != nil {
		if e.Type.NonNull {
			// Make type designation nullable
			s := e.Type.String()
			return s[:len(s)-1] + "(null)"
		}
		return e.Type.String() + "(null)"
	}
	return "null"
}

func (e *Enum) TypeDesignation() string {
	if e.TypeDef != nil {
		return e.TypeDef.Name
	}
	return "enum"
}

func (e *Array) TypeDesignation() string {
	if e.Type != nil {
		// Determine type designation based on schema
		return e.Type.String()
	}
	if len(e.Items) < 1 {
		return "array"
	}
	firstItem := e.Items[0].TypeDesignation()
	for _, i := range e.Items {
		if d := i.TypeDesignation(); d == firstItem || d == "null" {
			continue
		}
		// Array item type isn't homogeneous, compose designation
		var b strings.Builder
		b.WriteString("array{")
		for x, i := range e.Items {
			b.WriteString(i.TypeDesignation())
			if x+1 < len(e.Items) {
				b.WriteByte(',')
			}
		}
		b.WriteByte('}')
		return b.String()
	}
	// Determine type designation based on
	// constraint expression of the first item
	if firstItem == "null" {
		return "array"
	}
	return "[" + firstItem + "]"
}

func (e *Object) TypeDesignation() string {
	if e.TypeDef != nil {
		return e.TypeDef.Name
	}
	var b strings.Builder
	b.WriteByte('{')
	for i, f := range e.Fields {
		b.WriteString(f.Name.Name)
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
	LocRange

	Name string

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

// GetInfo returns the schema-type and constraint expression of the value
// behind the variable.
func (v *VariableDeclaration) GetInfo() (
	schemaType *ast.Type,
	constr Expression,
) {
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

// Parser is a GQT parser.
type Parser struct {
	schema   *ast.Schema
	enumVal  map[string]*ast.Definition
	varDecls map[string]*VariableDeclaration
	varRefs  []*Variable
	errors   []Error
}

// Source is a GraphQL schema source file.
type Source struct {
	Name    string
	Content string
}

func newParser() *Parser {
	return &Parser{
		enumVal:  make(map[string]*ast.Definition, 0),
		varDecls: make(map[string]*VariableDeclaration),
	}
}

// NewParser reads a GraphQL schema from given sources (if any)
// and returns a new Parser instance that will parse in schema-aware mode,
// according to the specified schema.
// Returns schemaless parser instance if no sources are provided.
func NewParser(schema []Source) (*Parser, error) {
	p := newParser()
	if len(schema) < 1 {
		return p, nil
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

	p.schema = s
	for _, t := range s.Types {
		for _, v := range t.EnumValues {
			p.enumVal[v.Name] = t
		}
	}

	return p, nil
}

// Parse parses the template in schemaless mode
// and returns its abstract syntax tree.
func Parse(src []byte) (
	operation *Operation,
	variables map[string]*VariableDeclaration,
	errors []Error,
) {
	return newParser().Parse(src)
}

// Parse parses the template and returns its abstract syntax tree.
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
	o := &Operation{LocRange: locRange(s.Location)}

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
	if s, o.SelectionSet = p.parseSelectionSet(s); s.stop() {
		return nil, nil, p.errors
	}
	o.LocationEnd = o.SelectionSet.LocationEnd

	s = s.consumeIgnored()
	if !s.isEOF() {
		p.errUnexpTok(s, "expected end of file")
	}

	for _, sel := range o.Selections {
		setParent(sel, o)
	}

	// Link variable declarations and references
	for _, r := range p.varRefs {
		if v, ok := p.varDecls[r.Name.Name]; !ok {
			p.newErr(r.LocRange, "undefined variable")
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
			s.Def = find(s.Name.Name)
			if s.Def != nil {
				if t := p.schema.Types[s.Def.Type.NamedType]; t != nil {
					p.setTypesSelSet(s.SelectionSet, t.Fields)
				}
				for _, a := range s.Arguments {
					a.Def = s.Def.Arguments.ForName(a.Name.Name)
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
		case *ConstrAny, *Variable:
		case *ConstrEquals:
			push(e.Value, exp)
		case *ConstrNotEquals:
			push(e.Value, exp)
		case *ConstrGreater:
			push(e.Value, nil)
		case *ConstrGreaterOrEqual:
			push(e.Value, nil)
		case *ConstrLess:
			push(e.Value, nil)
		case *ConstrLessOrEqual:
			push(e.Value, nil)
		case *ConstrLenEquals:
			push(e.Value, nil)
		case *ConstrLenNotEquals:
			push(e.Value, nil)
		case *ConstrLenGreater:
			push(e.Value, nil)
		case *ConstrLenGreaterOrEqual:
			push(e.Value, nil)
		case *ConstrLenLess:
			push(e.Value, nil)
		case *ConstrLenLessOrEqual:
			push(e.Value, nil)
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
			push(e.Expression, nil)
		case *ExprNumericNegation:
			push(e.Expression, nil)
		case *ExprLogicalOr:
			for _, e := range e.Expressions {
				push(e, exp)
			}
		case *ExprLogicalAnd:
			for _, e := range e.Expressions {
				push(e, exp)
			}
		case *ExprAddition:
			push(e.AddendLeft, nil)
			push(e.AddendRight, nil)
		case *ExprSubtraction:
			push(e.Minuend, nil)
			push(e.Subtrahend, nil)
		case *ExprMultiplication:
			push(e.Multiplicant, nil)
			push(e.Multiplicator, nil)
		case *ExprDivision:
			push(e.Dividend, nil)
			push(e.Divisor, nil)
		case *ExprModulo:
			push(e.Dividend, nil)
			push(e.Divisor, nil)
		case *ExprGreater:
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprGreaterOrEqual:
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprLess:
			push(e.Left, nil)
			push(e.Right, nil)
		case *ExprLessOrEqual:
			push(e.Left, nil)
			push(e.Right, nil)
		case *Int:
			if exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.Scalar &&
					t.Name != "ID" &&
					t.Name != "Int" &&
					t.Name != "Float" &&
					t.Name != "String" &&
					t.Name != "Boolean" {
					e.TypeDef = t
				}
			}
		case *Float:
			if exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.Scalar &&
					t.Name != "ID" &&
					t.Name != "Int" &&
					t.Name != "Float" &&
					t.Name != "String" &&
					t.Name != "Boolean" {
					e.TypeDef = t
				}
			}
		case *True:
			if exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.Scalar &&
					t.Name != "ID" &&
					t.Name != "Int" &&
					t.Name != "Float" &&
					t.Name != "String" &&
					t.Name != "Boolean" {
					e.TypeDef = t
				}
			}
		case *False:
			if exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.Scalar &&
					t.Name != "ID" &&
					t.Name != "Int" &&
					t.Name != "Float" &&
					t.Name != "String" &&
					t.Name != "Boolean" {
					e.TypeDef = t
				}
			}
		case *String:
			if exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.Scalar &&
					t.Name != "Int" &&
					t.Name != "Float" &&
					t.Name != "String" &&
					t.Name != "Boolean" {
					e.TypeDef = t
				}
			}
		case *Null:
			e.Type = exp
		case *Array:
			var expItem *ast.Type
			if exp != nil {
				expItem = exp.Elem
				if exp.Elem != nil {
					e.Type = exp
				} else {
					t := p.schema.Types[exp.NamedType]
					if t.Kind == ast.Scalar &&
						t.Name != "ID" &&
						t.Name != "Int" &&
						t.Name != "Float" &&
						t.Name != "String" &&
						t.Name != "Boolean" {
						e.Type = exp
					}
				}
			}
			for _, i := range e.Items {
				push(i, expItem)
			}
		case *Enum:
			e.TypeDef = p.enumVal[e.Value]
			if e.TypeDef == nil && exp != nil && exp.Elem == nil {
				t := p.schema.Types[exp.NamedType]
				if t.Kind == ast.Scalar &&
					t.Name != "ID" &&
					t.Name != "Int" &&
					t.Name != "Float" &&
					t.Name != "String" &&
					t.Name != "Boolean" {
					e.TypeDef = t
				}
			}
		case *Object:
			if exp != nil && exp.Elem == nil {
				if t := p.schema.Types[exp.NamedType]; t != nil &&
					(t.Kind == ast.InputObject || t.Kind == ast.Scalar &&
						t.Name != "ID" &&
						t.Name != "Int" &&
						t.Name != "Float" &&
						t.Name != "String" &&
						t.Name != "Boolean") {
					e.TypeDef = t
					for _, f := range e.Fields {
						var tp *ast.Type
						f.Def = t.Fields.ForName(f.Name.Name)
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
				p.newErr(o.LocRange, "type Query is undefined")
				return false
			}
			def = p.schema.Query
		}
		return p.validateSelSet(o, def)
	case OperationTypeMutation:
		if p.schema != nil {
			if p.schema.Mutation == nil {
				p.newErr(o.LocRange, "type Mutation is undefined")
				return false
			}
			def = p.schema.Mutation
		}
		return p.validateSelSet(o, def)
	case OperationTypeSubscription:
		if p.schema != nil {
			if p.schema.Subscription == nil {
				p.newErr(o.LocRange, "type Subscription is undefined")
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

	if s.Location.Index != 0 && len(s.Selections) < 1 {
		p.newErr(s.LocRange, "empty selection set")
		return false
	}

	parent, _ := host.(*SelectionField)
	if parent != nil &&
		expect != nil &&
		len(s.Selections) < 1 &&
		len(expect.Fields) > 0 {
		p.errMissingSelSet(locRange(l), parent.Name.Name, hostTypeName)
		return false
	}

	fields := map[string]struct{}{}
	typeConds := map[string]struct{}{}
	maxBlock := false

	for _, s := range s.Selections {
		switch s := s.(type) {
		case *SelectionField:
			if _, decl := fields[s.Name.Name]; decl {
				ok = false
				p.errRedeclField(s)
				continue
			}
			fields[s.Name.Name] = struct{}{}

			var def *ast.FieldDefinition
			if expect != nil {
				def = expect.Fields.ForName(s.Name.Name)
				if def == nil {
					if s.Name.Name == "__typename" {
						continue
					}
					p.errUndefFieldInType(
						s.LocRange, s.Name.Name, expect.Name,
					)
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
					if _, decl := fields[s.Name.Name]; decl {
						ok = false
						p.errRedeclField(s)
						continue
					}
					fields[s.Name.Name] = struct{}{}

					if s.Name.Name == "__typename" {
						ok = false
						p.errTypenameInMax(s)
						continue
					}
					var def *ast.FieldDefinition
					if expect != nil {
						def = expect.Fields.ForName(s.Name.Name)
						if def == nil {
							ok = false
							p.errUndefFieldInType(
								s.LocRange, s.Name.Name, expect.Name,
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
	if f.ArgumentList.Location.Index != 0 && len(f.Arguments) < 1 {
		p.newErr(f.ArgumentList.LocRange, "empty argument list")
		return false
	}

	byName := make(map[string]*Argument, len(f.Arguments))
	for _, a := range f.Arguments {
		if _, found := byName[a.Name.Name]; found {
			p.newErr(a.LocRange, fmt.Sprintf(
				"redeclared argument %q", a.Name.Name,
			))
			ok = false
			continue
		}
		byName[a.Name.Name] = a
	}

	if expect != nil {
		for _, a := range f.Arguments {
			// Check undefined arguments
			if ad := expect.Arguments.ForName(a.Name.Name); ad == nil {
				p.errUndefArg(a, f, host.Name)
				ok = false
			}
		}

		// Check required arguments
		for _, a := range expect.Arguments {
			if a.Type.NonNull && a.DefaultValue == nil {
				if _, found := byName[a.Name]; !found {
					l := f.ArgumentList.LocRange
					if len(f.Arguments) < 1 {
						l = f.LocRange
					}
					p.errMissingArg(l, a)
					ok = false
				}
			}
		}
	}

	// Check constraints
	for _, a := range byName {
		var exp *ast.Type
		if a.Def != nil {
			exp = a.Def.Type
		}

		if !p.validateExpr([]Expression{a}, a.Constraint, exp) {
			ok = false
		}
	}

	var def *ast.Definition
	if p.schema != nil {
		def = p.schema.Types[expect.Type.NamedType]
	}
	if !p.validateSelSet(f, def) {
		ok = false
	}
	return ok
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

			if expect != nil && !p.expectationIsNum(expect) {
				p.errCantApplyConstrRel(e, expect)
				ok = false
				break TYPESWITCH
			} else if _, found := find[*Null](val); found {
				p.errExpectedNumGotNull(val.GetLocation())
				ok = false
				break TYPESWITCH
			} else if !p.isNumeric(val) {
				p.errExpectedNum(val)
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

			if expect != nil &&
				!p.expectationIsString(expect) &&
				!p.expectationIsArray(expect) {
				// Only strings and arrays have a length
				p.errCantApplyConstrLen(e, expect)
				ok = false
				break TYPESWITCH
			} else if _, found := find[*Null](val); found {
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
			var loc LocRange
			switch e := e.(type) {
			case *ExprEqual:
				loc, left, right = e.LocRange, e.Left, e.Right
			case *ExprNotEqual:
				loc, left, right = e.LocRange, e.Left, e.Right
			}
			if !p.assumeComparableValue(left) {
				ok = false
			}
			if !p.assumeComparableValue(right) {
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
			if expect != nil && !p.expectationIsBool(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			} else if _, found := find[*Null](e.Expression); found {
				p.errExpectedBoolGotNull(e.Expression.GetLocation())
				ok = false
				break TYPESWITCH
			} else if !p.isBoolean(e.Expression) {
				p.errExpectedBool(e.Expression)
				ok = false
				break TYPESWITCH
			}
			push(e.Expression, nil)
		case *ExprNumericNegation:
			if expect != nil && !p.expectationIsNum(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			} else if _, found := find[*Null](e.Expression); found {
				p.errExpectedNumGotNull(e.Expression.GetLocation())
				ok = false
				break TYPESWITCH
			} else if !p.isNumeric(e.Expression) {
				p.errExpectedNum(e.Expression)
				ok = false
				break TYPESWITCH
			}
			push(e.Expression, mustMakeExpectNumNonNull(expect))
		case *ExprLogicalOr:
			objEncountered := false
			for _, e := range e.Expressions {
				o := getConstrEqValue[*Object](e)
				if o != nil && objEncountered {
					ok = false
					p.newErr(o.LocRange, "use single object "+
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
					p.newErr(o.LocRange, "use single object "+
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

			if expect != nil && !p.expectationIsNum(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			}

			if !p.isNumeric(left) {
				p.errExpectedNum(left)
				ok = false
			}
			if !p.isNumeric(right) {
				p.errExpectedNum(right)
				ok = false
			}
			if !ok {
				break TYPESWITCH
			}

			if _, found := find[*Null](left); found {
				p.errExpectedNumGotNull(left.GetLocation())
				ok = false
			}
			if _, found := find[*Null](right); found {
				p.errExpectedNumGotNull(right.GetLocation())
				ok = false
			}
			if !ok {
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

			if expect != nil && !p.expectationIsBool(expect) {
				p.errUnexpType(expect, e)
				ok = false
				break TYPESWITCH
			}

			if !p.isNumeric(left) {
				p.errExpectedNum(left)
				ok = false
			}
			if !p.isNumeric(right) {
				p.errExpectedNum(right)
				ok = false
			}
			if !ok {
				break TYPESWITCH
			}

			push(right, nil)
			push(left, nil)
		case *Int:
			if expect != nil && !p.expectationIsNum(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			}
		case *Float:
			if expect != nil && !p.expectationIsFloat(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			}
		case *True:
			if expect != nil && !p.expectationIsBool(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			}
		case *False:
			if expect != nil && !p.expectationIsBool(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			}
		case *Null:
			if expect != nil && expect.NonNull {
				ok = false
				p.errUnexpType(expect, e)
				break TYPESWITCH
			}
		case *Array:
			if expect != nil && !p.expectationIsArray(expect) {
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
				if expect.NonNull {
					if f, found := find[*Null](constr); found {
						p.errUnexpNull(e.LocRange, expect, f)
						return false
					}
				}
				if !areTypesCompatible(expect, tp) {
					p.errUnexpType(expect, e)
				}
			}
		case *Enum:
			if expect != nil && !p.expectationIsEnum(expect) {
				ok = false
				p.errUnexpType(expect, top.Expr)
				break TYPESWITCH
			}
			if p.schema != nil && e.TypeDef == nil {
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
			if expect != nil && !p.expectationIsString(expect) {
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

	if len(o.Fields) < 1 {
		p.newErr(o.LocRange, "empty input object")
		return false
	}

	if exp != nil && !p.expectationIsObject(exp) {
		ok = false
		// Make sure the object doesn't contain any
		// recursive variable references, otherwise
		// printing the object type will result in an endless loop.
		if p.checkObjectVarRefs(o) {
			p.errUnexpType(exp, o)
		}
		return
	}

	fieldNames := make(map[string]struct{}, len(o.Fields))
	validFields := o.Fields
	if p.schema != nil && exp != nil {
		validFields = make([]*ObjectField, 0, len(o.Fields))
		d := p.schema.Types[exp.NamedType]
		if d != nil && d.Kind != ast.Scalar {
			for _, f := range o.Fields {
				// Check undefined fields
				fd := d.Fields.ForName(f.Name.Name)
				if fd == nil {
					ok = false
					p.errUndefField(f, exp.NamedType)
					continue
				}
				fieldNames[f.Name.Name] = struct{}{}
				validFields = append(validFields, f)
			}

			// Check required fields
			for _, f := range d.Fields {
				if f.Type.NonNull && f.DefaultValue == nil {
					if _, exists := fieldNames[f.Name]; !exists {
						ok = false
						p.errMissingInputField(o.LocRange, f)
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

func (p *Parser) errUncompVal(e Expression) {
	p.errors = append(p.errors, Error{
		LocRange: e.GetLocation(),
		Msg:      "uncomparable value of type " + e.TypeDesignation(),
	})
}

func (p *Parser) errConstrSelf(v *Variable) {
	var on, name string
	switch v := v.Declaration.Parent.(type) {
	case *Argument:
		on, name = "argument", v.Name.Name
	case *ObjectField:
		on, name = "object field", v.Name.Name
	}
	p.errors = append(p.errors, Error{
		LocRange: v.LocRange,
		Msg: fmt.Sprintf(
			"illegal self-reference of %s %q through variable %q in constraint",
			on, name, v.Name.Name,
		),
	})
}

func (p *Parser) errUnexpTok(s source, msg string) {
	prefix := "unexpected token, "
	if s.Index >= len(s.s) {
		prefix = "unexpected end of file, "
	}
	p.errors = append(p.errors, Error{
		LocRange: locRange(s.Location),
		Msg:      prefix + msg,
	})
}

func (p *Parser) errTypenameInMax(f *SelectionField) {
	p.errors = append(p.errors, Error{
		LocRange: f.LocRange,
		Msg:      "avoid __typename in max sets",
	})
}

func (p *Parser) errNestedMaxSet(s *SelectionMax) {
	p.errors = append(p.errors, Error{
		LocRange: s.LocRange,
		Msg:      "nested max set",
	})
}

func (p *Parser) errRedeclMax(s *SelectionMax) {
	p.errors = append(p.errors, Error{
		LocRange: s.LocRange,
		Msg:      "redeclared max set",
	})
}

func (p *Parser) errUndefType(l LocRange, name string) {
	p.errors = append(p.errors, Error{
		LocRange: l,
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
		LocRange: f.LocRange,
		Msg: fmt.Sprintf(
			"field %q is undefined in type %s",
			f.Name.Name, hostTypeName,
		),
	})
}

func (p *Parser) errUndefArg(
	a *Argument,
	f *SelectionField,
	hostTypeName string,
) {
	p.errors = append(p.errors, Error{
		LocRange: a.LocRange,
		Msg: fmt.Sprintf(
			"argument %q is undefined on field %q in type %s",
			a.Name.Name, f.Name.Name, hostTypeName,
		),
	})
}

func (p *Parser) errMismatchingTypes(l LocRange, left, right Expression) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: "mismatching types " +
			left.TypeDesignation() +
			" and " +
			right.TypeDesignation(),
	})
}

func (p *Parser) errCompareWithNull(l LocRange, e, null Expression) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: "mismatching types " +
			e.TypeDesignation() +
			" and " +
			null.TypeDesignation(),
	})
}

func (p *Parser) errRedeclField(f *SelectionField) {
	p.newErr(f.LocRange, fmt.Sprintf("redeclared field %q", f.Name.Name))
}

func (p *Parser) errRedeclVar(v *VariableDeclaration) {
	p.newErr(v.LocRange, fmt.Sprintf("redeclared variable %q", v.Name))
}

func (p *Parser) errRedeclTypeCond(f *SelectionInlineFrag) {
	p.newErr(
		f.TypeCondition.LocRange,
		"redeclared condition for type "+f.TypeCondition.TypeName,
	)
}

func (p *Parser) newErr(l LocRange, msg string) {
	p.errors = append(p.errors, Error{
		LocRange: l,
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

func (p *Parser) parseSelectionSet(s source) (source, SelectionSet) {
	selset := SelectionSet{LocRange: locRange(s.Location)}
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
			selset.LocationEnd = locEnd(s)
			break
		}

		s = s.consumeIgnored()

		if _, ok = s.consume("..."); ok {
			var fragInline *SelectionInlineFrag
			if s, fragInline = p.parseInlineFrag(s); s.stop() {
				return stop(), SelectionSet{}
			}
			selset.Selections = append(selset.Selections, fragInline)
			continue
		}

		lBeforeName := s.Location
		sel := &SelectionField{LocRange: locRange(s.Location)}
		if s, name = s.consumeName(); name == nil {
			p.errUnexpTok(s, "expected selection")
			return stop(), SelectionSet{}
		}
		sel.Name = Name{
			LocRange: LocRange{
				Location:    lBeforeName,
				LocationEnd: locEnd(s),
			},
			Name: string(name),
		}
		sel.LocationEnd = locEnd(s)

		s = s.consumeIgnored()

		if sel.Name.Name == "max" {
			lBeforeMaxNum := s.Location
			var maxNum int64
			if s, maxNum, ok = s.consumeUnsignedInt(); ok {
				// Max
				s = s.consumeIgnored()

				if maxNum < 1 {
					p.newErr(
						locRange(lBeforeMaxNum),
						"limit of options must be an unsigned integer greater 0",
					)
					return stop(), SelectionSet{}
				}

				var options SelectionSet
				sBeforeOptionsBlock := s
				if s, options = p.parseSelectionSet(s); s.stop() {
					return stop(), SelectionSet{}
				}

				if len(options.Selections) < 2 {
					p.newErr(
						locRange(sBeforeOptionsBlock.Location),
						"max set must have at least 2 selection options",
					)
				} else if maxNum > int64(len(options.Selections)-1) {
					p.newErr(
						locRange(lBeforeMaxNum),
						"max limit exceeds number of options-1",
					)
				}

				e := &SelectionMax{
					LocRange: LocRange{
						Location:    lBeforeName,
						LocationEnd: options.LocationEnd,
					},
					Limit:   int(maxNum),
					Options: options,
				}
				selset.Selections = append(selset.Selections, e)
				continue
			}
		}

		if s.peek1('(') {
			if sel.Name.Name == "__typename" {
				p.newErr(locRange(s.Location), errFieldTypenameCantHaveArgs())
				return stop(), SelectionSet{}
			}
			if s, sel.ArgumentList = p.parseArguments(s); s.stop() {
				return stop(), SelectionSet{}
			}
			for _, arg := range sel.Arguments {
				setParent(arg, sel)
			}
			sel.LocationEnd = sel.ArgumentList.LocationEnd
		}

		s = s.consumeIgnored()

		if s.peek1('{') {
			if sel.Name.Name == "__typename" {
				p.newErr(locRange(s.Location), errFieldTypenameCantHaveSels())
				return stop(), SelectionSet{}
			}

			if s, sel.SelectionSet = p.parseSelectionSet(s); s.stop() {
				return stop(), SelectionSet{}
			}
			for _, sub := range sel.Selections {
				setParent(sub, sel)
			}
			sel.LocationEnd = sel.SelectionSet.LocationEnd
		}
		selset.Selections = append(selset.Selections, sel)
	}

	return s, selset
}

func (p *Parser) parseInlineFrag(s source) (source, *SelectionInlineFrag) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("..."); !ok {
		p.errUnexpTok(s, "expected '...'")
		return stop(), nil
	}

	inlineFrag := &SelectionInlineFrag{LocRange: locRange(l)}
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
		LocRange: LocRange{
			Location:    sBeforeName.Location,
			LocationEnd: locEnd(s),
		},
		TypeName: string(name),
	}

	s = s.consumeIgnored()
	var selset SelectionSet
	if s, selset = p.parseSelectionSet(s); s.stop() {
		return stop(), nil
	}
	for _, sel := range selset.Selections {
		setParent(sel, inlineFrag)
	}

	inlineFrag.SelectionSet = selset
	inlineFrag.LocationEnd = selset.LocationEnd
	return s, inlineFrag
}

func (p *Parser) parseArguments(s source) (source, ArgumentList) {
	si := s
	var ok bool
	if s, ok = s.consume("("); !ok {
		p.errUnexpTok(s, "expected opening parenthesis")
		return stop(), ArgumentList{}
	}

	list := ArgumentList{LocRange: locRange(si.Location)}
	for {
		s = s.consumeIgnored()
		var name []byte

		if s.isEOF() {
			p.errUnexpTok(s, "expected argument")
			return stop(), ArgumentList{}
		}

		if s, ok = s.consume(")"); ok {
			list.LocationEnd = locEnd(s)
			break
		}

		s = s.consumeIgnored()

		sBeforeName := s
		arg := &Argument{LocRange: locRange(s.Location)}
		if s, name = s.consumeName(); name == nil {
			p.errUnexpTok(s, "expected argument name")
			return stop(), ArgumentList{}
		}
		arg.Name = Name{
			LocRange: LocRange{
				Location:    sBeforeName.Location,
				LocationEnd: locEnd(s),
			},
			Name: string(name),
		}

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
				LocRange: LocRange{
					Location:    sBeforeDollar.Location,
					LocationEnd: locEnd(s),
				},
				Parent: arg,
				Name:   string(name),
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
		if s, expr = p.parseExprLogicalOr(
			s, expectConstraint,
		); s.stop() {
			return s, ArgumentList{}
		}
		setParent(expr, arg)
		arg.Constraint = expr
		arg.LocationEnd = expr.GetLocation().LocationEnd
		list.Arguments = append(list.Arguments, arg)
		s = s.consumeIgnored()

		if s, ok = s.consume(","); !ok {
			if s, ok = s.consume(")"); !ok {
				p.errUnexpTok(s, "expected comma or end of argument list")
				return stop(), ArgumentList{}
			}
			list.LocationEnd = locEnd(s)
			break
		}
		s = s.consumeIgnored()
	}

	return s, list
}

func (p *Parser) parseValue(
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
		e := &ExprParentheses{LocRange: locRange(l)}

		s = s.consumeIgnored()

		if s, e.Expression = p.parseExprLogicalOr(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Expression, e)

		if s, ok = s.consume(")"); !ok {
			p.errUnexpTok(s, "missing closing parenthesis")
			return stop(), nil
		}
		e.LocationEnd = locEnd(s)

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume("$"); ok {
		lBeforeName := s.Location
		var name []byte
		if s, name = s.consumeName(); name == nil {
			p.errUnexpTok(s, "expected variable name")
			return stop(), nil
		}

		v := &Variable{
			LocRange: LocRange{
				Location:    l,
				LocationEnd: locEnd(s),
			},
			Name: Name{
				LocRange: LocRange{
					Location:    lBeforeName,
					LocationEnd: locEnd(s),
				},
				Name: string(name),
			},
		}
		p.varRefs = append(p.varRefs, v)

		s = s.consumeIgnored()

		return s, v
	} else if s, str, ok = s.consumeString(); ok {
		return s, &String{
			LocRange: LocRange{
				Location:    l,
				LocationEnd: locEnd(s),
			},
			Value: string(str),
		}
	}

	if s, num = p.parseNumber(s); s.stop() {
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
		e := &Array{LocRange: locRange(l)}
		s = s.consumeIgnored()

		for {
			if s, ok = s.consume("]"); ok {
				e.LocationEnd = locEnd(s)
				break
			}

			var expr Expression
			if s, expr = p.parseExprLogicalOr(s, expectConstraint); s.stop() {
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
				e.LocationEnd = locEnd(s)
				break
			}
			s = s.consumeIgnored()
		}

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume("{"); ok {
		// Object
		o := &Object{LocRange: locRange(l)}
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
				o.LocationEnd = locEnd(s)
				break
			}

			s = s.consumeIgnored()

			sBeforeName := s
			fld := &ObjectField{
				LocRange: locRange(s.Location),
				Parent:   o,
			}
			if s, name = s.consumeName(); name == nil {
				p.errUnexpTok(s, "expected object field name")
				return stop(), nil
			}
			fld.Name = Name{
				LocRange: LocRange{
					Location:    sBeforeName.Location,
					LocationEnd: locEnd(s),
				},
				Name: string(name),
			}

			if _, ok := fieldNames[fld.Name.Name]; ok {
				p.newErr(locRange(sBeforeName.Location), "redeclared object field")
				return stop(), nil
			}

			fieldNames[fld.Name.Name] = struct{}{}

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
						locRange(sBeforeDollar.Location),
						"declaration of variables inside "+
							"map constraints is prohibited",
					)
					return stop(), nil
				}

				def := &VariableDeclaration{
					LocRange: LocRange{
						Location:    sBeforeDollar.Location,
						LocationEnd: locEnd(s),
					},
					Parent: fld,
					Name:   string(name),
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
			if s, expr = p.parseExprLogicalOr(s, expectConstraint); s.stop() {
				return stop(), nil
			}
			setParent(expr, fld)
			fld.Constraint = expr
			fld.LocationEnd = expr.GetLocation().LocationEnd
			o.Fields = append(o.Fields, fld)
			s = s.consumeIgnored()

			if s, ok = s.consume(","); !ok {
				if s, ok = s.consume("}"); !ok {
					p.errUnexpTok(s, "expected comma or end of object")
					return stop(), nil
				}
				o.LocationEnd = locEnd(s)
				break
			}
			s = s.consumeIgnored()
		}

		return s, o
	}

	s, str = s.consumeName()
	switch string(str) {
	case "true":
		return s, &True{
			LocRange: LocRange{
				Location:    l,
				LocationEnd: locEnd(s),
			},
		}
	case "false":
		return s, &False{
			LocRange: LocRange{
				Location:    l,
				LocationEnd: locEnd(s),
			},
		}
	case "null":
		return s, &Null{
			LocRange: LocRange{
				Location:    l,
				LocationEnd: locEnd(s),
			},
		}
	case "":
		p.errUnexpTok(s, "invalid value")
		return stop(), nil
	default:
		return s, &Enum{
			LocRange: LocRange{
				Location:    l,
				LocationEnd: locEnd(s),
			},
			Value: string(str),
		}
	}
}

func (p *Parser) parseExprUnary(
	s source,
	expect expect,
) (source, Expression) {
	l := s.Location
	var ok bool
	if s, ok = s.consume("!"); ok {
		e := &ExprLogicalNegation{LocRange: locRange(l)}

		s = s.consumeIgnored()

		if s, e.Expression = p.parseValue(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Expression, e)
		e.LocationEnd = e.Expression.GetLocation().LocationEnd

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume("-"); ok {
		e := &ExprNumericNegation{LocRange: locRange(l)}

		s = s.consumeIgnored()

		if s, e.Expression = p.parseValue(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Expression, e)
		e.LocationEnd = e.Expression.GetLocation().LocationEnd

		s = s.consumeIgnored()
		return s, e
	}
	return p.parseValue(s, expect)
}

func (p *Parser) parseExprMultiplicative(
	s source,
	expect expect,
) (source, Expression) {
	si := s

	var result Expression
	if s, result = p.parseExprUnary(s, expect); s.stop() {
		return stop(), nil
	}

	for {
		s = s.consumeIgnored()
		var selected int
		s, selected = s.consumeEitherOf3("*", "/", "%")
		switch selected {
		case 0:
			e := &ExprMultiplication{
				LocRange:     locRange(si.Location),
				Multiplicant: result,
			}
			setParent(e.Multiplicant, e)

			s = s.consumeIgnored()
			if s, e.Multiplicator = p.parseExprUnary(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Multiplicator, e)

			e.Float = e.Multiplicant.IsFloat() || e.Multiplicator.IsFloat()
			e.LocationEnd = e.Multiplicator.GetLocation().LocationEnd

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprDivision{
				LocRange: locRange(si.Location),
				Dividend: result,
			}
			setParent(e.Dividend, e)

			s = s.consumeIgnored()
			if s, e.Divisor = p.parseExprUnary(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Divisor, e)

			e.Float = e.Dividend.IsFloat() || e.Divisor.IsFloat()
			e.LocationEnd = e.Divisor.GetLocation().LocationEnd

			s = s.consumeIgnored()
			result = e
		case 2:
			e := &ExprModulo{
				LocRange: locRange(si.Location),
				Dividend: result,
			}
			setParent(e.Dividend, e)

			s = s.consumeIgnored()
			if s, e.Divisor = p.parseExprUnary(s, expect); s.stop() {
				return stop(), nil
			}
			setParent(e.Divisor, e)

			e.Float = e.Dividend.IsFloat() || e.Divisor.IsFloat()
			e.LocationEnd = e.Divisor.GetLocation().LocationEnd

			s = s.consumeIgnored()
			result = e
		default:
			s = s.consumeIgnored()
			return s, result
		}
	}
}

func (p *Parser) parseExprAdditive(
	s source,
	expect expect,
) (source, Expression) {
	si := s

	var result Expression
	if s, result = p.parseExprMultiplicative(s, expect); s.stop() {
		return stop(), nil
	}

	for {
		s = s.consumeIgnored()
		var selected int
		s, selected = s.consumeEitherOf3("+", "-", "")
		switch selected {
		case 0:
			e := &ExprAddition{
				LocRange:   locRange(si.Location),
				AddendLeft: result,
			}
			setParent(e.AddendLeft, e)

			s = s.consumeIgnored()
			s, e.AddendRight = p.parseExprMultiplicative(s, expect)
			if s.stop() {
				return stop(), nil
			}
			setParent(e.AddendRight, e)

			e.Float = e.AddendLeft.IsFloat() || e.AddendRight.IsFloat()
			e.LocationEnd = e.AddendRight.GetLocation().LocationEnd

			s = s.consumeIgnored()
			result = e
		case 1:
			e := &ExprSubtraction{
				LocRange: locRange(si.Location),
				Minuend:  result,
			}
			setParent(e.Minuend, e)

			s = s.consumeIgnored()
			s, e.Subtrahend = p.parseExprMultiplicative(s, expect)
			if s.stop() {
				return stop(), nil
			}
			setParent(e.Subtrahend, e)

			e.Float = e.Minuend.IsFloat() || e.Subtrahend.IsFloat()
			e.LocationEnd = e.Subtrahend.GetLocation().LocationEnd

			s = s.consumeIgnored()
			result = e
		default:
			s = s.consumeIgnored()
			return s, result
		}
	}
}

func (p *Parser) parseExprRelational(
	s source,
	expect expect,
) (source, Expression) {
	l := s.Location

	var left Expression
	if s, left = p.parseExprAdditive(s, expect); s.stop() {
		return stop(), nil
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("<="); ok {
		e := &ExprLessOrEqual{
			LocRange: locRange(l),
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.parseExprAdditive(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		e.LocationEnd = e.Right.GetLocation().LocationEnd

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume(">="); ok {
		e := &ExprGreaterOrEqual{
			LocRange: locRange(l),
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.parseExprAdditive(s, expect); s.stop() {
			return s, nil
		}
		setParent(e.Right, e)
		e.LocationEnd = e.Right.GetLocation().LocationEnd

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume("<"); ok {
		e := &ExprLess{
			LocRange: locRange(l),
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.parseExprAdditive(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		e.LocationEnd = e.Right.GetLocation().LocationEnd

		s = s.consumeIgnored()
		return s, e
	} else if s, ok = s.consume(">"); ok {
		e := &ExprGreater{
			LocRange: locRange(l),
			Left:     left,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.parseExprAdditive(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		e.LocationEnd = e.Right.GetLocation().LocationEnd

		s = s.consumeIgnored()
		return s, e
	}
	s = s.consumeIgnored()
	return s, left
}

func (p *Parser) parseExprEquality(
	s source,
	expect expect,
) (source, Expression) {
	si := s

	var exprLeft Expression
	if s, exprLeft = p.parseExprRelational(s, expect); s.stop() {
		return stop(), nil
	}

	s = s.consumeIgnored()

	var ok bool
	if s, ok = s.consume("=="); ok {
		e := &ExprEqual{
			LocRange: locRange(si.Location),
			Left:     exprLeft,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.parseExprRelational(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		e.LocationEnd = e.Right.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	}
	if s, ok = s.consume("!="); ok {
		e := &ExprNotEqual{
			LocRange: locRange(si.Location),
			Left:     exprLeft,
		}
		setParent(e.Left, e)

		s = s.consumeIgnored()

		if s, e.Right = p.parseExprRelational(s, expect); s.stop() {
			return stop(), nil
		}
		setParent(e.Right, e)
		e.LocationEnd = e.Right.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	}

	s = s.consumeIgnored()
	return s, exprLeft
}

func (p *Parser) parseExprLogicalOr(
	s source,
	expect expect,
) (source, Expression) {
	e := &ExprLogicalOr{LocRange: locRange(s.Location)}

	for {
		var expr Expression
		if s, expr = p.parseExprLogicalAnd(
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
			e.LocationEnd = e.Expressions[len(e.Expressions)-1].
				GetLocation().LocationEnd
			return s, e
		}

		s = s.consumeIgnored()
	}
}

func (p *Parser) parseExprLogicalAnd(
	s source,
	expect expect,
) (source, Expression) {
	e := &ExprLogicalAnd{LocRange: locRange(s.Location)}

	for {
		var expr Expression
		if s, expr = p.parseConstr(s, expect); s.stop() {
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
			e.LocationEnd = e.Expressions[len(e.Expressions)-1].
				GetLocation().LocationEnd
			return s, e
		}

		s = s.consumeIgnored()
	}
}

func (p *Parser) parseConstr(
	s source,
	expect expect,
) (source, Expression) {
	si := s
	var ok bool
	var expr Expression

	if expect == expectValue {
		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}

		s = s.consumeIgnored()

		return s, expr
	}

	if s, ok = s.consume("*"); ok {
		e := &ConstrAny{
			LocRange: LocRange{
				Location:    si.Location,
				LocationEnd: locEnd(s),
			},
		}
		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("!="); ok {
		e := &ConstrNotEquals{
			LocRange: locRange(si.Location),
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr
		e.LocationEnd = expr.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("<="); ok {
		e := &ConstrLessOrEqual{
			LocRange: locRange(si.Location),
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr
		e.LocationEnd = expr.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume(">="); ok {
		e := &ConstrGreaterOrEqual{
			LocRange: locRange(si.Location),
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr
		e.LocationEnd = expr.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("<"); ok {
		e := &ConstrLess{
			LocRange: locRange(si.Location),
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr
		e.LocationEnd = expr.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume(">"); ok {
		e := &ConstrGreater{
			LocRange: locRange(si.Location),
			Value:    expr,
		}
		s = s.consumeIgnored()

		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr
		e.LocationEnd = expr.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	} else if s, ok = s.consume("len"); ok {
		s = s.consumeIgnored()

		if s, ok = s.consume("!="); ok {
			e := &ConstrLenNotEquals{
				LocRange: locRange(si.Location),
				Value:    expr,
			}
			s = s.consumeIgnored()

			var expr Expression

			if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr
			e.LocationEnd = expr.GetLocation().LocationEnd

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume("<="); ok {
			e := &ConstrLenLessOrEqual{
				LocRange: locRange(si.Location),
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr
			e.LocationEnd = expr.GetLocation().LocationEnd

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume(">="); ok {
			e := &ConstrLenGreaterOrEqual{
				LocRange: locRange(si.Location),
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr
			e.LocationEnd = expr.GetLocation().LocationEnd

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume("<"); ok {
			e := &ConstrLenLess{
				LocRange: locRange(si.Location),
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr
			e.LocationEnd = expr.GetLocation().LocationEnd

			s = s.consumeIgnored()

			return s, e
		} else if s, ok = s.consume(">"); ok {
			e := &ConstrLenGreater{
				LocRange: locRange(si.Location),
				Value:    expr,
			}
			s = s.consumeIgnored()

			if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
				return stop(), nil
			}
			setParent(expr, e)
			e.Value = expr
			e.LocationEnd = expr.GetLocation().LocationEnd

			s = s.consumeIgnored()

			return s, e
		}

		e := &ConstrLenEquals{
			LocRange: locRange(si.Location),
			Value:    expr,
		}

		if s, expr = p.parseExprEquality(s, expectValue); s.stop() {
			return stop(), nil
		}
		setParent(expr, e)
		e.Value = expr
		e.LocationEnd = expr.GetLocation().LocationEnd

		s = s.consumeIgnored()

		return s, e
	}

	if s, ok = s.consume("["); ok {
		s = s.consumeIgnored()
		if s, ok = s.consume("..."); ok {
			e := &ConstrMap{LocRange: locRange(si.Location)}
			s = s.consumeIgnored()

			if s.isEOF() {
				p.errUnexpTok(s, "expected map constraint")
				return stop(), nil
			}

			var expr Expression
			if s, expr = p.parseExprLogicalOr(
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
			e.LocationEnd = locEnd(s)

			return s, e
		} else {
			s = si
		}
	}

	e := &ConstrEquals{
		LocRange: locRange(si.Location),
		Value:    expr,
	}

	if expect == expectConstraintInMap {
		expect = expectValueInMap
	} else {
		expect = expectValue
	}

	if s, expr = p.parseExprEquality(s, expect); s.stop() {
		return stop(), nil
	}
	setParent(expr, e)
	e.Value = expr
	e.LocationEnd = expr.GetLocation().LocationEnd

	s = s.consumeIgnored()

	return s, e
}

func (p *Parser) parseNumber(s source) (_ source, number Expression) {
	si := s
	if s.Index >= len(s.s) {
		return si, nil
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
				p.newErr(LocRange{
					Location: si.Location,
					LocationEnd: LocationEnd{
						IndexEnd:  s.Index,
						LineEnd:   s.Line,
						ColumnEnd: s.Column,
					},
				}, "exponent has no digits")
				return stop(), nil
			}
			s, _ = s.consumeEitherOf3("+", "-", "")
			if !s.isDigit() {
				p.newErr(LocRange{
					Location: si.Location,
					LocationEnd: LocationEnd{
						IndexEnd:  s.Index,
						LineEnd:   s.Line,
						ColumnEnd: s.Column,
					},
				}, "exponent has no digits")
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
			LocRange: LocRange{
				Location: si.Location,
				LocationEnd: LocationEnd{
					IndexEnd: s.Index, LineEnd: s.Line, ColumnEnd: s.Column,
				},
			},
			Value: v,
		}
	} else if v, err := strconv.ParseFloat(str, 64); err == nil {
		return s, &Float{
			LocRange: LocRange{
				Location: si.Location,
				LocationEnd: LocationEnd{
					IndexEnd: s.Index, LineEnd: s.Line, ColumnEnd: s.Column,
				},
			},
			Value: v,
		}
	}
	return si, nil
}

type Error struct {
	LocRange
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

func (p *Parser) isNumeric(e Expression) bool {
	switch e := e.(type) {
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
	case *ConstrAny, *Float, *Int,
		*ConstrGreater, *ConstrLess,
		*ConstrGreaterOrEqual, *ConstrLessOrEqual,
		*ExprAddition, *ExprSubtraction,
		*ExprMultiplication, *ExprDivision, *ExprModulo,
		*ExprNumericNegation:
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

func (p *Parser) isBoolean(e Expression) bool {
	switch e := e.(type) {
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

// isAny returns true if e is an any-constraint, otherwise returns false.
// Always returns false for variables in schemaless mode.
func (p *Parser) isAny(e Expression) bool {
	switch e := e.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			return false
		}
		// Check type based on constraint expression
		switch v := e.Declaration.Parent.(type) {
		case *Argument:
			return p.isAny(v.Constraint)
		case *ObjectField:
			return p.isAny(v.Constraint)
		}
	case *ConstrAny:
		return true
	case *ConstrEquals:
		return p.isAny(e.Value)
	case *ConstrNotEquals:
		return p.isAny(e.Value)
	case *ExprParentheses:
		return p.isAny(e.Expression)
	}
	return false
}

func (p *Parser) isString(e Expression) bool {
	switch e := e.(type) {
	case *Variable:
		if t, _ := e.Declaration.GetInfo(); t != nil {
			// Check type based on schema
			return t.Elem == nil &&
				(t.NamedType == "String" || t.NamedType == "ID")
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

func (p *Parser) isEnum(e Expression) bool {
	switch e := e.(type) {
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

func (p *Parser) isNull(e Expression) bool {
	switch e := e.(type) {
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

func (p *Parser) isArray(e Expression) bool {
	switch e := e.(type) {
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
	l LocRange, left, right Expression,
) (ok bool) {
	// Check for ineffectual comparison
	{
		varLeft, isLVar := find[*Variable](left)
		varRight, isRVar := find[*Variable](right)
		if isLVar && isRVar && varLeft.Name.Name == varRight.Name.Name {
			p.newErr(l, "ineffectual comparison")
		}
	}

	ok = true
	switch {
	case p.isAny(left):
		// Schemaless mode
		return true
	case p.isString(left):
		if f, found := find[*Null](right); found {
			p.errCompareWithNull(l, left, f)
			return false
		} else if _, found := find[*Object](right); found {
			p.errUncompVal(right)
			return false
		} else if !p.isString(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case p.isBoolean(left):
		if f, found := find[*Null](right); found {
			p.errCompareWithNull(l, left, f)
			return false
		} else if _, found := find[*Object](right); found {
			p.errUncompVal(right)
			return false
		} else if !p.isBoolean(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case p.isNumeric(left):
		if f, found := find[*Null](right); found {
			p.errCompareWithNull(l, left, f)
			return false
		} else if _, found := find[*Object](right); found {
			p.errUncompVal(right)
			return false
		} else if !p.isNumeric(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case p.isEnum(left):
		if f, found := find[*Null](right); found {
			p.errCompareWithNull(l, left, f)
			return false
		} else if _, found := find[*Object](right); found {
			p.errUncompVal(right)
			return false
		} else if !p.isEnum(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case p.isArray(left):
		if f, found := find[*Null](right); found {
			p.errCompareWithNull(l, left, f)
			return false
		} else if _, found := find[*Object](right); found {
			p.errUncompVal(right)
			return false
		} else if !p.isArray(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
		ld, rd := left.TypeDesignation(), right.TypeDesignation()
		if ld != rd && !(ld == "array" || rd == "array" ||
			ld == "*" || rd == "*") {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	case p.isNull(left):
		if _, found := find[*Object](right); found {
			p.errUncompVal(right)
			return false
		} else if !p.isNull(right) {
			p.errMismatchingTypes(l, left, right)
			return false
		}
	}
	return ok
}

// assumeComparableValue return true if expression e doesn't contain
// any non-ConstrEquals constraints recursively. Otherwise,
// returns false and pushes errors onto the error stack.
func (p *Parser) assumeComparableValue(e Expression) (ok bool) {
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
		return p.assumeComparableValue(v.Value)
	case *ExprParentheses:
		return p.assumeComparableValue(v.Expression)
	case *Array:
		for _, i := range v.Items {
			if !p.assumeComparableValue(i) {
				ok = false
			}
		}
	case *Object:
		p.errUncompVal(e)
		ok = false
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
	case *ExprNumericNegation:
		v.Parent = parent
	case *SelectionInlineFrag:
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
			frag.TypeCondition.LocRange,
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
						par.TypeCondition.TypeName,
						frag.TypeCondition.TypeName,
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
		LocRange: s.TypeCondition.LocRange,
		Msg: "fragment can't condition on scalar type " +
			s.TypeCondition.TypeName,
	})
}

func (p *Parser) errCondOnEnumType(s *SelectionInlineFrag) {
	p.errors = append(p.errors, Error{
		LocRange: s.TypeCondition.LocRange,
		Msg: "fragment can't condition on enum type " +
			s.TypeCondition.TypeName,
	})
}

func (p *Parser) errCondOnInputType(s *SelectionInlineFrag) {
	p.errors = append(p.errors, Error{
		LocRange: s.TypeCondition.LocRange,
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
		LocRange: locRange(l),
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
	p.newErr(c.GetLocation(), "relational constraint "+
		designation+" ("+description+") "+
		"only supports type Float and type Int, "+
		"it can't be applied to type "+expect.String())
}

func (p *Parser) errCantApplyConstrLen(c Expression, expect *ast.Type) {
	var designation, description string
	switch c.(type) {
	case *ConstrLenEquals:
		designation, description = "'len'", "length equal"
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
	p.newErr(c.GetLocation(), "length constraint "+
		designation+" ("+description+") "+
		"only supports arrays and type String, "+
		"it can't be applied to type "+expect.String())
}

func (p *Parser) errMissingArg(
	l LocRange, missingArgument *ast.ArgumentDefinition,
) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: fmt.Sprintf(
			"argument %q of type %s is required but missing",
			missingArgument.Name, missingArgument.Type,
		),
	})
}

func (p *Parser) errMissingInputField(
	l LocRange,
	missingField *ast.FieldDefinition,
) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: fmt.Sprintf(
			"field %q of type %q is required but missing",
			missingField.Name, missingField.Type,
		),
	})
}

func (p *Parser) errMissingSelSet(
	l LocRange, fieldName, hostTypeName string,
) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: fmt.Sprintf(
			"missing selection set for field %q of type %s",
			fieldName, hostTypeName,
		),
	})
}

func (p *Parser) errExpectedNum(actual Expression) {
	td := actual.TypeDesignation()
	p.errors = append(p.errors, Error{
		LocRange: actual.GetLocation(),
		Msg:      "expected number but received " + td,
	})
}

func (p *Parser) errExpectedBool(actual Expression) {
	td := actual.TypeDesignation()
	p.errors = append(p.errors, Error{
		LocRange: actual.GetLocation(),
		Msg:      "expected type Boolean but received " + td,
	})
}

func (p *Parser) errExpectedNumGotNull(l LocRange) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg:      "expected number but received null",
	})
}

func (p *Parser) errExpectedBoolGotNull(l LocRange) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg:      "expected type Boolean but received null",
	})
}

func (p *Parser) errUnexpType(
	expected *ast.Type,
	actual Expression,
) {
	p.errors = append(p.errors, Error{
		LocRange: actual.GetLocation(),
		Msg: "expected type " +
			expected.String() +
			" but received " +
			actual.TypeDesignation(),
	})
}

func (p *Parser) errUnexpNull(
	l LocRange,
	expected *ast.Type,
	null Expression,
) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: "expected type " +
			expected.String() +
			" but received " +
			null.TypeDesignation(),
	})
}

func (p *Parser) errUndefEnumVal(e *Enum) {
	p.errors = append(p.errors, Error{
		LocRange: e.LocRange,
		Msg:      fmt.Sprintf("undefined enum value %q", e.Value),
	})
}

func (p *Parser) errUndefFieldInType(
	l LocRange, fieldName, typeName string,
) {
	p.errors = append(p.errors, Error{
		LocRange: l,
		Msg: fmt.Sprintf(
			"field %q is undefined in type %s",
			fieldName, typeName,
		),
	})
}

func (p *Parser) expectationIsNum(t *ast.Type) bool {
	if t.Elem != nil {
		return false
	}
	if p.schema == nil {
		return t.NamedType == "Int" || t.NamedType == "Float"
	}
	x := p.schema.Types[t.NamedType]
	// Is custom scalar type?
	return x != nil && (x.Kind == ast.Scalar &&
		x.Name != "ID" &&
		x.Name != "String" &&
		x.Name != "Boolean")
}

func (p *Parser) expectationIsFloat(t *ast.Type) bool {
	if t.Elem != nil {
		return false
	}
	if p.schema == nil {
		return t.NamedType == "Float"
	}
	x := p.schema.Types[t.NamedType]
	// Is custom scalar type?
	return x != nil && (x.Kind == ast.Scalar &&
		x.Name != "ID" &&
		x.Name != "Int" &&
		x.Name != "String" &&
		x.Name != "Boolean")
}

func (p *Parser) expectationIsBool(t *ast.Type) bool {
	if t.Elem != nil {
		return false
	}
	if p.schema == nil {
		return t.NamedType == "Boolean"
	}
	x := p.schema.Types[t.NamedType]
	// Is custom scalar type?
	return x != nil && (x.Kind == ast.Scalar &&
		x.Name != "ID" &&
		x.Name != "Int" &&
		x.Name != "Float" &&
		x.Name != "String")
}

func (p *Parser) expectationIsString(t *ast.Type) bool {
	if t.Elem != nil {
		return false
	}
	if p.schema == nil {
		return t.NamedType == "String" || t.NamedType == "ID"
	}
	x := p.schema.Types[t.NamedType]
	// Is custom scalar type?
	return x != nil && (x.Kind == ast.Scalar &&
		x.Name != "Int" &&
		x.Name != "Float" &&
		x.Name != "Boolean")
}

func (p *Parser) expectationIsEnum(t *ast.Type) bool {
	if t.Elem != nil {
		return false
	}
	if p.schema == nil {
		return false
	}
	x := p.schema.Types[t.NamedType]
	// Is enum type or custom scalar type?
	return x != nil && (x.Kind == ast.Enum ||
		(x.Kind == ast.Scalar &&
			x.Name != "ID" &&
			x.Name != "Int" &&
			x.Name != "Float" &&
			x.Name != "String" &&
			x.Name != "Boolean"))
}

func (p *Parser) expectationIsObject(t *ast.Type) bool {
	if t.Elem != nil {
		return false
	}
	if p.schema == nil {
		return t.NamedType != "ID" &&
			t.NamedType != "Int" &&
			t.NamedType != "Float" &&
			t.NamedType != "String" &&
			t.NamedType != "Boolean"
	}
	x := p.schema.Types[t.NamedType]
	// Is input type or custom scalar type?
	return x != nil && (x.Kind == ast.InputObject ||
		(x.Kind == ast.Scalar &&
			x.Name != "ID" &&
			x.Name != "Int" &&
			x.Name != "Float" &&
			x.Name != "String" &&
			x.Name != "Boolean"))
}

func (p *Parser) expectationIsArray(t *ast.Type) bool {
	if t.Elem != nil {
		return true
	}
	if p.schema == nil {
		return true
	}
	x := p.schema.Types[t.NamedType]
	// Is custom scalar type?
	return x != nil && (x.Kind == ast.Scalar &&
		x.Name != "ID" &&
		x.Name != "Int" &&
		x.Name != "Float" &&
		x.Name != "String" &&
		x.Name != "Boolean")
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
		case *ExprNumericNegation:
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

// mustMakeExpectNumNonNull returns either typeIntNotNull or
// typeFloatNotNull if t is Int or Float respectively.
func mustMakeExpectNumNonNull(t *ast.Type) *ast.Type {
	if t != nil && t.Elem == nil {
		if t.NamedType == "Int" {
			return typeIntNotNull
		} else if t.NamedType == "Float" {
			return typeFloatNotNull
		}
	}
	return t
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

// find returns (instance, true) if e or any of its subexpressions
// contains an instance of T, otherwise returns (zero, false).
func find[T any](e Expression) (T, bool) {
	var zero T
	switch e := e.(type) {
	case T:
		return e, true
	case *ConstrAny:
		return zero, false
	case *ConstrEquals:
		return find[T](e.Value)
	case *ConstrNotEquals:
		return find[T](e.Value)
	case *ExprParentheses:
		return find[T](e.Expression)
	case *ExprLogicalAnd:
		for _, i := range e.Expressions {
			if t, ok := find[T](i); ok {
				return t, true
			}
		}
	case *ExprLogicalOr:
		for _, i := range e.Expressions {
			if t, ok := find[T](i); ok {
				return t, true
			}
		}
	case *Variable:
		switch p := e.Declaration.Parent.(type) {
		case *Argument:
			return find[T](p.Constraint)
		case *ObjectField:
			return find[T](p.Constraint)
		}
	}
	return zero, false
}

// areTypesCompatible returns true if a and b are compatible
// not taking nullability into account.
func areTypesCompatible(a, b *ast.Type) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

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

func locRange(l Location) LocRange { return LocRange{Location: l} }

func locEnd(s source) LocationEnd {
	return LocationEnd{
		IndexEnd: s.Index, LineEnd: s.Line, ColumnEnd: s.Column,
	}
}
