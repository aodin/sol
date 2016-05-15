package sol

import (
	"reflect"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// By default, an INSERT without values will assume a single entry
	expect.SQL(
		`INSERT INTO "contacts" ("id", "user_id", "key", "value") VALUES ($1, $2, $3, $4)`,
		contacts.Insert(),
		nil, nil, nil, nil,
	)

	expect.SQL(
		`INSERT INTO "users" ("name", "password") VALUES ($1, $2)`,
		Insert(users.C("name"), users.C("password")),
		nil, nil,
	)

	expect.SQL(
		`INSERT INTO "users" ("email", "name") VALUES ($1, $2)`,
		users.Insert().Values(user{Name: "admin", Email: "admin@example.com"}),
		"admin@example.com", "admin",
	)

	// Use sql.Values
	expect.SQL(
		`INSERT INTO "users" ("id", "name") VALUES ($1, $2)`,
		users.Insert().Values(Values{"id": 1, "name": "user"}),
		1, "user",
	)
}

var emptyValues = []bool{
	isEmptyValue(reflect.ValueOf(0)),
	isEmptyValue(reflect.ValueOf("")),
	isEmptyValue(reflect.ValueOf(false)),
	isEmptyValue(reflect.ValueOf(0.0)),
	isEmptyValue(reflect.ValueOf(time.Time{})),
}

var nonEmptyValues = []bool{
	isEmptyValue(reflect.ValueOf(1)),
	isEmptyValue(reflect.ValueOf("h")),
	isEmptyValue(reflect.ValueOf(true)),
	isEmptyValue(reflect.ValueOf(0.1)),
	isEmptyValue(reflect.ValueOf(time.Now())),
}

func TestIsEmptyValue(t *testing.T) {
	for i, isEmpty := range emptyValues {
		if !isEmpty {
			t.Errorf("Value %d should be empty", i)
		}
	}
	for i, isNotEmpty := range nonEmptyValues {
		if isNotEmpty {
			t.Errorf("Value %d should not be empty", i)
		}
	}
}
