package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	conf, err := Parse("./example.db.json")
	require.Nil(t, err, "Parsing the example DB config should not error")

	driver, credentials := conf.Credentials()
	assert.Equal(t, "postgres", driver)
	assert.Equal(t,
		`host=localhost port=5432 dbname=aspect_test user=postgres sslmode=disable`,
		credentials,
	)
}
