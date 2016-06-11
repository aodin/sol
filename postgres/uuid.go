package postgres

import (
	"github.com/aodin/sol/types"
)

const GenerateV4 = `uuid_generate_v4()`

type uuid struct {
	types.BaseType
}

// uuid must implement the Type interface
var _ types.Type = uuid{}

func (t uuid) Default(value string) uuid {
	t.BaseType.SetDefault(value)
	return t
}

func (t uuid) NotNull() uuid {
	t.BaseType.SetNotNull()
	return t
}

func (t uuid) Unique() uuid {
	t.BaseType.SetUnique()
	return t
}

func UUID() (t uuid) {
	return uuid{BaseType: types.Base("UUID")}
}
