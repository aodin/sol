package sol

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// mock values
var (
	mockInt   = int64(1)
	mockStr   = "a"
	mockBool  = true
	mockFloat = 1.1
	mockTime  = time.Date(2015, 3, 1, 0, 0, 0, 0, time.UTC)
)

// mock is a mock Scanner used only for testing. It returns an example
// value for every supported type
type mock struct {
	columns        []string
	counter, total int
}

var _ Scanner = &mock{}

func (mock mock) Close() error { return nil }

func (mock mock) Columns() ([]string, error) {
	return mock.columns, nil
}

func (mock mock) Err() error { return nil }

func (mock *mock) Next() bool {
	if mock.counter < mock.total {
		mock.counter += 1
		return true
	}
	return false
}

func (mock mock) Scan(dests ...interface{}) error {
	if len(dests) != len(mock.columns) {
		return fmt.Errorf(
			"Unequal number of scanner destinations (%d) for columns (%d)",
			len(dests), len(mock.columns),
		)
	}
	for i, dest := range dests {
		v := reflect.Indirect(reflect.ValueOf(dest))
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(mockInt)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v.SetUint(uint64(mockInt))
		case reflect.String:
			v.SetString(mockStr)
		case reflect.Bool:
			v.SetBool(mockBool)
		case reflect.Float32, reflect.Float64:
			v.SetFloat(mockFloat)
		case reflect.Interface, reflect.Ptr:
			// Used by Values
			switch mock.columns[i] {
			case "int":
				v.Set(reflect.ValueOf(mockInt))
			case "str":
				v.Set(reflect.ValueOf(mockStr))
			case "bool":
				v.Set(reflect.ValueOf(mockBool))
			case "float":
				v.Set(reflect.ValueOf(mockFloat))
			case "time":
				v.Set(reflect.ValueOf(mockTime))
			}
		case reflect.Struct:
			if t, ok := dest.(*time.Time); ok {
				*t = mockTime
			}
		}
	}
	return nil
}

func mockResult(total int, columns ...string) Result {
	return Result{Scanner: &mock{total: total, columns: columns}}
}

func TestMock(t *testing.T) {
	example := mockResult(1, "int", "str", "bool", "float", "time")

	var num int64
	var str string
	var boolean bool
	var float float64
	var timestamp time.Time

	example.Scanner.Scan(&num, &str, &boolean, &float, &timestamp)

	if num != mockInt {
		t.Errorf("Unequal mock int: have %d, want %d", num, mockInt)
	}
	if str != mockStr {
		t.Errorf("Unequal mock str: have %s, want %s", str, mockStr)
	}
	if boolean != mockBool {
		t.Errorf("Unequal mock bool: have %t, want %t", boolean, mockBool)
	}
	if float != mockFloat {
		t.Errorf("Unequal mock float: have %f, want %f", float, mockFloat)
	}
	if !timestamp.Equal(mockTime) {
		t.Errorf("Unequal mock time: have %v, want %v", timestamp, mockTime)
	}
}

func TestResult_One(t *testing.T) {
	var zero, one, two Result // Example results

	var values Values
	zero = mockResult(0, "int")
	if err := zero.One(values); err == nil {
		t.Errorf("Zero results should error with Result.One")
	}

	one = mockResult(1, "int")
	if err := one.One(values); err == nil {
		t.Errorf("Results.One should error when given an uninitialized map")
	}

	values = Values{}
	one = mockResult(1, "int", "str") // Reset
	if err := one.One(values); err != nil {
		t.Errorf(
			"Result.One should not error when given a Values type: %s",
			err,
		)
	}
	expected := Values{"int": mockInt, "str": mockStr}
	if !reflect.DeepEqual(expected, values) {
		t.Errorf("Unequal Values: %+v != %+v", expected, values)
	}

	one = mockResult(1, "int") // Reset
	ID := struct {
		ID int64
	}{}
	if err := one.One(ID); err == nil {
		t.Errorf("Result.One should error when given a struct type")
	}

	one = mockResult(1, "int") // Reset
	if err := one.One(&ID); err != nil {
		t.Errorf("Result.One should not error when given a *struct type")
	}
	if ID.ID != mockInt {
		t.Errorf("Unequal int: have %d, want %d", ID.ID, mockInt)
	}

	// Match misaligned fields
	two = mockResult(2, "user_id", "is_admin")
	user := struct {
		UserID  int64
		Email   string
		IsAdmin bool
	}{}
	if err := two.One(&user); err != nil {
		t.Errorf("Result.One should not error when given a struct dest")
	}
	if user.UserID != mockInt {
		t.Errorf("Unequal int: have %d, want %d", user.UserID, mockInt)
	}
	if !user.IsAdmin {
		t.Errorf("Unequal bool: have %t, want %t", user.IsAdmin, mockBool)
	}

	// Single addr dest
	var id int64
	one = mockResult(1, "int") // Reset
	if err := one.One(&id); err != nil {
		t.Errorf(
			"Single column results should not error when given a single dest",
		)
	}
	if id != mockInt {
		t.Errorf("Unequal int from single addr: have %d, want %d", id, mockInt)
	}

	two = mockResult(2, "int", "str")
	if err := two.One(&id); err == nil {
		t.Errorf("Result with multiple columns should error when given a single dest")
	}
}

func TestResult_All(t *testing.T) {
	var zero, two Result // Example results

	var values []Values
	zero = mockResult(0, "int")
	if err := zero.All(&values); err != nil {
		t.Errorf("Zero results should not error with Result.All")
	}

	two = mockResult(2, "int", "str")
	if err := two.All(&values); err != nil {
		t.Errorf("Result.All should not error when scanned into []Values")
	}

	if len(values) != 2 {
		t.Fatalf("Unexpected length of values: want 2, have %d", len(values))
	}
	if len(values[0]) != 2 {
		t.Fatalf("Unexpected values: want 2, have %d", len(values[0]))
	}
	// TODO Test actual values
}
