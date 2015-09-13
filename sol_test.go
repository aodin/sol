package sol

import (
	"time"

	"github.com/aodin/sol/types"
)

// Valid schemas should not panic
var users = Table("users",
	Column("id", types.Integer()),
	Column("email", types.Varchar().Limit(256).NotNull()), // TODO unique
	Column("name", types.Varchar().Limit(32).NotNull()),
	Column("password", types.Varchar()),
	Column("created_at", types.Timestamp()),
	PrimaryKey("id"),
	Unique("email"),
)

var contacts = Table("contacts",
	Column("id", types.Integer()),
	Column("user_id", types.Integer()), // TODO FK
	Column("key", types.Varchar()),
	Column("value", types.Varchar()),
	PrimaryKey("id"),
	Unique("user_id", "key"),
)

type user struct {
	ID        uint64 `db:",omitempty"`
	Email     string
	Name      string
	CreatedAt time.Time `db:",omitempty"`
}

type contact struct {
	Id         int64 `db:"id"`
	UserId     int64 `db:"user_id"`
	Key, Value string
}
