package types

type numeric struct {
	BaseType
	min, max         int64
	precision, scale int
}

var _ Type = numeric{}

func (t numeric) NotNull() numeric {
	t.BaseType.NotNull()
	return t
}

func (t numeric) Unique() numeric {
	t.BaseType.Unique()
	return t
}

func Decimal(precision, scale int) numeric {
	// TODO variadic arguments, with all but the first two ignored?
	datatype := Numeric(precision, scale)
	datatype.BaseType.name = "DECIMAL"
	return datatype
}

func Integer() numeric {
	return numeric{
		BaseType: Base("INTEGER"),
		min:      -2147483648,
		max:      2147483647,
	}
}

func SmallInt() numeric {
	return numeric{
		BaseType: Base("SMALLINT"),
		min:      -32768,
		max:      32767,
	}
}

func BigInt() numeric {
	return numeric{
		BaseType: Base("BIGINT"),
		min:      -9223372036854775808,
		max:      9223372036854775807,
	}
}

func Numeric(precision, scale int) numeric {
	return numeric{
		BaseType:  Base("NUMERIC"),
		precision: precision,
		scale:     scale,
	}
}

func Float() numeric {
	return numeric{BaseType: Base("FLOAT")}
}

func Real() numeric {
	return numeric{BaseType: Base("REAL")}
}

func Double() numeric {
	return numeric{BaseType: Base("DOUBLE PRECISION")}
}
