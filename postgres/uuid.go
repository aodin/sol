package postgres

import (
	"fmt"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

// TODO Don't hardcode utc?
const GenerateV4 = `uuid_generate_v4()`

type uuid struct {
	name         string
	isNotNull    bool
	isUnique     bool
	defaultValue string // TODO Additional defaults?
}

// uuid must implement the Type interface
var _ types.Type = uuid{}

func (t uuid) Create(d dialect.Dialect) (string, error) {
	compiled := t.name
	if t.isNotNull {
		compiled += " NOT NULL"
	}
	if t.isUnique {
		compiled += " UNIQUE"
	}
	if t.defaultValue != "" {
		compiled += fmt.Sprintf(" DEFAULT (%s)", t.defaultValue)
	}
	return compiled, nil
}

func (t uuid) Default(value string) uuid {
	t.defaultValue = value
	return t
}

func (t uuid) NotNull() uuid {
	t.isNotNull = true
	return t
}

func (t uuid) Unique() uuid {
	t.isUnique = true
	return t
}

func UUID() (t uuid) {
	t.name = "uuid"
	return
}
