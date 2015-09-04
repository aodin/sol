package types

import (
	"github.com/aodin/sol/dialect"
)

type Type interface {
	Create(dialect.Dialect) (string, error)
}
