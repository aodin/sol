package sol

import (
	"fmt"
)

type FunctionElem struct {
	Columnar
	name  string
	alias string
}

// As sets an alias for this FunctionElem
func (f FunctionElem) As(alias string) Columnar {
	f.alias = alias
	return f
}

// Alias returns the Column's alias
func (f FunctionElem) Alias() string {
	return f.alias
}

// Columns returns the FunctionElem itself in a slice of Columnar. This
// method implements the Selectable interface.
func (f FunctionElem) Columns() []Columnar {
	return []Columnar{f}
}

func (f FunctionElem) FullName() string {
	return fmt.Sprintf(`%s(%s)`, f.name, f.Columnar.FullName())
}

// TODO an entire expression can be columnar
func Function(name string, expression Columnar) FunctionElem {
	return FunctionElem{Columnar: expression, name: name}
}

func Avg(expression Columnar) FunctionElem {
	return Function("avg", expression)
}

func Count(expression Columnar) FunctionElem {
	return Function("count", expression)
}

func Date(expression Columnar) FunctionElem {
	return Function("date", expression)
}

func Max(expression Columnar) FunctionElem {
	return Function("max", expression)
}

func Min(expression Columnar) FunctionElem {
	return Function("min", expression)
}

func StdDev(expression Columnar) FunctionElem {
	return Function("stddev", expression)
}

func Sum(expression Columnar) FunctionElem {
	return Function("sum", expression)
}

func Variance(expression Columnar) FunctionElem {
	return Function("variance", expression)
}
