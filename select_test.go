package sol

import (
	"testing"
)

// All schemas are declared in sol_test.go

func TestSelect(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	// Select statements without destination structs
	expect.SQL(
		`SELECT "users"."id", "users"."email", "users"."name", "users"."password", "users"."created_at" FROM "users"`,
		Select(users),
	)

	expect.SQL(
		`SELECT "users"."email" FROM "users"`,
		Select(users.C("email")),
	)

	expect.SQL(
		`SELECT "users"."email", "contacts"."value" FROM "users", "contacts"`,
		Select(users.C("email"), contacts.C("value")),
	)

}
