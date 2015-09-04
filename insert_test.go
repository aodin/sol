package sol

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
}

func TestIsEmptyValue(t *testing.T) {
	assert.True(t, isEmptyValue(reflect.ValueOf(0)))
	assert.True(t, isEmptyValue(reflect.ValueOf("")))
	assert.True(t, isEmptyValue(reflect.ValueOf(false)))
	assert.True(t, isEmptyValue(reflect.ValueOf(0.0)))
	assert.True(t, isEmptyValue(reflect.ValueOf(time.Time{})))

	assert.False(t, isEmptyValue(reflect.ValueOf(1)))
	assert.False(t, isEmptyValue(reflect.ValueOf("h")))
	assert.False(t, isEmptyValue(reflect.ValueOf(true)))
	assert.False(t, isEmptyValue(reflect.ValueOf(0.1)))
	assert.False(t, isEmptyValue(reflect.ValueOf(time.Now())))
}
