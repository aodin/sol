package sol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO Unicode test cases?
var caseTests = []struct {
	In, Out string
}{
	{In: "SnakeCase", Out: "snake_case"},
	{In: "UserID", Out: "user_id"},
	{In: "UUID", Out: "uuid"},
}

func TestCamelToSnake(t *testing.T) {
	for i, test := range caseTests {
		out := camelToSnake(test.In)
		if out != test.Out {
			t.Errorf(
				"Unexpected camel to snake case conversion %d - %s: %s != %s",
				i, test.In, out, test.Out,
			)
		}
	}
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
	OneMore string `db:"text,omitempty"`
}
