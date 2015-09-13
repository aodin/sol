package sol

import (
	"fmt"
	"reflect"

	"github.com/aodin/sol/dialect"
	"github.com/aodin/sol/types"
)

type Columnar interface {
	Selectable
	Alias() string
	As(string) Columnar
	FullName() string
	IsInvalid() bool
	Name() string
	Table() *TableElem
}

type ColumnElem struct {
	name     string
	alias    string
	table    *TableElem
	datatype types.Type
	invalid  bool
}

// Alias returns the Column's alias
func (col ColumnElem) Alias() string {
	return col.alias
}

// As sets an alias for this ColumnElem
func (col ColumnElem) As(alias string) Columnar {
	col.alias = alias
	return col
}

// Columns returns the ColumnElem itself in a slice of ColumnElem. This
// method implements the Selectable interface.
func (col ColumnElem) Columns() []Columnar {
	return []Columnar{col}
}

// Compile produces the dialect specific SQL and adds any parameters
// in the clause to the given Parameters instance
func (col ColumnElem) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	return col.FullName(), nil
}

// Create implements the Creatable interface that outputs a column of a
// CREATE TABLE statement.
func (col ColumnElem) Create(d dialect.Dialect) (string, error) {
	compiled, err := col.datatype.Create(d)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`"%s" %s`, col.Name(), compiled), nil
}

// FullName prefixes the column name with the table name
func (col ColumnElem) FullName() string {
	return fmt.Sprintf(`"%s"."%s"`, col.table.name, col.name)
}

func (col ColumnElem) IsInvalid() bool {
	return col.invalid
}

func (col ColumnElem) Name() string {
	return fmt.Sprintf(`%s`, col.name)
}

// Modify implements the Modifier interface, allowing the ColumnElem to
// modify the given TableElem.
func (col ColumnElem) Modify(table *TableElem) error {
	if table == nil {
		return fmt.Errorf(
			"sol: column %s cannot modify a nil table",
			col.name,
		)
	}
	if err := isValidColumnName(col.name); err != nil {
		return err
	}

	// Add the table to the column
	if col.table != nil {
		return fmt.Errorf(
			"sol: column %s already belongs to table %s",
			col.name, col.table.name,
		)
	}
	col.table = table

	// Add the column to the table
	if err := table.columns.add(col); err != nil {
		return err
	}

	// Add the type to the table creates
	table.creates = append(table.creates, col)

	return nil
}

// Table returns the column's TableElem
func (col ColumnElem) Table() *TableElem {
	return col.table
}

// Conditionals
// ----
const (
	Equal              = "="
	NotEqual           = "<>"
	GreaterThan        = ">"
	GreaterThanOrEqual = ">="
	LessThan           = "<"
	LessThanOrEqual    = "<="
)

func (col ColumnElem) operator(op string, param interface{}) BinaryClause {
	return BinaryClause{
		Pre:  col,
		Post: &Parameter{param},
		Sep:  fmt.Sprintf(" %s ", op),
	}
}

// Equals creates an equals clause that can be used in conditional clauses.
//  table.Select().Where(table.C("id").Equals(3))
func (col ColumnElem) Equals(param interface{}) BinaryClause {
	return col.operator(Equal, param)
}

// DoesNotEqual creates a does not equal clause that can be used in
// conditional clauses.
//  table.Select().Where(table.C("id").DoesNotEqual(3))
func (col ColumnElem) DoesNotEqual(param interface{}) BinaryClause {
	return col.operator(NotEqual, param)
}

// LessThan creates a less than clause that can be used in conditional clauses.
//  table.Select().Where(table.C("id").LessThan(3))
func (col ColumnElem) LessThan(param interface{}) BinaryClause {
	return col.operator(LessThan, param)
}

// GreaterThan creates a greater than clause that can be used in conditional
// clauses.
//  table.Select().Where(table.C("id").GreaterThan(3))
func (col ColumnElem) GreaterThan(param interface{}) BinaryClause {
	return col.operator(GreaterThan, param)
}

// LTE creates a less than or equal to clause that can be used in conditional
// clauses.
//  table.Select().Where(table.C("id").LTE(3))
func (col ColumnElem) LTE(param interface{}) BinaryClause {
	return col.operator(LessThanOrEqual, param)
}

// GTE creates a greater than or equal to clause that can be used in
// conditional clauses.
//  table.Select().Where(table.C("id").GTE(3))
func (col ColumnElem) GTE(param interface{}) BinaryClause {
	return col.operator(GreaterThanOrEqual, param)
}

// IsNull creates a comparison clause that can be used for checking existence
// of NULLs in conditional clauses.
//  table.Select().Where(table.C("name").IsNull())
func (col ColumnElem) IsNull() UnaryClause {
	return UnaryClause{Pre: col, Sep: " IS NULL"}
}

// IsNotNull creates a comparison clause that can be used for checking absence
// of NULLs in conditional clauses.
//  table.Select().Where(table.C("name").IsNotNull())
func (col ColumnElem) IsNotNull() UnaryClause {
	return UnaryClause{Pre: col, Sep: " IS NOT NULL"}
}

// In creates a comparison clause with an IN operator that can be used in
// conditional clauses. An interface is used because the args may be of any
// type: ints, strings...
//  table.Select().Where(table.C("id").In([]int64{1, 2, 3}))
func (col ColumnElem) In(args interface{}) BinaryClause {
	// Create the inner array clause and parameters
	a := ArrayClause{clauses: make([]Clause, 0), sep: ", "}

	// Use reflect to get arguments from the interface only if it is a slice
	s := reflect.ValueOf(args)
	switch s.Kind() {
	case reflect.Slice:
		for i := 0; i < s.Len(); i++ {
			a.clauses = append(a.clauses, &Parameter{s.Index(i).Interface()})
		}
	}
	// TODO What if something other than a slice is given?
	// TODO This statement should be able to take clauses / subqueries
	return BinaryClause{
		Pre:  col,
		Post: FuncClause{Inner: a},
		Sep:  " IN ",
	}
}

func (col ColumnElem) Between(a, b interface{}) Clause {
	return AllOf(col.GTE(a), col.LTE(b))
}

func (col ColumnElem) NotBetween(a, b interface{}) Clause {
	return AnyOf(col.LessThan(a), col.GreaterThan(b))
}

// Ordering
// ----

// Orerable implements the Orderable interface that allows the column itself
// to be used in an OrderBy clause.
func (col ColumnElem) Orderable() OrderedColumn {
	return OrderedColumn{inner: col}
}

// Asc returns an OrderedColumn. It is the same as passing the column itself
// to an OrderBy clause.
func (col ColumnElem) Asc() OrderedColumn {
	return OrderedColumn{inner: col}
}

// Desc returns an OrderedColumn that will sort in descending order.
func (col ColumnElem) Desc() OrderedColumn {
	return OrderedColumn{inner: col, desc: true}
}

// NullsFirst returns an OrderedColumn that will sort NULLs first.
func (col ColumnElem) NullsFirst() OrderedColumn {
	return OrderedColumn{inner: col, nullsFirst: true}
}

// NullsLast returns an OrderedColumn that will sort NULLs last.
func (col ColumnElem) NullsLast() OrderedColumn {
	return OrderedColumn{inner: col, nullsLast: true}
}

// Column is the constructor for a ColumnElem
func Column(name string, datatype types.Type) ColumnElem {
	return ColumnElem{
		name:     name,
		datatype: datatype,
	}
}
