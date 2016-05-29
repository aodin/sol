package sol

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
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

// IntegrationTest runs a large, neutral dialect test
func IntegrationTest(t *testing.T, conn *DB) {
	// Perform all tests in a transaction
	// TODO What features should be testing outside of a transaction?
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Creating a new transaction should not error: %s", err)
	}
	defer tx.Rollback()

	// CREATE TABLE
	// TODO foreign keys
	testusers := Table("testusers",
		Column("id", types.Integer()),
		Column("email", types.Varchar().Limit(256).NotNull()),
		Column("is_admin", types.Boolean().NotNull()),
		Column("created_at", types.Timestamp()),
		PrimaryKey("id"),
		Unique("email"),
	)

	type testuser struct {
		ID        int64
		Email     string
		IsAdmin   bool
		CreatedAt time.Time
	}

	if err = tx.Query(testusers.Create()); err != nil {
		t.Fatalf("CREATE TABLE should not error: %s", err)
	}

	// INSERT by struct
	// Truncate the time.Time field to avoid significant digit errors
	admin := testuser{
		ID:        1,
		Email:     "admin@example.com",
		IsAdmin:   true,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	if err = tx.Query(testusers.Insert().Values(admin)); err != nil {
		t.Fatalf("INSERT by struct should not fail %s", err)
	}

	// SELECT
	var selected testuser
	if err = tx.Query(
		testusers.Select().Where(testusers.C("id").Equals(admin.ID)),
		&selected,
	); err != nil {
		t.Fatalf("SELECT should not fail: %s", err)
	}

	// TODO test with direct comparison: selected == admin
	// For now, test each field since DATETIME handling is terribly
	// inconsistent across databases
	if selected.ID != admin.ID {
		t.Errorf(
			"Unequal testusers id: have %d, want %d",
			selected.ID, admin.ID,
		)
	}
	if selected.Email != admin.Email {
		t.Errorf(
			"Unequal testusers email: have %s, want %s",
			selected.Email, admin.Email,
		)
	}
	if selected.IsAdmin != admin.IsAdmin {
		t.Errorf(
			"Unequal testusers is_admin: have %t, want %t",
			selected.IsAdmin, admin.IsAdmin,
		)
	}
	if !selected.CreatedAt.Equal(admin.CreatedAt) {
		t.Errorf(
			"Unequal testusers created_at: have %v, want %v",
			selected.CreatedAt, admin.CreatedAt,
		)
	}

	// UPDATE
	if err = tx.Query(
		testusers.Update().Values(
			Values{"is_admin": false},
		).Where(testusers.C("id").Equals(admin.ID)),
	); err != nil {
		t.Fatalf("UPDATE should not fail: %s", err)
	}

	var updated testuser
	if err = tx.Query(testusers.Select().Limit(1), &updated); err != nil {
		t.Fatalf("SELECT should not fail: %s", err)
	}

	selected.IsAdmin = false
	if updated != selected {
		t.Errorf(
			"Unequal testusers: have %+v, want %+v",
			updated, selected,
		)
	}

	// INSERT by values
	client := Values{
		"id":         2,
		"email":      "client@example.com",
		"is_admin":   false,
		"created_at": time.Now().UTC().Truncate(time.Second),
	}

	if err = tx.Query(testusers.Insert().Values(client)); err != nil {
		t.Fatalf("INSERT by values should not fail %s", err)
	}
	var list []testuser
	if err = tx.Query(
		testusers.Select().OrderBy(testusers.C("id").Desc()),
		&list,
	); err != nil {
		t.Fatalf("SELECT with ORDER BY should not fail: %s", err)
	}

	if len(list) != 2 {
		t.Fatalf("Unexpected length of list: want 2, have %d", len(list))
	}

	// The client should be first
	if list[0].Email != "client@example.com" {
		t.Errorf(
			"Unexpected email: want client@example.com, have %d",
			list[0].Email,
		)
	}

	// DELETE
	if err = tx.Query(
		testusers.Delete().Where(testusers.C("email").Equals(admin.Email)),
	); err != nil {
		t.Fatalf("DELETE should not fail: %s", err)
	}

	// DROP TABLE
	if err = tx.Query(testusers.Drop()); err != nil {
		t.Fatalf("DROP TABLE should not fail %s", err)
	}

	// Test a recover
	func() {
		defer func() {
			if panicked := recover(); panicked == nil {
				t.Errorf("Connection failed to panic on error")
			}
		}()
		conn.Must().Query(testusers.Select(), list)
	}()
}
