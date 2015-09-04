package sol

import (
	"database/sql"
	"reflect"

	"github.com/aodin/sol/dialect"
)

// Connection is an alias for Conn
type Connection interface {
	Conn
}

// Transaction is an alias for TX
type Transaction interface {
	TX
}

type Conn interface {
	Begin() TX
	Query(stmt Executable) error
}

type TX interface {
	Conn
	Commit()
	Rollback()
}

type conn struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func (c *conn) Begin() TX {
	return nil
}

func (c *conn) Close() error {
	return c.db.Close()
}

func (c *conn) Dialect() dialect.Dialect {
	return c.dialect
}

func (c *conn) compile(stmt Executable) (string, *Parameters, error) {
	// Initialize a list of empty parameters
	params := Params()

	// Compile with the database connection's current dialect
	compiled, err := stmt.Compile(c.dialect, params)
	return compiled, params, err
}

// information on rows affected and last ID inserted depending on the driver.
func (c *conn) Execute(stmt Executable) (sql.Result, error) {
	compiled, params, err := c.compile(stmt)
	if err != nil {
		return nil, err
	}
	return c.db.Exec(compiled, *params...)
}

// Query executes an Executable statement.
func (c *conn) Query(stmt Executable, dest ...interface{}) error {
	if len(dest) == 0 {
		_, err := c.Execute(stmt)
		return err
	}

	if len(dest) > 1 {
		return c.QueryAll(stmt, dest)
	}

	t := reflect.Indirect(reflect.ValueOf(dest[0]))
	if t.Kind() == reflect.Slice {
		return c.QueryAll(stmt, dest[0])
	}
	return c.QueryOne(stmt, dest[0])
}

// Query executes an Executable statement with the optional arguments. It
// returns a Result object, that can scan rows in various data types.
func (c *conn) query(stmt Executable) (*Result, error) {
	compiled, params, err := c.compile(stmt)
	if err != nil {
		return nil, err
	}

	rows, err := c.db.Query(compiled, *params...)
	if err != nil {
		return nil, err
	}
	// Wrap the sql rows in a result
	return &Result{Scanner: rows, stmt: compiled}, nil
}

// QueryAll will query the statement and populate the given destination
// interface with all results.
func (c *conn) QueryAll(stmt Executable, dest interface{}) error {
	result, err := c.query(stmt)
	if err != nil {
		return err
	}
	return result.All(dest)
}

// QueryOne will query the statement and populate the given destination
// interface with a single result.
func (c *conn) QueryOne(stmt Executable, dest interface{}) error {
	result, err := c.query(stmt)
	if err != nil {
		return err
	}
	// Close the result rows or sqlite3 will open another connection
	defer result.Close()
	return result.One(dest)
}

// String returns parameter-less SQL. If an error occurred during compilation,
// then the string output of the error will be returned.
func (c *conn) String(stmt Executable) string {
	compiled, err := stmt.Compile(c.dialect, Params())
	if err != nil {
		return err.Error()
	}
	return compiled
}

// Connect connects to the database using the given driver and credentials.
// It returns a database connection pool and an error if one occurred.
func Open(driver, credentials string) (*conn, error) {
	db, err := sql.Open(driver, credentials)
	if err != nil {
		return nil, err
	}

	// Get the dialect
	d, err := dialect.Get(driver)
	if err != nil {
		return nil, err
	}
	return &conn{db: db, dialect: d}, nil
}
