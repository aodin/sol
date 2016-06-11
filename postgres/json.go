package postgres

import "github.com/aodin/sol/types"

// json is a JSON column type
type json struct {
	types.BaseType
}

// json must implement the Type interface
var _ types.Type = json{}

func (t json) NotNull() json {
	t.BaseType.SetNotNull()
	return t
}

func (t json) Unique() json {
	t.BaseType.SetUnique()
	return t
}

func JSON() json {
	return json{BaseType: types.Base("JSON")}
}
