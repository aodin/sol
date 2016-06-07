package types

import (
	"strings"

	"github.com/aodin/sol/dialect"
)

// BaseType is foundational datatype that includes fields that nearly all
// datatypes implement
type BaseType struct {
	name      string
	isNotNull bool
	isUnique  bool
}

var _ Type = BaseType{}

// Create generates the
func (base BaseType) Create(d dialect.Dialect) (string, error) {
	clauses := append([]string{base.name}, base.Options()...)
	return strings.Join(clauses, " "), nil
}

// Unique sets the BaseType to UNIQUE
func (base *BaseType) Unique() {
	base.isUnique = true
}

// NotNull sets the BaseType to NOT NULL
func (base *BaseType) NotNull() {
	base.isNotNull = true
}

// Options returns the BaseTYPE options as a slice of strings
func (base BaseType) Options() (compiled []string) {
	if base.isNotNull {
		compiled = append(compiled, "NOT NULL")
	}
	if base.isUnique {
		compiled = append(compiled, "UNIQUE")
	}
	return
}

// Base creates a new BaseType
func Base(name string) BaseType {
	return BaseType{name: name}
}
