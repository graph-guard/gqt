package gqt

import (
	"fmt"
	"math"
	"strconv"
)

// Optimize recursively reduces all constant expressions inside e.
// For example, an addition expression where the left and right addends are
// constant or reducable expressions will be reduced to, or in other words,
// replaced with a single constant containing the resulting sum.
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
		if n, ok := e.Expression.(*Number); ok {
			n.vf, n.vi = -n.vf, -n.vi
			setParent(n, e.Parent)
			return n
		}
		return e
	case *ExprLogicalOr:
		for i, x := range e.Expressions {
			e.Expressions[i] = Optimize(x)
		}
		var r Expression = e
		if isBoolean(e.Expressions[0]) {
			r = func() Expression {
				for _, x := range e.Expressions {
					switch x.(type) {
					case *True:
						return &True{
							LocRange: e.LocRange,
							Parent:   e.Parent,
						}
					case *False:
						continue
					default:
						return e
					}
				}
				return &False{
					LocRange: e.LocRange,
					Parent:   e.Parent,
				}
			}()
		}
		if or, ok := r.(*ExprLogicalOr); ok {
			if len(or.Expressions) < 2 {
				setParent(or.Expressions[0], or.Parent)
				setLocRange(or.Expressions[0], or.LocRange)
				return or.Expressions[0]
			}
		}
		return r
	case *ExprLogicalAnd:
		for i, x := range e.Expressions {
			e.Expressions[i] = Optimize(x)
		}
		var r Expression = e
		if isBoolean(e.Expressions[0]) {
			r = func() Expression {
				for _, x := range e.Expressions {
					switch x.(type) {
					case *True:
						continue
					case *False:
						return &False{
							LocRange: e.LocRange,
							Parent:   e.Parent,
						}
					default:
						return e
					}
				}
				return &True{
					LocRange: e.LocRange,
					Parent:   e.Parent,
				}
			}()
		}
		if or, ok := r.(*ExprLogicalAnd); ok {
			if len(or.Expressions) < 2 {
				setParent(or.Expressions[0], or.Parent)
				setLocRange(or.Expressions[0], or.LocRange)
				return or.Expressions[0]
			}
		}
		return r
	case *ExprAddition:
		e.AddendLeft = Optimize(e.AddendLeft)
		e.AddendRight = Optimize(e.AddendRight)

		if isFloatOrInt(e.AddendLeft) && isFloatOrInt(e.AddendRight) {
			if e.AddendLeft.IsFloat() || e.AddendRight.IsFloat() {
				res := getFloat(e.AddendLeft) + getFloat(e.AddendRight)
				return &Number{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value:    fmt.Sprintf("%f", res),
					isFloat:  true,
					vf:       res,
				}
			}
			res := e.AddendLeft.(*Number).vi + e.AddendRight.(*Number).vi
			return &Number{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value:    strconv.Itoa(res),
				vi:       res,
			}
		}
		return e
	case *ExprSubtraction:
		e.Minuend = Optimize(e.Minuend)
		e.Subtrahend = Optimize(e.Subtrahend)

		if isFloatOrInt(e.Minuend) && isFloatOrInt(e.Subtrahend) {
			if e.Minuend.IsFloat() || e.Subtrahend.IsFloat() {
				res := getFloat(e.Minuend) - getFloat(e.Subtrahend)
				return &Number{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value:    fmt.Sprintf("%f", res),
					isFloat:  true,
					vf:       res,
				}
			}
			res := e.Minuend.(*Number).vi - e.Subtrahend.(*Number).vi
			return &Number{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value:    strconv.Itoa(res),
				vi:       res,
			}
		}
		return e
	case *ExprMultiplication:
		e.Multiplicant = Optimize(e.Multiplicant)
		e.Multiplicator = Optimize(e.Multiplicator)

		if isFloatOrInt(e.Multiplicant) && isFloatOrInt(e.Multiplicator) {
			if e.Multiplicant.IsFloat() || e.Multiplicator.IsFloat() {
				res := getFloat(e.Multiplicant) * getFloat(e.Multiplicator)
				return &Number{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value:    fmt.Sprintf("%f", res),
					isFloat:  true,
					vf:       res,
				}
			}
			res := e.Multiplicant.(*Number).vi * e.Multiplicator.(*Number).vi
			return &Number{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value:    strconv.Itoa(res),
				vi:       res,
			}
		}
		return e
	case *ExprDivision:
		e.Dividend = Optimize(e.Dividend)
		e.Divisor = Optimize(e.Divisor)

		// Check for division by 0
		if isZeroNum(e.Divisor) {
			return e
		}

		if isFloatOrInt(e.Dividend) && isFloatOrInt(e.Divisor) {
			if e.Dividend.IsFloat() || e.Divisor.IsFloat() {
				res := getFloat(e.Dividend) / getFloat(e.Divisor)
				return &Number{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value:    fmt.Sprintf("%f", res),
					isFloat:  true,
					vf:       res,
				}
			}
			res := e.Dividend.(*Number).vi / e.Divisor.(*Number).vi
			return &Number{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value:    strconv.Itoa(res),
				vi:       res,
			}
		}
		return e
	case *ExprModulo:
		e.Dividend = Optimize(e.Dividend)
		e.Divisor = Optimize(e.Divisor)

		// Check for division by 0
		if isZeroNum(e.Divisor) {
			return e
		}

		if isFloatOrInt(e.Dividend) && isFloatOrInt(e.Divisor) {
			if e.Dividend.IsFloat() || e.Divisor.IsFloat() {
				res := math.Mod(getFloat(e.Dividend), getFloat(e.Divisor))
				return &Number{
					Parent:   e.Parent,
					LocRange: e.LocRange,
					Value:    fmt.Sprintf("%f", res),
					isFloat:  true,
					vf:       res,
				}
			}
			res := e.Dividend.(*Number).vi % e.Divisor.(*Number).vi
			return &Number{
				Parent:   e.Parent,
				LocRange: e.LocRange,
				Value:    strconv.Itoa(res),
				vi:       res,
			}
		}
		return e
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
	case *Number:
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
		if c, ok := getVarDeclConstraint(e).(*ConstrEquals); ok {
			cv := Optimize(c.Value)
			return cv
		}
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
	_, ok := e.(*Number)
	return ok
}

func isBoolean(e Expression) bool {
	switch e.(type) {
	case *True, *False:
		return true
	}
	return false
}

func getFloat(e Expression) float64 {
	if n, ok := e.(*Number); ok {
		if n.isFloat {
			return n.vf
		}
		return float64(n.vi)
	}
	return 0
}

func comparableAndEqual(left, right Expression) (comparable, equal bool) {
	if v, ok := right.(*Variable); ok {
		if _, ok := getVarDeclConstraint(v).(*ConstrEquals); !ok {
			return false, false
		}
	}

	switch l := left.(type) {
	case *String:
		return true, l.Value == right.(*String).Value
	case *Number:
		return true, getFloat(l) == getFloat(right)
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

func isZeroNum(e Expression) bool {
	if n, ok := e.(*Number); ok {
		if n.isFloat {
			return n.vf == 0
		}
		return n.vi == 0
	}
	return false
}
