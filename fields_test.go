package sol

import (
	"time"
)

type Serial struct {
	ID uint64 `db:"id,omitempty"`
}

type ignored struct {
	manager struct{} `db:"-"`
	name    string
}

type embedded struct {
	Serial
	Name      string
	Timestamp struct {
		CreatedAt time.Time `db:",omitempty"`
		UpdatedAt *time.Time
		isActive  bool
	}
	manager *struct{}
}

type nested struct {
	Another string
	embedded
}

type moreNesting struct {
	nested
	OneMore string `db:"text,omitempty"`
}
