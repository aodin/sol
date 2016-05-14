package postgres

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aodin/config"
	sql "github.com/aodin/sol"
	"github.com/aodin/sol/types"
)

var travisCI = config.Database{
	Driver:  "postgres",
	Host:    "localhost",
	Port:    5432,
	Name:    "travis_ci_test",
	User:    "postgres",
	SSLMode: "disable",
}

// getConfigOrUseTravis returns the parsed db.json if it exists or the
// travisCI config if it does not
func getConfigOrUseTravis() (config.Database, error) {
	conf, err := config.ParseDatabasePath("./db.json")
	if os.IsNotExist(err) {
		return travisCI, nil
	}
	return conf, err
}

var things = sql.Table("things",
	sql.Column("name", types.Varchar()),
	sql.Column("created_at", Timestamp().Default(Now)),
)

type thing struct {
	Name      string
	CreatedAt time.Time `db:",omitempty"`
}

var itemsA = Table("items_a",
	sql.Column("id", Serial()),
	sql.Column("name", types.Varchar()),
)

var itemsB = Table("items_b",
	sql.Column("id", Serial()),
	sql.Column("name", types.Varchar()),
	sql.PrimaryKey("id"),
)

var itemsFK = Table("items_fk",
	sql.ForeignKey("id", itemsB, types.Integer().NotNull()),
	sql.Column("name", types.Varchar()),
)

type item struct {
	ID   uint64 `db:",omitempty"`
	Name string
}

func (i item) Exists() bool {
	return i.ID != 0
}

var meetings = Table("meetings",
	sql.Column("uuid", UUID().NotNull().Unique().Default(GenerateV4)),
	sql.Column("time", TimestampRange()),
)

// Connect to an PostGres instance and execute some statements.
func TestPostGres(t *testing.T) {
	conf, err := getConfigOrUseTravis()
	if err != nil {
		t.Fatalf("Failed to parse database config: %s", err)
	}

	// TODO in-memory postgres only?
	conn, err := sql.Open(conf.Credentials())
	if err != nil {
		t.Fatalf("Failed to connect to a PostGres instance: %s", err)
	}
	defer conn.Close()

	require.Nil(t, conn.Query(things.Drop().IfExists()))
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

	// Start a transaction and roll it back
	tx, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")

	beta := thing{Name: "Beta"}

	require.Nil(t,
		tx.Query(things.Insert().Values(beta)),
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
}

func TestPostGres_Create(t *testing.T) {
	expect := sql.NewTester(t, &PostGres{})

	expect.SQL(
		`CREATE TABLE "items_fk" (
  "id" INTEGER NOT NULL REFERENCES items_b("id"),
  "name" VARCHAR
);`,
		itemsFK.Create(),
	)
}

func TestPostGres_CRUD(t *testing.T) {
	conf, err := getConfigOrUseTravis()
	require.Nil(t, err, "Failed to parse database config")

	conn, err := sql.Open(conf.Credentials())
	require.Nil(t, err, `Failed to connect to a PostGres instance`)
	defer conn.Close()

	tx, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")
	defer tx.Rollback()

	if err = tx.Query(itemsA.Create().Temporary().IfNotExists()); err != nil {
		t.Fatalf("Create table %s should not error: %s", itemsA.Name(), err)
	}

	google := item{Name: "Google"}

	if err = tx.Query(
		itemsA.Insert().Values(google).Returning(),
		&google,
	); err != nil {
		t.Fatalf("INSERT should not fail %s", err)
	}

	// Update
	tx.Query(
		itemsA.Update().Values(
			sql.Values{"name": "Alphabet"},
		).Where(itemsA.C("id").Equals(google.ID)),
	)

	var alpha item
	if err = tx.Query(
		itemsA.Select().Where(itemsA.C("id").Equals(google.ID)),
		&alpha,
	); err != nil {
		t.Fatalf("Select should not fail: %s", err)
	}

	if google.ID != alpha.ID {
		t.Errorf(
			"Unexpected IDs of google and alphabet: %d != %d",
			google.ID, alpha.ID,
		)
	}
	if alpha.Name != "Alphabet" {
		t.Errorf("Unexpected name for alpha: %s", alpha.Name)
	}

	// Delete
	if err = tx.Query(
		itemsA.Delete().Where(itemsA.C("name").Equals("Alphabet")),
	); err != nil {
		t.Fatalf("Delete should not fail: %s", err)
	}
}

// TestPostGres_Select tests a variety of SelectStmt features against the
// postgres database
func TestPostGres_Select(t *testing.T) {
	conf, err := getConfigOrUseTravis()
	require.Nil(t, err, "Failed to parse database config")

	conn, err := sql.Open(conf.Credentials())
	require.Nil(t, err, `Failed to connect to a PostGres instance`)
	defer conn.Close()

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
	conf, err := getConfigOrUseTravis()
	require.Nil(t, err, "Failed to parse database config")

	// TODO in-memory postgres only?
	conn, err := sql.Open(conf.Credentials())
	require.Nil(t, err, `Failed to connect to a PostGres instance`)
	defer conn.Close()

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
