package sol

import (
	"database/sql"
	"log"

	"github.com/aodin/sol/dialect"
)

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

// Connection is an alias for Conn
type Connection interface {
	Conn
}

// Transaction is an alias for TX
type Transaction interface {
	TX
}

// DB is a database connection pool. Most functions should use the
// Conn interface instead of this type.
type DB struct {
	*sql.DB
	dialect dialect.Dialect
	panicky bool
}

var _ Conn = &DB{}

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

var _ Conn = &transaction{}
var _ TX = &transaction{}

// Begin simply returns the transaction itself
// TODO database/sql does not support nested transactions, more detail
// here: https://github.com/golang/go/issues/7898
func (tx *transaction) Begin() (TX, error) {
	return tx, nil
}

// Close will commit the transaction unless it has failed
func (tx *transaction) Close() (err error) {
	if tx.successful {
		err = tx.Tx.Commit()
	} else {
		err = tx.Tx.Rollback()
	}
	if tx.panicky && err != nil {
		log.Panic(err)
	}
	return
}

// Commit will attempt to commit the transaction
func (tx *transaction) Commit() error {
	err := tx.Tx.Commit()
	if tx.panicky && err != nil {
		log.Panic(err)
	}
	return err
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

// Rollback will attempt to roll back the transaction
func (tx *transaction) Rollback() error {
	err := tx.Tx.Rollback()
	if tx.panicky && err != nil {
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
