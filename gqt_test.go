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
			d, vars, err := gqt.Parse([]byte(input))
			require.Equal(
				t, expectation, err.Error(),
				"input: %q", input,
			)
			require.Zero(t, d)
			require.Nil(t, vars)
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

			o, vars, err := gqt.Parse([]byte(input))
			require.False(t, err.IsErr(), "unexpected error: %v", err)
			require.NotNil(t, vars)

			{
				var b bytes.Buffer
				_, err := gqt.WriteYAML(&b, o)
				require.NoError(t, err)
				require.Equal(t, expectation, b.String())
			}
		},
	)
}

func TestParseEmpty(t *testing.T) {
	input := `query{x}`
	opr, vars, err := gqt.Parse([]byte(input))
	require.False(t, err.IsErr(), "unexpected error: %v", err)
	require.NotNil(t, opr)
	require.Len(t, vars, 0)
}

func TestParseVariables(t *testing.T) {
	input := `query {
		f1(a: $b+$x, c=$c: $b) {
			f2(b=$b, x=$unused: $c+$b, d: {o1=$x})
		}
	}`
	opr, vars, err := gqt.Parse([]byte(input))
	require.False(t, err.IsErr(), "unexpected error: %v", err)
	require.NotNil(t, opr)

	require.Len(t, vars, 4)
	require.Contains(t, vars, "b")
	require.Contains(t, vars, "c")
	require.Contains(t, vars, "x")
	require.Contains(t, vars, "unused")

	// Check names
	require.Equal(t, "b", vars["b"].Name)
	require.Equal(t, "c", vars["c"].Name)
	require.Equal(t, "x", vars["x"].Name)
	require.Equal(t, "unused", vars["unused"].Name)

	// Check parents
	require.IsType(t, &gqt.Argument{}, vars["b"].Parent)
	require.IsType(t, &gqt.Argument{}, vars["c"].Parent)
	require.IsType(t, &gqt.Argument{}, vars["unused"].Parent)
	require.IsType(t, &gqt.ObjectField{}, vars["x"].Parent)

	require.Equal(t, "b", vars["b"].Parent.(*gqt.Argument).Name)
	require.Equal(t, "c", vars["c"].Parent.(*gqt.Argument).Name)
	require.Equal(t, "x", vars["unused"].Parent.(*gqt.Argument).Name)
	require.Equal(t, "o1", vars["x"].Parent.(*gqt.ObjectField).Name)

	// Check references
	require.Len(t, vars["b"].References, 3)
	require.Len(t, vars["c"].References, 1)
	require.Len(t, vars["unused"].References, 0)
	require.Len(t, vars["x"].References, 1)

	require.Equal(
		t, gqt.Location{Index: 16, Line: 2, Column: 9},
		vars["b"].References[0].(*gqt.Variable).Location,
	)
	require.Equal(
		t, gqt.Location{Index: 29, Line: 2, Column: 22},
		vars["b"].References[1].(*gqt.Variable).Location,
	)
	require.Equal(
		t, gqt.Location{Index: 61, Line: 3, Column: 27},
		vars["b"].References[2].(*gqt.Variable).Location,
	)

	require.Equal(
		t, gqt.Location{Index: 58, Line: 3, Column: 24},
		vars["c"].References[0].(*gqt.Variable).Location,
	)

	require.Equal(
		t, gqt.Location{Index: 19, Line: 2, Column: 12},
		vars["x"].References[0].(*gqt.Variable).Location,
	)
}
