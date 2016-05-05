package sol

import "testing"

func TestForeignKey(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`CREATE TABLE "contacts" (
  "id" INTEGER,
  "user_id" INTEGER REFERENCES users("id"),
  "key" VARCHAR,
  "value" VARCHAR,
  PRIMARY KEY ("id"),
  UNIQUE ("user_id", "key")
);`,
		contacts.Create(),
	)

	expect.SQL(
		`CREATE TABLE "messages" (
  "id" INTEGER,
  "user_id" INTEGER REFERENCES users("id"),
  "text" TEXT
);`,
		messages.Create(),
	)

	// The messages table should reference users
	if len(messages.ForeignKeys()) != 1 {
		t.Fatalf(
			"unexpected length of messages foreign keys: %d != 1",
			len(messages.ForeignKeys()),
		)
	}

	if messages.ForeignKeys()[0].references != users {
		t.Errorf("messages foreign key should reference table users")
	}

	// The users table should be referenced by messages and contacts
	if len(users.ReferencedBy()) != 2 {
		t.Fatalf(
			"unexpected length of users references: %d != 2",
			len(users.ReferencedBy()),
		)
	}
}
