package sol

import (
	"database/sql"
	"log"
	"reflect"

	"github.com/aodin/sol/dialect"
)

// executer is a common interface that database/sql *DB and *Tx can share
type executer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

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

// Connection is an alias for Conn
type Connection interface {
	Conn
}

// Transaction is an alias for TX
type Transaction interface {
	TX
}

type Conn interface {
	Begin() (TX, error)
	Close() error
	Query(stmt Executable, dest ...interface{}) error
	String(stmt Executable) string
}

type TX interface {
	Conn
	Commit() error
	IsSuccessful()
	Rollback() error
}

type conn struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// Begin will start a new transaction on the current connection pool
func (c *conn) Begin() (TX, error) {
	tx, err := c.db.Begin()
	return &transaction{Tx: tx, dialect: c.dialect}, err
}

// Close will make the current connection pool unusable
func (c *conn) Close() error {
	return c.db.Close()
}

// Dialect returns the current connection pool's dialect, e.g. sqlite3
func (c *conn) Dialect() dialect.Dialect {
	return c.dialect
}

// Query executes an Executable statement.
func (c *conn) Query(stmt Executable, dest ...interface{}) error {
	return perform(c.db, c.dialect, stmt, dest...)
}

// String returns parameter-less SQL. If an error occurred during compilation,
// then the string output of the error will be returned.
// TODO Common string function
func (c *conn) String(stmt Executable) string {
	compiled, err := stmt.Compile(c.dialect, Params())
	if err != nil {
		return err.Error()
	}
	return compiled
}

// PanicOnError returns a panicConn, which implements the Connection
// interface but will also panic on any error
func (c *conn) PanicOnError() panicConn {
	return panicConn{c}
}

// Must is an alias for PanicOnError
func (c *conn) Must() panicConn {
	return c.PanicOnError()
}

// Open connects to the database using the given driver and credentials.
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

// panicConn is a connection that panics on error
type panicConn struct {
	*conn
}

// Begin calls the internal conn Begin(), but will panic on any error
// and returns a panicTx instead of a transaction
func (c panicConn) Begin() (TX, error) {
	tx, err := c.conn.Begin()
	if err != nil {
		log.Panic(err)
	}
	return panicTx{tx.(*transaction)}, err
}

// Close calls the internal conn Close(), but will panic on any error
func (c panicConn) Close() error {
	err := c.conn.Close()
	if err != nil {
		log.Panic(err)
	}
	return err
}

// Query calls the internal conn Query(), but will panic on any error
func (c panicConn) Query(stmt Executable, dest ...interface{}) error {
	err := c.conn.Query(stmt, dest...)
	if err != nil {
		log.Panic(err)
	}
	return err
}

type transaction struct {
	*sql.Tx
	dialect    dialect.Dialect
	successful bool
}

// Begin simply returns the transaction itself
// TODO Are nested transactions possible? Or should this error?
func (tx *transaction) Begin() (TX, error) {
	return tx, nil
}

// Close will commit the transaction unless it has failed
func (tx *transaction) Close() error {
	if tx.successful {
		return tx.Commit()
	}
	return tx.Rollback()
}

func (tx *transaction) IsSuccessful() {
	tx.successful = true
}

// Query executes an Executable statement
func (tx *transaction) Query(stmt Executable, dest ...interface{}) error {
	return perform(tx.Tx, tx.dialect, stmt, dest...)
}

func (tx *transaction) String(stmt Executable) string {
	compiled, err := stmt.Compile(tx.dialect, Params())
	if err != nil {
		return err.Error()
	}
	return compiled
}

type panicTx struct {
	*transaction
}

// Begin returns the panicTx itself
// TODO Are nested transactions possible? Or should this error?
func (tx panicTx) Begin() (TX, error) {
	return tx, nil
}

// Close will call the internal transaction Close(). It will panic on error.
func (tx panicTx) Close() error {
	err := tx.transaction.Close()
	if err != nil {
		log.Panic(err)
	}
	return err
}

// Query will call the internal transaction Query(). It will panic on error.
func (tx panicTx) Query(stmt Executable, dest ...interface{}) error {
	err := tx.transaction.Query(stmt, dest...)
	if err != nil {
		log.Panic(err)
	}
	return err
}
