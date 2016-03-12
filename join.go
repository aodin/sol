package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

// JoinClause implements a variety of joins
type JoinClause struct {
	ArrayClause
	method string
	table  *TableElem
}

// String returns a default string representation of the JoinClause
func (j JoinClause) String() string {
	compiled, _ := j.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile compiles a JoinClause
func (j JoinClause) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	// Ignore clauses if CROSS
	if j.method == "CROSS JOIN" {
		return fmt.Sprintf(` CROSS JOIN "%s"`, j.table.Name()), nil
	}

	// If no clauses were given, assume the join is NATURAL
	if len(j.ArrayClause.clauses) == 0 {
		return fmt.Sprintf(
			` NATURAL %s "%s"`, j.method, j.table.Name(),
		), nil
	}

	// TODO Pass the joining table and auto-create the condition?
	// Compile the clauses of the join statement
	clauses, err := j.ArrayClause.Compile(d, ps)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		` %s "%s" ON %s`, j.method, j.table.Name(), clauses,
	), nil
}
