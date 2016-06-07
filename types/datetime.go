package types

type datetime struct {
	BaseType
}

var _ Type = datetime{}

func (t datetime) NotNull() datetime {
	t.BaseType.NotNull()
	return t
}

func (t datetime) Unique() datetime {
	t.BaseType.Unique()
	return t
}

func Date() datetime {
	return datetime{BaseType: Base("DATE")}
}

func Datetime() datetime {
	return datetime{BaseType: Base("DATETIME")}
}

func Timestamp() datetime {
	return datetime{BaseType: Base("TIMESTAMP")}
}
