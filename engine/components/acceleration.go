package components

const (
	// ComponentTypeAcceleration represents an acceleration component's type.
	ComponentTypeAcceleration = "acceleration"
)

// Acceleration represents the acceleration of an entity.
type Acceleration interface {
	Component
}

type acceleration struct {
}

// NewAcceleration creates a new Acceleration component.
func NewAcceleration() Acceleration {
	a := acceleration{}
	return &a
}

// ComponentType retrieves the type name of this component.
func (a *acceleration) ComponentType() string {
	return ComponentTypeAcceleration
}
