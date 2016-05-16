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

// Connection is an alias for Conn
type Connection interface {
	Conn
}

// Transaction is an alias for TX
type Transaction interface {
	TX
}

// Conn is the common database connection interface. It can perform
// queries and dialect specific compilation. Since transactions also
// implement this interface, it is highly recommended to pass the Conn
// interface to functions that do not need to modify transactional state.
type Conn interface {
	Begin() (TX, error)
	Close() error
	Query(stmt Executable, dest ...interface{}) error
	String(stmt Executable) string
}

// TX is the interface for a transaction. In addition to standard COMMIT
// and ROLLBACK behavior, the interface includes a generic Close method
// that can be used with defer. Close will rollback the transaction
// unless the IsSuccessful method has been called.
type TX interface {
	Conn
	Commit() error
	IsSuccessful()
	Rollback() error
}

// DB is a database connection pool. Most functions should use the
// Conn interface instead of this type.
type DB struct {
	*sql.DB
	dialect dialect.Dialect
	panicky bool
}

// Begin will start a new transaction on the current connection pool
func (c *DB) Begin() (TX, error) {
	tx, err := c.DB.Begin()
	if c.panicky && err != nil {
		log.Panic(err)
	}
	return &transaction{Tx: tx, dialect: c.dialect, panicky: c.panicky}, err
}

// Close will make the current connection pool unusable
func (c *DB) Close() error {
	err := c.DB.Close()
	if c.panicky && err != nil {
		log.Panic(err)
	}
	return err
}

// Dialect returns the current connection pool's dialect, e.g. sqlite3
func (c *DB) Dialect() dialect.Dialect {
	return c.dialect
}

// Query executes an Executable statement
func (c *DB) Query(stmt Executable, dest ...interface{}) error {
	err := perform(c.DB, c.dialect, stmt, dest...)
	if c.panicky && err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}
	return err
}

// String returns the compiled Executable using the DB's dialect.
// If an error is encountered during compilation, it will return the
// error instead.
func (c *DB) String(stmt Executable) string {
	compiled, err := stmt.Compile(c.dialect, Params())
	if err != nil {
		return err.Error()
	}
	return compiled
}

// PanicOnError will create a new connection that will panic on any error
func (c DB) PanicOnError() *DB {
	c.panicky = true
	return &c
}

// Must is an alias for PanicOnError
func (c DB) Must() *DB {
	return c.PanicOnError()
}

// Open connects to the database using the given driver and credentials.
// It returns a database connection pool and an error if one occurred.
func Open(driver, credentials string) (*DB, error) {
	db, err := sql.Open(driver, credentials)
	if err != nil {
		return nil, err
	}

	// Get the dialect
	d, err := dialect.Get(driver)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, dialect: d}, nil
}

type transaction struct {
	*sql.Tx
	dialect    dialect.Dialect
	successful bool
	panicky    bool
}

// Begin simply returns the transaction itself
// TODO database/sql does not support nested transactions, more detail
// here: https://github.com/golang/go/issues/7898
func (tx *transaction) Begin() (TX, error) {
	return tx, nil
}

// Close will commit the transaction unless it has failed
func (tx *transaction) Close() (err error) {
	if tx.successful {
		err = tx.Commit()
	} else {
		err = tx.Rollback()
	}
	if tx.panicky && err != nil {
		log.Panic(err)
	}
	return
}

// IsSuccessful will mark the transaction as successful, changing
// the behavior of Close()
func (tx *transaction) IsSuccessful() {
	tx.successful = true
}

// Query executes an Executable statement
func (tx *transaction) Query(stmt Executable, dest ...interface{}) error {
	err := perform(tx.Tx, tx.dialect, stmt, dest...)
	if tx.panicky && err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}
	return err
}

// String returns the compiled Executable using the transaction's dialect.
// If an error is encountered during compilation, it will return the
// error instead.
func (tx *transaction) String(stmt Executable) string {
	compiled, err := stmt.Compile(tx.dialect, Params())
	if err != nil {
		return err.Error()
	}
	return compiled
}
