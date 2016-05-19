package sol

import "testing"

func TestDelete(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Test a complete delete
	// TODO Require an All when no clauses are given to prevent mass deletes?
	expect.SQL(
		`DELETE FROM "users"`,
		users.Delete(),
	)

	// Test a delete with a WHERE
	expect.SQL(
		`DELETE FROM "users" WHERE "users"."id" = $1`,
		users.Delete().Where(users.C("id").Equals(1)),
		1,
	)

	expect.SQL(
		`DELETE FROM "users" WHERE ("users"."id" = $1 AND "users"."name" = $2)`,
		Delete(users).Where(users.C("id").Equals(1), users.C("name").Equals("admin")),
		1, "admin",
	)
}
