package postgres

import "github.com/aodin/sol"

// TableElem is a postgres specific implementation of a table. Major
// differences include:
// * Column and C methods returns postgres specific columns
// * Insert return postgres specific INSERT statements with RETURNING syntax
type TableElem struct {
	*sol.TableElem
}

var _ sol.Tabular = &TableElem{}

// Column will return a postgres specific ColumnElem rather than a generic
// ColumnElem
func (table TableElem) Column(name string) ColumnElem {
	if table.Has(name) {
		switch elem := table.GetColumn(name).(type) {
		case ColumnElem:
			return elem
		case sol.ColumnElem:
			return ColumnElem{ColumnElem: elem}
		}
		// TODO invalid column? Prevent the mixing of column types?
	}
	return ColumnElem{ColumnElem: sol.InvalidColumn(name, table)}
}

// C is an alias for Column
func (table TableElem) C(name string) ColumnElem {
	return table.Column(name)
}

// Insert creates a postgres.InsertStmt from the table
func (table *TableElem) Insert() InsertStmt {
	return Insert(table)
}

// Table creates a new table element. It will panic on any errors.
func Table(name string, modifiers ...sol.Modifier) *TableElem {
	return &TableElem{TableElem: sol.Table(name, modifiers...)}
}
