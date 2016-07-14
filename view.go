package sol

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

// ViewElem is a dialect neutral implementation of a SQL view
type ViewElem struct {
	*TableElem
	stmt SelectStmt // The statement that builds the view
}

// Create
func (view ViewElem) Create() (stmt CreateViewStmt) {
	stmt.view = view
	return
}

// View returns a new ViewElem that can be created or queried
func View(name string, stmt SelectStmt) (ViewElem, error) {
	var view ViewElem
	// Are the columns of the SELECT stmt unique?
	columns, err := stmt.columns.MakeUnique()
	if err != nil {
		return view, err
	}

	// Create a new table from the stmt
	view.TableElem = &TableElem{
		name:    name,
		columns: columns,
	}
	view.stmt = stmt
	return view, nil
}

// CreateViewStmt is the internal representation of a CREATE VIEW statement.
type CreateViewStmt struct {
	Stmt
	view        ViewElem
	isTemporary bool
	orReplace   bool
}

// String outputs the parameter-less CREATE View statement in a neutral
// dialect.
func (stmt CreateViewStmt) String() string {
	c, _ := stmt.Compile(&defaultDialect{}, Params())
	return c
}

func (stmt CreateViewStmt) Temporary() CreateViewStmt {
	stmt.isTemporary = true
	return stmt
}

func (stmt CreateViewStmt) OrReplace() CreateViewStmt {
	stmt.orReplace = true
	return stmt
}

// Compile outputs the CREATE VIEW statement using the given dialect and
// parameters. An error may be returned because of a pre-existing error or
// because an error occurred during compilation.
func (stmt CreateViewStmt) Compile(d dialect.Dialect, p *Parameters) (string, error) {
	name := "CREATE"
	if stmt.orReplace {
		name += " OR REPLACE"
	}
	if stmt.isTemporary {
		name = "TEMPORARY"
	}
	name += " VIEW"

	// TODO column aliases

	selectStmt, err := stmt.view.stmt.Compile(d, p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s %s AS (%s)", name, stmt.view.Name(), selectStmt,
	), nil
}
