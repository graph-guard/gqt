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
				y1: any
				t1: type = T
				t2: type != T
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
				a1: is = []
				a2: is = [ are < 10 ]
				a3: is = [ is = "a", is = "b" ]
				o1: is = {
					of1: is = true
					of2: is = false
					oo2: is = {
						oa1: is != null
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
					Constraint: gqt.ConstraintTypeEqual{TypeName: "T"},
				}, {
					Name:       "t2",
					Constraint: gqt.ConstraintTypeNotEqual{TypeName: "T"},
				}, {
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
				}, {
					Name:       "a1",
					Constraint: gqt.ConstraintIsEqual{Value: gqt.ValueArray{}},
				}, {
					Name: "a2",
					Constraint: gqt.ConstraintIsEqual{Value: gqt.ValueArray{
						Items: []gqt.Constraint{
							gqt.ConstraintAreLess{Value: 10},
						},
					}},
				}, {
					Name: "a3",
					Constraint: gqt.ConstraintIsEqual{Value: gqt.ValueArray{
						Items: []gqt.Constraint{
							gqt.ConstraintIsEqual{Value: "a"},
							gqt.ConstraintIsEqual{Value: "b"},
						},
					}},
				}, {
					Name: "o1",
					Constraint: gqt.ConstraintIsEqual{Value: gqt.ValueObject{
						Fields: []gqt.ObjectField{{
							Name:  "of1",
							Value: gqt.ConstraintIsEqual{Value: true},
						}, {
							Name:  "of2",
							Value: gqt.ConstraintIsEqual{Value: false},
						}, {
							Name: "oo2",
							Value: gqt.ConstraintIsEqual{
								Value: gqt.ValueObject{
									Fields: []gqt.ObjectField{{
										Name: "oa1",
										Value: gqt.ConstraintIsNotEqual{
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
				o1: is = {
					f1: is = true
					f2: is = false
				}
			) {b}
		}`,
		gqt.DocMutation{
			Selections: []gqt.Selection{{
				Name: "a",
				InputConstraints: []gqt.InputConstraint{{
					Name: "o1",
					Constraint: gqt.ConstraintIsEqual{Value: gqt.ValueObject{
						Fields: []gqt.ObjectField{
							{
								Name:  "f1",
								Value: gqt.ConstraintIsEqual{Value: true},
							},
							{
								Name:  "f2",
								Value: gqt.ConstraintIsEqual{Value: false},
							},
						},
					}},
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
			`query { x(a: is = 3 a: is = 4) }`,
			"error at 20: redundant constraint",
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
			"error at 10: expected constraint name",
		},
		{13,
			`query{f(x:  `,
			"error at 12: expected constraint name",
		},
		{14,
			`query{f(x:is`,
			"error at 12: expected constraint operator",
		},
		{15,
			`query{f(x:is  `,
			"error at 14: expected constraint operator",
		},
		{16,
			`query{f(x:is=`,
			"error at 13: expected value",
		},
		{17,
			`query{f(x:is=  `,
			"error at 15: expected value",
		},
		{18,
			`query{f(x:is=1`,
			"error at 14: expected parameter name",
		},
		{19,
			`query{f(x:is=1  `,
			"error at 16: expected parameter name",
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
