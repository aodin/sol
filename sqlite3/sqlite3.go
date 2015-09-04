package sqlite3

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/aodin/sol/dialect"
)

// Sqlite3 implements the Dialect interface for sqlite3 databases.
type Sqlite3 struct{}

// Param returns the sqlite3 specific parameterization scheme.
func (d *Sqlite3) Param(i int) string {
	return `?`
}

// Add the sqlite3 dialect to the dialect registry
func init() {
	dialect.Register("sqlite3", &Sqlite3{})
}
