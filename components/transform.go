package components

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// TypeTransform represents a transform component's type.
	TypeTransform = "transform"
)

// Transform represents the world position of an entity.
type Transform interface {
	Component
	// Set sets the transform to a specific matrix.
	Set(mgl32.Mat4)
	// Data retrieves the transforms matrix.
	Data() mgl32.Mat4
	// Rotate rotates the transform by the angles passed into the method.
	Rotate(mgl32.Vec3)
	// Translate translates the transform by the value passed into the method.
	Translate(mgl32.Vec3)
	// Update translates the transform based on the first argument then rotates it using the second argument.
	Update(mgl32.Vec3, mgl32.Vec3)
}

// NewTransform creates a new transform component.
func NewTransform() Transform {
	t := transform{
		modelView:   mgl32.Ident4(),
		rotation:    mgl32.Vec3{0, 0, 0},
		translation: mgl32.Vec3{0, 0, 0},
	}
	return &t
}

// transform represents the data of the Transform component.
type transform struct {
	modelView   mgl32.Mat4
	rotation    mgl32.Vec3
	translation mgl32.Vec3
	dataLock    sync.RWMutex
}

// Type retrieves the type of this component.
func (t *transform) Type() string {
	return TypeTransform
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

// Rotate rotates the transform by the angles passed into the method.
func (t *transform) Rotate(rotate mgl32.Vec3) {
	t.dataLock.Lock()
	defer t.dataLock.Unlock()

	t.rotation = t.rotation.Add(rotate)
	total := t.rotation
	rotX := mgl32.HomogRotate3DX(total.X())
	rotY := mgl32.HomogRotate3DY(total.Y())
	rotZ := mgl32.HomogRotate3DZ(total.Z())
	rotMatrix := rotZ.Mul4(rotY).Mul4(rotX)
	trans := t.translation
	t.modelView = mgl32.Ident4().Mul4(mgl32.Translate3D(trans.X(), trans.Y(), trans.Z())).Mul4(rotMatrix)
}

// Translate translates the transform by the value passed into the method.
func (t *transform) Translate(translate mgl32.Vec3) {
	t.dataLock.Lock()
	defer t.dataLock.Unlock()
	t.translation = t.translation.Add(translate)
	trans := t.translation
	t.modelView = t.modelView.Mul4(mgl32.Translate3D(trans.X(), trans.Y(), trans.Z()))
}

// Update translates and rotates the transform.
func (t *transform) Update(translate, rotate mgl32.Vec3) {

	t.dataLock.Lock()
	defer t.dataLock.Unlock()

	t.translation = t.translation.Add(translate)
	trans := t.translation

	t.rotation = t.rotation.Add(rotate)
	total := t.rotation
	rotX := mgl32.HomogRotate3DX(total.X())
	rotY := mgl32.HomogRotate3DY(total.Y())
	rotZ := mgl32.HomogRotate3DZ(total.Z())
	rotMatrix := rotZ.Mul4(rotY).Mul4(rotX)
	trans = t.translation
	t.modelView = mgl32.Ident4().Mul4(mgl32.Translate3D(trans.X(), trans.Y(), trans.Z())).Mul4(rotMatrix)
}
