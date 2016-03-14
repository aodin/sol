Sol [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/aodin/sol) [![Build Status](https://travis-ci.org/aodin/sol.svg?branch=master)](https://travis-ci.org/aodin/sol)
======

[![Join the chat at https://gitter.im/aodin/sol](https://badges.gitter.im/aodin/sol.svg)](https://gitter.im/aodin/sol?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

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

Happy Hacking!

aodin, 2015-2016
