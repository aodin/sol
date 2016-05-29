package sol

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aodin/sol/dialect"
)

// InsertStmt is the internal representation of an INSERT statement.
type InsertStmt struct {
	Stmt
	table       Tabular
	columns     ColumnSet
	multiValues []Values
}

// String outputs the parameter-less INSERT statement in a neutral dialect.
func (stmt InsertStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// Table returns the INSERT statement's table
func (stmt InsertStmt) Table() Tabular {
	return stmt.table
}

// Compile outputs the INSERT statement using the given dialect and parameters.
// An error may be returned because of a pre-existing error or because
// an error occurred during compilation.
func (stmt InsertStmt) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	// Check for delayed errors
	if err := stmt.Error(); err != nil {
		return "", err
	}

	cols := len(stmt.columns.order)
	// No columns? no statement!
	if cols == 0 {
		return "", ErrNoColumns
	}

	// TODO Bulk insert syntax is dialect specific
	// TODO what safety checks should happen here?
	// Every Values must have the expected number of columns with
	// matching column names
	var groups []string
	if len(stmt.multiValues) > 0 {
		groups = make([]string, len(stmt.multiValues))
		for g, values := range stmt.multiValues {
			// TODO method of ColumnSet - CompileWith?
			group := make([]string, len(stmt.columns.order))
			for i, column := range stmt.columns.order {
				param := NewParam(values[column.Name()])
				var err error
				if group[i], err = param.Compile(d, ps); err != nil {
					return "", err
				}
			}
			groups[g] = fmt.Sprintf(`(%s)`, strings.Join(group, ", "))
		}
	} else {
		group := make([]string, len(stmt.columns.order))
		for i, _ := range stmt.columns.order {
			param := NewParam(nil)
			var err error
			if group[i], err = param.Compile(d, ps); err != nil {
				return "", err
			}
		}
		groups = []string{fmt.Sprintf(`(%s)`, strings.Join(group, ", "))}
	}

	compiled := []string{
		INSERT,
		INTO,
		stmt.table.Name(),
		fmt.Sprintf("(%s)", strings.Join(stmt.columns.Names(), ", ")),
		VALUES,
		strings.Join(groups, ", "),
	}
	return strings.Join(compiled, WHITESPACE), nil
}

// Values adds parameters to the INSERT statement. Values can be given
// as structs, a single Values type, or slices of structs or Values.
// Both pointers and values are accepted.
func (stmt InsertStmt) Values(arg interface{}) InsertStmt {
	elem := reflect.Indirect(reflect.ValueOf(arg))

	switch elem.Kind() {
	case reflect.Map:
		values, ok := arg.(Values) // The only allowed map type is Values
		if !ok {
			stmt.AddMeta("sol: inserted values of type map must be Values")
			break
		}

		// Be friendly: take the intersection of the Values and the
		// currently selected columns
		// TODO how friendly - perform snake to case conversions here?
		stmt.columns = stmt.columns.Filter(values.Keys()...)
		stmt.multiValues = []Values{values.Filter(stmt.columns.Names()...)}
	default:
		stmt.AddMeta("sol: unaccepted Values type (for now)")
	}

	return stmt
}

// Insert creates an INSERT statement for the given columns and tables.
// There must be at least one column and all columns must belong to the
// same table.
func Insert(selections ...Selectable) (stmt InsertStmt) {
	columns := []ColumnElem{} // Holds columns until validated
	for _, selection := range selections {
		if selection == nil {
			stmt.AddMeta("sol: received a nil selectable in Insert()")
			return
		}
		columns = append(columns, selection.Columns()...)
	}

	if len(columns) < 1 {
		stmt.AddMeta("sol: Insert() must be given at least one column")
		return
	}

	// All columns must have the same table
	column := columns[0]
	if column.Table() == nil {
		stmt.AddMeta("sol: all columns given to Insert() must have a table")
		return
	}
	stmt.table = column.Table()

	// TODO inserted columns should be unique
	for _, column := range columns {
		if column.IsInvalid() {
			stmt.AddMeta("sol: column %s does not exist", column.FullName())
			return
		}

		if column.Table() != stmt.table {
			stmt.AddMeta(
				"sol: all columns in Insert() must belong to table %s",
				stmt.table.Name(),
			)
			return
		}
		stmt.columns.order = append(stmt.columns.order, column)
	}
	return
}
