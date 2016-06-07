package sol

import "testing"

// All schemas are declared in sol_test.go

func TestSelect(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	// Select statements without destination structs
	expect.SQL(
		Select(users),
		`SELECT users.id, users.email, users.name, users.password, users.created_at FROM users`,
	)

	expect.SQL(
		Select(users.C("email")),
		`SELECT users.email FROM users`,
	)

	expect.SQL(
		Select(users.C("email"), contacts.C("value")),
		`SELECT users.email, contacts.value FROM users, contacts`,
	)

	// SelectTable
	expect.SQL(
		SelectTable(users, contacts.C("value")),
		`SELECT users.id, users.email, users.name, users.password, users.created_at, contacts.value FROM users`,
	)

	// Add an alias
	expect.SQL(
		Select(users.C("email").As("Email")),
		`SELECT users.email AS "Email" FROM users`,
	)

	// Add an ORDER BY
	expect.SQL(
		Select(users.C("email")).OrderBy(users.C("email").Desc()),
		`SELECT users.email FROM users ORDER BY users.email DESC`,
	)

	// Mutiple conditionals will be merged with AND by default
	expect.SQL(
		Select(
			users.C("name"),
		).Where(
			users.C("id").DoesNotEqual(1),
			users.C("name").Equals("admin"),
		),
		`SELECT users.name FROM users WHERE (users.id <> $1 AND users.name = $2)`,
		1, "admin",
	)

	// Distinct
	expect.SQL(
		Select(users.C("name")).Distinct(),
		`SELECT DISTINCT users.name FROM users`,
	)

	expect.SQL(
		Select(users.C("name")).Distinct(users.C("id"), users.C("name")),
		`SELECT DISTINCT ON (users.id, users.name) users.name FROM users`,
	)

	// All is the default and will remove any existing Distinct clause
	expect.SQL(
		Select(users.C("name")).Distinct().All(),
		`SELECT users.name FROM users`,
	)

	// Build a GROUP BY statement using an aggregate
	expect.SQL(
		Select(
			contacts.C("user_id"),
			Count(contacts.C("id")),
		).GroupBy(
			contacts.C("user_id"),
		).Having(
			Count(contacts.C("id")).GTE(2),
		).OrderBy(
			Count(contacts.C("id")).Desc(),
		),
		`SELECT contacts.user_id, COUNT(contacts.id) FROM contacts GROUP BY contacts.user_id HAVING COUNT(contacts.id) >= $1 ORDER BY COUNT(contacts.id) DESC`,
		2,
	)

	// Test limit
	expect.SQL(
		Select(users.C("name")).Limit(1),
		`SELECT users.name FROM users LIMIT 1`,
	)

	// Test Offset
	expect.SQL(
		Select(users.C("name")).Offset(1),
		`SELECT users.name FROM users OFFSET 1`,
	)

	// Select zero columns
	expect.Error(Select())

	// Select a column that doesn't exist
	expect.Error(Select(users.C("what")))
}
