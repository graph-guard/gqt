package gqt

import (
	"fmt"
	"math"
)

// Optimize optimizes the operation o by reducing constant expressions.
// For example, an addition expression where the left and right
// addends are constants will be reduced to, or in other words,
// replaced with a single constant containing the result of the expression.
func Optimize(e Expression) Expression {
	switch e := e.(type) {
	case *Operation:
		for i, x := range e.Selections {
			e.Selections[i] = Optimize(x)
		}
		return e
	case *SelectionInlineFrag:
		for i, x := range e.Selections {
			e.Selections[i] = Optimize(x)
		}
		return e
	case *SelectionField:
		for i, x := range e.Arguments {
			e.Arguments[i].Constraint = Optimize(x.Constraint)
		}
		for i, x := range e.Selections {
			e.Selections[i] = Optimize(x)
		}
		return e
	case *ConstrAny:
		return e
	case *ConstrEquals:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrNotEquals:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrGreater:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLess:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrGreaterOrEqual:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLessOrEqual:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLenEquals:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLenNotEquals:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLenGreater:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLenLess:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLenGreaterOrEqual:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrLenLessOrEqual:
		e.Value = Optimize(e.Value)
		return e
	case *ConstrMap:
		e.Constraint = Optimize(e.Constraint)
		return e
	case *ExprParentheses:
		e.Expression = Optimize(e.Expression)
		setLocRange(e.Expression, e.LocRange)
		return e.Expression
	case *ExprEqual:
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)

		comparable, equal := comparableAndEqual(e.Left, e.Right)
		if !comparable {
			return e
		}
		if equal {
			return &True{
				LocRange: e.LocRange,
				Parent:   e.Parent,
			}
		}
		return &False{
			LocRange: e.LocRange,
			Parent:   e.Parent,
		}
	case *ExprNotEqual:
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)

		comparable, equal := comparableAndEqual(e.Left, e.Right)
		if !comparable {
			return e
		}
		if !equal {
			return &True{
				LocRange: e.LocRange,
				Parent:   e.Parent,
			}
		}
		return &False{
			LocRange: e.LocRange,
			Parent:   e.Parent,
		}
	case *ExprLogicalNegation:
		e.Expression = Optimize(e.Expression)
		switch e := e.Expression.(type) {
		case *True:
			return &False{
				LocRange: e.LocRange,
				Parent:   e.Parent,
			}
		case *False:
			return &True{
				LocRange: e.LocRange,
				Parent:   e.Parent,
			}
		}
		return e
	case *ExprNumericNegation:
		e.Expression = Optimize(e.Expression)
		switch e := e.Expression.(type) {
		case *Int:
			return &Int{
				LocRange: e.LocRange,
				Parent:   e.Parent,
				Value:    -e.Value,
			}
		case *Float:
			return &Float{
				LocRange: e.LocRange,
				Parent:   e.Parent,
				Value:    -e.Value,
			}
		}
		return e
	case *ExprLogicalOr:
		for i, x := range e.Expressions {
			e.Expressions[i] = Optimize(x)
		}
		for _, x := range e.Expressions {
			if eq, ok := x.(*ConstrEquals); ok {
				switch eq.Value.(type) {
				case *True:
					constrEq := &ConstrEquals{
						LocRange: e.LocRange,
						Parent:   e.Parent,
					}
					constrEq.Value = &True{
						LocRange: e.LocRange,
						Parent:   constrEq,
					}
					return constrEq
				case *False:
					continue
				default:
					return e
				}
			}
		}
		constrEq := &ConstrEquals{
			LocRange: e.LocRange,
			Parent:   e.Parent,
		}
		constrEq.Value = &False{
			LocRange: e.LocRange,
			Parent:   constrEq,
		}
		return constrEq
	case *ExprLogicalAnd:
		for i, x := range e.Expressions {
			e.Expressions[i] = Optimize(x)
		}
		for _, x := range e.Expressions {
			if eq, ok := x.(*ConstrEquals); ok {
				switch eq.Value.(type) {
				case *True:
					continue
				case *False:
					constrEq := &ConstrEquals{
						LocRange: e.LocRange,
						Parent:   e.Parent,
					}
					constrEq.Value = &False{
						LocRange: e.LocRange,
						Parent:   constrEq,
					}
					return constrEq
				default:
					return e
				}
			}
		}
		constrEq := &ConstrEquals{
			LocRange: e.LocRange,
			Parent:   e.Parent,
		}
		constrEq.Value = &True{
			LocRange: e.LocRange,
			Parent:   constrEq,
		}
		return constrEq
	case *ExprAddition:
		e.AddendLeft = Optimize(e.AddendLeft)
		e.AddendRight = Optimize(e.AddendRight)
		if isFloatOrInt(e.AddendLeft) && isFloatOrInt(e.AddendRight) {
			if e.AddendLeft.IsFloat() || e.AddendRight.IsFloat() {
				return &Float{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value: getFloat(e.AddendLeft) +
						getFloat(e.AddendRight),
				}
			}
			return &Int{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value: e.AddendLeft.(*Int).Value +
					e.AddendRight.(*Int).Value,
			}
		}
		return e
	case *ExprSubtraction:
		e.Minuend = Optimize(e.Minuend)
		e.Subtrahend = Optimize(e.Subtrahend)
		if isFloatOrInt(e.Minuend) && isFloatOrInt(e.Subtrahend) {
			if e.Minuend.IsFloat() || e.Subtrahend.IsFloat() {
				return &Float{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value: getFloat(e.Minuend) -
						getFloat(e.Subtrahend),
				}
			}
			return &Int{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value: e.Minuend.(*Int).Value -
					e.Subtrahend.(*Int).Value,
			}
		}
		return e
	case *ExprMultiplication:
		e.Multiplicant = Optimize(e.Multiplicant)
		e.Multiplicator = Optimize(e.Multiplicator)
		if isFloatOrInt(e.Multiplicant) && isFloatOrInt(e.Multiplicator) {
			if e.Multiplicant.IsFloat() || e.Multiplicator.IsFloat() {
				return &Float{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value: getFloat(e.Multiplicant) *
						getFloat(e.Multiplicator),
				}
			}
			return &Int{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value: e.Multiplicant.(*Int).Value *
					e.Multiplicator.(*Int).Value,
			}
		}
	case *ExprDivision:
		e.Dividend = Optimize(e.Dividend)
		e.Divisor = Optimize(e.Divisor)
		if isFloatOrInt(e.Dividend) && isFloatOrInt(e.Divisor) {
			if e.Dividend.IsFloat() || e.Divisor.IsFloat() {
				return &Float{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value: getFloat(e.Dividend) /
						getFloat(e.Divisor),
				}
			}
			return &Int{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value: e.Dividend.(*Int).Value /
					e.Divisor.(*Int).Value,
			}
		}
	case *ExprModulo:
		e.Dividend = Optimize(e.Dividend)
		e.Divisor = Optimize(e.Divisor)
		if isFloatOrInt(e.Dividend) && isFloatOrInt(e.Divisor) {
			if e.Dividend.IsFloat() || e.Divisor.IsFloat() {
				return &Float{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value: math.Mod(
						getFloat(e.Dividend), getFloat(e.Divisor),
					),
				}
			}
			return &Int{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value: e.Dividend.(*Int).Value %
					e.Divisor.(*Int).Value,
			}
		}
	case *ExprGreater:
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)
		if isFloatOrInt(e.Left) && isFloatOrInt(e.Right) {
			if getFloat(e.Left) > getFloat(e.Right) {
				return &True{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			} else {
				return &False{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			}
		}
		return e
	case *ExprLess:
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)
		if isFloatOrInt(e.Left) && isFloatOrInt(e.Right) {
			if getFloat(e.Left) < getFloat(e.Right) {
				return &True{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			} else {
				return &False{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			}
		}
		return e
	case *ExprGreaterOrEqual:
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)
		if isFloatOrInt(e.Left) && isFloatOrInt(e.Right) {
			if getFloat(e.Left) >= getFloat(e.Right) {
				return &True{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			} else {
				return &False{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			}
		}
		return e
	case *ExprLessOrEqual:
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)
		e.Left = Optimize(e.Left)
		e.Right = Optimize(e.Right)
		if isFloatOrInt(e.Left) && isFloatOrInt(e.Right) {
			if getFloat(e.Left) <= getFloat(e.Right) {
				return &True{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			} else {
				return &False{
					Parent:   e.Parent,
					LocRange: e.LocRange,
				}
			}
		}
		return e
	case *Int:
		return e
	case *Float:
		return e
	case *True:
		return e
	case *False:
		return e
	case *Null:
		return e
	case *Array:
		for i, x := range e.Items {
			e.Items[i] = Optimize(x)
		}
		return e
	case *Variable:
		return e
	case *Enum:
		return e
	case *String:
		return e
	case *Object:
		for i, x := range e.Fields {
			e.Fields[i].Constraint = Optimize(x.Constraint)
		}
		return e
	}
	panic(fmt.Errorf("unhandled type: %T", e))
}

func isFloatOrInt(e Expression) bool {
	_, f := e.(*Float)
	_, i := e.(*Int)
	return f || i
}

func getFloat(e Expression) float64 {
	if vf, ok := e.(*Float); ok {
		return vf.Value
	}
	if vf, ok := e.(*Int); ok {
		return float64(vf.Value)
	}
	return 0
}

func comparableAndEqual(left, right Expression) (comparable, equal bool) {
	switch l := left.(type) {
	case *String:
		return true, l.Value == right.(*String).Value
	case *Int:
		return true, float64(l.Value) == getFloat(right)
	case *Float:
		return true, l.Value == getFloat(right)
	case *True:
		_, ok := right.(*True)
		return true, ok
	case *False:
		_, ok := right.(*False)
		return true, ok
	case *Enum:
		return true, l.Value == right.(*Enum).Value
	case *Null:
		_, ok := right.(*Null)
		return true, ok
	case *Array:
		if r, ok := right.(*Array); ok {
			if len(l.Items) != len(r.Items) {
				return true, false
			}
			for i := range l.Items {
				l, r := l.Items[i].(*ConstrEquals), r.Items[i].(*ConstrEquals)
				comparable, equal = comparableAndEqual(l.Value, r.Value)
				if !comparable || !equal {
					return comparable, equal
				}
			}
			return true, true
		}
	}
	return false, false
}
