package dialect

import (
	"fmt"
	"log"
)

// Dialect is the common interface that all database drivers must implement.
type Dialect interface {
	Param(int) string
}

// Registry of available dialects
var dialects = make(map[string]Dialect)

// Register adds the given Dialect to the registry at the given name.
func Register(name string, d Dialect) {
	if d == nil {
		log.Panic("sol: unable to register a nil Dialect")
	}
	if _, duplicate := dialects[name]; duplicate {
		log.Panic("sol: a Dialect with this name already exists")
	}
	dialects[name] = d
}

// Get returns the Dialect in the registry with the given name. An
// error will be returned if no Dialect with that name exists.
func Get(name string) (Dialect, error) {
	d, ok := dialects[name]
	if !ok {
		return nil, fmt.Errorf(
			"sol: unknown Dialect %s (did you remember to import it?)", name,
		)
	}
	return d, nil
}

// MustGet returns the Dialect in the registry with the given name.
// It will panic if no Dialect with that name is found.
func MustGet(name string) Dialect {
	dialect, err := Get(name)
	if err != nil {
		log.Panic(err)
	}
	return dialect
}
