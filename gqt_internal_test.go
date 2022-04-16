package gqt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestF64ToUint(t *testing.T) {
	for _, td := range []struct {
		input    float64
		expect   uint
		expectOK bool
	}{
		{0, 0, true},
		{1, 1, true},
		{4294967295, 4294967295, true},
		{1234.0, 1234, true},
		{9_007_199_254_740_992, 9_007_199_254_740_992, true},

		{-1, 0, false},
		{-0.1, 0, false},
		{0.1, 0, false},
		{0.00000000001, 0, false},
	} {
		t.Run(fmt.Sprintf("%f_%t", td.input, td.expectOK), func(t *testing.T) {
			ui, ok := f64ToUint(td.input)
			r := require.New(t)
			r.Equal(td.expect, ui)
			r.Equal(td.expectOK, ok)
		})
	}
}

func TestConsumeIrrelevant(t *testing.T) {
	for _, td := range []struct {
		input       source
		expectAfter source
	}{
		{src(""), src("")},
		{src("  "), src("  ", "")},
		{src("xyz"), src("xyz", "xyz")},
		{src(" \t\n\r,xyz"), src(" \t\n\r,xyz", "xyz")},
		{
			src(" \t\n\r,#comment\n\n#another comment \nxyz "),
			src(" \t\n\r,#comment\n\n#another comment \nxyz ", "xyz "),
		},
	} {
		t.Run("", func(t *testing.T) {
			a := td.input.consumeIrrelevant()
			r := require.New(t)
			r.Equal(string(td.expectAfter.original), string(a.original))
			r.Equal(string(td.expectAfter.s), string(a.s))
			r.Equal(td.expectAfter.index(), a.index())
		})
	}
}

func src(s ...string) source {
	if len(s) < 1 || len(s) > 2 {
		panic(fmt.Errorf("invalid inputs: %#v", s))
	}
	o := []byte(s[0])
	src := source{original: o, s: o}
	if len(s) > 1 {
		src.s = []byte(s[1])
	}
	return src
}

func TestConsume(t *testing.T) {
	for _, td := range []struct {
		input       source
		consume     string
		expectOK    bool
		expectAfter source
	}{
		{src(""), "", true, src("")},
		{src("x"), "x", true, src("x", "")},
		{src("xyz"), "xyz", true, src("xyz", "")},
		{src("abc"), "xyz", false, src("abc")},
	} {
		t.Run("", func(t *testing.T) {
			a, ok := td.input.consume([]byte(td.consume))
			r := require.New(t)
			r.Equal(td.expectOK, ok)
			r.Equal(string(td.expectAfter.original), string(a.original))
			r.Equal(string(td.expectAfter.s), string(a.s))
			r.Equal(td.expectAfter.index(), a.index())
		})
	}
}

func TestConsumeToken(t *testing.T) {
	for _, td := range []struct {
		input       source
		expectToken string
		expectAfter source
	}{
		{src(""), "", src("")},
		{src("x"), "x", src("x", "")},
		{src("xyz "), "xyz", src("xyz ", " ")},
		{src("xyz\n"), "xyz", src("xyz\n", "\n")},
		{src("xyz\r"), "xyz", src("xyz\r", "\r")},
		{src("xyz\t"), "xyz", src("xyz\t", "\t")},
		{src("xyz,"), "xyz", src("xyz,", ",")},
		{src("xyz{"), "xyz", src("xyz{", "{")},
		{src("xyz("), "xyz", src("xyz(", "(")},
		{
			src("abc12_3456789#"),
			"abc12_3456789",
			src("abc12_3456789#", "#"),
		},
	} {
		t.Run("", func(t *testing.T) {
			a, k := td.input.consumeToken()
			r := require.New(t)
			r.Equal(td.expectToken, string(k))
			r.Equal(string(td.expectAfter.original), string(a.original))
			r.Equal(string(td.expectAfter.s), string(a.s))
			r.Equal(td.expectAfter.index(), a.index())
		})
	}
}

