package sol

func Function(name string, col Columnar) ColumnElem {
	return col.Column().AddOperator(name)
}

func Avg(expression Columnar) ColumnElem {
	return Function("AVG", expression)
}

func Count(expression Columnar) ColumnElem {
	return Function("COUNT", expression)
}

func Date(expression Columnar) ColumnElem {
	return Function("DATE", expression)
}

func Max(expression Columnar) ColumnElem {
	return Function("MAX", expression)
}

func Min(expression Columnar) ColumnElem {
	return Function("MIN", expression)
}

func StdDev(expression Columnar) ColumnElem {
	return Function("STDDEV", expression)
}

func Sum(expression Columnar) ColumnElem {
	return Function("SUM", expression)
}

func Variance(expression Columnar) ColumnElem {
	return Function("VARIANCE", expression)
}
