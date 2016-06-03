package sol

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Scanner is used for building mock result rows for testing
type Scanner interface {
	Close() error
	Columns() ([]string, error)
	Err() error
	Next() bool
	Scan(...interface{}) error
}

// Result is returned by a database query - it embeds a Scanner
type Result struct {
	stmt string
	Scanner
}

// One returns a single row from Result. The destination must be a pointer
// to a struct or Values type. If there is only a single result column,
// then the destination can be a single native type.
func (r Result) One(obj interface{}) error {
	// Confirm that there is at least one row to return
	if ok := r.Next(); !ok {
		return sql.ErrNoRows
	}

	columns, err := r.Columns()
	if err != nil {
		return fmt.Errorf("sol: error returning columns from result: %s", err)
	}

	// Since maps are already pointers, they can be used as destinations
	// no matter what - as long as they are of type Values
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Map:
		if value.IsNil() {
			// TODO this might be a lie...
			return fmt.Errorf("sol: map types must be initialized before being used as destinations")
		}
		values, ok := obj.(Values)
		if !ok {
			return fmt.Errorf("sol: map types can be destinations only if they are of type Values")
		}

		// TODO scan directly into values?
		addr := make([]interface{}, len(columns))
		dest := make([]interface{}, len(columns))
		for i := range addr {
			dest[i] = &addr[i]
		}

		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error while scanning map: %s", err)
		}

		for i, name := range columns {
			values[name] = addr[i]
		}
		return r.Err()
	case reflect.Ptr: // Do nothing here
	default:
		return fmt.Errorf(
			"sol: received a non-pointer destination for Result.One",
		)
	}

	// Other types must be given as pointers to be valid destinations
	elem := reflect.Indirect(value)
	switch elem.Kind() {
	case reflect.Struct:
		fields := DeepFields(obj)
		aligned := AlignFields(columns, fields)

		// Create an interface pointer for each column's destination.
		// Unmatched scanner values will be discarded
		dest := make([]interface{}, len(columns))

		// If nothing matched and the number of fields equals the number
		// columns, then blindly align columns and fields
		// TODO This may be too friendly of a feature
		if NoMatchingFields(aligned) && len(fields) == len(columns) {
			aligned = fields
		}
		for i, field := range aligned {
			if field.Exists() {
				dest[i] = field.Value.Addr().Interface()
			} else {
				dest[i] = &dest[i] // Discard
			}
		}

		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error while scanning struct: %s", err)
		}

	case reflect.Slice:
		return fmt.Errorf("sol: Result.One cannot scan into slices")
	default: // TODO enumerate types?
		if len(columns) != 1 {
			return fmt.Errorf(
				"sol: unsupported type %T for multi-column Result.One", obj,
			)
		}
		return r.Scan(elem.Addr().Interface()) // Scan directly into the elem
	}
	return r.Err()
}

// All returns all result rows scanned into the given interface, which
// must be a pointer to a slice of either structs or values. If there
// is only a single result column, then the destination can be a
// slice of a single native type (e.g. []int).
func (r Result) All(obj interface{}) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"sol: received a non-pointer destination for result.All",
		)
	}

	list := value.Elem()
	if list.Kind() != reflect.Slice {
		return fmt.Errorf(
			"sol: received a non-slice destination for result.All",
		)
	}

	elem := list.Type().Elem()
	columns, err := r.Columns()
	if err != nil {
		return fmt.Errorf("sol: error returning columns from result: %s", err)
	}

	switch elem.Kind() {
	case reflect.Struct:
		err = r.allStruct(columns, elem, list)
	case reflect.Map:
		err = r.allMap(columns, obj, list)
	default: // TODO enumerate types?
		err = r.allNative(columns, elem, list)
	}
	if err != nil {
		return err
	}
	return r.Err() // Check for delayed scan errors
}

// allStruct scanes the results into a slice of struct types
func (r Result) allStruct(columns []string, elem reflect.Type, list reflect.Value) error {
	fields := DeepFields(reflect.New(elem).Interface())
	aligned := AlignFields(columns, fields)

	// If nothing matched and the number of fields equals the number
	// columns, then blindly align columns and fields
	// TODO This may be too friendly of a feature
	if NoMatchingFields(aligned) && len(fields) == len(columns) {
		aligned = fields
	}

	// How many elements already exist? Merge scanned fields instead of
	// overwriting an entire new element
	existingElements := list.Len()
	index := 0
	dest := make([]interface{}, len(columns))
	for r.Next() {
		var newElem reflect.Value
		if index < existingElements {
			newElem = list.Index(index) // Merge with the existing element
		} else {
			newElem = reflect.New(elem).Elem() // Create a new element
		}

		for i, field := range aligned {
			if field.Exists() {
				dest[i] = newElem.FieldByIndex(field.Type.Index).Addr().Interface()
			} else {
				dest[i] = &dest[i] // Discard
			}
		}

		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error scanning struct: %s", err)
		}

		if index >= existingElements {
			list.Set(reflect.Append(list, newElem))
		}

		index += 1
	}
	return nil
}

// allMap scans the results into a slice of Values
func (r Result) allMap(columns []string, obj interface{}, list reflect.Value) error {
	// TODO support scaning into existing or partially populated slices?
	_, ok := obj.(*[]Values)
	if !ok {
		return fmt.Errorf(
			"sol: slices of maps must have an element type of sol.Values",
		)
	}

	// TODO How to scan directly into values?
	addr := make([]interface{}, len(columns))
	dest := make([]interface{}, len(columns))
	for i := range addr {
		dest[i] = &addr[i]
	}

	for r.Next() {
		if err := r.Scan(dest...); err != nil {
			return fmt.Errorf("sol: error scanning map slice: %s", err)
		}

		values := Values{}
		for i, name := range columns {
			values[name] = addr[i]
		}

		list.Set(reflect.Append(list, reflect.ValueOf(values)))
	}
	return nil
}

// allNative scans the results into a slice of a single native type,
// such as []int
func (r Result) allNative(columns []string, elem reflect.Type, list reflect.Value) error {
	// TODO support scaning into existing or partially populated slices?
	if len(columns) != 1 {
		return fmt.Errorf(
			"sol: unsupported destination for multi-column result: %s",
			elem.Kind(),
		)
	}
	for r.Next() {
		newElem := reflect.New(elem).Elem()
		if err := r.Scan(newElem.Addr().Interface()); err != nil {
			return fmt.Errorf("sol: error scanning native slice: %s", err)
		}
		list.Set(reflect.Append(list, newElem))
	}
	return nil
}
