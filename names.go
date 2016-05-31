package sol

import (
	"fmt"
	"unicode"
)

// TODO What are the actual rules for column and table names?

func isValidColumnName(name string) error {
	if name == "" {
		return fmt.Errorf("sol: column names cannot be blank")
	}
	return nil
}

func isValidTableName(name string) error {
	if name == "" {
		return fmt.Errorf("sol: table names cannot be blank")
	}
	return nil
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

// Aliases track column names during camel to snake case conversion
type Aliases map[string]string

// Keys returns the keys of the map in unspecified order
func (aliases Aliases) Keys() []string {
	keys := make([]string, len(aliases))
	var i int
	for key := range aliases {
		keys[i] = key
		i += 1
	}
	return keys
}
