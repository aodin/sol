package sol

import (
	"strings"

	"github.com/aodin/sol/dialect"
)

type Clause interface {
	Compiles
}

// ArrayClause is any number of clauses with a column join
type ArrayClause struct {
	clauses []Clause
	sep     string
}

func (c ArrayClause) String() string {
	compiled, _ := c.Compile(&defaultDialect{}, Params())
	return compiled
}

func (c ArrayClause) Compile(d dialect.Dialect, params *Parameters) (string, error) {
	compiled := make([]string, len(c.clauses))
	var err error
	for i, clause := range c.clauses {
		compiled[i], err = clause.Compile(d, params)
		if err != nil {
			return "", err
		}
	}
	return strings.Join(compiled, c.sep), nil
}
