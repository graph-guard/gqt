package gqt_test

import (
	"bytes"
	"embed"
	"testing"

	"github.com/graph-guard/gqt"
	"github.com/graph-guard/gqt/internal/test"

	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

//go:embed test/error
var errorFS embed.FS

//go:embed test/ast
var astFS embed.FS

func TestParseErr(t *testing.T) {
	test.ExecDirMD(
		t, errorFS, "test/error", "ERR: ",
		func(t *testing.T, input, expectation string) {
			d, err := gqt.Parse([]byte(input))
			require.Equal(
				t, expectation, err.Error(),
				"input: %q", input,
			)
			require.Zero(t, d)
		},
	)
}

func TestParse(t *testing.T) {
	test.ExecDirMD(
		t, astFS, "test/ast", "ERR: ",
		func(t *testing.T, input, expectation string) {
			var discard map[string]any
			require.NoError(
				t, yaml.Unmarshal([]byte(expectation), &discard),
				"invalid expectation YAML",
			)

			o, err := gqt.Parse([]byte(input))
			require.False(t, err.IsErr(), "unexpected error: %v", err)

			{
				var b bytes.Buffer
				_, err := gqt.WriteYAML(&b, o)
				require.NoError(t, err)
				require.Equal(t, expectation, b.String())
			}
		},
	)
}
