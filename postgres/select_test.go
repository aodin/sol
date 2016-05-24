package postgres

import (
	"testing"

	"github.com/aodin/sol"
)

func TestSelect(t *testing.T) {
	expect := sol.NewTester(t, &PostGres{})

	// Build a GROUP BY statement using an aggregate
	expect.SQL(
		sol.Select(
			things.C("name"),
			sol.Max(things.C("created_at")),
		).GroupBy(
			things.C("name"),
		).OrderBy(
			sol.Max(things.C("created_at")).Desc(),
		),
		`SELECT things.name, MAX(things.created_at) FROM things GROUP BY things.name ORDER BY MAX(things.created_at) DESC`,
	)
}
