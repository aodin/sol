package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

// DropStmt is the internal representation of an DROP TABLE statement.
type DropStmt struct {
	table    *TableElem
	ifExists bool
}

// IfExists adds the IF EXISTS modifier to a DROP TABLE statement.
func (stmt DropStmt) IfExists() DropStmt {
	stmt.ifExists = true
	return stmt
}

// String outputs the parameter-less CREATE TABLE statement in a neutral
// dialect.
func (stmt DropStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

// Compile outputs the DROP TABLE statement using the given dialect and
// parameters.
func (stmt DropStmt) Compile(d dialect.Dialect, p *Parameters) (string, error) {
	if stmt.ifExists {
		return fmt.Sprintf(`DROP TABLE IF EXISTS %s`, stmt.table.Name()), nil
	}
	return fmt.Sprintf(`DROP TABLE %s`, stmt.table.Name()), nil
}
