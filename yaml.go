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

func (s LocRange) MarshalYAML() (any, error) {
	if s.ColumnEnd == 0 {
		// Not a range
		return fmt.Sprintf("%d:%d:%d", s.Index, s.Line, s.Column), nil
	}
	return fmt.Sprintf(
		"%d:%d:%d-%d:%d:%d",
		s.Index, s.Line, s.Column,
		s.IndexEnd, s.LineEnd, s.ColumnEnd,
	), nil
}

func (c TypeCondition) MarshalYAML() (any, error) {
	var t string
	if c.TypeDef != nil {
		t = c.TypeDef.Name
	}
	return struct {
		Location LocRange `yaml:"location"`
		TypeName string   `yaml:"typeName"`
		Type     string   `yaml:"type,omitempty"`
	}{
		Location: c.LocRange,
		TypeName: c.TypeName,
		Type:     t,
	}, nil
}

func (l ArgumentList) MarshalYAML() (any, error) {
	return struct {
		Location  LocRange    `yaml:"location"`
		Arguments []*Argument `yaml:"arguments,omitempty"`
	}{
		Location:  l.LocRange,
		Arguments: l.Arguments,
	}, nil
}

func (s SelectionSet) MarshalYAML() (any, error) {
	return struct {
		Location   LocRange    `yaml:"location"`
		Selections []Selection `yaml:"selections,omitempty"`
	}{
		Location:   s.LocRange,
		Selections: s.Selections,
	}, nil
}

func (o *Operation) MarshalYAML() (any, error) {
	return struct {
		Location      LocRange     `yaml:"location"`
		OperationType string       `yaml:"operationType"`
		SelectionSet  SelectionSet `yaml:"selectionSet,omitempty"`
	}{
		Location:      o.LocRange,
		OperationType: o.Type.String(),
		SelectionSet:  o.SelectionSet,
	}, nil
}

func (c *ConstrAny) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange `yaml:"location"`
		ConstraintType string   `yaml:"constraintType"`
	}{
		Location:       c.LocRange,
		ConstraintType: "any",
	}, nil
}

func (c *ConstrEquals) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "equals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrNotEquals) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "notEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLess) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lessThan",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLessOrEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lessThanOrEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrGreater) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "greaterThan",
		Value:          c.Value,
	}, nil
}

func (c *ConstrGreaterOrEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "greaterThanOrEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenEquals) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lengthEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenNotEquals) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lengthNotEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenLess) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lengthLessThan",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenLessOrEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lengthLessThanOrEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenGreater) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lengthGreaterThan",
		Value:          c.Value,
	}, nil
}

func (c *ConstrLenGreaterOrEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Value          Expression `yaml:"value"`
	}{
		Location:       c.LocRange,
		ConstraintType: "lengthGreaterThanOrEquals",
		Value:          c.Value,
	}, nil
}

func (c *ConstrMap) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ConstraintType string     `yaml:"constraintType"`
		Constraint     Expression `yaml:"constraint"`
	}{
		Location:       c.LocRange,
		ConstraintType: "map",
		Constraint:     c.Constraint,
	}, nil
}

func (e *ExprParentheses) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Expression     Expression `yaml:"expression"`
	}{
		Location:       e.LocRange,
		ExpressionType: "parentheses",
		Expression:     e.Expression,
	}, nil
}

func (e *ExprModulo) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Float          bool       `yaml:"float"`
		Dividend       Expression `yaml:"dividend"`
		Divisor        Expression `yaml:"divisor"`
	}{
		Location:       e.LocRange,
		ExpressionType: "modulo",
		Float:          e.Float,
		Dividend:       e.Dividend,
		Divisor:        e.Divisor,
	}, nil
}

func (e *ExprDivision) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Float          bool       `yaml:"float"`
		Dividend       Expression `yaml:"dividend"`
		Divisor        Expression `yaml:"divisor"`
	}{
		Location:       e.LocRange,
		ExpressionType: "division",
		Float:          e.Float,
		Dividend:       e.Dividend,
		Divisor:        e.Divisor,
	}, nil
}

func (e *ExprMultiplication) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Float          bool       `yaml:"float"`
		Multiplicant   Expression `yaml:"multiplicant"`
		Multiplicator  Expression `yaml:"multiplicator"`
	}{
		Location:       e.LocRange,
		ExpressionType: "multiplication",
		Float:          e.Float,
		Multiplicant:   e.Multiplicant,
		Multiplicator:  e.Multiplicator,
	}, nil
}

