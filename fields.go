package sol

import (
	"reflect"
	"time"
	"unicode"

	"database/sql"
)

const (
	tagLabel  = "db"
	ignoreTag = "-"
)

var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

// Field holds value and type info on a struct field
type Field struct {
	Value reflect.Value
	Type  reflect.StructField
}

// DeepFields returns value and type info on struct types
func DeepFields(obj interface{}) (fields []Field) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	if typ != nil && typ.Kind() == reflect.Ptr {
		typ = val.Elem().Type()
		val = reflect.Indirect(val)
	}
	if typ == nil || typ.Kind() != reflect.Struct {
		return // Do not iterate on non-struct types
	}

	for i := 0; i < typ.NumField(); i++ {
		field := Field{Value: val.Field(i), Type: typ.Field(i)}

		// If the field has an ignore tag, skip it and any descendants
		if field.Type.Tag.Get(tagLabel) == ignoreTag {
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

type field struct {
	names  []string // struct field names - with possible embedding
	column string   // SQL column name
	table  string   // SQL table name
	options
}

// Exists returns true if the field contains a valid recursive field name
func (f field) Exists() bool {
	return len(f.names) > 0
}

type fields []field

// Empty returns true if none of the fields exist
func (f fields) Empty() bool {
	for _, field := range f {
		if field.Exists() {
			return false
		}
	}
	return true
}

// Has returns true if the given column exists in the fields
func (f fields) Has(column string) bool {
	for _, field := range f {
		// TODO what to do with table name?
		if field.column == column {
			return true
		}
	}
	return false
}

// camelToSnake converts camel case (FieldName) to snake case (field_name)
func camelToSnake(camel string) string {
	if camel == "" {
		return camel
	}
	runes := []rune(camel)
	lowered := unicode.ToLower(runes[0])
	prev := (runes[0] != lowered)
	snake := []rune{lowered}
	for _, char := range runes[1:] {
		lowered := unicode.ToLower(char)
		if !prev && (char != lowered) {
			snake = append(snake, []rune("_")...)
		}
		snake = append(snake, lowered)
		prev = (char != lowered)
	}
	return string(snake)
}

func recurse(names []string, elem reflect.Type) (matches fields) {
	if elem.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < elem.NumField(); i += 1 {
		f := elem.Field(i)

		// Check the tag first to see if this field should be ignored
		tag := f.Tag.Get(tagLabel)
		if tag == ignoreTag {
			continue
		}

		// Continue to build the fields recursively if the field is a struct
		// which does not implement the scanner interface
		if f.Type.Kind() == reflect.Struct && !reflect.PtrTo(f.Type).Implements(scannerType) {
			switch f.Type.String() { // TODO switch on the actual type
			case "time.Time": // TODO confirm this is actually time.Time
			default:
				matches = append(matches, recurse(append(names, f.Name), f.Type)...)
				continue
			}
		}

		// Check the db tag for options
		col, opts := parseTag(tag)

		// Fallback to the field name if no name was given in the tag
		if col == "" {
			col = f.Name
		}

		// A new array will not actually be allocated during every
		// append because capacity is being increased by 2 - make sure to
		// perform a copy to allocate new memory
		namesCopy := make([]string, len(names))
		copy(namesCopy, names)
		field := field{
			names:   append(names, f.Name),
			options: opts,
		}
		field.table, field.column = splitName(col)

		matches = append(matches, field)
	}
	return
}

// AlignColumns will reorder the given fields array to match the columns.
// Columns that do not match fields will be given empty field structs.
func AlignColumns(columns []string, fields []field) fields {
	aligned := make([]field, len(columns))
	// TODO aliases? tables? check if the columns first matches the fields?
	for i, column := range columns {
		for _, field := range fields {
			// Allow snake case columns to be declared in camel case
			alias := camelToSnake(field.column)
			if field.column == column || alias == column {
				field.column = column
				aligned[i] = field
				break
			}
		}
	}
	return aligned
}

// SelectFields returns the ordered list of fields from the given interface.
func SelectFields(v interface{}) fields {
	return recurse([]string{}, reflect.TypeOf(v).Elem())
}

// FieldsFromElem returns the list of fields from the given reflect.Type
func FieldsFromElem(elem reflect.Type) fields {
	return recurse([]string{}, elem)
}
