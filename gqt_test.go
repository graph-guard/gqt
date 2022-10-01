package gqt_test

import (
	"bytes"
	"embed"
	"strings"
	"testing"

	"github.com/graph-guard/gqt"
	"github.com/graph-guard/gqt/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

//go:embed test/error
var errorFS embed.FS

//go:embed test/ast
var astFS embed.FS

//go:embed test/schema_aware
var schemaFS embed.FS

//go:embed test/error_schema
var errorSchemaFS embed.FS

func TestParseErr(t *testing.T) {
	test.ExecDirMD2(
		t, errorFS, "test/error", "ERR: ",
		func(t *testing.T, input, expectation string) {
			expectedErrors := strings.Split(expectation, "\n")
			_, _, errs := gqt.Parse([]byte(input))
			for i, e := range expectedErrors {
				if i >= len(errs) {
					t.Errorf("missing error: %v", e)
					continue
				}
				assert.Equal(t, e, errs[i].Error())
			}
			if d := len(errs) - len(expectedErrors); d > 0 {
				for _, act := range errs[d:] {
					t.Errorf("unexpected error: %v", act)
				}
			}
		},
	)
}

func TestParse(t *testing.T) {
	test.ExecDirMD2(
		t, astFS, "test/ast", "",
		func(t *testing.T, input, expectation string) {

			var discard map[string]any
			require.NoError(
				t, yaml.Unmarshal([]byte(expectation), &discard),
				"invalid expectation YAML",
			)

			o, vars, errs := gqt.Parse([]byte(input))
			require.Nil(t, errs)
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
	opr, vars, errs := gqt.Parse([]byte(input))
	require.Len(t, errs, 0, "unexpected errors: %v", errs)
	require.NotNil(t, opr)
	require.Len(t, vars, 0)
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
		vars["b"].References[0].(*gqt.VariableRef).Location,
	)
	require.Equal(
		t, gqt.Location{Index: 29, Line: 2, Column: 22},
		vars["b"].References[1].(*gqt.VariableRef).Location,
	)
	require.Equal(
		t, gqt.Location{Index: 63, Line: 3, Column: 29},
		vars["b"].References[2].(*gqt.VariableRef).Location,
	)

	require.Equal(
		t, gqt.Location{Index: 60, Line: 3, Column: 26},
		vars["c"].References[0].(*gqt.VariableRef).Location,
	)

	require.Equal(
		t, gqt.Location{Index: 19, Line: 2, Column: 12},
		vars["x"].References[0].(*gqt.VariableRef).Location,
	)
}

func TestParserSchema(t *testing.T) {
	test.ExecDirMD3(
		t, schemaFS, "test/schema_aware", "",
		func(t *testing.T, schema, input, expectation string) {
			p, errp := gqt.NewParser([]gqt.Source{
				{Name: "test_schema_aware", Content: schema},
			})
			require.NoError(t, errp, "unexpected GraphQL schema parser error")
			o, vars, errs := p.Parse([]byte(input))
			require.Len(t, errs, 0, "unexpected errors: %v", errs)
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

func TestParserSchemaErr(t *testing.T) {
	test.ExecDirMD3(
		t, errorSchemaFS, "test/error_schema", "ERR: ",
		func(t *testing.T, schema, input, expectation string) {

			p, errp := gqt.NewParser([]gqt.Source{
				{Name: "testschema", Content: schema},
			})
			require.NoError(t, errp, "unexpected GraphQL schema parser error")

			expectedErrors := strings.Split(expectation, "\n")
			d, vars, errs := p.Parse([]byte(input))
			for i, e := range expectedErrors {
				if i >= len(errs) {
					t.Errorf("missing error: %v", e)
					continue
				}
				assert.Equal(t, e, errs[i].Error())
			}
			if d := len(errs) - len(expectedErrors); d > 0 {
				for _, act := range errs[d:] {
					t.Errorf("unexpected error: %v", act)
				}
			}

			require.Zero(t, d)
			require.Nil(t, vars)
		},
	)
}

// const schema = `
// type Query {
// 	id: ID!
// 	string: String!
// 	boolean: Boolean!
// 	float: Float!
// 	integer: Int!
// 	enum: Color!
// 	type: Type!
// 	array: [Int!]!

// 	idOpt: ID
// 	stringOpt: String
// 	booleanOpt: Boolean
// 	floatOpt: Float
// 	integerOpt: Int
// 	enumOpt: Color
// 	typeOpt: Type
// 	arrayOpt: [Int!]
// }

// type Mutation {
// 	withArgs(input: Input! inputOpt: Input!): Boolean!
// }

// type Type {
// 	id: ID!
// 	string: String!
// 	boolean: Boolean!
// 	float: Float!
// 	integer: Int!
// 	enum: Color!
// 	type: Type!
// 	array: [Int!]!

// 	idOpt: ID
// 	stringOpt: String
// 	booleanOpt: Boolean
// 	floatOpt: Float
// 	integerOpt: Int
// 	enumOpt: Color
// 	typeOpr: Type!
// 	arrayOpt: [Int!]
// }

// input Input {
// 	id: ID!
// 	string: String!
// 	boolean: Boolean!
// 	float: Float!
// 	integer: Int!
// 	enum: Color!
// 	array: [Int!]!

// 	idOpt: ID
// 	stringOpt: String
// 	booleanOpt: Boolean
// 	floatOpt: Float
// 	integerOpt: Int
// 	enumOpt: Color
// 	inputOpt: Input
// 	arrayOpt: [Int!]
// }

// enum Color {
// 	RED
// 	GREEN
// 	BLUE
// }
// `

// func TestValueType(t *testing.T) {
// 	type T struct {
// 		Schema string
// 		E      gqt.Expression
// 		Expect string
// 	}
// 	test := test.New(t, func(t *testing.T, x T) {
// 		p, err := gqt.NewParser([]gqt.Source{
// 			{
// 				Name:    "schema.graphqls",
// 				Content: x.Schema,
// 			},
// 		})
// 		require.NoError(t, err)
// 		actual := x.E.ValueType(p)
// 		require.Equal(t, x.Expect, actual)
// 	})

// 	test(T{
// 		E: &gqt.ExprLogicalAnd{
// 			Expressions: []gqt.Expression{
// 				&gqt.Int{},
// 				&gqt.Null{},
// 			},
// 		},
// 	})
// }
