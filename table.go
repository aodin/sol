package sol

import (
	"fmt"
	"log"

	"github.com/aodin/sol/types"
)

// Tabular is the interface that all dialects of a SQL table must implement
type Tabular interface {
	// Two methods for neutral SQL element interfaces:
	// (1) Require the interface to return the neutral implementation
	// (2) Enumerate all the methods an implmentation would require
	// Columnar and Tabular both use method (1)
	// Name has been left as a legacy shortcut but may be removed
	Selectable
	Name() string
	Table() *TableElem
}

// TableElem is a dialect neutral implementation of a SQL table
type TableElem struct {
	name         string
	alias        string
	columns      UniqueColumnSet
	pk           PKArray // Table's primary key
	uniques      []UniqueArray
	fks          []FKElem // This table's foreign keys
	referencedBy []FKElem // Foreign keys that reference this table
	creates      []types.Type
}

var _ Tabular = &TableElem{}

// Column returns the column as a ColumnElem. If the column does not exist
// it will return the ColumnElem in an invalid state that will be used to
// construct an error message
func (table TableElem) Column(name string) ColumnElem {
	col := table.columns.Get(name)
	// If the column is invalid add the current table in order
	// to construct a better error message
	if col.IsInvalid() {
		col.table = &table
	}
	return col
}

// C is an alias for Column
func (table TableElem) C(name string) ColumnElem {
	return table.Column(name)
}

// Columns returns all the table columns in the original schema order
func (table TableElem) Columns() []ColumnElem {
	return table.columns.All()
}

// Create generates the table's CREATE statement.
func (table *TableElem) Create() CreateStmt {
	return CreateStmt{table: table}
}

// Delete is an alias for Delete(table). It will generate a DELETE statement
// for the entire table
func (table *TableElem) Delete() DeleteStmt {
	return Delete(table)
}

// Create generates the table's DROP statement.
func (table *TableElem) Drop() DropStmt {
	return DropStmt{table: table}
}

// ForeignKeys returns the table's foreign keys
func (table *TableElem) ForeignKeys() []FKElem {
	return table.fks
}

// Has returns true if the column exists in this table
func (table *TableElem) Has(name string) bool {
	return table.columns.Has(name)
}

// Insert is an alias for Insert(table). It will create an INSERT statement
// for the entire table. Specify the insert values with the method Values().
func (table *TableElem) Insert() InsertStmt {
	return Insert(table)
}

// Name returns the table name without escaping
func (table *TableElem) Name() string {
	return fmt.Sprintf(`%s`, table.name)
}

// PrimaryKey returns the primary key array
func (table TableElem) PrimaryKey() PKArray {
	return table.pk
}

// ReferencedBy returns the foreign keys that reference this table
func (table *TableElem) ReferencedBy() []FKElem {
	return table.referencedBy
}

// Select returns a SelectStmt for the entire table
func (table *TableElem) Select(selections ...Selectable) (stmt SelectStmt) {
	return SelectTable(table, selections...)
}

// Table returns the table itself
func (table *TableElem) Table() *TableElem {
	return table
}

// Update is an alias for Update(table). It will create an UPDATE statement
// for the entire table. Specify the update values with the method Values().
func (table *TableElem) Update() UpdateStmt {
	return Update(table)
}

// Table creates a new dialect netural table. It will panic on any errors.
func Table(name string, modifiers ...Modifier) *TableElem {
	if err := isValidTableName(name); err != nil {
		log.Panic(err)
	}
	table := &TableElem{
		name:    name,
		columns: UniqueColumns(),
	}
	for _, modifier := range modifiers {
		if err := modifier.Modify(table); err != nil {
			log.Panic(err)
		}
	}
	return table
}
