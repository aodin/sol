package sol

import "testing"

func TestText(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	expect.SQL(
		`SELECT * FROM users WHERE id > $1 OR name LIKE $2`,
		Text(
			`SELECT * FROM users WHERE id > :id OR name LIKE :name`,
			Values{"name": "A", "id": 2},
		),
		2, "A",
	)
}
