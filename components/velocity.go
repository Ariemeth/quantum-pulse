package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// TypeVelocity represents a velocity component's type.
	TypeVelocity = "velocity"
)

// Velocity represents the velocity of an entity.
type Velocity interface {
	Component
	// Rotational retrieves the rotational velocity in radians per second.
	Rotational() mgl32.Vec3
	// Translational retrieves the translational velocity.
	Translational() mgl32.Vec3
	// SetRotational sets the rotational velocity.
	SetRotational(mgl32.Vec3)
	// SetTranslational sets the translational velocity.
	SetTranslational(mgl32.Vec3)
	// Set sets the rotational and translational velocities.
	Set(rotational mgl32.Vec3, translational mgl32.Vec3)
}

type velocity struct {
	rotational    mgl32.Vec3
	translational mgl32.Vec3
	dataLock      sync.RWMutex
}

// NewVelocity creates a new Velocity component.
func NewVelocity() Velocity {
	v := velocity{
		rotational:    mgl32.Vec3{0, 0, 0},
		translational: mgl32.Vec3{0, 0, 0},
	}

	return &v
}

// Type retrieves the type name of this component.
func (v *velocity) Type() string {
	return TypeVelocity
}

// Rotational retrieves the rotational velocity in radians per second.
func (v *velocity) Rotational() mgl32.Vec3 {
	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	return v.rotational
}

// Translational retrieves the translational velocity.
func (v *velocity) Translational() mgl32.Vec3 {
	v.dataLock.RLock()
	defer v.dataLock.RUnlock()
	return v.translational
}

// SetRotational sets the rotational velocity.
func (v *velocity) SetRotational(rotate mgl32.Vec3) {
	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	v.rotational = rotate
}

// SetTranslational sets the translational velocity.
func (v *velocity) SetTranslational(translate mgl32.Vec3) {
	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	v.translational = translate
}

// Set sets the rotational and translational velocities.
func (v *velocity) Set(rotate, translate mgl32.Vec3) {
	v.dataLock.Lock()
	defer v.dataLock.Unlock()
	v.rotational = rotate
	v.translational = translate
}
