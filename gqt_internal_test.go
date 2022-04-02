package gqt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for _, td := range []struct {
		input    float64
		expect   uint
		expectOK bool
	}{
		{0, 0, true},
		{1, 1, true},
		{4294967295, 4294967295, true},
		{1234.0, 1234, true},
		{18446744073709551615, 18446744073709551615, true},

		{-1, 0, false},
		{-0.1, 0, false},
		{0.1, 0, false},
		{0.00000000001, 0, false},
	} {
		t.Run(fmt.Sprintf("%f_%t", td.input, td.expectOK), func(t *testing.T) {
			ui, ok := f64ToUint(td.input)
			require.Equal(t, td.expect, ui)
			require.Equal(t, td.expectOK, ok)
		})
	}
}

func TestSkipIrrelevant(t *testing.T) {
	src := func(text string, index int) source {
		return source{s: []byte(text), index: index}
	}

	for _, td := range []struct {
		input       source
		expectAfter source
	}{
		{src("", 0), src("", 0)},
		{src("xyz", 0), src("xyz", 0)},
		{src(" \t\n\r,xyz", 0), src("xyz", 5)},
		{
			src(" \t\n\r,#comment\n\n#another comment \nxyz", 0),
			src("xyz", len(" \t\n\r,#comment\n\n#another comment \n")),
		},
	} {
		t.Run("", func(t *testing.T) {
			a := td.input.skipIrrelevant()
			require.Equal(t, string(td.expectAfter.s), string(a.s))
			require.Equal(t, td.expectAfter.index, a.index)
		})
	}
}

func TestConsume(t *testing.T) {
	src := func(text string, index int) source {
		return source{s: []byte(text), index: index}
	}

	for _, td := range []struct {
		input       source
		consume     string
		expectOK    bool
		expectAfter source
	}{
		{src("", 0), "", true, src("", 0)},
		{src("x", 0), "x", true, src("", len("x"))},
		{src("xyz", 0), "xyz", true, src("", len("xyz"))},
		{src("abc", 0), "xyz", false, src("abc", 0)},
	} {
		t.Run("", func(t *testing.T) {
			a, ok := td.input.consume([]byte(td.consume))
			require.Equal(t, td.expectOK, ok)
			require.Equal(t, string(td.expectAfter.s), string(a.s))
			require.Equal(t, td.expectAfter.index, a.index)
		})
	}
}

func TestConsumeToken(t *testing.T) {
	src := func(text string, index int) source {
		return source{s: []byte(text), index: index}
	}

	for _, td := range []struct {
		input       source
		expectToken string
		expectAfter source
	}{
		{src("", 0), "", src("", 0)},
		{src("x", 0), "x", src("", len("x"))},
		{src("xyz ", 0), "xyz", src(" ", len("xyz"))},
		{src("xyz\n", 0), "xyz", src("\n", len("xyz"))},
		{src("xyz\r", 0), "xyz", src("\r", len("xyz"))},
		{src("xyz\t", 0), "xyz", src("\t", len("xyz"))},
		{src("xyz,", 0), "xyz", src(",", len("xyz"))},
		{src("xyz{", 0), "xyz", src("{", len("xyz"))},
		{src("xyz(", 0), "xyz", src("(", len("xyz"))},
		{src("abc12_3456789#", 0), "abc12_3456789", src("#", len("abc12_3456789"))},
	} {
		t.Run("", func(t *testing.T) {
			a, k := td.input.consumeToken()
			require.Equal(t, td.expectToken, string(k))
			require.Equal(t, string(td.expectAfter.s), string(a.s))
			require.Equal(t, td.expectAfter.index, a.index)
		})
	}
}

