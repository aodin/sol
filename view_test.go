package sol

import "testing"

func TestView(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	// Create a new view from a table selection
	view, err := View("user_emails", users.Select().OrderBy(users.C("email")))
	if err != nil {
		t.Fatalf("Unexpected error creating View: %s", err)
	}

	if view.name != "user_emails" {
		t.Errorf("Unexpected table name: %s != user_emails", view.name)
	}

	expect.SQL(
		view.Create(),
		`CREATE VIEW user_emails AS (SELECT users.id, users.email, users.name, users.password, users.created_at FROM users ORDER BY users.email)`,
	)

	expect.SQL(
		view.Select(),
		`SELECT users.id, users.email, users.name, users.password, users.created_at FROM user_emails`,
	)
}
