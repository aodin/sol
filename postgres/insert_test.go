package postgres

import (
	"testing"

	"github.com/aodin/sol"
)

func TestInsert(t *testing.T) {
	expect := sol.NewTester(t, &PostGres{})

	// By default, an INSERT without values will assume a single entry
	expect.SQL(
		`INSERT INTO "meetings" ("uuid", "time") VALUES ($1, $2) RETURNING "meetings"."uuid", "meetings"."time"`,
		meetings.Insert().Returning(meetings),
		nil, nil,
	)

	// If no parameters are given to Returning(), it will default to the
	// INSERT statement's table
	expect.SQL(
		`INSERT INTO "meetings" ("uuid", "time") VALUES ($1, $2) RETURNING "meetings"."uuid", "meetings"."time"`,
		meetings.Insert().Returning(),
		nil, nil,
	)

	// Selecting a column or table that is not part of the insert table
	// should produce an error
	expect.Error(meetings.Insert().Returning(things))
	expect.Error(meetings.Insert().Returning(things.C("id")))
}
