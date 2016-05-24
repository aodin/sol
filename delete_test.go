package sol

import "testing"

func TestDelete(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Test a complete delete
	// TODO Require an All when no clauses are given to prevent mass deletes?
	expect.SQL(
		users.Delete(),
		`DELETE FROM users`,
	)

	// Test a delete with a WHERE
	expect.SQL(
		users.Delete().Where(users.C("id").Equals(1)),
		`DELETE FROM users WHERE users.id = $1`,
		1,
	)

	expect.SQL(
		Delete(users).Where(users.C("id").Equals(1), users.C("name").Equals("admin")),
		`DELETE FROM users WHERE (users.id = $1 AND users.name = $2)`,
		1, "admin",
	)
}
