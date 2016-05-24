package sol

import "testing"

func TestParamsRegex(t *testing.T) {
	example := `SELECT
    generate_series::date AS "Date"
FROM generate_series(
    :start_date::timestamp,
    :end_date::timestamp,
    :interval1
)`
	matches := paramsRegex.FindAllString(example, -1)
	if len(matches) != 6 { // This includes false matches
		t.Fatalf(
			"unexpected number of matches: 6 != %d (%v)",
			len(matches), matches,
		)
	}
}

func TestText(t *testing.T) {
	expect := NewTester(t, defaultDialect{})

	expect.SQL(
		Text(
			`SELECT * FROM users WHERE id > :id OR name::varchar LIKE :name`,
			Values{"name": "A", "id": 2},
		),
		`SELECT * FROM users WHERE id > $1 OR name::varchar LIKE $2`,
		2, "A",
	)

	// Missing values
	expect.Error(Text(`SELECT * FROM users WHERE id > :id`))
}
