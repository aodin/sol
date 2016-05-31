package sol

// TODO Merge with Clause?
type Operator interface {
	Wrap(string) string // TODO errors?
}

// Function adds a generic function to the column
func Function(name string, col Columnar) ColumnElem {
	return col.Column().AddOperator(FuncClause{Name: name})
}

// Avg returns a column wrapped in the AVG() function
func Avg(col Columnar) ColumnElem {
	return Function(AVG, col)
}

// Count returns a column wrapped in the COUNT() function
func Count(col Columnar) ColumnElem {
	return Function(COUNT, col)
}

// Date returns a column wrapped in the DATE() function
func Date(col Columnar) ColumnElem {
	return Function(DATE, col)
}

// DatePart returns a column wrapped in the DATE_PART() function
// TODO This method is unsafe - it should not accept direct user input
func DatePart(part string, col Columnar) ColumnElem {
	return col.Column().AddOperator(
		FuncClause{Name: DATEPART},
	).AddOperator(
		ArrayClause{clauses: []Clause{String(part)}, post: true, sep: ", "},
	)
}

// Max returns a column wrapped in the MAX() function
func Max(col Columnar) ColumnElem {
	return Function(MAX, col)
}

// Min returns a column wrapped in the MIN() function
func Min(col Columnar) ColumnElem {
	return Function(MIN, col)
}

// StdDev returns a column wrapped in the STDDEV() function
func StdDev(col Columnar) ColumnElem {
	return Function(STDDEV, col)
}

// Sum returns a column wrapped in the SUM() function
func Sum(col Columnar) ColumnElem {
	return Function(SUM, col)
}

// Variance returns a column wrapped in the VARIANCE() function
func Variance(col Columnar) ColumnElem {
	return Function(VARIANCE, col)
}
