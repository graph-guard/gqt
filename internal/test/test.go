package test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	mdast "github.com/gomarkdown/markdown/ast"
	mdparser "github.com/gomarkdown/markdown/parser"
	"github.com/stretchr/testify/require"
)

func New[T any](t *testing.T, fn func(t *testing.T, x T)) func(T) {
	return func(x T) {
		_, filename, line, _ := runtime.Caller(1)
		name := fmt.Sprintf("%s:%d", filepath.Base(filename), line)
		t.Run(name, func(t *testing.T) {
			fn(t, x)
		})
	}
}

// ExecDirMD applies fn to each .md file in dirName recursively.
// Expects the markdown file to contain two code blocks,
// one for the input and another one for the expected output.
func ExecDirMD(
	t *testing.T,
	filesystem fs.FS,
	dirName string,
	expectationPrefix string,
	fn func(t *testing.T, input, expectation string),
) {
	t.Helper()
	var execDirectory func(
		t *testing.T,
		filesystem fs.FS,
		fullPath, name string,
		dir []fs.DirEntry,
		expectationPrefix string,
		fn func(t *testing.T, input, expectation string),
	)
	execDirectory = func(
		t *testing.T,
		filesystem fs.FS,
		fullPath, name string,
		dir []fs.DirEntry,
		expectationPrefix string,
		fn func(t *testing.T, input, expectation string),
	) {
		t.Run(name, func(t *testing.T) {
			for _, e := range dir {
				en := e.Name()
				n := filepath.Join(fullPath, en)
				if e.IsDir() {
					e, err := fs.ReadDir(filesystem, n)
					require.NoError(t, err)
					execDirectory(t, filesystem, n, en, e, expectationPrefix, fn)
					continue
				}

				if !strings.HasSuffix(en, ".md") {
					t.Fatalf("unsupported file: %q", n)
				}

				t.Run(en, func(t *testing.T) {
					c, err := fs.ReadFile(filesystem, n)
					require.NoError(t, err)

					p := mdparser.New()
					ast := p.Parse(c)

					children := ast.GetChildren()

					if len(children) != 2 ||
						children[0].AsLeaf() == nil ||
						children[1].AsLeaf() == nil {
						t.Fatal("expected 2 code blocks, " +
							"one for the input " +
							"and another one for the expected output")
					}

					input := string(children[0].AsLeaf().Literal)
					output := string(children[1].AsLeaf().Literal)

					if input[len(input)-1] == '\n' {
						input = input[:len(input)-1]
					}
					if output != "" && output[len(output)-1] == '\n' {
						output = output[:len(output)-1]
					}

					fn(t, input, output)
				})
			}
		})
	}

	e, err := fs.ReadDir(filesystem, dirName)
	require.NoError(t, err)

	execDirectory(
		t, filesystem, dirName, dirName, e,
		expectationPrefix,
		fn,
	)
}

// Markdown applies fn to each .txt file in dirName recursively.
// Expects the expectation line to begin with expectationPrefix.
func Markdown(
	t *testing.T,
	filesystem fs.FS,
	dirName string,
	expectationPrefix string,
	fn func(t *testing.T, input, expectation string),
) {
	var execDirectory func(
		t *testing.T,
		filesystem fs.FS,
		fullPath, name string,
		dir []fs.DirEntry,
		expectationPrefix string,
		fn func(t *testing.T, input, expectation string),
	)
	execDirectory = func(
		t *testing.T,
		filesystem fs.FS,
		fullPath, name string,
		dir []fs.DirEntry,
		expectationPrefix string,
		fn func(t *testing.T, input, expectation string),
	) {
		t.Run(name, func(t *testing.T) {
			for _, e := range dir {
				en := e.Name()
				n := filepath.Join(fullPath, en)
				if e.IsDir() {
					e, err := fs.ReadDir(filesystem, n)
					require.NoError(t, err)
					execDirectory(t, filesystem, n, en, e, expectationPrefix, fn)
					continue
				}

				if !strings.HasSuffix(en, ".md") {
					t.Errorf("unsupported file type: %q", n)
					continue
				}

				t.Run(en, func(t *testing.T) {
					c, err := fs.ReadFile(filesystem, n)
					require.NoError(t, err)

					p := mdparser.New()
					ast := p.Parse(c)

					children := ast.GetChildren()

					var inputBlock *mdast.Container
					var outputBlock *mdast.Container
					if len(children) == 2 {
						inputBlock = children[0].AsContainer()
						outputBlock = children[1].AsContainer()
					}

					require.Len(t, children, 2)
					require.NotNil(t, inputBlock)
					require.NotNil(t, outputBlock)
					require.Len(t, inputBlock.Children, 1)
					require.Len(t, outputBlock.Children, 1)

					require.IsType(t, inputBlock.Children[0], &mdast.Text{})
					require.IsType(t, outputBlock.Children[0], &mdast.Text{})

					input := string(inputBlock.Children[0].(*mdast.Text).Leaf.Literal)
					output := string(outputBlock.Children[0].(*mdast.Text).Leaf.Literal)

					if len(children) != 2 ||
						inputBlock == nil ||
						outputBlock == nil ||
						len(inputBlock.Children) != 1 ||
						len(outputBlock.Children) != 1 {
						t.Fatal("expected 2 code blocks, " +
							"one for the input " +
							"and another one for the expected output")
					}

					fn(t, input, output)
				})
			}
		})
	}

	e, err := fs.ReadDir(filesystem, dirName)
	require.NoError(t, err)

	execDirectory(
		t, filesystem, dirName, dirName, e,
		expectationPrefix,
		fn,
	)
}
