package sol

import "testing"

func TestFunctions(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`SELECT count("users"."id") FROM "users"`,
		Select(Count(users.C("id"))),
	)

	expect.SQL(
		`SELECT count("users"."id") AS "Count" FROM "users"`,
		Select(Count(users.C("id")).As("Count")),
	)
}
