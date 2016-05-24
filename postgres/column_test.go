package postgres

import (
	"testing"

	"github.com/aodin/sol"
)

func TestColumn(t *testing.T) {
	expect := sol.NewTester(t, Dialect())

	expect.SQL(
		meetings.Select().Where(meetings.C("time").Contains("today")),
		`SELECT meetings.uuid, meetings.time FROM meetings WHERE meetings.time @> $1`,
		"today",
	)
}
