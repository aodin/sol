package postgres

import (
	"fmt"

	_ "github.com/lib/pq" // Register the PostGres driver

	"github.com/aodin/sol/dialect"
)

// PostGres implements the Dialect interface for postgres databases.
type PostGres struct{}

// The PostGres dialect must implement the dialect.Dialect interface
var _ dialect.Dialect = &PostGres{}

// Param returns the postgres specific parameterization scheme.
func (d *PostGres) Param(i int) string {
	return fmt.Sprintf(`$%d`, i+1)
}

// Dialect is a constructor for the PostGres Dialect
func Dialect() *PostGres {
	return &PostGres{}
}

// Add the postgres dialect to the dialect registry
func init() {
	dialect.Register("postgres", Dialect())
}
