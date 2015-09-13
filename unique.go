package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

// UniqueArray is a list of columns representing a table UNIQUE constraint.
type UniqueArray []string

var _ types.Type = UniqueArray{}
var _ Modifier = UniqueArray{}

// Create returns the proper syntax for CREATE TABLE commands.
func (unique UniqueArray) Create(d dialect.Dialect) (string, error) {
	columns := make([]string, len(unique))
	for i, col := range unique {
		columns[i] = fmt.Sprintf(`"%s"`, col)
	}
	return fmt.Sprintf("UNIQUE (%s)", strings.Join(columns, ", ")), nil
}

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
	table.creates = append(table.creates, unique)
	return nil
}

func Unique(names ...string) UniqueArray {
	return UniqueArray(names)
}
