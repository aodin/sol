package postgres

import (
	"fmt"

	"github.com/aodin/sol"
)

const (
	Contains                  = "@>"
	ContainedBy               = "<@"
	Overlap                   = "&&"
	StrictlyLeftOf            = "<<"
	StrictlyRightOf           = ">>"
	DoesNotExtendToTheRightOf = "&<"
	DoesNotExtendToTheLeftOf  = "&>"
	IsAdjacentTo              = "-|-"
	Union                     = "+"
	Intersection              = "*"
	Difference                = "-"
)

// ColumnElem is the postgres dialect's implementation of a SQL column
type ColumnElem struct {
	sol.ColumnElem
}

var _ sol.Columnar = ColumnElem{}

func (col ColumnElem) operator(op string, param interface{}) sol.BinaryClause {
	return sol.BinaryClause{
		Pre:  col,
		Post: &sol.Parameter{Value: param},
		Sep:  fmt.Sprintf(" %s ", op),
	}
}

func (col ColumnElem) Contains(param interface{}) sol.BinaryClause {
	return col.operator(Contains, param)
}

func (col ColumnElem) ContainedBy(param interface{}) sol.BinaryClause {
	return col.operator(ContainedBy, param)
}

func (col ColumnElem) Overlap(param interface{}) sol.BinaryClause {
	return col.operator(Overlap, param)
}

func (col ColumnElem) StrictlyLeftOf(param interface{}) sol.BinaryClause {
	return col.operator(StrictlyLeftOf, param)
}

func (col ColumnElem) StrictlyRightOf(param interface{}) sol.BinaryClause {
	return col.operator(StrictlyRightOf, param)
}

func (col ColumnElem) DoesNotExtendToTheRightOf(param interface{}) sol.BinaryClause {
	return col.operator(DoesNotExtendToTheRightOf, param)
}

func (col ColumnElem) DoesNotExtendToTheLeftOf(param interface{}) sol.BinaryClause {
	return col.operator(DoesNotExtendToTheLeftOf, param)
}

func (col ColumnElem) IsAdjacentTo(param interface{}) sol.BinaryClause {
	return col.operator(IsAdjacentTo, param)
}

func (col ColumnElem) Union(param interface{}) sol.BinaryClause {
	return col.operator(Union, param)
}

func (col ColumnElem) Intersection(param interface{}) sol.BinaryClause {
	return col.operator(Intersection, param)
}

func (col ColumnElem) Difference(param interface{}) sol.BinaryClause {
	return col.operator(Difference, param)
}
