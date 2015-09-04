package types

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// Provide a nullable boolean for Boolean type default values
var (
	internalTrue  bool = true
	internalFalse bool = false
	True               = &internalTrue
	False              = &internalFalse
)

type boolean struct {
	base
	value *bool
}

func (t boolean) Create(d dialect.Dialect) (string, error) {
	compiled := t.base.name
	if t.value != nil {
		compiled += strings.ToUpper(fmt.Sprintf(" DEFAULT %t", *t.value))
	}
	compiled += t.base.suffix()
	return compiled, nil
}

// Default will set the default value of the boolean
func (t boolean) Default(value bool) boolean {
	if value {
		t.value = True
	} else {
		t.value = False
	}
	return t
}

func (t boolean) NotNull() boolean {
	t.base.NotNull()
	return t
}

func (t boolean) Unique() boolean {
	t.base.Unique()
	return t
}

func Boolean() (t boolean) {
	t.name = "BOOLEAN"
	return
}
