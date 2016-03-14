package sol

import "testing"

func TestOrder(t *testing.T) {
	expect := NewTester(t, &defaultDialect{})

	// Asc is implied
	ord := OrderedColumn{inner: users.C("id")}
	expect.SQL(`"users"."id"`, ord)

	// Desc
	expect.SQL(`"users"."id" DESC`, ord.Desc())

	// Desc, nulls first
	expect.SQL(
		`"users"."id" DESC NULLS FIRST`,
		ord.Desc().NullsFirst(),
	)

	// Asc, Nulls last
	expect.SQL(`"users"."id" NULLS LAST`, ord.Asc().NullsLast())

	// Calling Orderable on an OrderableColumn should return a copy of itself
	if ord.inner.Name() != ord.Orderable().inner.Name() {
		t.Errorf(
			"Unexpected name of Orderable inner field: %s != %s",
			ord.Orderable().inner.Name(),
			ord.inner.Name(),
		)
	}
}
