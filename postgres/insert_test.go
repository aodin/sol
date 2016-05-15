package postgres

import (
	"testing"
	"time"

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

	// UPSERT
	now := time.Now()
	expect.SQL(
		`INSERT INTO "meetings" ("uuid", "time") VALUES ($1, $2) ON CONFLICT DO UPDATE SET "time" = $3 WHERE "meetings"."time" >= $4`,
		meetings.Insert().OnConflict().DoUpdate(sol.Values{"time": now}).Where(meetings.C("time").GTE(now)),
		nil, nil, now, now,
	)

	expect.SQL(
		`INSERT INTO "meetings" ("uuid", "time") VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		meetings.Insert().OnConflict().DoNothing(),
		nil, nil,
	)

	// Selecting a column or table that is not part of the insert table
	// should produce an error
	expect.Error(meetings.Insert().Returning(things))
	expect.Error(meetings.Insert().Returning(things.C("id")))
}
