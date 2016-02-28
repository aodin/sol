package config

import "testing"

func TestDatabase(t *testing.T) {
	conf, err := Parse("./example.db.json")
	if err != nil {
		t.Fatalf("Parsing the example DB config should not error: %s", err)
	}

	driver, credentials := conf.Credentials()
	if driver != "postgres" {
		t.Errorf("Unexpected config driver %s != postgres", driver)
	}
	expected := `host=localhost port=5432 dbname=aspect_test user=postgres sslmode=disable`
	if credentials != expected {
		t.Errorf(
			"Unexpected config credentials %s != %s", credentials, expected,
		)
	}
}
