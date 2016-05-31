package mysql

import (
	_ "github.com/go-sql-driver/mysql" // Register the MySQL driver

	"github.com/aodin/sol/dialect"
)

// MySQL implements the Dialect interface for MySQL databases.
type MySQL struct{}

// The MySQL dialect must implement the dialect.Dialect interface
var _ dialect.Dialect = &MySQL{}

// Param returns the MySQL specific parameterization scheme.
func (d *MySQL) Param(i int) string {
	return `?`
}

// Dialect is a constructor for the MySQL Dialect
func Dialect() *MySQL {
	return &MySQL{}
}

// Add the MySQL dialect to the dialect registry
func init() {
	dialect.Register("mysql", Dialect())
}
