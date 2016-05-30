package sol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Serial struct {
	ID uint64 `db:"id,omitempty"`
}

func TestSelectFields_Serial(t *testing.T) {
	var elem Serial
	fields := SelectFields(&elem)

	if len(fields) != 1 {
		t.Fatalf("Unexpected length of fields: %d", len(fields))
	}

	assert.Equal(t,
		field{
			column:  "id",
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
	Serial
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
