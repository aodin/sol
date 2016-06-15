package postgres

import (
	"fmt"

	"github.com/aodin/sol/dialect"
)

const (
	NowUTC = "now() at time zone 'utc'"
	Now    = "now()"
)

type timestamp struct {
	name         string
	isNotNull    bool
	isUnique     bool
	withTimezone bool
	defaultValue string // TODO Additional defaults?
}

func (t timestamp) Create(d dialect.Dialect) (string, error) {
	compiled := t.name
	if t.withTimezone {
		compiled += " with time zone"
	}
	if t.isNotNull {
		compiled += " NOT NULL"
	}
	if t.isUnique {
		compiled += " UNIQUE"
	}
	if t.defaultValue != "" {
		compiled += fmt.Sprintf(" DEFAULT (%s)", t.defaultValue)
	}
	return compiled, nil
}

func (t timestamp) Default(value string) timestamp {
	t.defaultValue = value
	return t
}

func (t timestamp) NotNull() timestamp {
	t.isNotNull = true
	return t
}

func (t timestamp) Unique() timestamp {
	t.isUnique = true
	return t
}

func (t timestamp) WithoutTimezone() timestamp {
	t.withTimezone = false
	return t
}

func (t timestamp) WithTimezone() timestamp {
	// TODO specify timezone?
	t.withTimezone = true
	return t
}

// TODO Date cannot have a time zone
func Date() (t timestamp) {
	t.name = "date"
	return
}

func Time() (t timestamp) {
	t.name = "time"
	return
}

func Timestamp() (t timestamp) {
	t.name = "timestamp"
	return
}
