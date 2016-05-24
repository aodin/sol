package sol

import (
	"testing"

	"github.com/aodin/sol/types"
)

var tableA = Table("a",
	Column("id", types.Integer()),
	Column("value", types.Varchar()),
)

var tableB = Table("b",
	Column("id", types.Integer()),
	Column("value", types.Varchar()),
)

var relations = Table("relations",
	Column("a_id", types.Integer()),
	Column("b_id", types.Integer()),
	Unique("a_id", "b_id"),
)

func TestJoinClause(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	expect.SQL(
		`SELECT "a"."id", "a"."value" FROM "a" CROSS JOIN "relations"`,
		Select(tableA).CrossJoin(relations),
	)

	expect.SQL(
		`SELECT "a"."id", "a"."value" FROM "a" NATURAL INNER JOIN "relations"`,
		Select(tableA).InnerJoin(relations),
	)

	expect.SQL(
		`SELECT "a"."id", "a"."value" FROM "a" LEFT OUTER JOIN "relations" ON "a"."id" = "relations"."a_id" AND "a"."id" = $1 LEFT OUTER JOIN "b" ON "b"."id" = "relations"."b_id"`,
		Select(tableA).LeftOuterJoin(
			relations,
			tableA.C("id").Equals(relations.C("a_id")),
			tableA.C("id").Equals(2),
		).LeftOuterJoin(
			tableB,
			tableB.C("id").Equals(relations.C("b_id")),
		),
		2,
	)

	// TODO self join with alias
}
