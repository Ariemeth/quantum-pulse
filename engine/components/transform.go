package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// ComponentTypeTransform represents a transform component's type.
	ComponentTypeTransform = "transform"
)

// Transform represents the world position of an entity.
type Transform interface {
	Component
	// Set sets the transform to a specific matrix.
	Set(mgl32.Mat4)
	// Data retrieves the transforms matrix.
	Data() mgl32.Mat4
}

// NewTransform creates a new transform component.
func NewTransform() Transform {
	t := transform{
		modelView: mgl32.Ident4(),
	}
	return &t
}

// transform represents the data of the Transform component.
type transform struct {
	modelView mgl32.Mat4
	dataLock  sync.RWMutex
}

// ComponentType retrieves the type of this component.
func (t *transform) ComponentType() string {
	return ComponentTypeTransform
}

// Set sets the transform to a specific matrix.
func (t *transform) Set(modelView mgl32.Mat4) {
	t.dataLock.Lock()
	defer t.dataLock.Unlock()
	t.modelView = modelView
}

// Data retrieves the transforms matrix.
func (t *transform) Data() mgl32.Mat4 {
	t.dataLock.RLock()
	defer t.dataLock.RUnlock()
	return t.modelView
}
