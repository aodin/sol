package postgres

import (
	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

// serial is a postgres-specific auto-increment type. It implies NOT NULL.
type serial struct {
	name     string
	isUnique bool
}

// serial must implement the Type interface
var _ types.Type = serial{}

func (t serial) Create(d dialect.Dialect) (string, error) {
	compiled := t.name + " NOT NULL"
	if t.isUnique {
		compiled += " UNIQUE"
	}
	return compiled, nil
}

// Unique sets the serial type as unique
func (t serial) Unique() serial {
	t.isUnique = true
	return t
}

// Serial creates a new serial type
func Serial() (t serial) {
	t.name = "serial"
	return
}
