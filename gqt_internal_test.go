package gqt

import (
	"testing"

	"github.com/graph-guard/gqt/v3/internal/test"

	"github.com/stretchr/testify/require"
)

func TestConsumeIgnored(t *testing.T) {
	type T struct {
		Location Location
		Input    string
		Expect   Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a := source{
			s:        []byte(x.Input),
			Location: x.Location,
		}.consumeIgnored()
		expect := source{
			s:        []byte(x.Input),
			Location: x.Expect,
		}
		require.Equal(t, expect, a)
	})

	run(T{
		startLoc(),
		"",
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"  ",
		Location{2, 1, 3},
	})
	run(T{
		startLoc(),
		"xyz",
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"#xyz",
		Location{4, 1, 5},
	})
	run(T{
		startLoc(),
		"#xyz\n",
		Location{5, 2, 1},
	})
	run(T{
		startLoc(),
		"#xyz\nabc",
		Location{5, 2, 1},
	})
	run(T{
		startLoc(),
		" \t\n\rxyz",
		Location{4, 2, 2},
	})
	run(T{
		startLoc(),
		" \t\n\r#comment\n\n#another comment \nxyz ",
		Location{32, 5, 1},
	})
	run(T{
		startLoc(),
		string([]byte{0x00, '1', '2', '3'}),
		Location{0, 1, 1},
	})
	run(T{
		Location{14, 4, 1},
		" \t\n\r#comment\n\n#another comment \nxyz ",
		Location{32, 5, 1},
	})
}

func TestConsume(t *testing.T) {
	type T struct {
		Location    Location
		Input       string
		Consume     string
		ExpectOK    bool
		ExpectAfter Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a, ok := source{
			s:        []byte(x.Input),
			Location: x.Location,
		}.consume(x.Consume)
		expect := source{
			s:        []byte(x.Input),
			Location: x.ExpectAfter,
		}
		require.Equal(t, x.ExpectOK, ok)
		if ok {
			require.Equal(t, expect, a)
		}
	})

	run(T{
		startLoc(),
		"",
		"",
		true,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"",
		"a",
		false,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"a",
		"ab",
		false,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"x",
		"x",
		true,
		Location{1, 1, 2},
	})
	run(T{
		startLoc(),
		"xyz",
		"xyz",
		true,
		Location{3, 1, 4},
	})
	run(T{
		startLoc(),
		"abc",
		"xyz",
		false,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"aaa",
		"aab",
		false,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"abcdef",
		"abc",
		true,
		Location{3, 1, 4},
	})
	run(T{
		startLoc(),
		"abc\ndef",
		"abc\n",
		true,
		Location{4, 2, 1},
	})
	run(T{
		Location{4, 2, 1},
		"abc\ndefg",
		"def",
		true,
		Location{7, 2, 4},
	})
}

func TestConsumeEitherOf3(t *testing.T) {
	type T struct {
		Location       Location
		Input          string
		Consume        [3]string
		ExpectSelected int
		ExpectAfter    Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a, selected := source{
			s:        []byte(x.Input),
			Location: x.Location,
		}.consumeEitherOf3(x.Consume[0], x.Consume[1], x.Consume[2])
		expect := source{s: []byte(x.Input), Location: x.ExpectAfter}
		require.Equal(t, x.ExpectSelected, selected)
		require.Equal(t, expect, a)
	})

	run(T{
		startLoc(),
		"",
		[3]string{"", "", ""},
		-1,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"",
		[3]string{"a", "", ""},
		-1,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"",
		[3]string{"a", "b", "c"},
		-1,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"a",
		[3]string{"a", "", ""},
		0,
		Location{1, 1, 2},
	})
	run(T{
		startLoc(),
		"apple",
		[3]string{"applepie", "honey", "smoke"},
		-1,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"apple",
		[3]string{"tomato", "applepie", "app"},
		2,
		Location{3, 1, 4},
	})
	run(T{
		startLoc(),
		"\n\n\n",
		[3]string{"\n", "\n\n", "\n\n\n"},
		0,
		Location{1, 2, 1},
	})
	run(T{
		startLoc(),
		"\n\n\n",
		[3]string{"", "\n", ""},
		1,
		Location{1, 2, 1},
	})
	run(T{
		Location{4, 3, 2},
		"\nx\nx\nx",
		[3]string{"", "", "\nx"},
		2,
		Location{6, 4, 2},
	})
}

