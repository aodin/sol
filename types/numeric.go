package types

type NumericType struct {
	BaseType
	min, max         int64
	precision, scale int
}

var _ Type = NumericType{}

func (t NumericType) NotNull() NumericType {
	t.BaseType.SetNotNull()
	return t
}

func (t NumericType) Unique() NumericType {
	t.BaseType.SetUnique()
	return t
}

func Decimal(precision, scale int) NumericType {
	// TODO variadic arguments, with all but the first two ignored?
	datatype := Numeric(precision, scale)
	datatype.BaseType.name = "DECIMAL"
	return datatype
}

func Integer() NumericType {
	return NumericType{
		BaseType: Base("INTEGER"),
		min:      -2147483648,
		max:      2147483647,
	}
}

func SmallInt() NumericType {
	return NumericType{
		BaseType: Base("SMALLINT"),
		min:      -32768,
		max:      32767,
	}
}

func BigInt() NumericType {
	return NumericType{
		BaseType: Base("BIGINT"),
		min:      -9223372036854775808,
		max:      9223372036854775807,
	}
}

func Numeric(precision, scale int) NumericType {
	return NumericType{
		BaseType:  Base("NUMERIC"),
		precision: precision,
		scale:     scale,
	}
}

func Float() NumericType {
	return NumericType{BaseType: Base("FLOAT")}
}

func Real() NumericType {
	return NumericType{BaseType: Base("REAL")}
}

func Double() NumericType {
	return NumericType{BaseType: Base("DOUBLE PRECISION")}
}
