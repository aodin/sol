package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/sol"
	"github.com/aodin/sol/dialect"
)

// InsertStmt is the internal representation of an INSERT ... RETURNING
// statement.
type InsertStmt struct {
	sol.InsertStmt
	returning []sol.Columnar
}

// String outputs the parameter-less INSERT ... RETURNING statement in a
// neutral dialect.
func (stmt InsertStmt) String() string {
	compiled, _ := stmt.Compile(&PostGres{}, sol.Params())
	return compiled
}

// Compile outputs the INSERT ... RETURNING statement using the given dialect
// and parameters. An error may be returned because of a pre-existing error
// or because an error occurred during compilation.
func (stmt InsertStmt) Compile(d dialect.Dialect, ps *sol.Parameters) (string, error) {
	compiled, err := stmt.InsertStmt.Compile(d, ps)
	if err != nil {
		return "", err
	}
	if len(stmt.returning) > 0 {
		compiled += fmt.Sprintf(
			" RETURNING %s",
			strings.Join(sol.CompileColumns(stmt.returning), ", "),
		)
	}
	return compiled, nil
}

// Returning adds a RETURNING clause to the statement.
// TODO How to remove a returning?
func (stmt InsertStmt) Returning(selections ...sol.Selectable) InsertStmt {
	// TODO An INSERT ... RETURING for all columns of the inserted row can
	// also use the syntax RETURNING *, see:
	// http://www.postgresql.org/docs/devel/static/sql-insert.html

	// If no selections were provided, default to the table
	if len(selections) == 0 {
		for _, column := range stmt.Table().Columns() {
			stmt.returning = append(stmt.returning, column)
		}
		return stmt
	}

	// If selections have been specified, use those instead
	for _, selection := range selections {
		if selection == nil {
			stmt.AddMeta(
				"postgres: received a nil selectable in Returning() - do the columns or tables you selected exist?",
			)
			return stmt
		}

		// All selected columns must belong to the INSERT table
		for _, column := range selection.Columns() {
			if column.Table() != stmt.Table() {
				stmt.AddMeta(
					"postgres: the column '%s' in Returning() does not belong to the inserted table '%s'",
					column.Name(), stmt.Table().Name(),
				)
				break
			}
			stmt.returning = append(stmt.returning, column)
		}
	}
	return stmt
}

// Values proxies to the inner InsertStmt's Values method
func (stmt InsertStmt) Values(args interface{}) InsertStmt {
	stmt.InsertStmt = stmt.InsertStmt.Values(args)
	return stmt
}

// Insert creates an INSERT ... RETURNING statement for the given columns.
// There must be at least one column and all columns must belong to the
// same table.
func Insert(selections ...sol.Selectable) InsertStmt {
	return InsertStmt{
		InsertStmt: sol.Insert(selections...),
	}
}
