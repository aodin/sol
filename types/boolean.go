package types

import (
	"fmt"
)

type BooleanType struct {
	BaseType
}

var _ Type = Boolean()

// Default will set the default value of the boolean
func (t BooleanType) Default(value bool) BooleanType {
	t.BaseType.SetDefault(fmt.Sprintf("%t", value))
	return t
}

func (t BooleanType) NotNull() BooleanType {
	t.BaseType.SetNotNull()
	return t
}

func (t BooleanType) Unique() BooleanType {
	t.BaseType.SetUnique()
	return t
}

// Boolean creats a new BOOLEAN datatype
func Boolean() (t BooleanType) {
	return BooleanType{BaseType: Base("BOOLEAN")}
}
