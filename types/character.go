package types

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

type character struct {
	BaseType
	limit int
}

func (t character) Create(d dialect.Dialect) (string, error) {
	name := t.BaseType.name
	if t.limit != 0 {
		name += fmt.Sprintf("(%d)", t.limit)
	}
	compiled := append([]string{name}, t.BaseType.Options()...)
	return strings.Join(compiled, " "), nil
}

func (t character) Limit(n int) character {
	t.limit = n
	return t
}

func (t character) NotNull() character {
	t.BaseType.NotNull()
	return t
}

func (t character) Unique() character {
	t.BaseType.Unique()
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
