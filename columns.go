package sol

import (
	"fmt"
)

// Columns maps column to name to a ColumnElem. It also maintains an order
// of columns.
type Columns struct {
	order []Columnar
	c     map[string]Columnar
}

func (columns *Columns) add(column Columnar) error {
	// Tables cannot have duplicate column names
	// TODO the column should already be assigned a table
	if columns.Has(column.Name()) {
		return fmt.Errorf(
			"sol: table '%s' already has a column named '%s'",
			column.Table().name,
			column.Name(),
		)
	}
	columns.order = append(columns.order, column)
	columns.c[column.Name()] = column
	return nil
}

// All returns all columns as ColumnElems in their default order
func (columns Columns) All() []Columnar {
	return columns.order
}

func (columns Columns) Get(name string) Columnar {
	return columns.c[name]
}

// Has returns true if there is a column with the given name in Columns
func (columns Columns) Has(name string) bool {
	_, ok := columns.c[name]
	return ok
}
