package types

type DatetimeType struct {
	BaseType
}

var _ Type = DatetimeType{}

func (t DatetimeType) NotNull() DatetimeType {
	t.BaseType.SetNotNull()
	return t
}

func (t DatetimeType) Unique() DatetimeType {
	t.BaseType.SetUnique()
	return t
}

func Date() DatetimeType {
	return DatetimeType{BaseType: Base("DATE")}
}

func Datetime() DatetimeType {
	return DatetimeType{BaseType: Base("DATETIME")}
}

func Timestamp() DatetimeType {
	return DatetimeType{BaseType: Base("TIMESTAMP")}
}
