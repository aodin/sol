package types

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

type CharacterType struct {
	BaseType
	limit int
}

var _ Type = Varchar()

func (t CharacterType) Create(d dialect.Dialect) (string, error) {
	name := t.BaseType.name
	if t.limit != 0 {
		name += fmt.Sprintf("(%d)", t.limit)
	}
	compiled := append([]string{name}, t.BaseType.Options()...)
	return strings.Join(compiled, " "), nil
}

func (t CharacterType) Limit(n int) CharacterType {
	t.limit = n
	return t
}

func (t CharacterType) NotNull() CharacterType {
	t.BaseType.SetNotNull()
	return t
}

func (t CharacterType) Unique() CharacterType {
	t.BaseType.SetUnique()
	return t
}

func Char(n int) CharacterType {
	return Character(n)
}

func Character(n int) (t CharacterType) {
	t.name = "CHAR"
	t.limit = n
	return
}

// Varchar creates a new VARCHAR. Limit will be set if an argument is given -
// all subsequent arguments will be ignored
func Varchar(limit ...int) (t CharacterType) {
	t.name = "VARCHAR"
	if limit != nil {
		t.limit = limit[0]
	}
	return
}

func Text() (t CharacterType) {
	t.name = "TEXT"
	return
}
