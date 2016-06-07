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
	table      Tabular
	columns    UniqueColumnSet
	valuesList []Values
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

	// There must be values, and there must be more than one value in the
	// first element of the values list - otherwise create nil values
	// TODO Create a ColumnValuesSet to handle the following?
	aliases := Aliases{}
	if len(stmt.valuesList) != 0 && len(stmt.valuesList[0]) != 0 {
		// Be friendly: take the intersection of the Values and the
		// currently selected columns. The Values keys will be matched
		// with the precedence:
		// 1. Exact match
		// 2. Camel to snake case conversion (case sensitive)
		// Only the first Values element will be matched
		for _, key := range stmt.valuesList[0].Keys() {
			column := stmt.columns.Get(key)
			if column.IsInvalid() {
				column = stmt.columns.Get(camelToSnake(key))
			}
			if column.IsValid() {
				aliases[column.Name()] = key
			}
		}

		// Remove the unmatched columns
		stmt.columns = stmt.columns.Filter(aliases.Keys()...)
	} else {
		stmt.valuesList = []Values{stmt.columns.EmptyValues()}
		for _, column := range stmt.columns.order {
			aliases[column.Name()] = column.Name() // ugly
		}
	}

	// No columns? no statement!
	if len(stmt.columns.order) == 0 {
		return "", ErrNoColumns
	}

	// TODO Bulk insert syntax is dialect specific
	// TODO must all Values must have the same keys?
	var groups []string
	if len(stmt.valuesList) != 0 {
		groups = make([]string, len(stmt.valuesList))
		for g, values := range stmt.valuesList {
			group := make([]string, len(stmt.columns.order))
			for i, column := range stmt.columns.order {
				param := NewParam(values[aliases[column.Name()]])
				var err error
				if group[i], err = param.Compile(d, ps); err != nil {
					return "", err
				}
			}
			groups[g] = fmt.Sprintf(`(%s)`, strings.Join(group, ", "))
		}
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

// Values adds parameters to the INSERT statement. Accepted types:
// struc, Values, or a slice of either. Both pointers and values are accepted.
func (stmt InsertStmt) Values(obj interface{}) InsertStmt {
	elem := reflect.Indirect(reflect.ValueOf(obj))
	stmt.valuesList = make([]Values, 1)

	// Examine allowed types
	var unsupported bool
	switch elem.Kind() {
	case reflect.Map:
		switch converted := obj.(type) {
		case Values:
			stmt.valuesList[0] = converted
		case *Values:
			stmt.valuesList[0] = *converted
		default:
			unsupported = true
		}
	case reflect.Struct:
		var err error
		if stmt.valuesList[0], err = ValuesOf(obj); err != nil {
			stmt.AddMeta(err.Error())
			return stmt
		}
	case reflect.Slice:
		if elem.Len() == 0 {
			stmt.AddMeta("sol: cannot insert values from an empty slice")
			return stmt
		}

		// Slices of structs or Values are acceptable
		if elem.Index(0).Kind() == reflect.Struct {
			stmt.valuesList = make([]Values, elem.Len())
			var err error
			for i := range stmt.valuesList {
				obj := elem.Index(i).Interface()
				if stmt.valuesList[i], err = ValuesOf(obj); err != nil {
					stmt.AddMeta(err.Error())
					return stmt
				}
			}
			break
		}

		switch converted := obj.(type) {
		// TODO []*Values, *[]*Values are unsupported
		case []Values:
			stmt.valuesList = converted
		case *[]Values:
			stmt.valuesList = *converted
		default:
			unsupported = true
		}
	default:
		unsupported = true
	}

	if unsupported {
		stmt.AddMeta(
			"sol: unsupported type %T for inserted values - accepted types: struct, Values, or a slice of either",
			obj,
		)
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

	if len(columns) == 0 {
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

	for _, column := range columns {
		if column.IsInvalid() {
			stmt.AddMeta("sol: column %s does not exist", column.FullName())
			continue
		}

		if column.Table() != stmt.table {
			stmt.AddMeta(
				"sol: all columns in Insert() must belong to table %s",
				stmt.table.Name(),
			)
			continue
		}
		var err error
		if stmt.columns, err = stmt.columns.Add(column); err != nil {
			stmt.AddMeta(err.Error())
		}
	}
	return
}
