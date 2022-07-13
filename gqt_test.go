package gqt_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/graph-guard/gqt"
	"github.com/stretchr/testify/require"
)

type KeyValueConstraintInterface interface {
	Key() string
	Content() gqt.Constraint
}

func TestConstraintKeyAndValue(t *testing.T) {
	for _, td := range []struct {
		input KeyValueConstraintInterface
		key   string
		value gqt.Constraint
	}{
		{
			input: gqt.InputConstraint{
				Name: "a",
				Constraint: gqt.ConstraintValEqual{
					Value: 88.0,
				},
			},
			key: "a",
			value: gqt.ConstraintValEqual{
				Value: 88.0,
			},
		},
		{
			input: gqt.ObjectField{
				Name: "a",
				Value: gqt.ConstraintValNotEqual{
					Value: 88.0,
				},
			},
			key: "a",
			value: gqt.ConstraintValNotEqual{
				Value: 88.0,
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			key := td.input.Key()
			value := td.input.Content()
			require.Equal(t, td.key, key)
			require.Equal(t, td.value, value)
		})
	}
}

var tests = []ExpectDoc{
	Expect(`query {
		a {
			b(
				x1: val = 1
			) {
				c
				d
			}
		}
	}`, gqt.DocQuery{
		Selections: []gqt.Selection{
			gqt.SelectionField{
				Name: "a",
				Selections: []gqt.Selection{
					gqt.SelectionField{
						Name: "b",
						InputConstraints: []gqt.InputConstraint{{
							Name: "x1",
							Constraint: gqt.ConstraintValEqual{
								Value: int64(1),
							},
						}},
						Selections: []gqt.Selection{
							gqt.SelectionField{
								Name: "c",
							},
							gqt.SelectionField{
								Name: "d",
							},
						},
					},
				},
			},
		},
	}),
	Expect(`query {
		filesystemObject(id: any) {
			name
			... on File {
				format
			}
			... on Directory {
				objects {
					... on File {
						format
					}
				}
			}
		}
	}`, gqt.DocQuery{
		Selections: []gqt.Selection{
			gqt.SelectionField{
				Name: "filesystemObject",
				InputConstraints: []gqt.InputConstraint{{
					Name:       "id",
					Constraint: gqt.ConstraintAny{},
				}},
				Selections: []gqt.Selection{
					gqt.SelectionField{
						Name: "name",
					},
					gqt.SelectionInlineFragment{
						TypeName: "File",
						Selections: []gqt.Selection{
							gqt.SelectionField{
								Name: "format",
							},
						},
					},
					gqt.SelectionInlineFragment{
						TypeName: "Directory",
						Selections: []gqt.Selection{
							gqt.SelectionField{
								Name: "objects",
								Selections: []gqt.Selection{
									gqt.SelectionInlineFragment{
										TypeName: "File",
										Selections: []gqt.Selection{
											gqt.SelectionField{
												Name: "format",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}),
	Expect(`query#comment after token
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
	#comment after token`, gqt.DocQuery{
		Selections: []gqt.Selection{
			gqt.SelectionField{
				Name: "a",
				Selections: []gqt.Selection{
					gqt.SelectionField{
						Name: "b",
						InputConstraints: []gqt.InputConstraint{
							{
								Name: "x1",
								Constraint: gqt.ConstraintValEqual{
									Value: int64(1),
								},
							},
						},
						Selections: []gqt.Selection{
							gqt.SelectionField{
								Name: "c",
							},
							gqt.SelectionField{
								Name: "d",
							},
						},
					},
				},
			},
		},
	}),
	Expect(`mutation {
		a(
			y1: any
			i1: val = 42
			i2: val != 42
			i3: val > 42
			i4: val < 42
			i5: val >= 42
			i6: val <= 42
			i7: val = 42.0
			i8: val != 42.0
			i9: val > 42.0
			i10: val < 42.0
			i11: val >= 42.0
			i12: val <= 42.0
			i13: val = "text"
			i14: val != "text"
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
	}`, gqt.DocMutation{
		Selections: []gqt.Selection{
			gqt.SelectionField{
				Name: "a",
				InputConstraints: []gqt.InputConstraint{{
					Name:       "y1",
					Constraint: gqt.ConstraintAny{},
				}, {
					Name: "i1",
					Constraint: gqt.ConstraintValEqual{
						Value: int64(42),
					},
				}, {
					Name: "i2",
					Constraint: gqt.ConstraintValNotEqual{
						Value: int64(42),
					},
				}, {
					Name: "i3",
					Constraint: gqt.ConstraintValGreater{
						Value: int64(42),
					},
				}, {
					Name: "i4",
					Constraint: gqt.ConstraintValLess{
						Value: int64(42),
					},
				}, {
					Name: "i5",
					Constraint: gqt.ConstraintValGreaterOrEqual{
						Value: int64(42),
					},
				}, {
					Name: "i6",
					Constraint: gqt.ConstraintValLessOrEqual{
						Value: int64(42),
					},
				}, {
					Name: "i7",
					Constraint: gqt.ConstraintValEqual{
						Value: float64(42),
					},
				}, {
					Name: "i8",
					Constraint: gqt.ConstraintValNotEqual{
						Value: float64(42),
					},
				}, {
					Name: "i9",
					Constraint: gqt.ConstraintValGreater{
						Value: float64(42),
					},
				}, {
					Name: "i10",
					Constraint: gqt.ConstraintValLess{
						Value: float64(42),
					},
				}, {
					Name: "i11",
					Constraint: gqt.ConstraintValGreaterOrEqual{
						Value: float64(42),
					},
				}, {
					Name: "i12",
					Constraint: gqt.ConstraintValLessOrEqual{
						Value: float64(42),
					},
				}, {
					Name: "i13",
					Constraint: gqt.ConstraintValEqual{
						Value: "text",
					},
				}, {
					Name: "i14",
					Constraint: gqt.ConstraintValNotEqual{
						Value: "text",
					},
				}, {
					Name: "l1",
					Constraint: gqt.ConstraintLenEqual{
						Value: uint(42),
					},
				}, {
					Name: "l2",
					Constraint: gqt.ConstraintLenNotEqual{
						Value: uint(42),
					},
				}, {
					Name: "l3",
					Constraint: gqt.ConstraintLenGreater{
						Value: uint(42),
					},
				}, {
					Name: "l4",
					Constraint: gqt.ConstraintLenLess{
						Value: uint(42),
					},
				}, {
					Name: "l5",
					Constraint: gqt.ConstraintLenGreaterOrEqual{
						Value: uint(42),
					},
				}, {
					Name: "l6",
					Constraint: gqt.ConstraintLenLessOrEqual{
						Value: uint(42),
					},
				}, {
					Name: "b1",
					Constraint: gqt.ConstraintBytelenEqual{
						Value: uint(42),
					},
				}, {
					Name: "b2",
					Constraint: gqt.ConstraintBytelenNotEqual{
						Value: uint(42),
					},
				}, {
					Name: "b3",
					Constraint: gqt.ConstraintBytelenGreater{
						Value: uint(42),
					},
				}, {
					Name: "b4",
					Constraint: gqt.ConstraintBytelenLess{
						Value: uint(42),
					},
				}, {
					Name: "b5",
					Constraint: gqt.ConstraintBytelenGreaterOrEqual{
						Value: uint(42),
					},
				}, {
					Name: "b6",
					Constraint: gqt.ConstraintBytelenLessOrEqual{
						Value: uint(42),
					},
				}, {
					Name: "a1",
					Constraint: gqt.ConstraintValEqual{
						Value: gqt.ValueArray{},
					},
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
					Constraint: gqt.ConstraintValEqual{
						Value: gqt.ValueObject{
							Fields: []gqt.ObjectField{
								{
									Name:  "of1",
									Value: gqt.ConstraintValEqual{Value: true},
								}, {
									Name:  "of2",
									Value: gqt.ConstraintValEqual{Value: false},
								}, {
									Name: "oo2",
									Value: gqt.ConstraintValEqual{
										Value: gqt.ValueObject{
											Fields: []gqt.ObjectField{
												{
													Name: "oa1",
													Value: gqt.ConstraintValNotEqual{
														Value: nil,
													},
												}, {
													Name:  "oa2",
													Value: gqt.ConstraintAny{},
												},
											},
										},
									},
								},
							},
						},
					},
				}},
				Selections: []gqt.Selection{
					gqt.SelectionField{
						Name: "b",
					},
				},
			},
		},
	}),
	Expect(`mutation {
		a(
			a: val > 0.0 && val < 3.0
			b: val > 0.0 && val < 9.0 && val != 5.0
			c: val = 1.0 || val = 2.0
			d: val = 1.0 || val = 2.0 || val = 3.0
			e: bytelen > 3 && bytelen <= 10 ||
				val = true ||
				val = 1.0
		) {b}
	}`, gqt.DocMutation{
		Selections: []gqt.Selection{
			gqt.SelectionField{
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
							gqt.ConstraintBytelenGreater{Value: uint(3)},
							gqt.ConstraintBytelenLessOrEqual{Value: uint(10)},
						},
						gqt.ConstraintValEqual{Value: true},
						gqt.ConstraintValEqual{Value: float64(1)},
					},
				}},
				Selections: []gqt.Selection{
					gqt.SelectionField{
						Name: "b",
					},
				},
			},
		},
	})}

func TestParse(t *testing.T) {
	for _, td := range tests {
		t.Run(td.Decl, func(t *testing.T) {
			d, err := gqt.Parse([]byte(td.Input))
			require.Zero(t, err.Error())
			require.False(t, err.IsErr())
			require.Equal(t, td.Expect, d)
		})
	}
}

var testsSemanticErr = []ExpectErr{
	// Redundant field selection
	SyntaxErr(
		`query { x x }`,
		"error at 10: redundant field selection",
	),
	// Redundant type condition
	SyntaxErr(
		`query {
			... on T {x}
			... on T {y}
		}`,
		"error at 27: redundant type condition",
	),
	// Redundant constraint
	SyntaxErr(
		`query { x(a: val = 3 a: val = 4) }`,
		"error at 21: redundant constraint",
	),
	// Map multiple constraints
	SyntaxErr(
		`query {x(a: val = [ ... val=0 val=1 ])}`,
		"error at 30: expected right square bracket",
	),
}

func TestSemanticErr(t *testing.T) {
	for _, td := range testsSemanticErr {
		t.Run(td.Decl, func(t *testing.T) {
			d, err := gqt.Parse([]byte(td.Input))
			require.True(t, err.IsErr())
			require.Equal(t, td.ExpectedError, err.Error())
			require.Nil(t, d)
		})
	}
}

var testsSyntaxErr = []ExpectErr{
	SyntaxErr(
		``,
		"error at 0: expected definition",
	),
	SyntaxErr(
		`  `,
		"error at 2: expected definition",
	),
	SyntaxErr(
		`query`,
		"error at 5: expected selection set",
	),
	SyntaxErr(
		`query  `,
		"error at 7: expected selection set",
	),
	SyntaxErr(
		`query{`,
		"error at 6: expected field name",
	),
	SyntaxErr(
		`query{  `,
		"error at 8: expected field name",
	),
	SyntaxErr(
		`query{f`,
		"error at 7: expected field name",
	),
	SyntaxErr(
		`query{f  `,
		"error at 9: expected field name",
	),
	SyntaxErr(
		`query{f(`,
		"error at 8: expected parameter name",
	),
	SyntaxErr(
		`query{f(  `,
		"error at 10: expected parameter name",
	),
	SyntaxErr(
		`query{f(x`,
		"error at 9: expected column after parameter name",
	),
	SyntaxErr(
		`query{f(x  `,
		"error at 11: expected column after parameter name",
	),
	SyntaxErr(
		`query{f(x:`,
		"error at 10: expected constraint subject",
	),
	SyntaxErr(
		`query{f(x:  `,
		"error at 12: expected constraint subject",
	),
	SyntaxErr(
		`query{f(x:val`,
		"error at 13: expected constraint operator",
	),
	SyntaxErr(
		`query{f(x:val  `,
		"error at 15: expected constraint operator",
	),
	SyntaxErr(
		`query{f(x:val=`,
		"error at 14: expected value",
	),
	SyntaxErr(
		`query{f(x:val=  `,
		"error at 16: expected value",
	),
	SyntaxErr(
		`query{f(x:val=1`,
		"error at 15: expected parameter name",
	),
	SyntaxErr(
		`query{f(x:val=1  `,
		"error at 17: expected parameter name",
	),
	SyntaxErr(
		`query{f(x:val=[`,
		"error at 15: expected constraint subject",
	),
	SyntaxErr(
		`query{f(x:val=[.`,
		"error at 15: expected constraint subject",
	),
	SyntaxErr(
		`query{f(x:val=[..`,
		"error at 15: expected constraint subject",
	),
	SyntaxErr(
		`query{f(x:val=[...`,
		"error at 18: expected constraint subject",
	),
	SyntaxErr(
		`query{f(x:val=[...  `,
		"error at 20: expected constraint subject",
	),
	SyntaxErr(
		`query{f{...`,
		"error at 11: expected keyword 'on'",
	),
	SyntaxErr(
		`query{f{... `,
		"error at 12: expected keyword 'on'",
	),
	SyntaxErr(
		`query{f{...on`,
		"error at 13: expected type name",
	),
	SyntaxErr(
		`query{f{...on `,
		"error at 14: expected type name",
	),
	SyntaxErr(
		`query{f{...on T`,
		"error at 15: expected selection set",
	),
	SyntaxErr(
		`query{f{...on T `,
		"error at 16: expected selection set",
	),
	SyntaxErr(
		`query{f{...onT{x}}}`,
		"error at 11: expected keyword 'on'",
	),
}

func TestSyntaxErr(t *testing.T) {
	for _, td := range testsSyntaxErr {
		t.Run(td.Decl, func(t *testing.T) {
			r := require.New(t)
			d, err := gqt.Parse([]byte(td.Input))
			r.True(err.IsErr())
			r.Equal(td.ExpectedError, err.Error())
			r.Nil(d)
		})
	}
}

func decl(skipFrames int) string {
	_, filename, line, _ := runtime.Caller(skipFrames)
	return fmt.Sprintf("%s:%d", filepath.Base(filename), line)
}

type ExpectErr struct {
	Decl          string
	Input         string
	ExpectedError string
}

func SyntaxErr(input, expectedError string) ExpectErr {
	return ExpectErr{
		Decl:          decl(2),
		Input:         input,
		ExpectedError: expectedError,
	}
}

type ExpectDoc struct {
	Decl   string
	Input  string
	Expect gqt.Doc
}

func Expect(input string, doc gqt.Doc) ExpectDoc {
	return ExpectDoc{
		Decl:   decl(2),
		Input:  input,
		Expect: doc,
	}
}