func TestConsumeToken(t *testing.T) {
	type T struct {
		InputLocation Location
		Input         string
		ExpectToken   string
		ExpectAfter   Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a, token := source{
			s:        []byte(x.Input),
			Location: x.InputLocation,
		}.consumeToken()
		r := require.New(t)
		r.Equal(x.ExpectToken, string(token))
		expect := source{
			s:        []byte(x.Input),
			Location: x.ExpectAfter,
		}
		require.Equal(t, expect, a)
	})

	run(T{
		Location{Line: 1, Column: 1},
		"", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"x", "x", Location{1, 1, 2},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz ", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz\n", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz\r", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz\t", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz,abc", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz{foo", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz(bar", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"123[\"baz\"", "123", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz}", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"xyz)", "xyz", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"123]", "123", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"123#asd", "123", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"abc12_3456789#", "abc12_3456789", Location{13, 1, 14},
	})
	run(T{
		Location{Line: 1, Column: 1},
		string([]byte{'1', '2', '3', 0x00}), "123", Location{3, 1, 4},
	})
	run(T{
		Location{Line: 1, Column: 1},
		string([]byte{0x00, '1', '2', '3'}), "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		" xyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"\nxyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"\rxyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"\txyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		",abc", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"{foo", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"(bar", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"[\"baz\"", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"}xyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		")xyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"]xyz", "", Location{0, 1, 1},
	})
	run(T{
		Location{Line: 1, Column: 1},
		"#asd", "", Location{0, 1, 1},
	})
	run(T{
		Location{Index: 4, Line: 2, Column: 3},
		"a\nb asd ", "asd", Location{7, 2, 6},
	})

}

func TestParseNumber(t *testing.T) {
	type T struct {
		Location    Location
		Input       string
		ExpectNum   any
		ExpectErr   []Error
		ExpectAfter Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		p := &Parser{}
		a, n := p.parseNumber(source{
			s:        []byte(x.Input),
			Location: x.Location,
		})
		require.Equal(t, x.ExpectNum, n)
		if len(x.ExpectErr) > 0 {
			require.Equal(t, x.ExpectErr, p.errors)
			require.Equal(t, source{}, a)
		} else {
			require.Len(t, p.errors, 0)
			require.Equal(t, source{
				s:        []byte(x.Input),
				Location: x.ExpectAfter,
			}, a)
		}

	})

	run(T{
		startLoc(),
		"",
		nil,
		nil,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"x",
		nil,
		nil,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"0.1234x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0",
		&Int{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{1, 1, 2},
			},
			Value: 0,
		},
		nil,
		Location{1, 1, 2},
	})
	run(T{
		startLoc(),
		"10.0",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{4, 1, 5},
			},
			Value: 10.0,
		},
		nil,
		Location{4, 1, 5},
	})
	run(T{
		startLoc(),
		"0.1234",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"6E+5",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{4, 1, 5},
			},
			Value: 6e+5,
		},
		nil,
		Location{4, 1, 5},
	})
	run(T{
		startLoc(),
		"6E",
		nil,
		[]Error{
			{
				LocRange: LocRange{
					Location:    startLoc(),
					LocationEnd: LocationEnd{2, 1, 3},
				},
				Msg: "exponent has no digits",
			},
		},
		Location{},
	})
	run(T{
		startLoc(),
		"6E ",
		nil,
		[]Error{
			{
				LocRange: LocRange{
					Location:    startLoc(),
					LocationEnd: LocationEnd{2, 1, 3},
				},
				Msg: "exponent has no digits",
			},
		},
		Location{},
	})
	run(T{
		startLoc(),
		"6Ee",
		nil,
		[]Error{
			{
				LocRange: LocRange{
					Location:    startLoc(),
					LocationEnd: LocationEnd{2, 1, 3},
				},
				Msg: "exponent has no digits",
			},
		},
		startLoc(),
	})
	run(T{
		startLoc(),
		"6E-",
		nil,
		[]Error{
			{
				LocRange: LocRange{
					Location:    startLoc(),
					LocationEnd: LocationEnd{3, 1, 4},
				},
				Msg: "exponent has no digits",
			},
		},
		startLoc(),
	})
	run(T{
		startLoc(),
		"6E- ",
		nil,
		[]Error{
			{
				LocRange: LocRange{
					Location:    startLoc(),
					LocationEnd: LocationEnd{3, 1, 4},
				},
				Msg: "exponent has no digits",
			},
		},
		startLoc(),
	})
	run(T{
		startLoc(),
		"6E-e",
		nil,
		[]Error{
			{
				LocRange: LocRange{
					Location:    startLoc(),
					LocationEnd: LocationEnd{3, 1, 4},
				},
				Msg: "exponent has no digits",
			},
		},
		startLoc(),
	})
	run(T{
		startLoc(),
		"0.1234}x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234)x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234]x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234||",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234&&",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234#x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234,x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234 x",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234\tx",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234\nx",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"0.1234\rx",
		&Float{
			LocRange: LocRange{
				Location:    startLoc(),
				LocationEnd: LocationEnd{6, 1, 7},
			},
			Value: 0.1234,
		},
		nil,
		Location{6, 1, 7},
	})
	run(T{
		Location{2, 1, 3},
		"  2-321 ",
		&Int{
			LocRange: LocRange{
				Location:    Location{2, 1, 3},
				LocationEnd: LocationEnd{3, 1, 4},
			},
			Value: 2,
		},
		nil,
		Location{3, 1, 4},
	})
	run(T{
		Location{2, 1, 3},
		"  2+321 ",
		&Int{
			LocRange: LocRange{
				Location:    Location{2, 1, 3},
				LocationEnd: LocationEnd{3, 1, 4},
			},
			Value: 2,
		},
		nil,
		Location{3, 1, 4},
	})
	run(T{
		Location{2, 1, 3},
		"  2*321",
		&Int{
			LocRange: LocRange{
				Location:    Location{2, 1, 3},
				LocationEnd: LocationEnd{3, 1, 4},
			},
			Value: 2,
		},
		nil,
		Location{3, 1, 4},
	})
	run(T{
		Location{2, 1, 3},
		"  2/321",
		&Int{
			LocRange: LocRange{
				Location:    Location{2, 1, 3},
				LocationEnd: LocationEnd{3, 1, 4},
			},
			Value: 2,
		},
		nil,
		Location{3, 1, 4},
	})
	run(T{
		Location{2, 1, 3},
		"  2%321",
		&Int{
			LocRange: LocRange{
				Location:    Location{2, 1, 3},
				LocationEnd: LocationEnd{3, 1, 4},
			},
			Value: 2,
		},
		nil,
		Location{3, 1, 4},
	})
}

