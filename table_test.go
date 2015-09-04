package sol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// All schemas are declared in sol_test.go

func TestTable(t *testing.T) {
	assert.Equal(t, "users", users.name)

	// Confirm the fields created by Modifiers
	assert.Equal(t, PKArray{"id"}, users.pk)
	assert.Equal(t, []UniqueArray{{"email"}}, users.uniques)
}

func TestTable_Select(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	// Select statements without destination structs
	expect.SQL(
		`SELECT "users"."id", "users"."email", "users"."name", "users"."password", "users"."created_at" FROM "users"`,
		users.Select(),
	)
}
