package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// Selectable is an interface that allows both tables and columns to be
// selected. It is implemented by TableElem and ColumnElem.
type Selectable interface {
	Columns() []ColumnElem
}

// SelectStmt is the internal representation of an SQL SELECT statement.
type SelectStmt struct {
	ConditionalStmt
	tables  []*TableElem
	columns []ColumnElem
	limit   int
	offset  int
}

func (stmt SelectStmt) compileColumns() []string {
	names := make([]string, len(stmt.columns))
	for i, col := range stmt.columns {
		names[i] = col.FullName()
	}
	return names
}

func (stmt SelectStmt) compileTables() []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = table.Name()
	}
	return names
}

func (stmt SelectStmt) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	compiled := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(stmt.compileColumns(), ", "),
		strings.Join(stmt.compileTables(), ", "),
	)
	if stmt.limit != 0 {
		compiled += fmt.Sprintf(" LIMIT %d", stmt.limit)
	}
	if stmt.offset != 0 {
		compiled += fmt.Sprintf(" OFFSET %d", stmt.offset)
	}
	return compiled, nil
}

func (stmt SelectStmt) hasTable(name string) bool {
	for _, table := range stmt.tables {
		if table.Name() == name {
			return true
		}
	}
	return false
}

// Limit sets the limit of the SELECT statement.
func (stmt SelectStmt) Limit(limit int) SelectStmt {
	// TODO Error (or warning) if limit was already set
	stmt.limit = limit
	return stmt
}

// Offset sets the offset of the SELECT statement.
func (stmt SelectStmt) Offset(offset int) SelectStmt {
	// TODO Error (or warning) if offset was already set
	stmt.offset = offset
	return stmt
}

// TODO SelectColumn

func SelectTable(table *TableElem, dest ...interface{}) (stmt SelectStmt) {
	stmt.tables = []*TableElem{table}

	// Add the columns from the alias
	stmt.columns = table.columns.order
	return
}

func Select(selections ...Selectable) (stmt SelectStmt) {
	columns := make([]ColumnElem, 0)
	for _, selection := range selections {
		if selection == nil {
			stmt.AddMeta("sol: received a nil selectable in Select()")
			return
		}
		columns = append(columns, selection.Columns()...)
	}

	if len(columns) < 1 {
		stmt.AddMeta("sol: Select() must be given at least one column")
		return
	}

	for _, column := range columns {
		if column.invalid {
			// TODO field error
			stmt.AddMeta("sol: selected column does not exist")
			return
		}
		stmt.columns = append(stmt.columns, column)

		// Add the table to the stmt tables if it does not already exist
		if !stmt.hasTable(column.Table().Name()) {
			stmt.tables = append(stmt.tables, column.Table())
		}
	}
	return
}
