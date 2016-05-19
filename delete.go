package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

// DeleteStmt is the internal representation of a DELETE statement.
type DeleteStmt struct {
	ConditionalStmt
	table *TableElem
}

// String outputs the parameter-less DELETE statement in a neutral dialect.
func (stmt DeleteStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile outputs the DELETE statement using the given dialect and parameters.
// An error may be returned because of a pre-existing error or because
// an error occurred during compilation.
func (stmt DeleteStmt) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	if err := stmt.Error(); err != nil {
		return "", err
	}
	compiled := fmt.Sprintf(`DELETE FROM "%s"`, stmt.table.Name())

	if stmt.where != nil {
		cc, err := stmt.where.Compile(d, ps)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	return compiled, nil
}

// Where adds a conditional WHERE clause to the DELETE statement.
func (stmt DeleteStmt) Where(clauses ...Clause) DeleteStmt {
	if len(clauses) > 1 {
		// By default, multiple where clauses will be joined will AllOf
		stmt.where = AllOf(clauses...)
	} else if len(clauses) == 1 {
		stmt.where = clauses[0]
	}
	return stmt
}

// Delete creates a DELETE statement for the given table.
func Delete(table *TableElem) (stmt DeleteStmt) {
	if table == nil {
		stmt.AddMeta("sol: attempting to DELETE a nil table")
		return
	}
	stmt.table = table
	return
}
