/*
Sol is a SQL toolkit for Go - in the style of Python's SQLAlchemy Core:

- Build complete database schemas

- Create reusable and cross-dialect SQL statements

- Allow struct instances and slices to be directly populated by the database

- Support for MySQL, PostGres, and SQLite3

*/
package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

// Test dialect - uses postgres style parameterization
type defaultDialect struct{}

// The default dialect must implement the dialect.Dialect interface
var _ dialect.Dialect = &defaultDialect{}

func (dialect defaultDialect) Param(i int) string {
	return fmt.Sprintf(`$%d`, i+1)
}
