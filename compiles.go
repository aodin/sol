package sol

import (
	"github.com/aodin/sol/dialect"
)

// Compiles in the main SQL statement interface.
type Compiles interface {
	Compile(dialect.Dialect, *Parameters) (string, error)
}
