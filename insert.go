package sol

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aodin/sol/dialect"
)

var (
	// ErrNoColumns is returned when a query for one columns returns none
	ErrNoColumns = errors.New(
		"sol: attempt to create a statement with zero columns",
	)
)

// InsertStmt is the internal representation of an INSERT statement.
type InsertStmt struct {
	Stmt
	table   Tabular
	columns []Columnar
	args    []interface{}
	fields  fields
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

	cols := len(stmt.columns)
	// No columns? no statement!
	if cols == 0 {
		return "", ErrNoColumns
	}

	columns := make([]string, len(stmt.columns))
	for i, column := range stmt.columns {
		columns[i] = fmt.Sprintf(`"%s"`, column.Name())
	}

	// args must be divisable by cols without remainder
	if len(stmt.args)%cols != 0 {
		return "", fmt.Errorf(
			`sol: size mismatch between arguments and columns during INSERT: %d is not a multiple of %d`,
			len(stmt.args),
			cols,
		)
	}

	// Determine the number of rows that will be inserted
	rows := len(stmt.args) / cols
	// If there are no arguments, default to one rows and create
	// placeholder values
	if rows == 0 {
		rows = 1
		stmt.args = make([]interface{}, cols)
		for i := range stmt.args {
			stmt.args[i] = nil
		}
	}
	parameters := make([]string, rows)

	var param int
	for i := 0; i < rows; i += 1 {
		row := make([]string, cols)
		for j := 0; j < cols; j += 1 {
			// Parameters are dialect specific
			// TODO Why isn't compile being used here?
			row[j] = d.Param(param)
			ps.Add(stmt.args[param])
			param += 1
		}
		// TODO Parameters compilation?
		parameters[i] = fmt.Sprintf(`(%s)`, strings.Join(row, ", "))
	}

	// TODO Bulk insert syntax is dialect specific
	return fmt.Sprintf(
		`INSERT INTO "%s" (%s) VALUES %s`,
		stmt.table.Name(),
		strings.Join(columns, ", "),
		strings.Join(parameters, ", "),
	), nil
}

// TODO []ColumnElem should be a custom type in order to attach getter/setters
func (stmt InsertStmt) Has(name string) bool {
	for _, column := range stmt.columns {
		if column.Name() == name {
			return true
		}
	}
	return false
}

// Values adds parameters to the INSERT statement. If the given values do not
// match the statement's current columns, the columns will be updated.
// Valid values include structs, Values maps, or slices of structs or Values.
// It accepts pointers or values.
func (stmt InsertStmt) Values(arg interface{}) InsertStmt {
	// TODO If auto-updating fields are required, they will need pointers
	elem := reflect.Indirect(reflect.ValueOf(arg))

	switch elem.Kind() {
	case reflect.Struct:
		// Inspect the fields of the given struct
		unaligned := SelectFieldsFromElem(elem.Type())

		// TODO function to return names of columns
		columns := make([]string, len(stmt.columns))
		for i, column := range stmt.columns {
			columns[i] = column.Name()
		}

		stmt.fields = AlignColumns(columns, unaligned)

		// If no fields were found and the number of fields matches the
		// columns requested, then insert the struct's values as is.
		if stmt.fields.Empty() && len(unaligned) == len(stmt.columns) {
			stmt.fields = unaligned
			stmt.argsFromElem(elem)
			return stmt
		}

		// Remove unmatched columns and empty values from fields
		// with the 'omitempty' option
		stmt.updateColumns()
		stmt.trimFields(elem)

		// If no fields remain after trimming, abort
		if len(stmt.fields) == 0 {
			stmt.AddMeta("sol: could not match fields for INSERT - are the 'db' tags correct?")
			return stmt
		}

		// Collect the parameters
		stmt.argsFromElem(elem)

	case reflect.Slice:
		if elem.Len() == 0 {
			stmt.AddMeta("sol: args cannot be set by empty slices")
			return stmt
		}
		// Slices of structs or Values are acceptable
		// TODO check kind of elem directly?
		elem0 := elem.Index(0)
		if elem0.Kind() == reflect.Struct {
			unaligned := SelectFieldsFromElem(elem.Type().Elem())

			// TODO function to return names of columns
			columns := make([]string, len(stmt.columns))
			for i, column := range stmt.columns {
				columns[i] = column.Name()
			}
			stmt.fields = AlignColumns(columns, unaligned)

			// If no fields were found and the number of fields matches the
			// columns requested, then insert the struct's values as is.
			if stmt.fields.Empty() && len(unaligned) == len(stmt.columns) {
				stmt.fields = unaligned
				for i := 0; i < elem.Len(); i++ {
					stmt.argsFromElem(elem.Index(i))
				}
				return stmt
			}

			// Remove unmatched columns and empty values from fields
			// with the 'omitempty' option
			stmt.updateColumns()
			stmt.trimFields(elem0)

			// If no fields remain after trimming, abort
			if len(stmt.fields) == 0 {
				stmt.AddMeta(
					"sol: could not match any fields for INSERT - do the field names or 'db' match the table columns?",
				)
				return stmt
			}

			// Add the parameters for each element
			for i := 0; i < elem.Len(); i++ {
				stmt.argsFromElem(elem.Index(i))
			}

			return stmt
		}

		valuesSlice, ok := arg.([]Values)
		if ok {
			if len(valuesSlice) == 0 {
				stmt.AddMeta(
					"sol: cannot insert []Values of length zero",
				)
				return stmt
			}

			// Set the table columns according to the values
			var err error
			if stmt.fields, err = valuesMap(stmt, valuesSlice[0]); err != nil {
				stmt.AddMeta(err.Error())
				return stmt
			}

			// TODO Column names should match in each values element!
			stmt.updateColumns()

			// Add the args in the values
			for _, v := range valuesSlice {
				stmt.argsFromValues(v)
			}

			return stmt
		}

		stmt.AddMeta(
			"sol: unsupported type %T for INSERT %s - values must be of type struct, Values, or a slice of either",
			arg, stmt,
		)

	case reflect.Map:
		// The only allowed map type is Values
		values, ok := arg.(Values)
		if !ok {
			stmt.AddMeta(
				"sol: inserted maps must be of type Values",
			)
			return stmt
		}

		// Set the table columns according to the values
		var err error
		if stmt.fields, err = valuesMap(stmt, values); err != nil {
			stmt.AddMeta(err.Error())
			return stmt
		}

		// Remove unmatched columns and add args from the values
		stmt.updateColumns()
		stmt.argsFromValues(values)
	}
	return stmt
}

