package sol

import "testing"

func TestFunctions(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		Select(Count(users.C("id"))),
		`SELECT COUNT(users.id) FROM users`,
	)

	expect.SQL(
		Select(Count(users.C("id")).As("Count")),
		`SELECT COUNT(users.id) AS "Count" FROM users`,
	)
}
