package sol

import (
	"time"
)

type Serial struct {
	ID int64 `db:"id,omitempty"`
}

type metadata struct{ Attrs map[string]string }

type customScanner struct{ info []byte }

func (cs customScanner) Scan(src interface{}) error { return nil }

type Nested struct {
	Level2 struct {
		Level3 struct {
			Value bool
		}
	}
}

type embedded struct {
	Serial
	Name      string
	Timestamp struct {
		CreatedAt time.Time `db:",omitempty"`
		UpdatedAt *time.Time
		isActive  bool
	}
	manager  *struct{}
	Custom   customScanner
	Metadata metadata `db:"-"`
	Deep     Nested
}
