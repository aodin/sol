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
