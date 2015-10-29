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

	// Mutiple conditionals will be merged with AND by default
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

	// Distinct
	expect.SQL(
		`SELECT DISTINCT "users"."name" FROM "users"`,
		Select(users.C("name")).Distinct(),
	)

	expect.SQL(
		`SELECT DISTINCT ON ("users"."id", "users"."name") "users"."name" FROM "users"`,
		Select(users.C("name")).Distinct(users.C("id"), users.C("name")),
	)

	// All is the default and will remove any existing Distinct clause
	expect.SQL(
		`SELECT "users"."name" FROM "users"`,
		Select(users.C("name")).Distinct().All(),
	)

	// Build a GROUP BY statement using an aggregate
	expect.SQL(
		`SELECT "contacts"."user_id", count("contacts"."id") FROM "contacts" GROUP BY "contacts"."user_id" ORDER BY count("contacts"."id") DESC`,
		Select(
			contacts.C("user_id"),
			Count(contacts.C("id")),
		).GroupBy(
			contacts.C("user_id"),
		).OrderBy(
			Count(contacts.C("id")).Desc(),
		),
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