func (e *ExprAddition) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Float          bool       `yaml:"float"`
		AddendLeft     Expression `yaml:"addendLeft"`
		AddendRight    Expression `yaml:"addendRight"`
	}{
		Location:       e.LocRange,
		ExpressionType: "addition",
		Float:          e.Float,
		AddendLeft:     e.AddendLeft,
		AddendRight:    e.AddendRight,
	}, nil
}

func (e *ExprSubtraction) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Float          bool       `yaml:"float"`
		Minuend        Expression `yaml:"minuend"`
		Subtrahend     Expression `yaml:"subtrahend"`
	}{
		Location:       e.LocRange,
		ExpressionType: "subtraction",
		Float:          e.Float,
		Minuend:        e.Minuend,
		Subtrahend:     e.Subtrahend,
	}, nil
}

func (e *ExprLogicalNegation) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Expression     Expression `yaml:"expression"`
	}{
		Location:       e.LocRange,
		ExpressionType: "logicalNegation",
		Expression:     e.Expression,
	}, nil
}

func (e *ExprEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		Location:       e.LocRange,
		ExpressionType: "equals",
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprNotEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		Location:       e.LocRange,
		ExpressionType: "notEquals",
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprLess) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		Location:       e.LocRange,
		ExpressionType: "lessThan",
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprLessOrEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		Location:       e.LocRange,
		ExpressionType: "lessThanOrEquals",
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprGreater) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		Location:       e.LocRange,
		ExpressionType: "greaterThan",
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprGreaterOrEqual) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange   `yaml:"location"`
		ExpressionType string     `yaml:"expressionType"`
		Left           Expression `yaml:"left"`
		Right          Expression `yaml:"right"`
	}{
		Location:       e.LocRange,
		ExpressionType: "greaterThanOrEquals",
		Left:           e.Left,
		Right:          e.Right,
	}, nil
}

func (e *ExprLogicalAnd) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange     `yaml:"location"`
		ExpressionType string       `yaml:"expressionType"`
		Expressions    []Expression `yaml:"expressions"`
	}{
		Location:       e.LocRange,
		ExpressionType: "logicalAND",
		Expressions:    e.Expressions,
	}, nil
}

func (e *ExprLogicalOr) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange     `yaml:"location"`
		ExpressionType string       `yaml:"expressionType"`
		Expressions    []Expression `yaml:"expressions"`
	}{
		Location:       e.LocRange,
		ExpressionType: "logicalOR",
		Expressions:    e.Expressions,
	}, nil
}

func (e *True) MarshalYAML() (any, error) {
	var t string
	if e.TypeDef != nil {
		t = e.TypeDef.Name
	}
	return struct {
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Type           string   `yaml:"type,omitempty"`
	}{
		Location:       e.LocRange,
		ExpressionType: "true",
		Type:           t,
	}, nil
}

func (e *False) MarshalYAML() (any, error) {
	var t string
	if e.TypeDef != nil {
		t = e.TypeDef.Name
	}
	return struct {
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Type           string   `yaml:"type,omitempty"`
	}{
		Location:       e.LocRange,
		ExpressionType: "false",
		Type:           t,
	}, nil
}

func (i *Int) MarshalYAML() (any, error) {
	var t string
	if i.TypeDef != nil {
		t = i.TypeDef.Name
	}
	return struct {
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Type           string   `yaml:"type,omitempty"`
		Value          int64    `yaml:"value"`
	}{
		Location:       i.LocRange,
		ExpressionType: "int",
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
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Type           string   `yaml:"type,omitempty"`
		Value          float64  `yaml:"value"`
	}{
		Location:       f.LocRange,
		ExpressionType: "float",
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
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Type           string   `yaml:"type,omitempty"`
		Value          string   `yaml:"value"`
	}{
		Location:       s.LocRange,
		ExpressionType: "string",
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
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Type           string   `yaml:"type,omitempty"`
	}{
		Location:       n.LocRange,
		ExpressionType: "null",
		Type:           t,
	}, nil
}

func (e *Enum) MarshalYAML() (any, error) {
	var t string
	if e.TypeDef != nil {
		t = e.TypeDef.Name
	}
	return struct {
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Value          string   `yaml:"value"`
		Type           string   `yaml:"type,omitempty"`
	}{
		Location:       e.LocRange,
		ExpressionType: "enum",
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
		Location       LocRange     `yaml:"location"`
		ExpressionType string       `yaml:"expressionType"`
		ItemType       string       `yaml:"itemType,omitempty"`
		Items          []Expression `yaml:"items"`
	}{
		Location:       a.LocRange,
		ExpressionType: "array",
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
		Location       LocRange       `yaml:"location"`
		ExpressionType string         `yaml:"expressionType"`
		Type           string         `yaml:"type,omitempty"`
		Fields         []*ObjectField `yaml:"fields"`
	}{
		Location:       o.LocRange,
		ExpressionType: "object",
		Type:           t,
		Fields:         o.Fields,
	}, nil
}