func (stmt *InsertStmt) argsFromValues(values Values) {
	for _, column := range stmt.columns {
		stmt.args = append(stmt.args, values[column.Name()])
	}
}

func (stmt *InsertStmt) argsFromElem(elem reflect.Value) {
	for _, field := range stmt.fields {
		// TODO FieldByNameFunc?
		var fieldElem reflect.Value = elem
		for _, name := range field.names {
			fieldElem = fieldElem.FieldByName(name)
		}
		stmt.args = append(stmt.args, fieldElem.Interface())
	}
}

// A field marked omitempty can cause the removal of a column, only
// to have another value not have an empty value for that field
func (stmt *InsertStmt) trimFields(elem reflect.Value) {
	// TODO this function could be skipped if it was known that the given
	// struct has no omitempty fields
	validFields := fields{}
	for _, field := range stmt.fields {
		if !field.Exists() {
			continue
		}
		if field.options.Has(OmitEmpty) {
			var fieldElem reflect.Value = elem
			for _, name := range field.names {
				fieldElem = fieldElem.FieldByName(name)
			}
			if isEmptyValue(fieldElem) {
				// Remove the column
				stmt.columns = removeColumn(stmt.columns, field.column)
				continue
			}
		}
		// Keep the field
		validFields = append(validFields, field)
	}
	stmt.fields = validFields
}

// updateColumns removes any columns that weren't matched by the fields.
func (stmt *InsertStmt) updateColumns() {
	// TODO keep actual target columns separate from requested in case
	// the statement is updated?
	for _, column := range stmt.columns {
		if !stmt.fields.Has(column.Name()) {
			// TODO pass an index to prevent further nesting of iteration
			stmt.columns = removeColumn(stmt.columns, column.Name())
		}
	}
}

// isEmptyValue is from Go's encoding/json package: encode.go
// Copyright 2010 The Go Authors. All rights reserved.
// TODO what about pointer fields?
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		t, ok := v.Interface().(time.Time)
		if ok {
			return t.IsZero()
		}
	}
	return false
}

// TODO error if there was no match?
func removeColumn(columns []Columnar, name string) []Columnar {
	for i, col := range columns {
		if col.Name() == name {
			return append(columns[:i], columns[i+1:]...)
		}
	}
	return columns
}

// TODO better way to pass columns than by using the whole statement?
// TODO if this is better generalized then it can be used with UPDATE and
// DELETE statements.
// TODO use a column set - that's all it needs - maybe move to fields?
func valuesMap(stmt InsertStmt, values Values) (fields, error) {
	matches := fields{}
	for column := range values {
		if !stmt.Has(column) {
			return nil, fmt.Errorf(
				"sol: cannot INSERT a value with column '%s' as it has no corresponding column in the INSERT statement",
				column,
			)
		}
		matches = append(matches, field{column: column}) // TODO set index?
	}
	return matches, nil
}

// Insert creates an INSERT statement for the given columns. There must be at
// least one column and all columns must belong to the same table.
func Insert(selections ...Selectable) (stmt InsertStmt) {
	var columns []Columnar
	for _, selection := range selections {
		if selection == nil {
			stmt.AddMeta("sol: INSERT received a nil selectable - do the columns or tables you selected exist?")
			return
		}
		columns = append(columns, selection.Columns()...)
	}

	if len(columns) < 1 {
		stmt.AddMeta("sol: no columns were selected for INSERT")
		return
	}

	// The table is set from the first column
	column := columns[0]
	if column.Table() == nil {
		stmt.AddMeta(
			"sol: attempting to INSERT to a column unattached to a table",
		)
		return
	}
	stmt.table = column.Table()

	// Prepend the first column
	for _, column := range columns {
		// TODO Check column validity

		if column.Table() != stmt.table {
			stmt.AddMeta("sol: columns of an INSERT must all belong to the same table")
			return
		}
		stmt.columns = append(stmt.columns, column)
	}
	return
}
