package sol

func Function(name string, col Columnar) ColumnElem {
	return col.AddOperator(name)
}

func Avg(expression Columnar) ColumnElem {
	return Function("avg", expression)
}

func Count(expression Columnar) ColumnElem {
	return Function("count", expression)
}

func Date(expression Columnar) ColumnElem {
	return Function("date", expression)
}

func Max(expression Columnar) ColumnElem {
	return Function("max", expression)
}

func Min(expression Columnar) ColumnElem {
	return Function("min", expression)
}

func StdDev(expression Columnar) ColumnElem {
	return Function("stddev", expression)
}

func Sum(expression Columnar) ColumnElem {
	return Function("sum", expression)
}

func Variance(expression Columnar) ColumnElem {
	return Function("variance", expression)
}
