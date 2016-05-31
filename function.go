package sol

// Function adds a generic function to the column
func Function(name string, col Columnar) ColumnElem {
	return col.Column().AddOperator(name)
}

// Avg returns a column wrapped in the AVG() function
func Avg(expression Columnar) ColumnElem {
	return Function("AVG", expression)
}

// Count returns a column wrapped in the COUNT() function
func Count(expression Columnar) ColumnElem {
	return Function("COUNT", expression)
}

// Date returns a column wrapped in the DATE() function
func Date(expression Columnar) ColumnElem {
	return Function("DATE", expression)
}

// Max returns a column wrapped in the MAX() function
func Max(expression Columnar) ColumnElem {
	return Function("MAX", expression)
}

// Min returns a column wrapped in the MIN() function
func Min(expression Columnar) ColumnElem {
	return Function("MIN", expression)
}

// StdDev returns a column wrapped in the STDDEV() function
func StdDev(expression Columnar) ColumnElem {
	return Function("STDDEV", expression)
}

// Sum returns a column wrapped in the SUM() function
func Sum(expression Columnar) ColumnElem {
	return Function("SUM", expression)
}

// Variance returns a column wrapped in the VARIANCE() function
func Variance(expression Columnar) ColumnElem {
	return Function("VARIANCE", expression)
}
