package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

// UpdateStmt is the internal representation of an SQL UPDATE statement.
type UpdateStmt struct {
	ConditionalStmt
	table  *TableElem
	values Values
}

// String outputs the parameter-less UPDATE statement in a neutral dialect.
func (stmt UpdateStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Compile outputs the UPDATE statement using the given dialect and parameters.
func (stmt UpdateStmt) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	// Check for delayed errors
	if err := stmt.Error(); err != nil {
		return "", err
	}

	// If no values were attached, then create a default values map
	if stmt.values == nil {
		stmt.values = Values{}
		for _, column := range stmt.table.Columns() {
			stmt.values[column.Name()] = nil
		}
	}

	// Error if there are no values
	if len(stmt.values) == 0 {
		return "", fmt.Errorf("sol: UPDATE has no values")
	}

	// Compile the values
	compiledValues, err := stmt.values.Compile(d, ps)
	if err != nil {
		return "", fmt.Errorf("sol: failed to compile UPDATE values: %s", err)
	}

	// Begin building the UPDATE statement
	compiled := fmt.Sprintf(
		`UPDATE %s SET %s`,
		stmt.table.Name(),
		compiledValues,
	)

	// Add a conditional statement if it exists
	if stmt.where != nil {
		cc, err := stmt.where.Compile(d, ps)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", cc)
	}
	return compiled, nil
}

// Values attaches the given values to the statement. The keys of values
// must match columns in the table.
// TODO Allow structs to be used if a primary key is specified in the schema
func (stmt UpdateStmt) Values(values Values) UpdateStmt {
	// There must be some columns to update!
	if len(values) == 0 {
		stmt.AddMeta("sol: there must be at least one value to update")
		return stmt
	}

	// Confirm that all values' keys are columns in the table
	// TODO perform column alias matching? e.g. UUID > uuid or ItemID > item_id
	for key := range values {
		if !stmt.table.Has(key) {
			stmt.AddMeta(
				"sol: no column '%s' exists in the table '%s'",
				key, stmt.table.Name(),
			)
		}
	}

	stmt.values = values
	return stmt
}

// Where adds a conditional WHERE clause to the UPDATE statement.
func (stmt UpdateStmt) Where(clauses ...Clause) UpdateStmt {
	if len(clauses) > 1 {
		// By default, multiple where clauses will be joined will AllOf
		stmt.where = AllOf(clauses...)
	} else if len(clauses) == 1 {
		stmt.where = clauses[0]
	}
	return stmt
}

// Update creates an UPDATE statement for the given table.
func Update(table *TableElem) (stmt UpdateStmt) {
	if table == nil {
		stmt.AddMeta("sol: attempting to UPDATE a nil table")
		return
	}
	stmt.table = table
	return
}
