package sol

import "github.com/aodin/sol/dialect"

// Compiles is the main interface that all SQL statements and clauses
// must implement in order to be queried.
type Compiles interface {
	Compile(dialect.Dialect, *Parameters) (string, error)
}
