package sol

import (
	"fmt"
	"log"

	"github.com/aodin/sol/types"
)

type Tabular interface {
	Name() string
}

type TableElem struct {
	name    string
	alias   string
	columns Columns
	pk      PKArray // Table's primary key
	uniques []UniqueArray
	fks     []FKElem // This table's foreign keys
	reverse []FKElem // Tables that link to this table
	creates []types.Type
}

// Column returns the column as a ColumnElem. If the column does not exist
// it will return the ColumnElem in an invalid state that will be used to
// construct an error message
func (table TableElem) Column(name string) ColumnElem {
	if table.Has(name) {
		switch elem := table.GetColumn(name).(type) {
		case ColumnElem:
			return elem
		}
		// TODO invalid column?
	}
	return InvalidColumn(name, &table)
}

// C is an alias for Column
func (table TableElem) C(name string) ColumnElem {
	return table.Column(name)
}

// Columns returns all the table columns in the original schema order
func (table TableElem) Columns() []Columnar {
	return table.columns.order
}

// Create generates the table's CREATE statement.
func (table *TableElem) Create() CreateStmt {
	return CreateStmt{table: table}
}

// Delete is an alias for Delete(table). It will generate a DELETE statement
// for the entire table. Conditionals joined with AND can be passed as
// parameters or later added with the Where() method
func (table *TableElem) Delete(clauses ...Clause) DeleteStmt {
	return Delete(table, clauses...)
}

// Create generates the table's DROP statement.
func (table *TableElem) Drop() DropStmt {
	return DropStmt{table: table}
}

func (table TableElem) GetColumn(name string) Columnar {
	return table.columns.Get(name)
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

// Select returns a SelectStmt for the entire table
func (table *TableElem) Select(selections ...Selectable) (stmt SelectStmt) {
	return SelectTable(table, selections...)
}

// Update is an alias for Update(table). It will create an UPDATE statement
// for the entire table. Specify the update values with the method Values().
func (table *TableElem) Update() UpdateStmt {
	return Update(table)
}

// Table creates a new table element. It will panic on any errors.
func Table(name string, modifiers ...Modifier) *TableElem {
	if err := isValidTableName(name); err != nil {
		log.Panic(err)
	}
	table := &TableElem{
		name:    name,
		columns: Columns{c: make(map[string]Columnar)},
	}
	for _, modifier := range modifiers {
		if err := modifier.Modify(table); err != nil {
			log.Panic(err)
		}
	}
	return table
}
