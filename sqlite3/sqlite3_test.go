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
	CreatedAt *time.Time `db:",omitempty"`
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

	require.Nil(t,
		conn.Query(things.Insert().Values(thing{Name: "Alphabet"})),
		`Insert into table "things" should not error`,
	)

	var company thing
	conn.Query(things.Select().Limit(1), &company)

	assert.Equal(t, "Alphabet", company.Name)

	// Start a transaction and roll it back
	tx, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")

	require.Nil(t,
		tx.Query(things.Insert().Values(thing{Name: "Beta"})),
		`Insert into table "things" within a transaction should not error`,
	)

	var all []thing
	require.Nil(t,
		tx.Query(things.Select(), &all),
		`Select from table "things" within a transaction should not error`,
	)
	assert.Equal(t, 2, len(all))

	// Rolling back the transaction should remove the second insert
	tx.Rollback()

	var one []thing
	conn.Query(things.Select(), &one)
	assert.Equal(t, 1, len(one))

	// Test the panicConn
	var another thing
	conn.Must().Query(
		things.Select().OrderBy(things.C("name")).Limit(1), &another,
	)
	assert.Equal(t, "Alphabet", another.Name)

	assert.Panics(t, func() {
		// A non-pointer receiver will error, and with Must(), will panic
		conn.Must().Query(things.Select(), all)
	})
}

// TestSqlite3_Transaction tests the transactional operations of Sqlite3,
// including Commit, Rollback, and Close
func TestSqlite3_Transaction(t *testing.T) {
	conn, err := sql.Open("sqlite3", ":memory:")
	require.Nil(t, err, `Failed to connect to in-memory sqlite3 instance`)
	defer conn.Close()

	require.Nil(t,
		conn.Query(things.Create()),
		`Create table "things" should not error`,
	)

	// Start a transaction with the intent to commit
	tx1, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")

	require.Nil(t,
		tx1.Query(things.Insert().Values(thing{Name: "A"})),
		`Insert into table "things" within a transaction should not error`,
	)

	require.Nil(t, tx1.Commit(), "Committing a transaction should not error")

	var first []thing
	require.Nil(t, conn.Query(things.Select(), &first))
	assert.Equal(t, 1, len(first), "Thing A should have been committed")

	// Start a transaction with the intent to rollback
	tx2, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")

	require.Nil(t, tx2.Query(things.Insert().Values(thing{Name: "B"})))
	require.Nil(t, tx2.Rollback(), "Transaction Rollback should not error")

	var second []thing
	require.Nil(t, conn.Query(things.Select(), &second))
	assert.Equal(t, 1, len(second), "Thing B should not have been committed")

	// Create functions that fail, panic, and are successful
	func() {
		tx3, err := conn.Begin()
		require.Nil(t, err)
		defer tx3.Close()

		require.Nil(t, tx3.Query(things.Insert().Values(thing{Name: "C"})))
		// Without a call to IsSuccessful, the transaction should rollback
	}()

	var third []thing
	require.Nil(t, conn.Query(things.Select(), &third))
	assert.Equal(t, 1, len(third), "Thing C should not have been committed")

	func() {
		tx4, err := conn.Begin()
		require.Nil(t, err)
		defer tx4.Close()

		require.Nil(t, tx4.Query(things.Insert().Values(thing{Name: "D"})))
		tx4.IsSuccessful()
	}()

	var fourth []thing
	require.Nil(t, conn.Query(things.Select(), &fourth))
	assert.Equal(t, 2, len(fourth), "Thing D should have been committed")

	assert.Panics(t,
		func() {
			tx5, err := conn.Begin()
			require.Nil(t, err)
			defer tx5.Close()

			require.Nil(t, tx5.Query(things.Insert().Values(thing{Name: "E"})))
			panic("I'm panicking")
		},
	)

	var fifth []thing
	require.Nil(t, conn.Query(things.Select(), &fifth))
	assert.Equal(t, 2, len(fifth), "Thing E should have been committed")

	// Create a panicTx from the panicConnection
	panicTx, _ := conn.Must().Begin()

	assert.Panics(t, func() {
		// A non-pointer receiver will error, and with Must(), will panic
		panicTx.Query(things.Select(), fifth)
	})

	// A valid transaction will still commit
	var one thing
	panicTx.Query(things.Select(), &one)
	assert.NotEqual(t, "", one.Name)
}
