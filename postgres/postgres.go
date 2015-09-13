package postgres

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/aodin/sol/dialect"
)

// PostGres implements the Dialect interface for postgres databases.
type PostGres struct{}

// Parameterize returns the postgres specific parameterization scheme.
func (d *PostGres) Param(i int) string {
	return fmt.Sprintf(`$%d`, i+1)
}

// Add the postgres dialect to the dialect registry
func init() {
	dialect.Register("postgres", &PostGres{})
}
