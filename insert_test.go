package sol

import "testing"

func TestInsert(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// By default, an INSERT without values will assume a single entry
	expect.SQL(
		contacts.Insert(),
		`INSERT INTO contacts (id, user_id, key, value) VALUES ($1, $2, $3, $4)`,
		nil, nil, nil, nil,
	)

	expect.SQL(
		Insert(users.C("name"), users.C("password")),
		`INSERT INTO users (name, password) VALUES ($1, $2)`,
		nil, nil,
	)

	// Use structs
	admin := user{Name: "admin", Email: "admin@example.com"}
	expect.SQL(
		users.Insert().Values(admin),
		`INSERT INTO users (email, name) VALUES ($1, $2)`,
		"admin@example.com", "admin",
	)
	expect.SQL(
		users.Insert().Values(&admin),
		`INSERT INTO users (email, name) VALUES ($1, $2)`,
		"admin@example.com", "admin",
	)

	exampleUsers := []user{
		admin,
		user{Name: "client", Email: "client@example.com"},
	}

	expect.SQL(
		users.Insert().Values(exampleUsers),
		`INSERT INTO users (email, name) VALUES ($1, $2), ($3, $4)`,
		"admin@example.com", "admin", "client@example.com", "client",
	)
	expect.SQL(
		users.Insert().Values(&exampleUsers),
		`INSERT INTO users (email, name) VALUES ($1, $2), ($3, $4)`,
		"admin@example.com", "admin", "client@example.com", "client",
	)

	// Use sql.Values
	expect.SQL(
		users.Insert().Values(Values{"id": 1, "name": "user"}),
		`INSERT INTO users (id, name) VALUES ($1, $2)`,
		1, "user",
	)

	github := Values{"UserID": 1, "KEY": "github"}
	expect.SQL(
		contacts.Insert().Values(github),
		`INSERT INTO contacts (user_id, key) VALUES ($1, $2)`,
		1, "github",
	)
	expect.SQL(
		contacts.Insert().Values(&github),
		`INSERT INTO contacts (user_id, key) VALUES ($1, $2)`,
		1, "github",
	)

	exampleContacts := []Values{
		github,
		Values{"UserID": 1, "KEY": "bitbucket"},
	}
	expect.SQL(
		contacts.Insert().Values(exampleContacts),
		`INSERT INTO contacts (user_id, key) VALUES ($1, $2), ($3, $4)`,
		1, "github", 1, "bitbucket",
	)
	expect.SQL(
		contacts.Insert().Values(&exampleContacts),
		`INSERT INTO contacts (user_id, key) VALUES ($1, $2), ($3, $4)`,
		1, "github", 1, "bitbucket",
	)

	// Handle errors
	expect.Error(users.Insert().Values("a"))
	expect.Error(users.Insert().Values([]int{1}))
	expect.Error(users.Insert().Values([]struct{}{}))
	expect.Error(users.Insert().Values(nil))
}
