package gqt_test

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/graph-guard/gqt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

//go:embed tests
var testsFS embed.FS

func TestParse(t *testing.T) {
	type T struct {
		Schema                 string         `yaml:"schema"`
		Template               string         `yaml:"template"`
		ExpectAST              map[string]any `yaml:"expect-ast"`
		ExpectASTSchemaless    map[string]any `yaml:"expect-ast(schemaless)"`
		ExpectErrors           []string       `yaml:"expect-errors"`
		ExpectErrorsSchemaless []string       `yaml:"expect-errors(schemaless)"`
	}

	d, err := fs.ReadDir(testsFS, "tests")
	require.NoError(t, err)

	for _, do := range d {
		fileName := do.Name()
		if do.IsDir() {
			t.Run(fileName, func(t *testing.T) {
				t.Skipf("ignoring directory %q", fileName)
			})
			continue
		}
		if !strings.HasSuffix(fileName, ".yml") {
			t.Run(fileName, func(t *testing.T) {
				t.Skipf("ignoring file %q", fileName)
			})
			continue
		}
		f, err := testsFS.ReadFile(filepath.Join("tests", fileName))
		require.NoError(t, err, "reading YAML test file")
		t.Run(fileName[:len(fileName)-len(".yml")], func(t *testing.T) {
			var ts T
			{
				d := yaml.NewDecoder(bytes.NewReader(f))
				d.KnownFields(true)
				if err := d.Decode(&ts); err != nil {
					t.Fatal("parsing YAML test definition", err)
				}
			}
			if ts.ExpectAST != nil && ts.ExpectErrors != nil {
				t.Fatal("expecting both AST and errors in schema-aware mode")
			} else if ts.ExpectASTSchemaless != nil &&
				ts.ExpectErrorsSchemaless != nil {
				t.Fatal("expecting both AST and errors in schema-less mode")
			}
			t.Run("schema", func(t *testing.T) {
				p, err := gqt.NewParser([]gqt.Source{
					{Name: "schema.graphqls", Content: ts.Schema},
				})
				require.NoError(t, err, "unexpected error while parsing schema")
				opr, vars, errs := p.Parse([]byte(ts.Template))
				compareErrors(t, ts.ExpectErrors, errs)
				if len(ts.ExpectErrors) > 0 {
					// Expect failure
					require.Zero(t, vars)
					require.Zero(t, opr)
				} else {
					// Expect success
					var j bytes.Buffer
					d := yaml.NewEncoder(&j)
					d.SetIndent(2)
					err := d.Encode(opr)
					require.NoError(t, err)
					var decoded map[string]any
					require.NoError(t, yaml.Unmarshal(j.Bytes(), &decoded))
					if !assert.ObjectsAreEqual(ts.ExpectAST, decoded) {
						fmt.Println("actual:")
						fmt.Println(j.String())
					}
					require.Equal(t, ts.ExpectAST, decoded)
				}
			})
			t.Run("schemaless", func(t *testing.T) {
				opr, vars, errs := gqt.Parse([]byte(ts.Template))
				compareErrors(t, ts.ExpectErrorsSchemaless, errs)
				if len(ts.ExpectErrorsSchemaless) > 0 {
					// Expect failure
					require.Zero(t, vars)
					require.Zero(t, opr)
				} else {
					// Expect success
					var j bytes.Buffer
					d := yaml.NewEncoder(&j)
					d.SetIndent(2)
					err := d.Encode(opr)
					require.NoError(t, err)
					var decoded map[string]any
					require.NoError(t, yaml.Unmarshal(j.Bytes(), &decoded))
					if !assert.ObjectsAreEqual(ts.ExpectASTSchemaless, decoded) {
						fmt.Println("actual(schemaless):")
						fmt.Println(j.String())
					}
					require.Equal(t, ts.ExpectASTSchemaless, decoded)
				}
			})
		})
	}
}

func compareErrors(t *testing.T, expected []string, actual []gqt.Error) {
	if len(expected) < 1 {
		for _, act := range actual {
			t.Errorf("unexpected error: %v", act)
		}
		return
	}
	for i, e := range expected {
		if i >= len(actual) {
			t.Errorf("missing error: %v", e)
			continue
		}
		assert.Equal(t, e, actual[i].Error(), "at index %d", i)
	}
	if d := len(actual) - len(expected); d > 0 {
		for _, act := range actual[d:] {
			t.Errorf("unexpected error: %v", act)
		}
	}
}

func TestParseEmpty(t *testing.T) {
	opr, vars, errs := gqt.Parse([]byte(""))
	require.Equal(t, []gqt.Error{
		{
			Location: gqt.Location{Index: 0, Line: 1, Column: 1},
			Msg: "unexpected end of file, expected " +
				"query, mutation, or subscription operation definition",
		},
	}, errs)
	require.Nil(t, opr)
	require.Zero(t, vars)
}

func TestParseVariables(t *testing.T) {
	input := `query {
		f1(a: $b+$x, c=$c: $b) {
			f2(b=$b:*, x=$unused: $c+$b, d: {o1=$x:*})
		}
	}`
	opr, vars, errs := gqt.Parse([]byte(input))
	require.Len(t, errs, 0, "unexpected errors: %v", errs)
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
		t, gqt.Location{Index: 63, Line: 3, Column: 29},
		vars["b"].References[2].(*gqt.Variable).Location,
	)

	require.Equal(
		t, gqt.Location{Index: 60, Line: 3, Column: 26},
		vars["c"].References[0].(*gqt.Variable).Location,
	)

	require.Equal(
		t, gqt.Location{Index: 19, Line: 2, Column: 12},
		vars["x"].References[0].(*gqt.Variable).Location,
	)
}
