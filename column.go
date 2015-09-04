package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

type ColumnElem struct {
	name     string
	table    *TableElem
	datatype types.Type
	invalid  bool
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

// Column is the constructor for a ColumnElem
func Column(name string, datatype types.Type) ColumnElem {
	return ColumnElem{
		name:     name,
		datatype: datatype,
	}
}
