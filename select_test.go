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

	// Add an alias
	expect.SQL(
		`SELECT "users"."email" AS "Email" FROM "users"`,
		Select(users.C("email").As("Email")),
	)

	// Add an ORDER BY
	expect.SQL(
		`SELECT "users"."email" FROM "users" ORDER BY "users"."email" DESC`,
		Select(users.C("email")).OrderBy(users.C("email").Desc()),
	)

	// Mutiple conditionals will joined with AND by default
	expect.SQL(
		`SELECT "users"."name" FROM "users" WHERE ("users"."id" <> $1 AND "users"."name" = $2)`,
		Select(
			users.C("name"),
		).Where(
			users.C("id").DoesNotEqual(1),
			users.C("name").Equals("admin"),
		),
		1, "admin",
	)

	// Test limit
	expect.SQL(
		`SELECT "users"."name" FROM "users" LIMIT 1`,
		Select(users.C("name")).Limit(1),
	)

	// Test Offset
	expect.SQL(
		`SELECT "users"."name" FROM "users" OFFSET 1`,
		Select(users.C("name")).Offset(1),
	)

	// Select a column that doesn't exist
	expect.Error(Select(users.C("what")))
}
