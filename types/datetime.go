package types

type datetime struct {
	base
}

func (t datetime) NotNull() datetime {
	t.base.NotNull()
	return t
}

func Timestamp() (t datetime) {
	t.name = "TIMESTAMP"
	return
}
