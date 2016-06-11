package postgres

import (
	"github.com/aodin/sol/types"
)

// serial is a postgres-specific auto-increment type. It implies NOT NULL.
type serial struct {
	types.BaseType
}

// serial must implement the Type interface
var _ types.Type = serial{}

func (t serial) Unique() serial {
	t.BaseType.SetUnique()
	return t
}

// Serial creates a new serial type. Serial implies NOT NULL
func Serial() (t serial) {
	base := types.Base("SERIAL")
	base.SetNotNull()
	return serial{BaseType: base}
}
