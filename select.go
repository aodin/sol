package sol

import (
	"fmt"
	"strings"

	"github.com/aodin/sol/dialect"
)

// Selectable is an interface that allows both tables and columns to be
// selected. It is implemented by TableElem and ColumnElem.
type Selectable interface {
	Columns() []Columnar
}

// SelectStmt is the internal representation of an SQL SELECT statement.
type SelectStmt struct {
	ConditionalStmt
	tables     []Tabular
	columns    []Columnar
	joins      []JoinClause
	groupBy    []Columnar
	having     Clause
	orderBy    []OrderedColumn
	isDistinct bool
	distincts  []Columnar
	limit      int
	offset     int
}

// String outputs the parameter-less SELECT statement in a neutral dialect.
func (stmt SelectStmt) String() string {
	compiled, _ := stmt.Compile(&defaultDialect{}, Params())
	return compiled
}

// TODO where should this function live? Also used in postgres.InsertStmt
func CompileColumns(columns []Columnar) []string {
	names := make([]string, len(columns))
	for i, col := range columns {
		// Ignore dialect, parameters and error?
		compiled, _ := col.Compile(nil, nil)
		if col.Alias() != "" {
			compiled += fmt.Sprintf(` AS "%s"`, col.Alias())
		}
		names[i] = compiled
	}
	return names
}

func (stmt SelectStmt) compileTables() []string {
	names := make([]string, len(stmt.tables))
	for i, table := range stmt.tables {
		names[i] = fmt.Sprintf(`"%s"`, table.Name())
	}
	return names
}

func (stmt SelectStmt) Compile(d dialect.Dialect, ps *Parameters) (string, error) {
	if err := stmt.Error(); err != nil {
		return "", err
	}

	compiled := "SELECT"

	// DISTINCT
	if stmt.isDistinct {
		compiled += " DISTINCT"
		if len(stmt.distincts) > 0 {
			distincts := make([]string, len(stmt.distincts))
			for i, col := range stmt.distincts {
				distincts[i] = fmt.Sprintf(`%s`, col.FullName())
			}
			compiled += fmt.Sprintf(
				" ON (%s)", strings.Join(distincts, ", "),
			)
		}
	}

	compiled += fmt.Sprintf(
		" %s FROM %s",
		strings.Join(CompileColumns(stmt.columns), ", "),
		strings.Join(stmt.compileTables(), ", "),
	)

	if len(stmt.joins) > 0 {
		for _, j := range stmt.joins {
			jc, err := j.Compile(d, ps)
			if err != nil {
				return "", err
			}
			compiled += jc
		}
	}

	if stmt.where != nil {
		conditional, err := stmt.where.Compile(d, ps)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" WHERE %s", conditional)
	}

	// GROUP BY ...
	if len(stmt.groupBy) > 0 {
		groupBy := make([]string, len(stmt.groupBy))
		for i, col := range stmt.groupBy {
			groupBy[i] = fmt.Sprintf(`%s`, col.FullName())
		}
		compiled += fmt.Sprintf(" GROUP BY %s", strings.Join(groupBy, ", "))
	}

	//  HAVING
	if stmt.having != nil {
		conditional, err := stmt.having.Compile(d, ps)
		if err != nil {
			return "", err
		}
		compiled += fmt.Sprintf(" HAVING %s", conditional)
	}

	if len(stmt.orderBy) > 0 {
		order := make([]string, len(stmt.orderBy))
		for i, ord := range stmt.orderBy {
			order[i], _ = ord.Compile(d, ps)
		}
		compiled += fmt.Sprintf(" ORDER BY %s", strings.Join(order, ", "))
	}

	if stmt.limit != 0 {
		compiled += fmt.Sprintf(" LIMIT %d", stmt.limit)
	}

	if stmt.offset != 0 {
		compiled += fmt.Sprintf(" OFFSET %d", stmt.offset)
	}
	return compiled, nil
}

func (stmt SelectStmt) hasTable(name string) bool {
	for _, table := range stmt.tables {
		if table.Name() == name {
			return true
		}
	}
	return false
}

// From manually specifies the SelectStmt's FROM clause
func (stmt SelectStmt) From(tables ...Tabular) SelectStmt {
	stmt.tables = tables
	return stmt
}

// All removes the DISTINCT clause from the SELECT statement.
func (stmt SelectStmt) All() SelectStmt {
	stmt.isDistinct = false
	stmt.distincts = nil
	return stmt
}

// Distinct adds a DISTINCT clause to the SELECT statement. If any
// column are provided, the clause will be compiled as a DISTINCT ON.
func (stmt SelectStmt) Distinct(columns ...Columnar) SelectStmt {
	stmt.isDistinct = true
	stmt.distincts = columns
	return stmt
}

func (stmt SelectStmt) join(table Tabular, method string, clauses ...Clause) SelectStmt {
	stmt.joins = append(
		stmt.joins,
		JoinClause{
			method:      method,
			table:       table,
			ArrayClause: ArrayClause{clauses: clauses, sep: " AND "},
		},
	)
	return stmt
}

// CrossJoin adds a CROSS JOIN ... clause to the SELECT statement.
func (stmt SelectStmt) CrossJoin(table Tabular) SelectStmt {
	return stmt.join(table, "CROSS JOIN")
}

