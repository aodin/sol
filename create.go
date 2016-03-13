package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// For the Postgres implementation:
// http://www.postgresql.org/docs/9.4/static/sql-createtable.html

// CreateStmt is the internal representation of an CREATE TABLE statement.
type CreateStmt struct {
	table       *TableElem
	ifNotExists bool
	isTemporary bool
}

// String outputs the parameter-less CREATE TABLE statement in a neutral
// dialect.
func (stmt CreateStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

func (stmt CreateStmt) IfNotExists() CreateStmt {
	stmt.ifNotExists = true
	return stmt
}

func (stmt CreateStmt) Temporary() CreateStmt {
	stmt.isTemporary = true
	return stmt
}

// Compile outputs the CREATE TABLE statement using the given dialect and
// parameters. An error may be returned because of a pre-existing error or
// because an error occurred during compilation.
func (stmt CreateStmt) Compile(d dialect.Dialect, p *Parameters) (string, error) {
	// Compiled elements
	compiled := make([]string, len(stmt.table.creates))

	var err error
	for i, create := range stmt.table.creates {
		if compiled[i], err = create.Create(d); err != nil {
			return "", err
		}
	}

	var name string
	if stmt.isTemporary {
		name = "CREATE TEMPORARY TABLE"
	} else {
		name = "CREATE TABLE"
	}
	if stmt.ifNotExists {
		name += " IF NOT EXISTS"
	}

	return fmt.Sprintf(
		"%s \"%s\" (\n  %s\n);",
		name,
		stmt.table.Name(),
		strings.Join(compiled, ",\n  "),
	), nil
}
