package sol

import "testing"

func TestUpdate(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Values do not need to be attached to produce an UPDATE statement. It
	// will default to all columns in the table with nil parameters.
	expect.SQL(
		messages.Update(),
		`UPDATE messages SET id = $1, parent_id = $2, text = $3, user_id = $4`,
		nil, nil, nil, nil,
	)

	expect.SQL(
		messages.Update().Values(Values{"text": "hello"}),
		`UPDATE messages SET text = $1`,
		"hello",
	)

	values := Values{"text": "goodbye", "user_id": 2}

	// With Where
	expect.SQL(
		Update(messages).Values(values).Where(messages.C("id").Equals(1)),
		`UPDATE messages SET text = $1, user_id = $2 WHERE messages.id = $3`,
		"goodbye", 2, 1,
	)

	expect.SQL(
		Update(messages).Values(Values{"text": "waka"}).Where(
			messages.C("id").Equals(1),
			messages.C("user_id").Equals(2),
		),
		`UPDATE messages SET text = $1 WHERE (messages.id = $2 AND messages.user_id = $3)`,
		"waka", 1, 2,
	)

	// The statement should have an error if the values map is empty
	expect.Error(messages.Update().Values(Values{}))

	// Attempt to update values with keys that do not correspond to columns
	expect.Error(Update(messages).Values(Values{"nope": "what"}))
}
