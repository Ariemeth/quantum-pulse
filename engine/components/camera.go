package components

//TODO: move project and view handling here.
const (
	// ComponentTypeCamera represents a Camera component's type.
	ComponentTypeCamera = "Camera"
)

// Camera represents the Camera of an entity.
type Camera interface {
	Component
}

type camera struct {
}

// NewCamera creates a new Camera component.
func NewCamera() Camera {
	c := camera{}

	return &c
}

// ComponentType retrieves the type name of this component.
func (c *camera) ComponentType() string {
	return ComponentTypeCamera
}
