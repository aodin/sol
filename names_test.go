package sol

import "testing"

// TODO Unicode test cases?
var caseTests = []struct {
	In, Out string
}{
	{In: "SnakeCase", Out: "snake_case"},
	{In: "UserID", Out: "user_id"},
	{In: "UUID", Out: "uuid"},
}

func TestCamelToSnake(t *testing.T) {
	for i, test := range caseTests {
		out := camelToSnake(test.In)
		if out != test.Out {
			t.Errorf(
				"Unexpected camel to snake case conversion %d - %s: %s != %s",
				i, test.In, out, test.Out,
			)
		}
	}
}
