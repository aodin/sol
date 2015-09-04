package sol

import (
	"fmt"
)

// UniqueArray is a list of columns representing a table UNIQUE constraint.
type UniqueArray []string

// Has returns true if the UniqueArray contains the given column name.
func (unique UniqueArray) Has(name string) bool {
	for _, col := range unique {
		if col == name {
			return true
		}
	}
	return false
}

// Modify implements the TableModifier interface. It confirms that every column
// given exists in the parent table.
func (unique UniqueArray) Modify(table *TableElem) error {
	for _, col := range unique {
		if !table.Has(col) {
			return fmt.Errorf(
				"sol: table '%s' does not have a column '%s'. Is it created after Unique()?",
				table.name,
				col,
			)
		}
	}
	table.uniques = append(table.uniques, unique)

	// TODO Add the unique to the create array
	return nil
}

func Unique(names ...string) UniqueArray {
	return UniqueArray(names)
}
