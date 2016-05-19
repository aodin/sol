package sol

import "fmt"

// ColumnSet maintains a []ColumnElem. It includes a variety of
// getters and setter. Optionally, it can force unique
type ColumnSet struct {
	// TODO what about uniqueness with aliases? across tables?
	unique bool
	order  []ColumnElem
}

// Add adds any number of ColumnElem types to the set and returns the new set.
// If the set is marked unique, adding a column with the same name
// as an existing column in the set will return an error.
func (set ColumnSet) Add(columns ...ColumnElem) (ColumnSet, error) {
	if set.unique {
		for _, column := range columns {
			for _, existing := range set.order {
				// TODO across tables? aliases?
				if existing.Name() == column.Name() {
					if existing.Table() == nil {
						return set, fmt.Errorf(
							"sol: this set already has a column named '%s'",
							existing.Name(),
						)
					}
					return set, fmt.Errorf(
						"sol: table '%s' already has a column named '%s'",
						existing.Table().Name(),
						existing.Name(),
					)
				}
			}
			set.order = append(set.order, column)
		}
	} else {
		set.order = append(set.order, columns...)
	}
	return set, nil
}

// All returns all columns in their default order
func (set ColumnSet) All() []ColumnElem {
	return set.order
}

// Get returns a ColumnElem - or an invalid ColumnElem if a column
// with the given name does not exist in the set
func (set ColumnSet) Get(name string) ColumnElem {
	for _, column := range set.order {
		// TODO What about table? aliases?
		if column.Name() == name {
			return column
		}
	}
	return InvalidColumn(name, nil)
}

// Has returns true if there is a column with the given name in the ColumnSet
func (set ColumnSet) Has(name string) bool {
	return set.Get(name).IsValid()
}

// UniqueColumns creates a new ColumnSet that can only hold columns
// with unique names
func UniqueColumns() ColumnSet {
	return ColumnSet{unique: true}
}

// Columns creates a new ColumnSet
func Columns() ColumnSet {
	return ColumnSet{}
}