func (r *Variable) MarshalYAML() (any, error) {
	return struct {
		Location       LocRange `yaml:"location"`
		ExpressionType string   `yaml:"expressionType"`
		Name           string   `yaml:"name"`
	}{
		Location:       r.LocRange,
		ExpressionType: "variableReference",
		Name:           r.Name.Name,
	}, nil
}

func (s *SelectionInlineFrag) MarshalYAML() (any, error) {
	return struct {
		Location      LocRange      `yaml:"location"`
		SelectionType string        `yaml:"selectionType"`
		TypeCondition TypeCondition `yaml:"typeCondition"`
		SelectionsSet SelectionSet  `yaml:"selectionSet,omitempty"`
	}{
		Location:      s.LocRange,
		SelectionType: "inlineFragment",
		TypeCondition: s.TypeCondition,
		SelectionsSet: s.SelectionSet,
	}, nil
}

func (f *ObjectField) MarshalYAML() (any, error) {
	var t string
	if f.Def != nil {
		t = f.Def.Type.String()
	}
	type Var struct {
		Location LocRange `yaml:"location"`
		Name     string   `yaml:"name"`
	}
	var v *Var
	if f.AssociatedVariable != nil {
		v = &Var{
			Location: f.AssociatedVariable.LocRange,
			Name:     f.AssociatedVariable.Name,
		}
	}
	return struct {
		Location   LocRange   `yaml:"location"`
		Name       Name       `yaml:"name"`
		Variable   *Var       `yaml:"variable,omitempty"`
		Type       string     `yaml:"type,omitempty"`
		Constraint Expression `yaml:"constraint"`
	}{
		Location:   f.LocRange,
		Name:       f.Name,
		Variable:   v,
		Type:       t,
		Constraint: f.Constraint,
	}, nil
}

func (n Name) MarshalYAML() (any, error) {
	return struct {
		Location LocRange `yaml:"location"`
		Name     string   `yaml:"name"`
	}{
		Location: n.LocRange,
		Name:     n.Name,
	}, nil
}

func (s *SelectionField) MarshalYAML() (any, error) {
	var t string
	if s.Def != nil {
		t = s.Def.Type.String()
	}
	return struct {
		Location      LocRange     `yaml:"location"`
		SelectionType string       `yaml:"selectionType"`
		Name          Name         `yaml:"name"`
		Type          string       `yaml:"type,omitempty"`
		ArgumentList  ArgumentList `yaml:"argumentList,omitempty"`
		SelectionSet  SelectionSet `yaml:"selectionSet,omitempty"`
	}{
		Location:      s.LocRange,
		SelectionType: "field",
		Name:          s.Name,
		Type:          t,
		ArgumentList:  s.ArgumentList,
		SelectionSet:  s.SelectionSet,
	}, nil
}

func (e *SelectionMax) MarshalYAML() (any, error) {
	return struct {
		Location      LocRange     `yaml:"location"`
		SelectionType string       `yaml:"selectionType"`
		Limit         int          `yaml:"limit"`
		Options       SelectionSet `yaml:"options"`
	}{
		Location:      e.LocRange,
		SelectionType: "max",
		Limit:         e.Limit,
		Options:       e.Options,
	}, nil
}

func (a *Argument) MarshalYAML() (any, error) {
	var t string
	if a.Def != nil {
		t = a.Def.Type.String()
	}
	type Var struct {
		Location LocRange `yaml:"location"`
		Name     string   `yaml:"name"`
	}
	var v *Var
	if a.AssociatedVariable != nil {
		v = &Var{
			Location: a.AssociatedVariable.LocRange,
			Name:     a.AssociatedVariable.Name,
		}
	}
	return struct {
		Location   LocRange   `yaml:"location"`
		Name       Name       `yaml:"name"`
		Variable   *Var       `yaml:"variable,omitempty"`
		Type       string     `yaml:"type,omitempty"`
		Constraint Expression `yaml:"constraint"`
	}{
		Location:   a.LocRange,
		Name:       a.Name,
		Variable:   v,
		Type:       t,
		Constraint: a.Constraint,
	}, nil
}