func TestConsumeName(t *testing.T) {
	type T struct {
		Location    Location
		Input       string
		ExpectName  string
		ExpectAfter Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a, n := source{
			s:        []byte(x.Input),
			Location: x.Location,
		}.consumeName()
		require.Equal(t, x.ExpectName, string(n))
		expect := source{
			s:        []byte(x.Input),
			Location: x.ExpectAfter,
		}
		require.Equal(t, expect, a)
	})

	run(T{startLoc(), "", "", Location{0, 1, 1}})
	run(T{startLoc(), "x", "x", Location{1, 1, 2}})
	run(T{startLoc(), "xyz", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz ", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz(", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz)", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz{", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz}", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz[", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz]", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz\n", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz\t", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz\r", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz,", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz#", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "xyz#abc", "xyz", Location{3, 1, 4}})
	run(T{startLoc(), "1name#abc", "", Location{0, 1, 1}})
	run(T{Location{2, 1, 3}, "{ f(a) }", "f", Location{3, 1, 4}})
}

func TestConsumeString(t *testing.T) {
	type T struct {
		Location     Location
		Input        string
		ExpectOK     bool
		ExpectString string
		ExpectAfter  Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a, n, ok := source{
			s:        []byte(x.Input),
			Location: x.Location,
		}.consumeString()
		require.Equal(t, x.ExpectOK, ok)
		require.Equal(t, x.ExpectString, string(n))
		expect := source{
			s:        []byte(x.Input),
			Location: x.ExpectAfter,
		}
		require.Equal(t, expect, a)
	})

	run(T{
		startLoc(),
		"",
		false,
		"",
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		`""`,
		true,
		"",
		Location{2, 1, 3},
	})
	run(T{
		startLoc(),
		`"",`,
		true,
		"",
		Location{2, 1, 3},
	})
	run(T{
		startLoc(),
		`"x" `,
		true,
		"x",
		Location{3, 1, 4},
	})
	run(T{
		startLoc(),
		`"this is a string" `,
		true,
		"this is a string",
		Location{18, 1, 19},
	})
	run(T{
		startLoc(),
		`"\\" `,
		true,
		`\\`,
		Location{4, 1, 5},
	})
	run(T{
		startLoc(),
		`"`,
		false,
		"",
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		`" `,
		false,
		"",
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		`"\" `,
		false,
		"",
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		string([]byte{'"', 0x00, '"'}),
		false,
		"",
		Location{0, 1, 1},
	})
	run(T{
		Location{2, 1, 3},
		"  \"okay\"abc",
		true,
		"okay",
		Location{8, 1, 9},
	})
}

// startLoc returns a location at the start of a file.
func startLoc() Location { return Location{0, 1, 1} }
