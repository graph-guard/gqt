package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
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
