package sol

import "testing"

func TestForeignKey(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		contacts.Create(),
		`CREATE TABLE contacts (
  id INTEGER,
  user_id INTEGER REFERENCES users(id),
  key VARCHAR,
  value VARCHAR,
  PRIMARY KEY (id),
  UNIQUE (user_id, key)
);`,
	)

	expect.SQL(
		messages.Create(),
		`CREATE TABLE messages (
  id INTEGER,
  user_id INTEGER REFERENCES users(id),
  parent_id INTEGER REFERENCES messages(id),
  text TEXT
);`,
	)

	if len(messages.ForeignKeys()) != 2 {
		t.Fatalf(
			"unexpected length of messages foreign keys: %d != 2",
			len(messages.ForeignKeys()),
		)
	}

	// The messages table's first reference should be users
	if messages.ForeignKeys()[0].references != users {
		t.Errorf("messages' first foreign key should reference table users")
	}

	// The users table should be referenced by messages and contacts
	if len(users.ReferencedBy()) != 2 {
		t.Fatalf(
			"unexpected length of users references: %d != 2",
			len(users.ReferencedBy()),
		)
	}
}
