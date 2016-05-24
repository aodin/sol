package sol

import "testing"

// All schemas are declared in sol_test.go

func TestTable(t *testing.T) {
	if users.name != "users" {
		t.Errorf("Unexpected table name: %s != users", users.name)
	}

	// Confirm the fields created by Modifiers
	if len(users.pk) != 1 {
		t.Fatalf("Unexpected length of primary keys: %d != 1", len(users.pk))
	}
	if users.pk[0] != "id" {
		t.Errorf("Unexpected primary key: %s != id", users.pk[0])
	}
	if len(users.uniques) != 1 {
		t.Fatalf("Unexpected length of uniques: %d != 1", len(users.uniques))
	}
	if len(users.uniques[0]) != 1 {
		t.Fatalf(
			"Unexpected length of unique array: %d != 1",
			len(users.uniques),
		)
	}
	if users.uniques[0][0] != "email" {
		t.Errorf("Unexpected unique: %s != email", users.uniques[0][0])
	}
}

func TestTable_Select(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	// Select statements without destination structs
	expect.SQL(
		users.Select(),
		`SELECT users.id, users.email, users.name, users.password, users.created_at FROM users`,
	)
}
