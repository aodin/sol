package sol

import (
	"testing"
)

// Valid schemas are declared in sol_test

func TestCreate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`CREATE TABLE "users" (
  "id" INTEGER,
  "email" VARCHAR(256) NOT NULL UNIQUE,
  "name" VARCHAR(32) NOT NULL,
  "password" VARCHAR,
  "created_at" TIMESTAMP,
  PRIMARY KEY ("id")
);`,
		users.Create(),
	)
}
