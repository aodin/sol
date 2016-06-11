package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aodin/sol"
	"github.com/aodin/sol/types"
)

var typetester = Table("typetester",
	sol.Column("id", Serial()),
	sol.Column("string", types.Varchar(30).NotNull()),
	sol.Column("uuid", UUID().NotNull().Unique().Default(GenerateV4)),
	sol.Column("timerange", TimestampRange()),
	sol.Column("timestamp", Timestamp()),
	sol.Column("timestampwith", Timestamp().WithTimezone().Default(Now)),
	sol.Column("boolean", types.Boolean().NotNull().Default(false)),
	sol.Column("text", types.Text()),
	sol.Column("json", JSON()),
	sol.Column("inet", Inet()),
	sol.Column("cidr", Cidr()),
	sol.Column("macaddr", Macaddr()),
	sol.PrimaryKey("id"),
	sol.Unique("string", "uuid"),
	// TODO All the types!
)

// TestP_Select tests a variety of SelectStmt features against the
// postgres database
func TestParseDDL(t *testing.T) {
	conn := getConn(t) // TODO close

	tx, err := conn.Begin()
	require.Nil(t, err, "Creating a new transaction should not error")
	defer tx.Rollback()

	require.Nil(t, tx.Query(typetester.Create().Temporary().IfNotExists()))
	tables, err := ParseDDL(tx)
	require.Nil(t, err)

	var parsed bool
	for _, table := range tables {
		if table.Name() == "typetester" {
			parsed = true
			if table.Create().String() == "" {
				t.Errorf("failed to create table from parsed DDL")
			}
		}
	}

	if !parsed {
		t.Errorf("failed to parse table from DDL")
	}
}
