package postgres

import (
	"fmt"

	"github.com/aodin/sol"
	"github.com/aodin/sol/dialect"
)

// InsertStmt is the internal representation of an INSERT ... RETURNING
// statement.
type InsertStmt struct {
	sol.InsertStmt
	onConflict      bool
	conflictTargets []string
	values          sol.Values
	where           sol.Clause
	returning       sol.ColumnSet
}

// String outputs the parameter-less INSERT ... RETURNING statement in the
// PostGres dialect.
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

	if stmt.onConflict {
		compiled += " ON CONFLICT"
		// TODO conflict targets
		if len(stmt.values) > 0 {
			compiledValues, err := stmt.values.Compile(d, ps)
			if err != nil {
				return "", fmt.Errorf("sol: failed to compile values: %s", err)
			}
			compiled += fmt.Sprintf(" DO UPDATE SET %s", compiledValues)

			// Add a WHERE clause if specified
			if stmt.where != nil {
				where, err := stmt.where.Compile(d, ps)
				if err != nil {
					return "", err
				}
				compiled += fmt.Sprintf(" WHERE %s", where)
			}
		} else {
			compiled += " DO NOTHING"
		}
	}

	if stmt.returning.Exists() {
		selections, err := stmt.returning.Compile(d, ps)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" RETURNING %s", selections)
	}
	return compiled, nil
}

// OnConflict adds UPSERT behavior to the INSERT. By Default, it will
// DO NOTHING.
func (stmt InsertStmt) OnConflict(targets ...string) InsertStmt {
	stmt.conflictTargets = targets
	stmt.onConflict = true
	return stmt
}

// Where should only be used alongside OnConflict. Only one WHERE
// is allowed per statement. Additional calls to Where will overwrite the
// existing WHERE clause.
func (stmt InsertStmt) Where(conditions ...sol.Clause) InsertStmt {
	if len(conditions) > 1 {
		// By default, multiple where clauses will be joined using AllOf
		stmt.where = sol.AllOf(conditions...)
	} else if len(conditions) == 1 {
		stmt.where = conditions[0]
	} else {
		// Clear the existing conditions
		stmt.where = nil
	}
	return stmt
}

// DoNothing sets the ON CONFLICT behavior to DO NOTHING
func (stmt InsertStmt) DoNothing() InsertStmt {
	stmt.onConflict = true
	stmt.values = sol.Values{}
	return stmt
}

// DoUpdate sets the ON CONFLICT behavior to DO UPDATE if at least
// one value is given
func (stmt InsertStmt) DoUpdate(values sol.Values) InsertStmt {
	stmt.onConflict = true
	stmt.values = values
	return stmt
}

// RemoveOnConflict will remove the ON CONFLICT behavior
func (stmt InsertStmt) RemoveOnConflict() InsertStmt {
	stmt.onConflict = false
	stmt.values = sol.Values{}
	stmt.conflictTargets = nil
	stmt.where = nil
	return stmt
}

// Returning adds a RETURNING clause to the statement.
// TODO How to remove a returning?
func (stmt InsertStmt) Returning(selections ...sol.Selectable) InsertStmt {
	// TODO An INSERT ... RETURING for all columns of the inserted row can
	// also use the syntax RETURNING *, see:
	// http://www.postgresql.org/docs/devel/static/sql-insert.html

	// If no selections were provided, default to the table
	if len(selections) == 0 && stmt.Table() != nil {
		for _, column := range stmt.Table().Columns() {
			stmt.returning, _ = stmt.returning.Add(column)
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
			stmt.returning, _ = stmt.returning.Add(column)
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
