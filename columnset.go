package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// ColumnSet maintains a []ColumnElem. It includes a variety of
// getters, setters, and filters.
type ColumnSet struct {
	order []ColumnElem
}

func (set ColumnSet) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	names := make([]string, len(set.order))
	for i, col := range set.order {
		compiled, err := col.Compile(d, ps)
		if err != nil {
			return "", err
		}
		if col.Alias() != "" {
			compiled += fmt.Sprintf(` AS "%s"`, col.Alias())
		}
		names[i] = compiled
	}
	return strings.Join(names, ", "), nil
}

// Add adds any number of Columnar types to the set and returns the new set.
func (set ColumnSet) Add(columns ...Columnar) (ColumnSet, error) {
	for _, column := range columns {
		set.order = append(set.order, column.Column())
	}
	return set, nil
}

// All returns all columns in their default order
func (set ColumnSet) All() []ColumnElem {
	return set.order
}

// Exists returns true if there is at least one column in the set
func (set ColumnSet) Exists() bool {
	return len(set.order) > 0
}

// Filters returns a ColumnSet of columns from the original set that
// match the given names
func (set ColumnSet) Filter(names ...string) (out ColumnSet) {
	for _, column := range set.order {
		for _, name := range names {
			if column.Name() == name {
				out.order = append(out.order, column)
				break
			}
		}
	}
	return
}

// FullNames returns the full names of the set's columns without alias
func (set ColumnSet) FullNames() []string {
	names := make([]string, len(set.order))
	for i, col := range set.order {
		names[i] = col.FullName()
	}
	return names
}

// Get returns a ColumnElem - or an invalid ColumnElem if a column
// with the given name does not exist in the set
func (set ColumnSet) Get(name string) ColumnElem {
	for _, column := range set.order {
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

// IsEmpty returns true if there are no columns in this set
func (set ColumnSet) IsEmpty() bool {
	return len(set.order) == 0
}

// Names returns the names of the set's columns without alias
func (set ColumnSet) Names() []string {
	names := make([]string, len(set.order))
	for i, col := range set.order {
		names[i] = col.Name()
	}
	return names
}

// Reject removes the given column names and returns a new ColumnSet
func (set ColumnSet) Reject(names ...string) (out ColumnSet) {
ColumnLoop:
	for _, column := range set.order {
		for _, name := range names {
			if column.Name() == name {
				continue ColumnLoop
			}
		}
		out.order = append(out.order, column)
	}
	return
}

// Columns creates a new ColumnSet
func Columns(columns ...ColumnElem) ColumnSet {
	return ColumnSet{order: columns}
}

// UniqueColumnSet is a ColumnSet, but column names must be unique
type UniqueColumnSet struct{ ColumnSet }

// Add adds any number of Columnar types to the set and returns the new set.
// Adding a column with the same name as an existing column in the set
// will return an error.
func (set UniqueColumnSet) Add(columns ...Columnar) (UniqueColumnSet, error) {
	for _, column := range columns {
		for _, existing := range set.order {
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
		set.order = append(set.order, column.Column())
	}
	return set, nil
}

// EmptyValues creates an empty Values from the UniqueColumnSet
func (set UniqueColumnSet) EmptyValues() Values {
	values := Values{}
	for _, column := range set.order {
		values[column.Name()] = nil
	}
	return values
}

// Filters returns a UniqueColumnSet of columns from the original
// set that match the given names
func (set UniqueColumnSet) Filter(names ...string) (out UniqueColumnSet) {
	return UniqueColumnSet{ColumnSet: set.ColumnSet.Filter(names...)}
}

// Reject removes the given column names and returns a new ColumnSet
func (set UniqueColumnSet) Reject(names ...string) (out UniqueColumnSet) {
	return UniqueColumnSet{ColumnSet: set.ColumnSet.Reject(names...)}
}

// UniqueColumns creates a new ColumnSet that can only hold columns
// with unique names
func UniqueColumns() (set UniqueColumnSet) {
	return
}
