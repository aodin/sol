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

	expect.SQL(
		Select(Date(users.C("created_at"))),
		`SELECT DATE(users.created_at) FROM users`,
	)

	expect.SQL(
		Select(DatePart("hour", users.C("created_at")).As("Hour")),
		`SELECT DATE_PART('hour', users.created_at) AS "Hour" FROM users`,
	)
}
