package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

type Columnar interface {
	Alias() string
	FullName() string
	Name() string
	Table() *TableElem
}

type ColumnElem struct {
	name     string
	alias    string
	table    *TableElem
	datatype types.Type
	invalid  bool
}

// Alias returns the Column's alias
func (col ColumnElem) Alias() string {
	return col.alias
}

// As sets an alias for this ColumnElem
func (col ColumnElem) As(alias string) ColumnElem {
	col.alias = alias
	return col
}

// Columns returns the ColumnElem itself in a slice of ColumnElem. This
// method implements the Selectable interface.
func (col ColumnElem) Columns() []ColumnElem {
	return []ColumnElem{col}
}

// Create implements the Creatable interface that outputs a column of a
// CREATE TABLE statement.
func (col ColumnElem) Create(d dialect.Dialect) (string, error) {
	compiled, err := col.datatype.Create(d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s %s`, col.Name(), compiled), nil
}

// FullName prefixes the column name with the table name
func (col ColumnElem) FullName() string {
	return fmt.Sprintf(`"%s"."%s"`, col.table.name, col.name)
}

func (col ColumnElem) Name() string {
	return fmt.Sprintf(`"%s"`, col.name)
}

// Modify implements the Modifier interface, allowing the ColumnElem to
// modify the given TableElem.
func (col ColumnElem) Modify(table *TableElem) error {
	if table == nil {
		return fmt.Errorf(
			"sol: column %s cannot modify a nil table",
			col.name,
		)
	}
	if err := isValidColumnName(col.name); err != nil {
		return err
	}

	// Add the table to the column
	if col.table != nil {
		return fmt.Errorf(
			"sol: column %s already belongs to table %s",
			col.name, col.table.name,
		)
	}
	col.table = table

	// Add the column to the table
	if err := table.columns.add(col); err != nil {
		return err
	}

	// Add the type to the table creates
	table.creates = append(table.creates, col)

	return nil
}

// Table returns the column's TableElem
func (col ColumnElem) Table() *TableElem {
	return col.table
}

// Ordering
// --------

// Orerable implements the Orderable interface that allows the column itself
// to be used in an OrderBy clause.
func (col ColumnElem) Orderable() OrderedColumn {
	return OrderedColumn{inner: col}
}

// Asc returns an OrderedColumn. It is the same as passing the column itself
// to an OrderBy clause.
func (col ColumnElem) Asc() OrderedColumn {
	return OrderedColumn{inner: col}
}

// Desc returns an OrderedColumn that will sort in descending order.
func (col ColumnElem) Desc() OrderedColumn {
	return OrderedColumn{inner: col, desc: true}
}

// NullsFirst returns an OrderedColumn that will sort NULLs first.
func (col ColumnElem) NullsFirst() OrderedColumn {
	return OrderedColumn{inner: col, nullsFirst: true}
}

// NullsLast returns an OrderedColumn that will sort NULLs last.
func (col ColumnElem) NullsLast() OrderedColumn {
	return OrderedColumn{inner: col, nullsLast: true}
}

// Column is the constructor for a ColumnElem
func Column(name string, datatype types.Type) ColumnElem {
	return ColumnElem{
		name:     name,
		datatype: datatype,
	}
}
