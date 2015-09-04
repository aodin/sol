package types

import (
	"github.com/aodin/sol/dialect"
)

type base struct {
	name      string
	isNotNull bool
	isUnique  bool
}

func (t base) Create(d dialect.Dialect) (string, error) {
	return t.name + t.suffix(), nil
}

func (t *base) Unique() {
	t.isUnique = true
}

func (t *base) NotNull() {
	t.isNotNull = true
}

func (t base) suffix() (compiled string) {
	if t.isNotNull {
		compiled += " NOT NULL"
	}
	if t.isUnique {
		compiled += " UNIQUE"
	}
	return
}
