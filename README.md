Sol
===

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/sol) [![Build Status](https://travis-ci.org/aodin/sol.svg?branch=master)](https://travis-ci.org/aodin/sol) [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/aodin/sol)

A relational database toolkit for Go - in the style of Python's [SQLAlchemy Core](http://docs.sqlalchemy.org/en/latest/core/):

* Build complete database schemas
* Create reusable and cross-dialect SQL statements
* Allow struct instances and slices to be directly populated by the database

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

// Create a database schema using sol's Table function
var Users = sol.Table("users",
	sol.Column("id", types.Integer().NotNull()),
	sol.Column("name", types.Varchar().Limit(32).NotNull()),
	sol.Column("password", types.Varchar().Limit(128).NotNull()),
	sol.PrimaryKey("id"),
)

// Structs are used to send and receive values to the database
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

	// Insert a user - they can be inserted by value or reference
	admin := User{ID: 1, Name: "admin", Password: "secret"}
	conn.Query(Users.Insert().Values(admin))

	// Select a user - query methods must be given a pointer
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

Import Sol and at least one database dialect.

```go
import (
    "github.com/aodin/sol"
    _ "github.com/aodin/sol/postgres"
    _ "github.com/aodin/sol/sqlite3"
)
```

Calling `Open` will return a `*DBConn` that implements the `Conn` interface. All queries can be performed the `Query` method, which will return an error if the query fails.

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

If you'd prefer to have the query panic on error, you can create a panicky version of the connection.

```go
panicky := conn.PanicOnError()
panicky.Query(Users.Create())
```

Transactions are started with `Begin` and include the standard methods `Rollback` and `Commit`. There is also a `Close` method which will rollback the transaction unless `IsSuccessful` is called.

```go
tx, _ := conn.PanicOnError().Begin()
defer tx.Close()

tx.Query(Users.Insert().Values(User{Name: "Zero"}))
tx.IsSuccessful()
```

### Statements

Sol includes a variety of SQL statements that can be constructed directly from declared table schemas.

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

It will output dialect neutral SQL from its `String()` method and a dialect specific version from conn.String().

```sql
CREATE TABLE "users" (
  "id" INTEGER NOT NULL,
  "name" VARCHAR(32) NOT NULL,
  "password" VARCHAR(128) NOT NULL,
  PRIMARY KEY ("id")
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
INSERT INTO "users" ("id", "name", "password") VALUES (?, ?, ?)
```

Values can be inserted to the database using custom struct types or the generic `sol.Values` type. If given a struct, Sol will attempt to match SQL column names to struct field names in a case sensitive manner that is aware of camel to snake case conversion. More complex names or aliases should specify db struct tags.

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

Both will produce the following SQL with the Sqlite3 dialect:

```sql
UPDATE "users" SET "id" = ?, "name" = ?, "password" = ?
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

#### SELECT

Results can be queried in a number of ways. Each of the following statements will produce the same SQL output:

```go
Users.Select()
sol.Select(Users)
sol.Select(Users.C("id"), Users.C("name"), Users.C("password"))
```

```sql
SELECT "users"."id", "users"."name", "users"."password" FROM "users"
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

Tables can be constructed with foreign keys, unique constraints, and composite primary keys. See the `sol_test.go` file for more examples.

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
CREATE TABLE "contacts" (
  "id" INTEGER,
  "user_id" INTEGER NOT NULL REFERENCES users("id"),
  "key" VARCHAR,
  "value" VARCHAR,
  PRIMARY KEY ("id"),
  UNIQUE ("user_id", "key")
);
```

Develop
-------

To run all tests, the `postgres` sub-package requires a configuration file at `db.json` within the `postgres/` directory. See the example file in that directory for the correct format.

Statements and clauses can be tested by creating a new dialect-specific tester; for example using the `postgres` package:

```go
expect := NewTester(t, &postgres.PostGres{})
```

The instance's `SQL` method will test expected output and parameterization:

```go
expect.SQL(`DELETE FROM "users"`, users.Delete())

expect.SQL(
    `INSERT INTO "users" ("name") VALUES ($1), ($2)`,
    users.Insert().Values([]sol.Values{{"name": "Totti"}, {"name": "De Rossi"}}),
    "Totti", "De Rossi",
)
```

And the `Error` method will test that an error occurred:

```go
expect.Error(sql.Select(users.C("does-not-exist")))
```

Happy Hacking!

aodin, 2015-2016
