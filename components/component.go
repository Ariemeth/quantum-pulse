package components

// Component presents the behavior all components must have.
type Component interface {
	// Type is expected to return a string representing the type of component.
	Type() string
}