func TestConsumeNumber(t *testing.T) {
	src := func(text string, index int) source {
		return source{s: []byte(text), index: index}
	}

	for _, td := range []struct {
		input       source
		expectNum   float64
		expectOK    bool
		expectAfter source
	}{
		{src("", 0), 0, false, src("", 0)},
		{src("x", 0), 0, false, src("x", 0)},
		{src("0", 0), 0, true, src("", len("0"))},
		{src("10", 0), 10, true, src("", len("10"))},
		{src("-1", 0), -1, true, src("", len("-1"))},
		{src("0.1234", 0), 0.1234, true, src("", len("0.1234"))},
		{src("-0.1234", 0), -0.1234, true, src("", len("-0.1234"))},
		{src("-0.1234{x", 0), -0.1234, true, src("{x", len("-0.1234"))},
		{src("-0.1234(x", 0), -0.1234, true, src("(x", len("-0.1234"))},
		{src("-0.1234#x", 0), -0.1234, true, src("#x", len("-0.1234"))},
		{src("-0.1234,x", 0), -0.1234, true, src(",x", len("-0.1234"))},
		{src("-0.1234 x", 0), -0.1234, true, src(" x", len("-0.1234"))},
		{src("-0.1234\tx", 0), -0.1234, true, src("\tx", len("-0.1234"))},
		{src("-0.1234\nx", 0), -0.1234, true, src("\nx", len("-0.1234"))},
		{src("-0.1234\rx", 0), -0.1234, true, src("\rx", len("-0.1234"))},
		{src("0.1234x", 0), 0, false, src("0.1234x", 0)},
	} {
		t.Run("", func(t *testing.T) {
			a, f, ok := td.input.consumeNumber()
			require.Equal(t, td.expectOK, ok)
			require.Equal(t, td.expectNum, f)
			require.Equal(t, string(td.expectAfter.s), string(a.s))
			require.Equal(t, td.expectAfter.index, a.index)
		})
	}
}

func TestConsumeName(t *testing.T) {
	src := func(text string, index int) source {
		return source{s: []byte(text), index: index}
	}

	for _, td := range []struct {
		input       source
		expectName  string
		expectAfter source
	}{
		{src("", 0), "", src("", 0)},
		{src("x", 0), "x", src("", len("x"))},
		{src("xyz", 0), "xyz", src("", len("xyz"))},
		{src("xyz ", 0), "xyz", src(" ", len("xyz"))},
		{src("xyz(", 0), "xyz", src("(", len("xyz"))},
		{src("xyz{", 0), "xyz", src("{", len("xyz"))},
		{src("xyz\n", 0), "xyz", src("\n", len("xyz"))},
		{src("xyz\t", 0), "xyz", src("\t", len("xyz"))},
		{src("xyz\r", 0), "xyz", src("\r", len("xyz"))},
		{src("xyz,", 0), "xyz", src(",", len("xyz"))},
		{src("xyz#", 0), "xyz", src("#", len("xyz"))},
	} {
		t.Run("", func(t *testing.T) {
			a, n := td.input.consumeName()
			require.Equal(t, td.expectName, string(n))
			require.Equal(t, string(td.expectAfter.s), string(a.s))
			require.Equal(t, td.expectAfter.index, a.index)
		})
	}
}

func TestConsumeString(t *testing.T) {
	src := func(text string, index int) source {
		return source{s: []byte(text), index: index}
	}

	for _, td := range []struct {
		input       source
		expectName  string
		expectOK    bool
		expectAfter source
	}{
		{src("", 0), "", false, src("", 0)},
		{src(`""`, 0), "", true, src("", len(`""`))},
		{src(`"",`, 0), "", true, src(",", len(`""`))},
		{src(`"x" `, 0), "x", true, src(" ", len(`"x"`))},
		{
			src(`"this is a string" `, 0),
			"this is a string",
			true,
			src(" ", len(`"this is a string"`)),
		},
		{src(`"\\" `, 0), `\\`, true, src(` `, len(`"\\"`))},

		{src(`"`, 0), "", false, src(``, len(`"`))},
		{src(`" `, 0), "", false, src(``, len(`" `))},
		{src(`"\" `, 0), "", false, src(``, len(`"\" `))},
	} {
		t.Run("", func(t *testing.T) {
			a, n, ok := td.input.consumeString()
			require.Equal(t, td.expectOK, ok)
			require.Equal(t, td.expectName, string(n))
			require.Equal(t, string(td.expectAfter.s), string(a.s))
			require.Equal(t, td.expectAfter.index, a.index)
		})
	}
}
