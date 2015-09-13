package postgres

import (
	"testing"

	"github.com/aodin/sol"
)

func TestColumn(t *testing.T) {
	expect := sol.NewTester(t, &PostGres{})

	expect.SQL(
		`SELECT "meetings"."uuid", "meetings"."time" FROM "meetings" WHERE "meetings"."time" @> $1`,
		meetings.Select().Where(meetings.C("time").Contains("today")),
		"today",
	)
}
