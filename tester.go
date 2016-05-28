package sol

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/aodin/sol/dialect"
)

// callerInfo returns a string containing the file and line number of the
// assert call that failed.
// https://github.com/stretchr/testify/blob/master/assert/assertions.go
// Copyright (c) 2012 - 2013 Mat Ryer and Tyler Bunnell
func callerInfo() string {
	file := ""
	line := 0
	ok := false

	for i := 0; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return ""
		}
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]

		// dir := parts[len(parts)-2]
		if file == "tester.go" {
			continue
		}
		break
	}
	return fmt.Sprintf("%s:%d", file, line)
}

type tester struct {
	t       *testing.T
	dialect dialect.Dialect
}

// Error tests that the given Compiles instances generates an error for the
// current dialect.
func (t *tester) Error(stmt Compiles) {
	// TODO Allow a specific error
	if _, err := stmt.Compile(t.dialect, Params()); err == nil {
		t.t.Errorf("%s: expected error, received nil", callerInfo())
	}
}

// SQL tests that the given Compiles instance matches the expected string for
// the current dialect.
func (t *tester) SQL(stmt Compiles, expect string, ps ...interface{}) {
	// Get caller information in case of failure
	caller := callerInfo()

	// Start a new parameters instance
	params := Params()

	// Compile the given stmt with the tester's dialect
	actual, err := stmt.Compile(t.dialect, params)
	if err != nil {
		t.t.Errorf("%s: unexpected error from Compile(): %s", caller, err)
		return
	}

	if expect != actual {
		t.t.Errorf(
			"%s: unexpected SQL from Compile(): \n - have: %s\n - want: %s",
			caller,
			actual,
			expect,
		)
	}
	// Test that the parameters are equal
	if len(*params) != len(ps) {
		t.t.Errorf(
			"%s: unexpected number of parameters for %s: \n - have %d, want %d",
			caller,
			actual,
			len(ps),
			len(*params),
		)
		return
	}

	// Examine individual parameters for equality
	for i, param := range *params {
		if !reflect.DeepEqual(ps[i], param) {
			t.t.Errorf(
				"%s: unequal parameters at index %d: \n - have %#v, want %#v",
				caller,
				i,
				param,
				ps[i],
			)
		}
	}
}

// NewTester creates a new SQL/Error tester that uses the given dialect
func NewTester(t *testing.T, d dialect.Dialect) *tester {
	return &tester{t: t, dialect: d}
}
