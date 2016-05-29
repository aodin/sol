Sol
===

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/sol) [![Build Status](https://travis-ci.org/aodin/sol.svg?branch=master)](https://travis-ci.org/aodin/sol) [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/aodin/sol)

A SQL toolkit for Go - in the style of Python's [SQLAlchemy Core](http://docs.sqlalchemy.org/en/latest/core/):

* Build complete database schemas
* Create reusable and cross-dialect SQL statements
* Allow struct instances and slices to be directly populated by the database
* Support for MySQL, PostGres, and SQLite3

Quickstart
----------

```go
package main

import (
	"log"

	"github.com/aodin/sol"
	_ "github.com/aodin/sol/sqlite3"
	"github.com/aodin/sol/types"
)

// Database schemas are created using sol's Table function
var Users = sol.Table("users",
	sol.Column("id", types.Integer().NotNull()),
	sol.Column("name", types.Varchar().Limit(32).NotNull()),
	sol.Column("password", types.Varchar().Limit(128).NotNull()),
	sol.PrimaryKey("id"),
)

// Structs can be used to send and receive values
type User struct {
	ID       int64
	Name     string
	Password string
}

func main() {
	// Connect to an in-memory sqlite3 instance
	conn, err := sol.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	// Create the users table
	conn.Query(Users.Create())

	// Insert a user by struct
	admin := User{ID: 1, Name: "admin", Password: "secret"}
	conn.Query(Users.Insert().Values(admin))

	// Select a user - query methods must be given a pointer receiver!
	var user User
	conn.Query(Users.Select(), &user)
	log.Println(user)
}
```

Install
-------

```
go get -u github.com/aodin/sol
```

Usage
-----

Import Sol and at least one database dialect:

```go
import (
    "github.com/aodin/sol"
    _ "github.com/aodin/sol/mysql"
    _ "github.com/aodin/sol/postgres"
    _ "github.com/aodin/sol/sqlite3"
)
```

Calling `Open` will return a `*DB` that implements Sol's `Conn` interface and embeds Go's `*sql.DB`. All queries use the `Query` method, which will return an error if the query fails:

```go
conn, err := sol.Open("sqlite3", ":memory:")
if err != nil {
    log.Panic(err)
}
defer conn.Close()

if err = conn.Query(Users.Create()); err != nil {
	log.Panic(err)
}

var user User
if err = conn.Query(Users.Select(), &user); err != nil {
	log.Panic(err)
}
```

If you'd prefer to have queries panic on error, you can create a panicky version of the connection:

```go
panicky := conn.PanicOnError() // or Must()
panicky.Query(Users.Create())
```

Transactions are started with `Begin` and include the standard methods `Rollback` and `Commit`. There is also a `Close` method which will rollback the transaction unless `IsSuccessful` is called:

```go
tx, _ := conn.PanicOnError().Begin()
defer tx.Close()

tx.Query(Users.Insert().Values(User{Name: "Zero"}))
tx.IsSuccessful()
```

### Statements

SQL can be handwritten using the `Text` function, which requires parameters to be written in a dialect neutral format and passed within a `Values` type:

```go
sol.Text(
    `SELECT * FROM users WHERE id = :id OR name = :name`,
    sol.Values{"name": "admin", "id": 1},
)
```

The parameters will be re-written for the current dialect:

```sql
SELECT * FROM users WHERE id = ? OR name = ?
```

Sol also includes a variety of statements that can be constructed directly from declared schemas.

#### CREATE TABLE

Once a schema has been specified with `Table`, such as:

```go
var Users = sol.Table("users",
	sol.Column("id", types.Integer().NotNull()),
	sol.Column("name", types.Varchar().Limit(32).NotNull()),
	sol.Column("password", types.Varchar().Limit(128).NotNull()),
	sol.PrimaryKey("id"),
)
```

A `CREATE TABLE` statement can be created with:

```go
Users.Create()
```

As with most statements, it will output dialect neutral SQL from its `String()` method. Dialect specific output is created with the `String()` method on the current connection.

```sql
CREATE TABLE users (
  id INTEGER NOT NULL,
  name VARCHAR(32) NOT NULL,
  password VARCHAR(128) NOT NULL,
  PRIMARY KEY (id)
);
```

#### DROP TABLE

Using the `Users` schema, a `DROP TABLE` statement can be created with:

```go
Users.Drop()
```

```sql
DROP TABLE "users"
```

#### INSERT

Insert statements can be created without specifying values. For instance, the method `Insert()` on a schema such as `Users` can be created with:

```go
Users.Insert()
```

When using the Sqlite3 dialect it will generate the following SQL:

```sql
INSERT INTO users (id, name, password) VALUES (?, ?, ?)
```

Values can be inserted to the database using custom struct types or the generic `sol.Values` type. If given a struct, Sol will attempt to match SQL column names to struct field names in a case insensitive manner that is also aware of camel to snake case conversion. More complex names or aliases should specify db struct tags:

```go
type User struct {
	DoesNotMatchColumn int64 `db:"id"`
	Name               string
	Password           string
}
```

#### UPDATE

Rows in a table can be updated using either of:

```go
Users.Update()
sol.Update(Users)
```

Both will produce the following SQL with the sqlite3 dialect:

```sql
UPDATE users SET id = ?, name = ?, password = ?
```

Conditionals can be specified using `Where`:

```go
conn.Query(Users.Update().Values(
	sol.Values{"password": "supersecret"},
).Where(Users.C("name").Equals("admin")))
```

#### DELETE

```go
Users.Delete().Where(Users.C("name").Equals("admin"))
```

```sql
DELETE FROM users WHERE users.name = ?
```

#### SELECT

Results can be queried in a number of ways. Each of the following statements will produce the same SQL output:

```go
Users.Select()
sol.Select(Users)
sol.Select(Users.C("id"), Users.C("name"), Users.C("password"))
```

```sql
SELECT users.id, users.name, users.password FROM users
```

Multiple results can be returned directly into slice of structs:

```go
var users []User
conn.Query(Users.Select(), &users)
```

Single column selections can select directly into a slice of the appropriate type:

```go
var ids []int64
conn.Query(sol.Select(Users.C("id")), &ids)
```

### Table Schema

Tables can be constructed with foreign keys, unique constraints, and composite primary keys. See the `sol_test.go` file for more examples:

```go
var Contacts = sol.Table("contacts",
	sol.Column("id", types.Integer()),
	sol.ForeignKey("user_id", Users),
	sol.Column("key", types.Varchar()),
	sol.Column("value", types.Varchar()),
	sol.PrimaryKey("id"),
	sol.Unique("user_id", "key"),
)
```

```sql
CREATE TABLE contacts (
  id INTEGER,
  user_id INTEGER NOT NULL REFERENCES users(id),
  key VARCHAR,
  value VARCHAR,
  PRIMARY KEY (id),
  UNIQUE (user_id, key)
);
```

Develop
-------

Some dialects require a configuration to be set via an environmental variable for testing, such as `SOL_TEST_POSTGRES` for the `postgres` dialect. Example variables and, if possible, [Docker](https://www.docker.com/) containers have been provided in the subpackages where these variables are required.

Statements and clauses can be tested by creating a new dialect-specific tester; for example using the `postgres` package:

```go
expect := NewTester(t, postgres.Dialect())
```

The instance's `SQL` method will test expected output and parameterization:

```go
expect.SQL(Users.Delete(), `DELETE FROM users`)

expect.SQL(
    Users.Insert().Values(Values{"id": 1, "name": "user"}),
    `INSERT INTO users (id, name) VALUES ($1, $2)`,
    1, "user",
)
```

And the `Error` method will test that an error occurred:

```go
expect.Error(sql.Select(Users.C("does-not-exist")))
```

Happy Hacking!

aodin, 2015-2016
