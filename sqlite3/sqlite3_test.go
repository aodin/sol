package sqlite3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sql "github.com/aodin/sol"
	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

// The sql dialect must implement the dialect interface
var _ dialect.Dialect = &Sqlite3{}

var things = sql.Table("things",
	sql.Column("name", types.Varchar()),
	sql.Column("created_at", types.Timestamp()), // TODO auto-timestamp?
)

type thing struct {
	Name      string
	CreatedAt time.Time `db:",omitempty"`
}

// Connect to an in-memory sqlite3 instance and execute some statements.
func TestSqlite3(t *testing.T) {
	conn, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err, `Failed to connect to in-memory sqlite3 instance`)
	defer conn.Close()

	require.Nil(t,
		conn.Query(things.Create()),
		`Create table "things" should not error`,
	)

	alphabet := thing{
		Name:      "Alphabet",
		CreatedAt: time.Now(),
	}
	require.Nil(t,
		conn.Query(things.Insert().Values(alphabet)),
		`Insert into table "things" should not error`,
	)

	var company thing
	conn.Query(things.Select().Limit(1), &company)

	assert.Equal(t, "Alphabet", company.Name)
	assert.False(t, company.CreatedAt.IsZero())
}
