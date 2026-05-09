package cadence

import (
	"time"

	"github.com/connordoman/cadence/internal/compiler"
)

func Compile(source string) ([]time.Time, error) {
	return compiler.Compile(source)
}

func CompileAsStringSlice(source string) ([]string, error) {
	return compiler.CompileAsStringSlice(source)
}

func CompileAsJSON(source string) (string, error) {
	return compiler.CompileAsJSON(source)
}
