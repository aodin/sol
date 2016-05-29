package sol

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNoColumns is returned when attempting to compile a query without
// any columns
var ErrNoColumns = errors.New(
	"sol: cannot compile a statement without columns",
)

type fieldError struct {
	column string
	table  string
	clause string
}

type stmtErrors struct {
	meta   []string
	fields map[fieldError]string
}

// Error implements the error interface
func (e stmtErrors) Error() string {
	errs := e.meta
	for field, err := range e.fields {
		errs = append(errs, fmt.Sprintf("%s (%s)", err, field))
	}
	return strings.Join(errs, "; ")
}

// Exist returns true if there are either meta or fields errors
func (e stmtErrors) Exist() bool {
	return len(e.meta) > 0 || len(e.fields) > 0
}
