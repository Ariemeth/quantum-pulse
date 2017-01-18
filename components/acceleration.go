package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// ComponentTypeAcceleration represents an acceleration component's type.
	ComponentTypeAcceleration = "acceleration"
)

// Acceleration represents the acceleration of an entity.
type Acceleration interface {
	Component
	// Rotational retrieves the rotational acceleration in radians per second.
	Rotational() mgl32.Vec3
	// Translational retrieves the translational acceleration.
	Translational() mgl32.Vec3
	// SetRotational sets the rotational acceleration.
	SetRotational(mgl32.Vec3)
	// SetTranslational sets the translational acceleration.
	SetTranslational(mgl32.Vec3)
}

type acceleration struct {
	rotational    mgl32.Vec3
	translational mgl32.Vec3
	dataLock      sync.RWMutex
}

// NewAcceleration creates a new Acceleration component.
func NewAcceleration() Acceleration {
	a := acceleration{
		rotational:    mgl32.Vec3{0, 0, 0},
		translational: mgl32.Vec3{0, 0, 0},
	}
	return &a
}

// ComponentType retrieves the type name of this component.
func (a *acceleration) ComponentType() string {
	return ComponentTypeAcceleration
}

// Rotational retrieves the rotational acceleration in radians per second.
func (a *acceleration) Rotational() mgl32.Vec3 {
	a.dataLock.RLock()
	defer a.dataLock.RUnlock()
	return a.rotational
}

// Translational retrieves the translational acceleration.
func (a *acceleration) Translational() mgl32.Vec3 {
	a.dataLock.RLock()
	defer a.dataLock.RUnlock()
	return a.translational
}

// SetRotational sets the rotational acceleration.
func (a *acceleration) SetRotational(r mgl32.Vec3) {
	a.dataLock.Lock()
	defer a.dataLock.Unlock()
	a.rotational = r
}

// SetTranslational sets the translational acceleration.
func (a *acceleration) SetTranslational(t mgl32.Vec3) {
	a.dataLock.Lock()
	defer a.dataLock.Unlock()
	a.translational = t
}
