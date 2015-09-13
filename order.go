package sol

import (
	"github.com/aodin/sol/dialect"
)

// Both ColumnElem and OrderedColumns will implement the Orderable interface
type Orderable interface {
	Orderable() OrderedColumn
}

// OrderedColumn represents a ColumnElem that will be used in an ORDER BY
// clause within SELECT statements. It provides additional sorting
// features, such as ASC, DESC, NULLS FIRST, and NULLS LAST.
// If not specified, ASC is assumed by default.
// In Postgres: the default behavior is NULLS LAST when ASC is specified or
// implied, and NULLS FIRST when DESC is specified
// http://www.postgresql.org/docs/9.2/static/sql-select.html#SQL-ORDERBY
type OrderedColumn struct {
	inner                       Columnar
	desc, nullsFirst, nullsLast bool
}

// OrderedColumn should implement the Orderable interface
var _ Orderable = OrderedColumn{}

func (ord OrderedColumn) String() string {
	compiled, _ := ord.Compile(&defaultDialect{}, Params())
	return compiled
}

func (ord OrderedColumn) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	compiled := ord.inner.FullName()
	if ord.desc {
		compiled += " DESC"
	}
	if ord.nullsFirst || ord.nullsLast {
		if ord.nullsFirst {
			compiled += " NULLS FIRST"
		} else {
			compiled += " NULLS LAST"
		}
	}
	return compiled, nil
}

func (ord OrderedColumn) Orderable() OrderedColumn {
	return ord
}

func (ord OrderedColumn) Asc() OrderedColumn {
	ord.desc = false
	return ord
}

func (ord OrderedColumn) Desc() OrderedColumn {
	ord.desc = true
	return ord
}

func (ord OrderedColumn) NullsFirst() OrderedColumn {
	ord.nullsFirst = true
	ord.nullsLast = false
	return ord
}

func (ord OrderedColumn) NullsLast() OrderedColumn {
	ord.nullsFirst = false
	ord.nullsLast = true
	return ord
}
