package sol

// Modifier is the interface that all Table elements - such as columns
// and constraints - must implement in order to be added to the Table()
// constructor
type Modifier interface {
	Modify(Tabular) error
}
