package postgres

import (
	"github.com/aodin/sol/dialect"
)

type json struct {
	name      string
	isNotNull bool
	isUnique  bool
}

func (t json) Create(d dialect.Dialect) (string, error) {
	compiled := t.name
	if t.isNotNull {
		compiled += " NOT NULL"
	}
	if t.isUnique {
		compiled += " UNIQUE"
	}
	return compiled, nil
}

func (t json) NotNull() json {
	t.isNotNull = true
	return t
}

func (t json) Unique() json {
	t.isUnique = true
	return t
}

func JSON() (t json) {
	t.name = "json"
	return
}
