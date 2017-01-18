package components

// Component presents the behavior all components must have.
type Component interface {
	// ComponentType is expected to return a string representing the type of component.
	ComponentType() string
}
