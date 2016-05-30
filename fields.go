package sol

import (
	"reflect"
	"time"

	"database/sql"
)

const (
	tagLabel  = "db"
	ignoreTag = "-"
)

var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

// Field holds value and type info on a struct field
type Field struct {
	Value   reflect.Value
	Type    reflect.StructField
	Name    string
	Options options
}

// Exists returns true if the field exists
func (field Field) Exists() bool {
	return field.Name != ""
}

// IsIgnorable returns true if the field can be ignored
func (field Field) IsIgnorable() bool {
	return field.Name == ignoreTag
}

// IsOmittable returns true if the field can be omitted
func (field Field) IsOmittable() bool {
	return field.Options.Has(OmitEmpty) && isEmptyValue(field.Value)
}

// NewField creates a Field from a reflect.Value and Type
func NewField(val reflect.Value, typ reflect.StructField) (field Field) {
	field.Value = val
	field.Type = typ
	field.Name, field.Options = parseTag(typ.Tag.Get(tagLabel))
	if field.Name == "" {
		field.Name = field.Type.Name // Fallback to struct field name
	}
	return
}

// DeepFields returns value and type info on struct types. It will return
// nothing if the given object is not a struct or *struct type.
func DeepFields(obj interface{}) (fields []Field) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	if typ != nil && typ.Kind() == reflect.Ptr {
		typ = val.Elem().Type()
		val = reflect.Indirect(val)
	}
	if typ == nil || typ.Kind() != reflect.Struct {
		return // Do not iterate over non-struct types
	}

	for i := 0; i < typ.NumField(); i++ {
		field := NewField(val.Field(i), typ.Field(i))

		// If the field has an ignore tag, skip it and any descendants
		if field.IsIgnorable() {
			continue
		}

		// Skip fields that cannot be interfaced
		if !field.Value.CanInterface() {
			continue
		}

		// Time structs have special handling
		switch field.Value.Interface().(type) {
		case time.Time, *time.Time:
			fields = append(fields, field)
			continue
		}

		// Scanners have special handling
		if reflect.PtrTo(field.Type.Type).Implements(scannerType) {
			fields = append(fields, field)
			continue
		}

		// Save the field or recurse further
		switch field.Value.Kind() {
		case reflect.Struct:
			fields = append(fields, DeepFields(field.Value.Interface())...)
		default:
			fields = append(fields, field)
		}
	}
	return
}

// AlignFields will reorder the given fields array to match the columns.
// Columns that do not match fields will be given empty field structs.
func AlignFields(columns []string, fields []Field) []Field {
	out := make([]Field, len(columns))

	for i, column := range columns {
		for _, field := range fields {
			// Match names either exactly and using camel to snake
			if field.Name == column || camelToSnake(field.Name) == column {
				out[i] = field
				break
			}
		}
	}
	return out
}

func NoMatchingFields(fields []Field) bool {
	for _, field := range fields {
		if field.Exists() {
			return false
		}
	}
	return true
}
