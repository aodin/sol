package types

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

type character struct {
	base
	limit int
}

func (t character) Create(d dialect.Dialect) (string, error) {
	compiled := t.base.name
	if t.limit != 0 {
		compiled += fmt.Sprintf("(%d)", t.limit)
	}
	compiled += t.base.suffix()
	return compiled, nil
}

func (t character) Limit(n int) character {
	t.limit = n
	return t
}

func (t character) NotNull() character {
	t.base.NotNull()
	return t
}

func (t character) Unique() character {
	t.base.Unique()
	return t
}

func Char(n int) character {
	return Character(n)
}

func Character(n int) (t character) {
	t.name = "CHAR"
	t.limit = n
	return
}

// Varchar creates a new VARCHAR. Limit will be set if an argument is given -
// all subsequent arguments will be ignored
func Varchar(limit ...int) (t character) {
	t.name = "VARCHAR"
	if limit != nil {
		t.limit = limit[0]
	}
	return
}

func Text() (t character) {
	t.name = "TEXT"
	return
}
