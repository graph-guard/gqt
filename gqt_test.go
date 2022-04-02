package gqt_test

import (
	"testing"

	"github.com/graph-guard/gqt"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for _, td := range []struct {
		input  string
		expect gqt.Doc
	}{{
		`query {
			a {
				b(
					x1: is = 1
				) {
					c
					d
				}
			}
		}`,
		gqt.DocQuery{
			Selections: []gqt.Selection{{
				Name: "a",
				Selections: []gqt.Selection{{
					Name: "b",
					InputConstraints: []gqt.InputConstraint{{
						Name: "x1",
						Constraint: gqt.ConstraintIsEqual{
							Value: float64(1),
						},
					}},
					Selections: []gqt.Selection{{
						Name: "c",
					}, {
						Name: "d",
					}},
				}},
			}},
		},
	}, { // Comments
		`query#comment after token
		#comment after token
		{#comment after token
		#comment after token
			a#comment after token
			#comment after token
			{#comment after token
			#comment after token
				b#comment after token
				#comment after token
				(#comment after token
				#comment after token
					x1#comment after token
					#comment after token
					:#comment after token
					#comment after token
					is#comment after token
					#comment after token
					=#comment after token
					#comment after token
					1#comment after token
					#comment after token
				)#comment after token
				#comment after token
				{#comment after token
				#comment after token
					c#comment after token
					#comment after token
					d#comment after token
					#comment after token
				}#comment after token
				#comment after token
			}#comment after token
			#comment after token
		}#comment after token
		#comment after token`,
		gqt.DocQuery{
			Selections: []gqt.Selection{{
				Name: "a",
				Selections: []gqt.Selection{{
					Name: "b",
					InputConstraints: []gqt.InputConstraint{{
						Name: "x1",
						Constraint: gqt.ConstraintIsEqual{
							Value: float64(1),
						},
					}},
					Selections: []gqt.Selection{{
						Name: "c",
					}, {
						Name: "d",
					}},
				}},
			}},
		},
	}, {
		`mutation {
			a(
				i1: is = 42
				i2: is != 42
				i3: is > 42
				i4: is < 42
				i5: is >= 42
				i6: is <= 42
				i7: is = "text"
				i8: is != "text"
				l1: len = 42
				l2: len != 42
				l3: len > 42
				l4: len < 42
				l5: len >= 42
				l6: len <= 42
				b1: bytelen = 42
				b2: bytelen != 42
				b3: bytelen > 42
				b4: bytelen < 42
				b5: bytelen >= 42
				b6: bytelen <= 42
			) {b}
		}`,
		gqt.DocMutation{
			Selections: []gqt.Selection{{
				Name: "a",
				InputConstraints: []gqt.InputConstraint{{
					Name:       "i1",
					Constraint: gqt.ConstraintIsEqual{Value: float64(42)},
				}, {
					Name:       "i2",
					Constraint: gqt.ConstraintIsNotEqual{Value: float64(42)},
				}, {
					Name:       "i3",
					Constraint: gqt.ConstraintIsGreater{Value: float64(42)},
				}, {
					Name:       "i4",
					Constraint: gqt.ConstraintIsLess{Value: float64(42)},
				}, {
					Name:       "i5",
					Constraint: gqt.ConstraintIsGreaterOrEqual{Value: float64(42)},
				}, {
					Name:       "i6",
					Constraint: gqt.ConstraintIsLessOrEqual{Value: float64(42)},
				}, {
					Name:       "i7",
					Constraint: gqt.ConstraintIsEqual{Value: "text"},
				}, {
					Name:       "i8",
					Constraint: gqt.ConstraintIsNotEqual{Value: "text"},
				}, {
					Name:       "l1",
					Constraint: gqt.ConstraintLenEqual{Value: uint(42)},
				}, {
					Name:       "l2",
					Constraint: gqt.ConstraintLenNotEqual{Value: uint(42)},
				}, {
					Name:       "l3",
					Constraint: gqt.ConstraintLenGreater{Value: uint(42)},
				}, {
					Name:       "l4",
					Constraint: gqt.ConstraintLenLess{Value: uint(42)},
				}, {
					Name:       "l5",
					Constraint: gqt.ConstraintLenGreaterOrEqual{Value: uint(42)},
				}, {
					Name:       "l6",
					Constraint: gqt.ConstraintLenLessOrEqual{Value: uint(42)},
				}, {
					Name:       "b1",
					Constraint: gqt.ConstraintBytelenEqual{Value: uint(42)},
				}, {
					Name:       "b2",
					Constraint: gqt.ConstraintBytelenNotEqual{Value: uint(42)},
				}, {
					Name:       "b3",
					Constraint: gqt.ConstraintBytelenGreater{Value: uint(42)},
				}, {
					Name:       "b4",
					Constraint: gqt.ConstraintBytelenLess{Value: uint(42)},
				}, {
					Name:       "b5",
					Constraint: gqt.ConstraintBytelenGreaterOrEqual{Value: uint(42)},
				}, {
					Name:       "b6",
					Constraint: gqt.ConstraintBytelenLessOrEqual{Value: uint(42)},
				}},
				Selections: []gqt.Selection{{
					Name: "b",
				}},
			}},
		},
	}} {
		t.Run("", func(t *testing.T) {
			d, err := gqt.Parse([]byte(td.input))
			require.Zero(t, err.Error())
			require.False(t, err.IsErr())
			require.Equal(t, td.expect, d)
		})
	}

}
