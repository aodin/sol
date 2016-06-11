package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

const Now = "TIMEZONE('utc'::TEXT, now())"

type timestamp struct {
	types.BaseType
	withTimezone bool
	defaultValue string // TODO Use BaseType default?
}

func (t timestamp) Create(d dialect.Dialect) (string, error) {
	name := t.BaseType.Name()
	if t.withTimezone {
		name += " WITH TIME ZONE"
	}
	compiled := append([]string{name}, t.BaseType.Options()...)

	if t.defaultValue != "" {
		compiled = append(
			compiled, fmt.Sprintf("DEFAULT (%s)", t.defaultValue),
		)
	}
	return strings.Join(compiled, " "), nil
}

func (t timestamp) Default(value string) timestamp {
	t.defaultValue = value
	return t
}

func (t timestamp) NotNull() timestamp {
	t.BaseType.SetNotNull()
	return t
}

func (t timestamp) Unique() timestamp {
	t.BaseType.SetUnique()
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
	return timestamp{BaseType: types.Base("DATE")}
}

func Time() (t timestamp) {
	return timestamp{BaseType: types.Base("TIME")}
}

func Timestamp() (t timestamp) {
	return timestamp{BaseType: types.Base("TIMESTAMP")}
}
