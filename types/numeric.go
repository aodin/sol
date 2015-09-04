package types

type numeric struct {
	base
	min, max int64
}

func (t numeric) NotNull() numeric {
	t.base.NotNull()
	return t
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

func Integer() numeric {
	return numeric{
		base: base{
			name: "INTEGER",
		},
		min: -2147483648,
		max: 2147483647,
	}
}
