package gqt

import (
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v3"
)

// WriteYAML writes any AST object as YAML to w.
func WriteYAML(w io.Writer, o *Operation) error {
	d := yaml.NewEncoder(w)
	d.SetIndent(2)
	return d.Encode(o)
}

func (s Location) MarshalYAML() (any, error) {
	return fmt.Sprintf("%d:%d:%d", s.Index, s.Line, s.Column), nil
}

func (c TypeCondition) MarshalYAML() (any, error) {
	var t string
	if c.TypeDef != nil {
		t = c.TypeDef.Name
	}
	return struct {
		Location Location `yaml:"location"`
		TypeName string   `yaml:"typeName"`
		Type     string   `yaml:"type,omitempty"`
	}{
		Location: c.Location,
		TypeName: c.TypeName,
		Type:     t,
	}, nil
}

func (l ArgumentList) MarshalYAML() (any, error) {
	return struct {
		Location  Location    `yaml:"location"`
		Arguments []*Argument `yaml:"arguments,omitempty"`
	}{
		Location:  l.Location,
		Arguments: l.Arguments,
	}, nil
}

func (s SelectionSet) MarshalYAML() (any, error) {
	return struct {
		Location   Location    `yaml:"location"`
		Selections []Selection `yaml:"selections,omitempty"`
	}{
		Location:   s.Location,
		Selections: s.Selections,
	}, nil
}

func (o *Operation) MarshalYAML() (any, error) {
	return struct {
		OperationType string       `yaml:"operationType"`
		Location      Location     `yaml:"location"`
		SelectionSet  SelectionSet `yaml:"selectionSet,omitempty"`
	}{
		OperationType: o.Type.String(),
		Location:      o.Location,
		SelectionSet:  o.SelectionSet,
	}, nil
}

func (c *ConstrAny) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string   `yaml:"constraintType"`
		Location       Location `yaml:"location"`
	}{
		ConstraintType: "any",
		Location:       c.Location,
	}, nil
}

func (c *ConstrEquals) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "equals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrNotEquals) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "notEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLess) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lessThan",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLessOrEqual) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lessThanOrEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrGreater) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "greaterThan",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrGreaterOrEqual) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "greaterThanOrEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenEquals) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lengthEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenNotEquals) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lengthNotEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenLess) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lengthLessThan",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenLessOrEqual) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lengthLessThanOrEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenGreater) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lengthGreaterThan",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenGreaterOrEqual) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Value          Expression `yaml:"value"`
	}{
		ConstraintType: "lengthGreaterThanOrEquals",
		Location:       c.Location,
		Value:          c.Value,
	}, nil
}

func (c *ConstrMap) MarshalYAML() (any, error) {
	return struct {
		ConstraintType string     `yaml:"constraintType"`
		Location       Location   `yaml:"location"`
		Constraint     Expression `yaml:"constraint"`
	}{
		ConstraintType: "map",
		Location:       c.Location,
		Constraint:     c.Constraint,
	}, nil
}

func (e *ExprParentheses) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Expression     Expression `yaml:"expression"`
	}{
		ExpressionType: "parentheses",
		Location:       e.Location,
		Expression:     e.Expression,
	}, nil
}

func (e *ExprModulo) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Float          bool       `yaml:"float"`
		Dividend       Expression `yaml:"dividend"`
		Divisor        Expression `yaml:"divisor"`
	}{
		ExpressionType: "modulo",
		Location:       e.Location,
		Float:          e.Float,
		Dividend:       e.Dividend,
		Divisor:        e.Divisor,
	}, nil
}

func (e *ExprDivision) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Float          bool       `yaml:"float"`
		Dividend       Expression `yaml:"dividend"`
		Divisor        Expression `yaml:"divisor"`
	}{
		ExpressionType: "division",
		Location:       e.Location,
		Float:          e.Float,
		Dividend:       e.Dividend,
		Divisor:        e.Divisor,
	}, nil
}

func (e *ExprMultiplication) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Float          bool       `yaml:"float"`
		Multiplicant   Expression `yaml:"multiplicant"`
		Multiplicator  Expression `yaml:"multiplicator"`
	}{
		ExpressionType: "multiplication",
		Location:       e.Location,
		Float:          e.Float,
		Multiplicant:   e.Multiplicant,
		Multiplicator:  e.Multiplicator,
	}, nil
}

func (e *ExprAddition) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Float          bool       `yaml:"float"`
		AddendLeft     Expression `yaml:"addendLeft"`
		AddendRight    Expression `yaml:"addendRight"`
	}{
		ExpressionType: "addition",
		Location:       e.Location,
		Float:          e.Float,
		AddendLeft:     e.AddendLeft,
		AddendRight:    e.AddendRight,
	}, nil
}

