package gqt

import (
	"testing"

	"github.com/graph-guard/gqt/internal/test"

	"github.com/stretchr/testify/require"
)

func TestConsumeIgnored(t *testing.T) {
	type T struct {
		Input  string
		Expect Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a := source{
			s:        []byte(x.Input),
			Location: startLoc(),
		}.consumeIgnored()
		expect := source{
			s:        []byte(x.Input),
			Location: x.Expect,
		}
		require.Equal(t, expect, a)
	})

	run(T{"", Location{0, 1, 1}})
	run(T{"  ", Location{2, 1, 3}})
	run(T{"xyz", Location{0, 1, 1}})
	run(T{"#xyz", Location{4, 1, 5}})
	run(T{"#xyz\n", Location{5, 2, 1}})
	run(T{"#xyz\nabc", Location{5, 2, 1}})
	run(T{" \t\n\rxyz", Location{4, 2, 2}})
	run(T{" \t\n\r#comment\n\n#another comment \nxyz ", Location{32, 5, 1}})
	run(T{string([]byte{0x00, '1', '2', '3'}), Location{0, 1, 1}})
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
			Location: startLoc(),
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

func TestConsumeNumber(t *testing.T) {
	type T struct {
		Location    Location
		Input       string
		ExpectOK    bool
		ExpectNum   any
		ExpectAfter Location
	}

	run := test.New(t, func(t *testing.T, x T) {
		a, n, ok := source{
			s:        []byte(x.Input),
			Location: x.Location,
		}.consumeNumber()
		require.Equal(t, x.ExpectOK, ok)
		require.Equal(t, x.ExpectNum, n)
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
		nil,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"x",
		false,
		nil,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"0.1234x",
		false,
		nil,
		Location{0, 1, 1},
	})
	run(T{
		startLoc(),
		"0",
		true,
		&ValueInt{
			Location: startLoc(),
			Value:    0,
		},
		Location{1, 1, 2},
	})
	run(T{
		startLoc(),
		"10.0",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    10.0,
		},
		Location{4, 1, 5},
	})
	run(T{
		startLoc(),
		"-1",
		true,
		&ValueInt{
			Location: startLoc(),
			Value:    -1,
		},
		Location{2, 1, 3},
	})
	run(T{
		startLoc(),
		"0.1234",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    0.1234,
		},
		Location{6, 1, 7},
	})
	run(T{
		startLoc(),
		"-0.1234",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234{x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234}x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234(x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234)x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234[x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234]x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234#x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234,x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234 x",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234\tx",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234\nx",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
	})
	run(T{
		startLoc(),
		"-0.1234\rx",
		true,
		&ValueFloat{
			Location: startLoc(),
			Value:    -0.1234,
		},
		Location{7, 1, 8},
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
