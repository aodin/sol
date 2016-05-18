package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// PKArray is a list of columns representing the table's primary key
// array. It is an array instead of a map because order is
// significant for primary keys.
type PKArray []string

// Create returns the proper syntax for CREATE TABLE commands.
func (pk PKArray) Create(d dialect.Dialect) (string, error) {
	cols := make([]string, len(pk))
	for i, col := range pk {
		cols[i] = fmt.Sprintf(`"%s"`, col)
	}
	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(cols, ", ")), nil
}

// Has returns true if the PKArray contains the given column name.
func (pk PKArray) Has(name string) bool {
	for _, col := range pk {
		if col == name {
			return true
		}
	}
	return false
}

// Modify implements the TableModifier interface. It confirms that every column
// given exists in the parent table.
func (pk PKArray) Modify(tabular Tabular) error {
	if tabular == nil || tabular.Table() == nil {
		return fmt.Errorf("sol: primary keys cannot modify a nil table")
	}
	table := tabular.Table() // Get the dialect neutral table
	for _, col := range pk {
		if !table.Has(col) {
			return fmt.Errorf(
				"sol: table '%s' does not have a column '%s'. Is it created after the PrimaryKey?",
				table.name,
				col,
			)
		}
	}
	table.pk = pk

	// Add the pk to the create array
	table.creates = append(table.creates, pk)
	return nil
}

// PrimaryKey creates a new PKArray. Only one primary key is allowed
// per table.
func PrimaryKey(names ...string) PKArray {
	return PKArray(names)
}
