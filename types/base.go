package types

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// BaseType is foundational datatype that includes fields that nearly all
// datatypes implement
type BaseType struct {
	name         string
	defaultValue string
	isNotNull    bool
	isUnique     bool
}

var _ Type = BaseType{}

// Create returns a string suitable for use in a CREATE TABLE statement
func (base BaseType) Create(d dialect.Dialect) (string, error) {
	clauses := append([]string{base.name}, base.Options()...)
	return strings.Join(clauses, " "), nil
}

func (base BaseType) Name() string {
	return base.name
}

func (base BaseType) Default(value string) BaseType {
	base.SetDefault(value)
	return base
}

func (base BaseType) NotNull() BaseType {
	base.SetNotNull()
	return base
}

func (base BaseType) Unique() BaseType {
	base.SetUnique()
	return base
}

// SetDefault sets the type's default value
func (base *BaseType) SetDefault(value string) {
	base.defaultValue = value
}

// SetName sets the type's name
func (base *BaseType) SetName(name string) {
	base.name = name
}

// SetNotNull sets the BaseType to NOT NULL
func (base *BaseType) SetNotNull() {
	base.isNotNull = true
}

// SetUnique sets the BaseType to UNIQUE
func (base *BaseType) SetUnique() {
	base.isUnique = true
}

// Options returns the BaseTYPE options as a slice of strings
func (base BaseType) Options() (compiled []string) {
	if base.isNotNull {
		compiled = append(compiled, "NOT NULL")
	}
	if base.isUnique {
		compiled = append(compiled, "UNIQUE")
	}
	if base.defaultValue != "" {
		compiled = append(compiled,
			fmt.Sprintf("DEFAULT %s", base.defaultValue),
		)
	}
	return
}

// Base creates a new BaseType
func Base(name string) BaseType {
	return BaseType{name: name}
}

// New is an alias for Base
func New(name string) BaseType {
	return Base(name)
}