func TestConsumeNumber(t *testing.T) {
	for _, td := range []struct {
		input       source
		expectNum   float64
		expectOK    bool
		expectAfter source
	}{
		{src(""), 0, false, src("")},
		{src("x"), 0, false, src("x")},
		{src("0"), 0, true, src("0", "")},
		{src("10"), 10, true, src("10", "")},
		{src("-1"), -1, true, src("-1", "")},
		{src("0.1234"), 0.1234, true, src("0.1234", "")},
		{src("-0.1234"), -0.1234, true, src("-0.1234", "")},
		{src("-0.1234{x"), -0.1234, true, src("-0.1234{x", "{x")},
		{src("-0.1234(x"), -0.1234, true, src("-0.1234(x", "(x")},
		{src("-0.1234#x"), -0.1234, true, src("-0.1234#x", "#x")},
		{src("-0.1234,x"), -0.1234, true, src("-0.1234,x", ",x")},
		{src("-0.1234 x"), -0.1234, true, src("-0.1234 x", " x")},
		{src("-0.1234\tx"), -0.1234, true, src("-0.1234\tx", "\tx")},
		{src("-0.1234\nx"), -0.1234, true, src("-0.1234\nx", "\nx")},
		{src("-0.1234\rx"), -0.1234, true, src("-0.1234\rx", "\rx")},
		{src("0.1234x"), 0, false, src("0.1234x")},
	} {
		t.Run("", func(t *testing.T) {
			a, f, ok := td.input.consumeNumber()
			r := require.New(t)
			r.Equal(td.expectOK, ok)
			r.Equal(td.expectNum, f)
			r.Equal(string(td.expectAfter.original), string(a.original))
			r.Equal(string(td.expectAfter.s), string(a.s))
			r.Equal(td.expectAfter.index(), a.index())
		})
	}
}

func TestConsumeName(t *testing.T) {
	for _, td := range []struct {
		input       source
		expectName  string
		expectAfter source
	}{
		{src(""), "", src("", "")},
		{src("x"), "x", src("x", "")},
		{src("xyz"), "xyz", src("xyz", "")},
		{src("xyz "), "xyz", src("xyz ", " ")},
		{src("xyz("), "xyz", src("xyz(", "(")},
		{src("xyz{"), "xyz", src("xyz{", "{")},
		{src("xyz\n"), "xyz", src("xyz\n", "\n")},
		{src("xyz\t"), "xyz", src("xyz\t", "\t")},
		{src("xyz\r"), "xyz", src("xyz\r", "\r")},
		{src("xyz,"), "xyz", src("xyz,", ",")},
		{src("xyz#"), "xyz", src("xyz#", "#")},
	} {
		t.Run("", func(t *testing.T) {
			a, n := td.input.consumeName()
			r := require.New(t)
			r.Equal(td.expectName, string(n))
			r.Equal(string(td.expectAfter.original), string(a.original))
			r.Equal(string(td.expectAfter.s), string(a.s))
			r.Equal(td.expectAfter.index(), a.index())
		})
	}
}

func TestConsumeString(t *testing.T) {
	for _, td := range []struct {
		input       source
		expectName  string
		expectOK    bool
		expectAfter source
	}{
		{src(""), "", false, src("")},
		{src(`""`), "", true, src(`""`, "")},
		{src(`"",`), "", true, src(`"",`, ",")},
		{src(`"x" `), "x", true, src(`"x" `, " ")},
		{
			src(`"this is a string" `),
			"this is a string",
			true,
			src(`"this is a string" `, " "),
		},
		{src(`"\\" `), `\\`, true, src(`"\\" `, " ")},

		{src(`"`), "", false, src(`"`)},
		{src(`" `), "", false, src(`" `)},
		{src(`"\" `), "", false, src(`"\" `)},
	} {
		t.Run("", func(t *testing.T) {
			a, n, ok := td.input.consumeString()
			r := require.New(t)
			r.Equal(td.expectOK, ok)
			r.Equal(td.expectName, string(n))
			r.Equal(string(td.expectAfter.original), string(a.original))
			r.Equal(string(td.expectAfter.s), string(a.s))
			r.Equal(td.expectAfter.index(), a.index())
		})
	}
}
