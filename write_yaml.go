package gqt

import (
	"fmt"
	"io"
)

// WriteYAML writes any AST object as YAML to w.
func WriteYAML(w io.Writer, o any) (written int, err error) {
	if written, err = writeYAML(w, 0, o); err == nil {
		var n int
		n, err = w.Write([]byte("\n"))
		written += n
	}
	return
}

func writeYAML(
	w io.Writer,
	indent int,
	o any,
) (written int, err error) {
	wrBR := func() (stop bool) {
		var n int
		if n, err = w.Write([]byte("\n")); err != nil {
			return true
		}
		written += n
		return false
	}
	wrf := func(format string, v ...any) (stop bool) {
		var n int
		if n, err = w.Write([]byte(fmt.Sprintf(format, v...))); err != nil {
			return true
		}
		written += n
		return false
	}
	writeIndent := func(indent int) {
		for i := 0; i < indent; i++ {
			var n int
			if n, err = w.Write([]byte("  ")); err != nil {
				return
			}
			written += n
		}
	}

	if indent != 0 {
		if wrBR() {
			return
		}
	}
	writeIndent(indent)

	var n int
	switch v := o.(type) {
	case *Operation:
		if wrf("Operation[%d:%d](%s):", v.Line, v.Column, v.Type) {
			return
		}
		for _, s := range v.Selections {
			if n, err = writeYAML(w, indent+1, s); err != nil {
				return
			}
			written += n
		}
	case *Int:
		if wrf("- Int[%d:%d](%v)", v.Line, v.Column, v.Value) {
			return
		}
	case *Float:
		if wrf("- Float[%d:%d](%v)", v.Line, v.Column, v.Value) {
			return
		}
	case *String:
		if wrf("- String[%d:%d](%q)", v.Line, v.Column, v.Value) {
			return
		}
	case *True:
		if wrf("- True[%d:%d]", v.Line, v.Column) {
			return
		}
	case *False:
		if wrf("- False[%d:%d]", v.Line, v.Column) {
			return
		}
	case *Null:
		if wrf("- Null[%d:%d]", v.Line, v.Column) {
			return
		}
	case *Enum:
		if wrf("- Enum[%d:%d](%s)", v.Line, v.Column, v.Value) {
			return
		}
	case *Array:
		colon := ""
		if len(v.Items) > 0 {
			colon = ":"
		}
		if wrf(
			"- Array[%d:%d](%d items)%s",
			v.Line, v.Column, len(v.Items), colon,
		) {
			return
		}
		for _, item := range v.Items {
			if n, err = writeYAML(w, indent+1, item); err != nil {
				return
			}
			written += n
		}
	case *Object:
		colon := ""
		if len(v.Fields) > 0 {
			colon = ":"
		}
		if wrf(
			"- Object[%d:%d](%d fields)%s",
			v.Line, v.Column, len(v.Fields), colon,
		) {
			return
		}
		for _, field := range v.Fields {
			if n, err = writeYAML(w, indent+1, field); err != nil {
				return
			}
			written += n
		}
	case *ObjectField:
		colon := ""
		if v.Constraint != nil {
			colon = ":"
		}
		if v.AssociatedVariable != nil {
			if wrf(
				"- ObjectField[%d:%d](%s=$%s)%s",
				v.Line, v.Column, v.Name, v.AssociatedVariable.Name, colon,
			) {
				return
			}
		} else {
			if wrf(
				"- ObjectField[%d:%d](%s)%s",
				v.Line, v.Column, v.Name, colon,
			) {
				return
			}
		}
		if v.Constraint != nil {
			if n, err = writeYAML(w, indent+1, v.Constraint); err != nil {
				return
			}
			written += n
		}
	case *ExprModulo:
		if wrf("- Modulo[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Dividend); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Divisor); err != nil {
			return
		}
		written += n
	case *ExprDivision:
		if wrf("- Division[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Dividend); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Divisor); err != nil {
			return
		}
		written += n
	case *ExprMultiplication:
		if wrf("- Multiplication[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Multiplicant); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Multiplicator); err != nil {
			return
		}
		written += n
	case *ExprAddition:
		if wrf("- Addition[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.AddendLeft); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.AddendRight); err != nil {
			return
		}
		written += n
	case *ExprSubtraction:
		if wrf("- Subtraction[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Minuend); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Subtrahend); err != nil {
			return
		}
		written += n
	case *ExprEqual:
		if wrf("- Equal[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Left); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Right); err != nil {
			return
		}
		written += n
	case *ExprNotEqual:
		if wrf("- NotEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Left); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Right); err != nil {
			return
		}
		written += n
	case *ExprLessOrEqual:
		if wrf("- LessOrEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Left); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Right); err != nil {
			return
		}
		written += n
	case *ExprGreaterOrEqual:
		if wrf("- GreaterOrEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Left); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Right); err != nil {
			return
		}
		written += n
	case *ExprLess:
		if wrf("- Less[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Left); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Right); err != nil {
			return
		}
		written += n
	case *ExprGreater:
		if wrf("- Greater[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Left); err != nil {
			return
		}
		written += n
		if n, err = writeYAML(w, indent+1, v.Right); err != nil {
			return
		}
		written += n
	case *ExprLogicalNegation:
		if wrf("- LogicalNegation[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Expression); err != nil {
			return
		}
		written += n
	case *Variable:
		if wrf("- Variable[%d:%d](%s)", v.Line, v.Column, v.Name) {
			return
		}
	case *ExprParentheses:
		if wrf("- Parentheses[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Expression); err != nil {
			return
		}
		written += n
	case *ExprConstrEquals:
		if wrf("- ConstrEquals[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrNotEquals:
		if wrf("- ConstrNotEquals[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLess:
		if wrf("- ConstrLess[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLessOrEqual:
		if wrf("- ConstrLessOrEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrGreater:
		if wrf("- ConstrGreater[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrGreaterOrEqual:
		if wrf("- ConstrGreaterOrEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLenEquals:
		if wrf("- ConstrLenEquals[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLenNotEquals:
		if wrf("- ConstrLenNotEquals[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLenLess:
		if wrf("- ConstrLenLess[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLenLessOrEqual:
		if wrf("- ConstrLenLessOrEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLenGreater:
		if wrf("- ConstrLenGreater[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprConstrLenGreaterOrEqual:
		if wrf("- ConstrLenGreaterOrEqual[%d:%d]:", v.Line, v.Column) {
			return
		}
		if n, err = writeYAML(w, indent+1, v.Value); err != nil {
			return
		}
		written += n
	case *ExprLogicalAnd:
		if wrf("- ConstrLogicalAND[%d:%d]:", v.Line, v.Column) {
			return
		}
		for _, e := range v.Expressions {
			if n, err = writeYAML(w, indent+1, e); err != nil {
				return
			}
			written += n
		}
	case *ExprLogicalOr:
		if wrf("- ConstrLogicalOR[%d:%d]:", v.Line, v.Column) {
			return
		}
		for _, e := range v.Expressions {
			if n, err = writeYAML(w, indent+1, e); err != nil {
				return
			}
			written += n
		}
	case *SelectionField:
		colon := ""
		if len(v.Selections) > 0 || len(v.Arguments) > 0 {
			colon = ":"
		}
		if wrf(
			"- SelectionField[%d:%d](%s)%s",
			v.Line, v.Column, v.Name, colon,
		) {
			return
		}
		if len(v.Arguments) > 0 {
			if wrBR() {
				return
			}
			writeIndent(indent + 1)
			if wrf("arguments:") {
				return
			}
			for _, a := range v.Arguments {
				if n, err = writeYAML(w, indent+2, a); err != nil {
					return
				}
				written += n
			}
		}
		if len(v.Selections) > 0 {
			if wrBR() {
				return
			}
			writeIndent(indent + 1)
			if wrf("selections:") {
				return
			}
			for _, a := range v.Selections {
				if n, err = writeYAML(w, indent+2, a); err != nil {
					return
				}
				written += n
			}
		}
	case *Argument:
		colon := ""
		if v.Constraint != nil {
			colon = ":"
		}
		if v.AssociatedVariable != nil {
			if wrf(
				"- Argument[%d:%d](%s=$%s)%s",
				v.Line, v.Column, v.Name, v.AssociatedVariable.Name, colon,
			) {
				return
			}
		} else {
			if wrf(
				"- Argument[%d:%d](%s)%s",
				v.Line, v.Column, v.Name, colon,
			) {
				return
			}
		}
		if v.Constraint != nil {
			if n, err = writeYAML(w, indent+1, v.Constraint); err != nil {
				return
			}
			written += n
		}
	case *SelectionInlineFrag:
		colon := ""
		if len(v.Selections) > 0 {
			colon = ":"
		}
		if wrf(
			"- SelectionInlineFrag[%d:%d](%s)%s",
			v.Line, v.Column, v.TypeCondition, colon,
		) {
			return
		}
		if len(v.Selections) > 0 {
			if wrBR() {
				return
			}
			writeIndent(indent + 1)
			if wrf("selections:") {
				return
			}
			for _, a := range v.Selections {
				if n, err = writeYAML(w, indent+2, a); err != nil {
					return
				}
				written += n
			}
		}
	default:
		return written, fmt.Errorf("unsupported value: %#v", o)
	}
	return
}
