package sol

type Modifier interface {
	Modify(*TableElem) error
}
