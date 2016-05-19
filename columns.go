package sol

import (
	"fmt"
)

// Columns maps column to name to a ColumnElem. It also maintains an order
// of columns.
type Columns struct {
	order []ColumnElem
	c     map[string]ColumnElem
}

func (columns *Columns) add(column Columnar) error {
	// Tables cannot have duplicate column names
	// TODO the column should already be assigned a table
	if columns.Has(column.Name()) {
		return fmt.Errorf(
			"sol: table '%s' already has a column named '%s'",
			column.Table().Name(),
			column.Name(),
		)
	}
	columns.order = append(columns.order, column.Column())
	columns.c[column.Name()] = column.Column()
	return nil
}

// All returns all columns as ColumnElems in their default order
func (columns Columns) All() []ColumnElem {
	return columns.order
}

// Get returns a ColumnElem - or an invalid ColumnElem if a column
// with the given name does not exist in Columns
func (columns Columns) Get(name string) ColumnElem {
	col, ok := columns.c[name]
	if !ok {
		return InvalidColumn(name, nil)
	}
	return col
}

// Has returns true if there is a column with the given name in Columns
func (columns Columns) Has(name string) bool {
	_, ok := columns.c[name]
	return ok
}
