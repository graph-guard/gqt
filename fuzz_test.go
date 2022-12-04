package gqt_test

import (
	"testing"

	"github.com/graph-guard/gqt/v4"
)

func FuzzParse(f *testing.F) {
	for _, tc := range []string{
		`query { foo }`,
		`query { foo bar bazz fuzz maz }`,
		"query {\n\tfoo\n\tbar\n\tbazz\n\tfuzz\n\tmaz\n}",
		`query { foo(string:"fuzz") }`,
		`query { foo(magic_number_int:42) }`,
		`query { foo(piNumberFloat:3.1415) }`,
		`query { foo(enum:SOME_ENUM_VALUE) }`,
		`query { foo(zero:0) }`,
		`query { bool(booleanTrue:true, booleanFalse:false) }`,
		`mutation { newUser(name: *) }`,
		`mutation { createFile(tags: len < 10 && [... len < 64]) }`,
		`mutation { orStatement(c: [1] || [1,2,3]) }`,
		`subscription { newEvents(category: "this" || "that" || "those") }`,
	} {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		_, _, _ = gqt.Parse([]byte(orig))
	})
}