// InnerJoin adds an INNER JOIN ... ON ... clause to the SELECT statement.
// If no clauses are given, it will assume the clause is NATURAL.
func (stmt SelectStmt) InnerJoin(table Tabular, clauses ...Clause) SelectStmt {
	return stmt.join(table, "INNER JOIN", clauses...)
}

// LeftOuterJoin adds a LEFT OUTER JOIN ... ON ... clause to the SELECT
// statement. If no clauses are given, it will assume the clause is NATURAL.
func (stmt SelectStmt) LeftOuterJoin(table Tabular, clauses ...Clause) SelectStmt {
	return stmt.join(table, "LEFT OUTER JOIN", clauses...)
}

// RightOuterJoin adds a RIGHT OUTER JOIN ... ON ... clause to the SELECT
// statement. If no clauses are given, it will assume the clause is NATURAL.
func (stmt SelectStmt) RightOuterJoin(table Tabular, clauses ...Clause) SelectStmt {
	return stmt.join(table, "RIGHT OUTER JOIN", clauses...)
}

// FullOuterJoin adds a FULL OUTER JOIN ... ON ... clause to the SELECT
// statement. If no clauses are given, it will assume the clause is NATURAL.
func (stmt SelectStmt) FullOuterJoin(table Tabular, clauses ...Clause) SelectStmt {
	return stmt.join(table, "FULL OUTER JOIN", clauses...)
}

// Where adds a conditional clause to the SELECT statement. Only one WHERE
// is allowed per statement. Additional calls to Where will overwrite the
// existing WHERE clause.
func (stmt SelectStmt) Where(conditions ...Clause) SelectStmt {
	if len(conditions) > 1 {
		// By default, multiple where clauses will be joined using AllOf
		stmt.where = AllOf(conditions...)
	} else if len(conditions) == 1 {
		stmt.where = conditions[0]
	} else {
		// Clear the existing conditions
		stmt.where = nil
	}
	return stmt
}

// GroupBy adds a GROUP BY to the SELECT statement. Only one GROUP BY
// is allowed per statement. Additional calls to GroupBy will overwrite the
// existing GROUP BY clause.
func (stmt SelectStmt) GroupBy(columns ...Columnar) SelectStmt {
	stmt.groupBy = columns
	return stmt
}

// Having adds a conditional clause to the SELECT statement. Only one HAVING
// is allowed per statement. Additional calls to Having will overwrite the
// existing HAVING clause.
func (stmt SelectStmt) Having(conditions ...Clause) SelectStmt {
	if len(conditions) > 1 {
		// By default, multiple having clauses will be joined using AllOf
		stmt.having = AllOf(conditions...)
	} else if len(conditions) == 1 {
		stmt.having = conditions[0]
	} else {
		// Clear the existing conditions
		stmt.having = nil
	}
	return stmt
}

// OrderBy adds an ORDER BY to the SELECT statement. Only one ORDER BY
// is allowed per statement. Additional calls to OrderBy will overwrite the
// existing ORDER BY clause.
func (stmt SelectStmt) OrderBy(ords ...Orderable) SelectStmt {
	stmt.orderBy = make([]OrderedColumn, len(ords))
	// Since columns may be given without an ordering method, perform the
	// orderable conversion whether or not it is already ordered
	for i, column := range ords {
		stmt.orderBy[i] = column.Orderable()
	}
	return stmt
}

// Limit sets the limit of the SELECT statement.
func (stmt SelectStmt) Limit(limit int) SelectStmt {
	// TODO Error (or warning) if limit was already set
	stmt.limit = limit
	return stmt
}

// Offset sets the offset of the SELECT statement.
func (stmt SelectStmt) Offset(offset int) SelectStmt {
	// TODO Error (or warning) if offset was already set
	stmt.offset = offset
	return stmt
}

// SelectTable creates a SELECT statement from the given table and its
// columns. Any additional selections will not have their table added to
// the SelectStmt's tables field - they must be added manually or through
// a join. To perform selections using cartesian logic, use Select() instead.
func SelectTable(table Tabular, selects ...Selectable) (stmt SelectStmt) {
	stmt.tables = []Tabular{table}

	// Add the columns from the alias
	stmt.columns = table.Columns()

	// Add any additional selections
	for _, selection := range selects {
		if selection == nil {
			stmt.AddMeta("sol: received a nil selectable in SelectTable()")
			return
		}
		stmt.columns = append(stmt.columns, selection.Columns()...)
	}

	for _, column := range stmt.columns {
		if column.IsInvalid() {
			stmt.AddMeta(
				"sol: the column %s does not exist", column.FullName(),
			)
			return
		}
	}
	return
}

// Select generates a new SELECT statement from the given columns and tables.
func Select(selections ...Selectable) (stmt SelectStmt) {
	columns := make([]Columnar, 0)
	for _, selection := range selections {
		if selection == nil {
			stmt.AddMeta("sol: received a nil selectable in Select()")
			return
		}
		columns = append(columns, selection.Columns()...)
	}

	if len(columns) < 1 {
		stmt.AddMeta("sol: Select() must be given at least one column")
		return
	}

	for _, column := range columns {
		if column.IsInvalid() {
			stmt.AddMeta(
				"sol: the column %s does not exist", column.FullName(),
			)
			return
		}
		stmt.columns = append(stmt.columns, column)

		// Add the table to the stmt tables if it does not already exist
		if !stmt.hasTable(column.Table().Name()) {
			stmt.tables = append(stmt.tables, column.Table())
		}
	}
	return
}
