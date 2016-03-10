package sol

import (
	"testing"
)

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
}