func (e *ExprSubtraction) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Float          bool       `yaml:"float"`
		Minuend        Expression `yaml:"minuend"`
		Subtrahend     Expression `yaml:"subtrahend"`
	}{
		ExpressionType: "subtraction",
		Location:       e.Location,
		Float:          e.Float,
		Minuend:        e.Minuend,
		Subtrahend:     e.Subtrahend,
	}, nil
}

func (e *ExprLogicalNegation) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Expression     Expression `yaml:"expression"`
	}{
		ExpressionType: "logicalNegation",
		Location:       e.Location,
		Expression:     e.Expression,
	}, nil
}

func (e *ExprEqual) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		ExpressionType: "equals",
		Location:       e.Location,
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprNotEqual) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		ExpressionType: "notEquals",
		Location:       e.Location,
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprLess) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		ExpressionType: "lessThan",
		Location:       e.Location,
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprLessOrEqual) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		ExpressionType: "lessThanOrEquals",
		Location:       e.Location,
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprGreater) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		ExpressionType: "greaterThan",
		Location:       e.Location,
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprGreaterOrEqual) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string     `yaml:"expressionType"`
		Location       Location   `yaml:"location"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		ExpressionType: "greaterThanOrEquals",
		Location:       e.Location,
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprLogicalAnd) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string       `yaml:"expressionType"`
		Location       Location     `yaml:"location"`
		Expressions    []Expression `yaml:"expressions"`
	}{
		ExpressionType: "logicalAND",
		Location:       e.Location,
		Expressions:    e.Expressions,
	}, nil
}

func (e *ExprLogicalOr) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string       `yaml:"expressionType"`
		Location       Location     `yaml:"location"`
		Expressions    []Expression `yaml:"expressions"`
	}{
		ExpressionType: "logicalOR",
		Location:       e.Location,
		Expressions:    e.Expressions,
	}, nil
}

func (e *True) MarshalYAML() (any, error) {
	var t string
	if e.TypeDef != nil {
		t = e.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Type           string   `yaml:"type,omitempty"`
	}{
		ExpressionType: "true",
		Location:       e.Location,
		Type:           t,
	}, nil
}

func (e *False) MarshalYAML() (any, error) {
	var t string
	if e.TypeDef != nil {
		t = e.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Type           string   `yaml:"type,omitempty"`
	}{
		ExpressionType: "false",
		Location:       e.Location,
		Type:           t,
	}, nil
}

func (i *Int) MarshalYAML() (any, error) {
	var t string
	if i.TypeDef != nil {
		t = i.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Type           string   `yaml:"type,omitempty"`
		Value          int64    `yaml:"value"`
	}{
		ExpressionType: "int",
		Location:       i.Location,
		Type:           t,
		Value:          i.Value,
	}, nil
}

func (f *Float) MarshalYAML() (any, error) {
	var t string
	if f.TypeDef != nil {
		t = f.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Type           string   `yaml:"type,omitempty"`
		Value          float64  `yaml:"value"`
	}{
		ExpressionType: "float",
		Location:       f.Location,
		Type:           t,
		Value:          f.Value,
	}, nil
}

func (s *String) MarshalYAML() (any, error) {
	var t string
	if s.TypeDef != nil {
		t = s.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Type           string   `yaml:"type,omitempty"`
		Value          string   `yaml:"value"`
	}{
		ExpressionType: "string",
		Location:       s.Location,
		Type:           t,
		Value:          s.Value,
	}, nil
}

func (n *Null) MarshalYAML() (any, error) {
	var t string
	if n.TypeDef != nil {
		t = n.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Type           string   `yaml:"type,omitempty"`
	}{
		ExpressionType: "null",
		Location:       n.Location,
		Type:           t,
	}, nil
}

func (e *Enum) MarshalYAML() (any, error) {
	var t string
	if e.TypeDef != nil {
		t = e.TypeDef.Name
	}
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Value          string   `yaml:"value"`
		Type           string   `yaml:"type,omitempty"`
	}{
		ExpressionType: "enum",
		Location:       e.Location,
		Value:          e.Value,
		Type:           t,
	}, nil
}

func (a *Array) MarshalYAML() (any, error) {
	var t string
	if a.ItemTypeDef != nil {
		t = a.ItemTypeDef.Name
	}
	return struct {
		ExpressionType string       `yaml:"expressionType"`
		Location       Location     `yaml:"location"`
		ItemType       string       `yaml:"itemType,omitempty"`
		Items          []Expression `yaml:"items"`
	}{
		ExpressionType: "array",
		Location:       a.Location,
		ItemType:       t,
		Items:          a.Items,
	}, nil
}

func (o *Object) MarshalYAML() (any, error) {
	var t string
	if o.TypeDef != nil {
		t = o.TypeDef.Name
	}
	return struct {
		ExpressionType string         `yaml:"expressionType"`
		Location       Location       `yaml:"location"`
		Type           string         `yaml:"type,omitempty"`
		Fields         []*ObjectField `yaml:"fields"`
	}{
		ExpressionType: "object",
		Location:       o.Location,
		Type:           t,
		Fields:         o.Fields,
	}, nil
}

func (r *Variable) MarshalYAML() (any, error) {
	return struct {
		ExpressionType string   `yaml:"expressionType"`
		Location       Location `yaml:"location"`
		Name           string   `yaml:"name"`
	}{
		ExpressionType: "variableReference",
		Location:       r.Location,
		Name:           r.Name,
	}, nil
}

func (s *SelectionInlineFrag) MarshalYAML() (any, error) {
	return struct {
		SelectionType string        `yaml:"selectionType"`
		TypeCondition TypeCondition `yaml:"typeCondition"`
		Location      Location      `yaml:"location"`
		SelectionsSet SelectionSet  `yaml:"selectionSet,omitempty"`
	}{
		SelectionType: "inlineFragment",
		TypeCondition: s.TypeCondition,
		Location:      s.Location,
		SelectionsSet: s.SelectionSet,
	}, nil
}

func (f *ObjectField) MarshalYAML() (any, error) {
	var t string
	if f.Def != nil {
		t = f.Def.Type.String()
	}
	type Var struct {
		Location Location `yaml:"location"`
		Name     string   `yaml:"name"`
	}
	var v *Var
	if f.AssociatedVariable != nil {
		v = &Var{
			Location: f.AssociatedVariable.Location,
			Name:     f.AssociatedVariable.Name,
		}
	}
	return struct {
		Name       string     `yaml:"name"`
		Location   Location   `yaml:"location"`
		Variable   *Var       `yaml:"variable,omitempty"`
		Type       string     `yaml:"type,omitempty"`
		Constraint Expression `yaml:"constraint"`
	}{
		Name:       f.Name,
		Location:   f.Location,
		Variable:   v,
		Type:       t,
		Constraint: f.Constraint,
	}, nil
}

func (s *SelectionField) MarshalYAML() (any, error) {
	var t string
	if s.Def != nil {
		t = s.Def.Type.String()
	}
	return struct {
		SelectionType string       `yaml:"selectionType"`
		Name          string       `yaml:"name"`
		Location      Location     `yaml:"location"`
		Type          string       `yaml:"type,omitempty"`
		ArgumentList  ArgumentList `yaml:"argumentList,omitempty"`
		SelectionSet  SelectionSet `yaml:"selectionSet,omitempty"`
	}{
		SelectionType: "field",
		Name:          s.Name,
		Location:      s.Location,
		Type:          t,
		ArgumentList:  s.ArgumentList,
		SelectionSet:  s.SelectionSet,
	}, nil
}

func (e *SelectionMax) MarshalYAML() (any, error) {
	return struct {
		SelectionType string       `yaml:"selectionType"`
		Location      Location     `yaml:"location"`
		Limit         int          `yaml:"limit"`
		Options       SelectionSet `yaml:"options"`
	}{
		SelectionType: "max",
		Limit:         e.Limit,
		Location:      e.Location,
		Options:       e.Options,
	}, nil
}

func (a *Argument) MarshalYAML() (any, error) {
	var t string
	if a.Def != nil {
		t = a.Def.Type.String()
	}
	type Var struct {
		Location Location `yaml:"location"`
		Name     string   `yaml:"name"`
	}
	var v *Var
	if a.AssociatedVariable != nil {
		v = &Var{
			Location: a.AssociatedVariable.Location,
			Name:     a.AssociatedVariable.Name,
		}
	}
	return struct {
		Name       string     `yaml:"name"`
		Location   Location   `yaml:"location"`
		Variable   *Var       `yaml:"variable,omitempty"`
		Type       string     `yaml:"type,omitempty"`
		Constraint Expression `yaml:"constraint"`
	}{
		Name:       a.Name,
		Location:   a.Location,
		Variable:   v,
		Type:       t,
		Constraint: a.Constraint,
	}, nil
}
