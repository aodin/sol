package sol

import "testing"

func TestOrder(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Asc is implied
	ord := OrderedColumn{inner: users.C("id")}
	expect.SQL(ord, `users.id`)

	// Desc
	expect.SQL(ord.Desc(), `users.id DESC`)

	// Desc, nulls first
	expect.SQL(
		ord.Desc().NullsFirst(),
		`users.id DESC NULLS FIRST`,
	)

	// Asc, Nulls last
	expect.SQL(
		ord.Asc().NullsLast(),
		`users.id NULLS LAST`,
	)

	// Calling Orderable on an OrderableColumn should return a copy of itself
	if ord.inner.Name() != ord.Orderable().inner.Name() {
		t.Errorf(
			"Unexpected name of Orderable inner field: %s != %s",
			ord.Orderable().inner.Name(),
			ord.inner.Name(),
		)
	}
}
