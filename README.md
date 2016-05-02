Sol
======

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/sol) [![Build Status](https://travis-ci.org/aodin/sol.svg?branch=master)](https://travis-ci.org/aodin/sol) [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/aodin/sol)

A relational database toolkit for Go:

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
go get github.com/aodin/sol
```


Develop
-------

To perform test, the `postgres` sub-package requires a configuration at `db.json` within the `postgres/` directory. See the example file in that directory for the correct format.

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
