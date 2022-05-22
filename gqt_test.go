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
					x1: val = 1
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
						Constraint: gqt.ConstraintValEqual{
							Value: int64(1),
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
					val#comment after token
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
						Constraint: gqt.ConstraintValEqual{
							Value: int64(1),
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
				y1: any
				t1: type = Boolean
				t2: type = Int
				t3: type = Float
				t4: type = String
				t5: type = ID
				t6: type = Input
				i1: val = 42.0
				i2: val != 42.0
				i3: val > 42
				i4: val < 42.0
				i5: val >= 42
				i6: val <= 42.0
				i7: val = "text"
				i8: val != "text"
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
				a1: val = []
				a2: val = [ ... val < 10.0 ]
				a3: val = [ val = "a", val = "b" ]
				o1: val = {
					of1: val = true
					of2: val = false
					oo2: val = {
						oa1: val != null
						oa2: any
					}
				}
			) {b}
		}`,
		gqt.DocMutation{
			Selections: []gqt.Selection{{
				Name: "a",
				InputConstraints: []gqt.InputConstraint{{
					Name:       "y1",
					Constraint: gqt.ConstraintAny{},
				}, {
					Name:       "t1",
					Constraint: gqt.ConstraintTypeEqual{Type: gqt.TypeKindBoolean},
				}, {
					Name:       "t2",
					Constraint: gqt.ConstraintTypeEqual{Type: gqt.TypeKindInt},
				}, {
					Name:       "t3",
					Constraint: gqt.ConstraintTypeEqual{Type: gqt.TypeKindFloat},
				}, {
					Name:       "t4",
					Constraint: gqt.ConstraintTypeEqual{Type: gqt.TypeKindString},
				}, {
					Name:       "t5",
					Constraint: gqt.ConstraintTypeEqual{Type: gqt.TypeKindID},
				}, {
					Name:       "t6",
					Constraint: gqt.ConstraintTypeEqual{Type: gqt.TypeKindInput},
				}, {
					Name:       "i1",
					Constraint: gqt.ConstraintValEqual{Value: float64(42)},
				}, {
					Name:       "i2",
					Constraint: gqt.ConstraintValNotEqual{Value: float64(42)},
				}, {
					Name:       "i3",
					Constraint: gqt.ConstraintValGreater{Value: int64(42)},
				}, {
					Name:       "i4",
					Constraint: gqt.ConstraintValLess{Value: float64(42)},
				}, {
					Name:       "i5",
					Constraint: gqt.ConstraintValGreaterOrEqual{Value: int64(42)},
				}, {
					Name:       "i6",
					Constraint: gqt.ConstraintValLessOrEqual{Value: float64(42)},
				}, {
					Name:       "i7",
					Constraint: gqt.ConstraintValEqual{Value: "text"},
				}, {
					Name:       "i8",
					Constraint: gqt.ConstraintValNotEqual{Value: "text"},
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
				}, {
					Name:       "a1",
					Constraint: gqt.ConstraintValEqual{Value: gqt.ValueArray{}},
				}, {
					Name: "a2",
					Constraint: gqt.ConstraintMap{
						Constraint: gqt.ConstraintValLess{
							Value: 10.0,
						},
					},
				}, {
					Name: "a3",
					Constraint: gqt.ConstraintValEqual{Value: gqt.ValueArray{
						Items: []gqt.Constraint{
							gqt.ConstraintValEqual{Value: "a"},
							gqt.ConstraintValEqual{Value: "b"},
						},
					}},
				}, {
					Name: "o1",
					Constraint: gqt.ConstraintValEqual{Value: gqt.ValueObject{
						Fields: []gqt.ObjectField{{
							Name:  "of1",
							Value: gqt.ConstraintValEqual{Value: true},
						}, {
							Name:  "of2",
							Value: gqt.ConstraintValEqual{Value: false},
						}, {
							Name: "oo2",
							Value: gqt.ConstraintValEqual{
								Value: gqt.ValueObject{
									Fields: []gqt.ObjectField{{
										Name: "oa1",
										Value: gqt.ConstraintValNotEqual{
											Value: nil,
										},
									}, {
										Name:  "oa2",
										Value: gqt.ConstraintAny{},
									}},
								},
							},
						}},
					}},
				}},
				Selections: []gqt.Selection{{
					Name: "b",
				}},
			}},
		},
	}, {
		`mutation {
			a(
				a: val > 0.0 && val < 3.0
				b: val > 0.0 && val < 9.0 && val != 5.0
				c: val = 1.0 || val = 2.0
				d: val = 1.0 || val = 2.0 || val = 3.0
				e: type = String && bytelen > 3 && bytelen <= 10 ||
					val = true ||
					val = 1.0
			) {b}
		}`,
		gqt.DocMutation{
			Selections: []gqt.Selection{{
				Name: "a",
				InputConstraints: []gqt.InputConstraint{{
					Name: "a",
					Constraint: gqt.ConstraintAnd{
						gqt.ConstraintValGreater{Value: float64(0)},
						gqt.ConstraintValLess{Value: float64(3)},
					},
				}, {
					Name: "b",
					Constraint: gqt.ConstraintAnd{
						gqt.ConstraintValGreater{Value: float64(0)},
						gqt.ConstraintValLess{Value: float64(9)},
						gqt.ConstraintValNotEqual{Value: float64(5)},
					},
				}, {
					Name: "c",
					Constraint: gqt.ConstraintOr{
						gqt.ConstraintValEqual{Value: float64(1)},
						gqt.ConstraintValEqual{Value: float64(2)},
					},
				}, {
					Name: "d",
					Constraint: gqt.ConstraintOr{
						gqt.ConstraintValEqual{Value: float64(1)},
						gqt.ConstraintValEqual{Value: float64(2)},
						gqt.ConstraintValEqual{Value: float64(3)},
					},
				}, {
					Name: "e",
					Constraint: gqt.ConstraintOr{
						gqt.ConstraintAnd{
							gqt.ConstraintTypeEqual{Type: gqt.TypeKindString},
							gqt.ConstraintBytelenGreater{Value: uint(3)},
							gqt.ConstraintBytelenLessOrEqual{Value: uint(10)},
						},
						gqt.ConstraintValEqual{Value: true},
						gqt.ConstraintValEqual{Value: float64(1)},
					},
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

func TestParseErr(t *testing.T) {
	for _, td := range []struct {
		name      string
		input     string
		expectErr string
	}{
		{"redundant_selection",
			`query { x x }`,
			"error at 10: redundant selection",
		},
		{"redundant_constraint",
			`query { x(a: val = 3 a: val = 4) }`,
			"error at 21: redundant constraint",
		},
		{"map_multiple_constraints",
			`query {x(a: val = [ ... val=0 val=1 ])}`,
			"error at 30: expected right square bracket",
		},
		{"unsupported_arg_type",
			`query {x(a: type = Foo}`,
			"error at 19: unsupported type kind",
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			d, err := gqt.Parse([]byte(td.input))
			require.True(t, err.IsErr())
			require.Equal(t, td.expectErr, err.Error())
			require.Nil(t, d)
		})
	}
}

func TestParseErrUnexpectedEOF(t *testing.T) {
	for ti, td := range []struct {
		index    int
		input    string
		expected string
	}{
		{0,
			``,
			"error at 0: expected definition",
		},
		{1,
			`  `,
			"error at 2: expected definition",
		},
		{2,
			`query`,
			"error at 5: expected selection set",
		},
		{3,
			`query  `,
			"error at 7: expected selection set",
		},
		{4,
			`query{`,
			"error at 6: expected field name",
		},
		{5,
			`query{  `,
			"error at 8: expected field name",
		},
		{6,
			`query{f`,
			"error at 7: expected field name",
		},
		{7,
			`query{f  `,
			"error at 9: expected field name",
		},
		{8,
			`query{f(`,
			"error at 8: expected parameter name",
		},
		{9,
			`query{f(  `,
			"error at 10: expected parameter name",
		},
		{10,
			`query{f(x`,
			"error at 9: expected column after parameter name",
		},
		{11,
			`query{f(x  `,
			"error at 11: expected column after parameter name",
		},
		{12,
			`query{f(x:`,
			"error at 10: expected constraint subject",
		},
		{13,
			`query{f(x:  `,
			"error at 12: expected constraint subject",
		},
		{14,
			`query{f(x:val`,
			"error at 13: expected constraint operator",
		},
		{15,
			`query{f(x:val  `,
			"error at 15: expected constraint operator",
		},
		{16,
			`query{f(x:val=`,
			"error at 14: expected value",
		},
		{17,
			`query{f(x:val=  `,
			"error at 16: expected value",
		},
		{18,
			`query{f(x:val=1`,
			"error at 15: expected parameter name",
		},
		{19,
			`query{f(x:val=1  `,
			"error at 17: expected parameter name",
		},
		{20,
			`query{f(x:val=[`,
			"error at 15: expected constraint subject",
		},
		{21,
			`query{f(x:val=[.`,
			"error at 15: expected constraint subject",
		},
		{22,
			`query{f(x:val=[..`,
			"error at 15: expected constraint subject",
		},
		{23,
			`query{f(x:val=[...`,
			"error at 18: expected constraint subject",
		},
		{24,
			`query{f(x:val=[...  `,
			"error at 20: expected constraint subject",
		},
	} {
		t.Run("", func(t *testing.T) {
			r := require.New(t)
			r.Equal(ti, td.index)
			d, err := gqt.Parse([]byte(td.input))
			r.True(err.IsErr())
			r.Equal(td.expected, err.Error())
			r.Nil(d)
		})
	}
}
