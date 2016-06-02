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
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Ptr {
		return fmt.Errorf(
			"sol: received a non-pointer destination for result.All",
		)
	}

	objElem := objVal.Elem()
	if objElem.Kind() != reflect.Slice {
		return fmt.Errorf(
			"sol: received a non-slice destination for result.All",
		)
	}

	// Get the type of the slice element
	elem := objElem.Type().Elem()

	columns, err := r.Columns()
	if err != nil {
		return fmt.Errorf("sol: error returning columns from result: %s", err)
	}

	total := objElem.Len() // Existing elements
	index := 0             // Current scanner index

	switch elem.Kind() {
	case reflect.Struct:
		fields := DeepFields(reflect.New(elem).Interface())
		aligned := AlignFields(columns, fields)

		// If nothing matched and the number of fields equals the number
		// columns, then blindly align columns and fields
		// TODO This may be too friendly of a feature
		if NoMatchingFields(aligned) && len(fields) == len(columns) {
			aligned = fields
		}

		dest := make([]interface{}, len(columns))
		for r.Next() {
			var newElem reflect.Value
			if index < total {
				// Match the values on the existing object
				newElem = objElem.Index(index)
			} else {
				// Create a new element
				newElem = reflect.New(elem).Elem()
			}

			for i, field := range aligned {
				if field.Exists() {
					dest[i] = newElem.FieldByIndex(field.Type.Index).Addr().Interface()
				} else {
					dest[i] = &dest[i] // Discard
				}
			}

			if err := r.Scan(dest...); err != nil {
				return fmt.Errorf("sol: error while scanning struct: %s", err)
			}

			if index >= total {
				objElem.Set(reflect.Append(objElem, newElem))
			}

			index += 1
		}
	case reflect.Map:
		_, ok := obj.(*[]Values)
		if !ok {
			return fmt.Errorf("sol: slices of maps are only allowed if they are of type sol.Values")
		}

		// TODO How to scan directly into values?
		addr := make([]interface{}, len(columns))
		dest := make([]interface{}, len(columns))
		for i := range addr {
			dest[i] = &addr[i]
		}

		for r.Next() {
			// TODO scan into existing elements?
			if err := r.Scan(dest...); err != nil {
				return fmt.Errorf("sol: error while scanning map: %s", err)
			}

			values := Values{}
			for i, name := range columns {
				values[name] = addr[i]
			}

			objElem.Set(reflect.Append(objElem, reflect.ValueOf(values)))
		}
	default: // TODO enumerate types?
		// Single column results can be scanned into native types
		// TODO scan into existing elements?
		if len(columns) != 1 {
			return fmt.Errorf(
				"sol: unsupported destination for multi-column result: %s",
				elem.Kind(),
			)
		}
		for r.Next() {
			newElem := reflect.New(elem).Elem()
			if err := r.Scan(newElem.Addr().Interface()); err != nil {
				return err
			}
			objElem.Set(reflect.Append(objElem, newElem))
		}
	}
	return r.Err()
}
