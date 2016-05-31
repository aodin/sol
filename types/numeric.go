package types

type numeric struct {
	base
	min, max         int64
	precision, scale int
}

var _ Type = numeric{}

func (t numeric) NotNull() numeric {
	t.base.NotNull()
	return t
}

func Decimal(precision, scale int) numeric {
	typ := Numeric(precision, scale)
	typ.base.name = "DECIMAL"
	return typ
}

func Integer() numeric {
	return numeric{
		base: base{
			name: "INTEGER",
		},
		min: -2147483648,
		max: 2147483647,
	}
}

func SmallInt() numeric {
	return numeric{
		base: base{
			name: "SMALLINT",
		},
		min: -32768,
		max: 32767,
	}
}

func BigInt() numeric {
	return numeric{
		base: base{
			name: "BIGINT",
		},
		min: -9223372036854775808,
		max: 9223372036854775807,
	}
}

func Numeric(precision, scale int) numeric {
	return numeric{
		base: base{
			name: "NUMERIC",
		},
		precision: precision,
		scale:     scale,
	}
}

func Float() numeric {
	return numeric{
		base: base{
			name: "FLOAT",
		},
	}
}

func Real() numeric {
	return numeric{
		base: base{
			name: "REAL",
		},
	}
}

func Double() numeric {
	return numeric{
		base: base{
			name: "DOUBLE PRECISION",
		},
	}
}
