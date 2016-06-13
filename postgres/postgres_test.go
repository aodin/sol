package postgres

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aodin/sol"
	"github.com/aodin/sol/types"
)

const travisCI = "host=localhost port=5432 dbname=sol_test user=postgres sslmode=disable"

var testconn *sol.DB
var once sync.Once

// getConn returns a postgres connection pool
func getConn(t *testing.T) *sol.DB {
	// Check if an ENV VAR has been set, otherwise, use travis
	credentials := os.Getenv("SOL_TEST_POSTGRES")
	if credentials == "" {
		credentials = travisCI
	}

	once.Do(func() {
		var err error
		if testconn, err = sol.Open("postgres", credentials); err != nil {
			t.Fatalf("Failed to open connection: %s", err)
		}
		testconn.SetMaxOpenConns(20)
	})
	return testconn
}

var things = sol.Table("things",
	sol.Column("name", types.Varchar()),
	sol.Column("created_at", Timestamp().Default(Now)),
)

var thingsWithoutNow = sol.Table("things",
	sol.Column("name", types.Varchar()),
	sol.Column("created_at", Timestamp()),
)

type thing struct {
	Name      string
	CreatedAt time.Time `db:",omitempty"`
}

var itemsA = Table("items_a",
	sol.Column("id", Serial()),
	sol.Column("name", types.Varchar()),
)

var itemsB = Table("items_b",
	sol.Column("id", Serial()),
	sol.Column("name", types.Varchar()),
	sol.PrimaryKey("id"),
)

var itemsFK = Table("items_fk",
	sol.ForeignKey("id", itemsB, types.Integer().NotNull()),
	sol.Column("name", types.Varchar()),
)

type item struct {
	ID   uint64 `db:",omitempty"`
	Name string
}

func (i item) Exists() bool {
	return i.ID != 0
}

var meetings = Table("meetings",
	sol.Column("uuid", UUID().NotNull().Unique().Default(GenerateV4)),
	sol.Column("time", TimestampRange()),
)

// TestPostGres performs the standard integration test
func TestPostGres(t *testing.T) {
	conn := getConn(t) // TODO close
	sol.IntegrationTest(t, conn, false)
}

func TestPostGres_NullTime(t *testing.T) {
	conn := getConn(t) // TODO close

	tx, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")
	defer tx.Rollback()

	// TODO temp tables
	require.Nil(t,
		tx.Query(thingsWithoutNow.Create().Temporary().IfNotExists()),
		`Create table "%s" should not error`, thingsWithoutNow.Name(),
	)

	type nullthing struct {
		Name      string
		CreatedAt pq.NullTime `db:"created_at"`
	}

	nonzero := nullthing{
		Name:      "a",
		CreatedAt: pq.NullTime{Valid: true, Time: time.Now()},
	}
	tx.Query(thingsWithoutNow.Insert().Values(nonzero))

	var things []nullthing
	tx.Query(thingsWithoutNow.Select(), &things)
	require.Equal(t, 1, len(things))
	assert.True(t, things[0].CreatedAt.Valid)
	assert.False(t, things[0].CreatedAt.Time.IsZero())
}

func TestPostGres_Create(t *testing.T) {
	expect := sol.NewTester(t, &PostGres{})

	expect.SQL(
		itemsFK.Create(),
		`CREATE TABLE items_fk (
  id INTEGER NOT NULL REFERENCES items_b(id),
  name VARCHAR
);`,
	)
}

// TestPostGres_Select tests a variety of SelectStmt features against the
// postgres database
func TestPostGres_Select(t *testing.T) {
	conn := getConn(t) // TODO close

	tx, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")
	defer tx.Rollback()

	// TODO temp tables
	require.Nil(t,
		tx.Query(itemsA.Create().Temporary().IfNotExists()),
		`Create table "%s" should not error`, itemsA.Name(),
	)

	require.Nil(t,
		tx.Query(itemsB.Create().Temporary().IfNotExists()),
		`Create table "%s" should not error`, itemsB.Name(),
	)

	a := item{Name: "A"}
	require.Nil(t,
		tx.Query(itemsA.Insert().Values(a).Returning(), &a),
		`Insert into table "%s" within a transaction should not error`,
		itemsA.Name(),
	)
	if a.ID == 0 {
		t.Fatal("Failed to set primary key of item during INSERT")
	}

	require.Nil(t,
		tx.Query(itemsB.Insert().Values([]item{{Name: "A"}, {Name: "B"}})),
		`Insert into table "%s" within a transaction should not error`,
		itemsB.Name(),
	)

	stmt := itemsB.Select().InnerJoin(
		itemsA,
		itemsB.C("name").Equals(itemsA.C("name")),
	).Where(
		itemsA.C("id").Equals(a.ID),
	)

	var selected item
	tx.Query(stmt, &selected)
	if selected.ID == 0 {
		t.Fatal("Failed to SELECT item through joined table")
	}
}

// TestPostGres_Transaction tests the transactional operations of PostGres,
// including Commit, Rollback, and Close
func TestPostGres_Transaction(t *testing.T) {
	conn := getConn(t) // TODO close

	require.Nil(t, conn.Query(things.Drop().IfExists()))
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
}
