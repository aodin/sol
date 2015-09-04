package sol

import ()

type Modifier interface {
	Modify(*TableElem) error
}
