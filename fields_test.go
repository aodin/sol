package sol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCamelToSnake(t *testing.T) {
	assert.Equal(t, "snake_case", camelToSnake("SnakeCase"))
	assert.Equal(t, "user_id", camelToSnake("UserID"))
	assert.Equal(t, "uuid", camelToSnake("UUID"))

	// TODO Unicode test cases?
}

type embeddedID struct {
	ID uint64 `db:",omitempty"`
}

func TestSelectFields_embeddedID(t *testing.T) {
	var elem embeddedID
	fields := SelectFields(&elem)

	if len(fields) != 1 {
		t.Fatalf("Unexpected length of fields: %d", len(fields))
	}

	assert.Equal(t,
		field{
			column:  "ID",
			names:   []string{"ID"},
			options: []string{OmitEmpty},
		},
		fields[0],
	)
}

type ignored struct {
	manager struct{} `db:"-"`
	name    string
}

func TestSelectFields_ignored(t *testing.T) {
	var elem ignored
	fields := SelectFields(&elem)

	if len(fields) != 1 {
		t.Fatalf("Unexpected length of fields: %d", len(fields))
	}

	assert.Equal(t,
		field{
			column:  "name",
			names:   []string{"name"},
			options: []string{},
		},
		fields[0],
	)
}

// TODO implement scanner

type embedded struct {
	embeddedID
	Name      string
	Timestamp struct {
		CreatedAt time.Time `db:",omitempty"`
		UpdatedAt *time.Time
		isActive  bool
	}
	manager *struct{}
}

type nested struct {
	Another string
	embedded
}

type moreNesting struct {
	nested
	OneMore string
}
