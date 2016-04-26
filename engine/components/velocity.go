package components

const (
	// ComponentTypeVelocity represents a velocity component's type.
	ComponentTypeVelocity = "velocity"
)

// Velocity represents the velocity of an entity.
type Velocity interface {
	Component
}

type velocity struct {
}

// NewVelocity creates a new Velocity component.
func NewVelocity() Velocity {
	v := velocity{}

	return &v
}

// ComponentType retrieves the type name of this component.
func (v *velocity) ComponentType() string {
	return ComponentTypeVelocity
}
