package sol

import "fmt"

// Test dialect - uses postgres style parameterization
type defaultDialect struct{}

func (dialect defaultDialect) Param(i int) string {
	return fmt.Sprintf(`$%d`, i+1)
}
