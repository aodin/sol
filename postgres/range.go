package postgres

import (
	"github.com/aodin/sol/dialect"
)

type rangeType struct {
	name      string
	isNotNull bool
	isUnique  bool
}

func (t rangeType) Create(d dialect.Dialect) (string, error) {
	compiled := t.name
	if t.isNotNull {
		compiled += " NOT NULL"
	}
	if t.isUnique {
		compiled += " UNIQUE"
	}
	return compiled, nil
}

func (t rangeType) NotNull() rangeType {
	t.isNotNull = true
	return t
}

func (t rangeType) Unique() rangeType {
	t.isUnique = true
	return t
}

func Int4Range() (t rangeType) {
	t.name = "INT4RANGE"
	return
}

func Int8Range() (t rangeType) {
	t.name = "INT8RANGE"
	return
}

func NumRange() (t rangeType) {
	t.name = "NUMRANGE"
	return
}

func TimestampRange() (t rangeType) {
	t.name = "TSRANGE"
	return
}

func TimestampWithTimezoneRange() (t rangeType) {
	t.name = "TSTZRANGE"
	return
}

func DateRange() (t rangeType) {
	t.name = "DATERANGE"
	return
}
