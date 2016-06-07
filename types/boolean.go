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
	BaseType
	value *bool
}

var _ Type = boolean{}

func (t boolean) Create(d dialect.Dialect) (string, error) {
	compiled, err := t.BaseType.Create(d)
	if err != nil {
		return "", err
	}
	if t.value != nil {
		compiled += strings.ToUpper(fmt.Sprintf(" DEFAULT %t", *t.value))
	}
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
	t.BaseType.NotNull()
	return t
}

func (t boolean) Unique() boolean {
	t.BaseType.Unique()
	return t
}

// Boolean creats a new BOOLEAN datatype
func Boolean() (t boolean) {
	return boolean{BaseType: Base("BOOLEAN")}
}
