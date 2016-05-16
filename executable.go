package sol

import (
	"database/sql"
	"reflect"

	"github.com/aodin/sol/dialect"
)

// Executable is - for now - an alias of Compiles
type Executable interface {
	Compiles
}

// executer is a common interface that database/sql *DB and *Tx can share
type executer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

var _ executer = &sql.DB{}
var _ executer = &sql.Tx{}

func compile(d dialect.Dialect, stmt Executable) (string, *Parameters, error) {
	// Initialize a list of empty parameters
	params := Params()

	// Compile with the database connection's current dialect
	compiled, err := stmt.Compile(d, params)
	return compiled, params, err
}

func execute(exec executer, d dialect.Dialect, stmt Executable) (sql.Result, error) {
	compiled, params, err := compile(d, stmt)
	if err != nil {
		return nil, err
	}
	return exec.Exec(compiled, *params...)
}

func perform(exec executer, d dialect.Dialect, stmt Executable, dest ...interface{}) error {
	if len(dest) == 0 {
		_, err := execute(exec, d, stmt)
		return err
	}

	if len(dest) > 1 {
		return queryAll(exec, d, stmt, dest)
	}

	t := reflect.Indirect(reflect.ValueOf(dest[0]))
	if t.Kind() == reflect.Slice {
		return queryAll(exec, d, stmt, dest[0])
	}
	return queryOne(exec, d, stmt, dest[0])
}

func query(exec executer, d dialect.Dialect, stmt Executable) (*Result, error) {
	compiled, params, err := compile(d, stmt)
	if err != nil {
		return nil, err
	}

	rows, err := exec.Query(compiled, *params...)
	if err != nil {
		return nil, err
	}
	// Wrap the sql rows in a result
	return &Result{Scanner: rows, stmt: compiled}, nil
}

// QueryAll will query the statement and populate the given destination
// interface with all results.
func queryAll(exec executer, d dialect.Dialect, stmt Executable, dest interface{}) error {
	result, err := query(exec, d, stmt)
	if err != nil {
		return err
	}
	return result.All(dest)
}

// QueryOne will query the statement and populate the given destination
// interface with a single result.
func queryOne(exec executer, d dialect.Dialect, stmt Executable, dest interface{}) error {
	result, err := query(exec, d, stmt)
	if err != nil {
		return err
	}
	// Close the result rows or sqlite3 will open another connection
	defer result.Close()
	return result.One(dest)
}
