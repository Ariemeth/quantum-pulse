package systems

import (
	"github.com/Ariemeth/quantum-pulse/engine/components"
	"github.com/Ariemeth/quantum-pulse/engine/entity"
)

const (
	// SystemTypeMovement is the name of the movement system.
	SystemTypeMovement = "renderer"
)

// Movement represents a system that knows how to alter an Entity's position based on its velocities.
type Movement interface {
	System
	// Process updates entities position based on their velocities.
	Process(elapsed float64)
}

type movement struct {
	entities map[string]movable
}

// NewMovement creates a new Movement system.
func NewMovement() Movement {
	m := movement{
		entities: make(map[string]movable, 0),
	}

	return &m
}

// SystemType retrieves the type of system such as renderer, mover, etc.
func (m *movement) SystemType() string {
	return SystemTypeMovement
}

// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (m *movement) AddEntity(e entity.Entity) {
	velocity, isVelocity := e.Component(components.ComponentTypeVelocity).(components.Velocity)
	acceleration, isAcceleration := e.Component(components.ComponentTypeAcceleration).(components.Acceleration)
	transform, isTransform := e.Component(components.ComponentTypeTransform).(components.Transform)

	if isVelocity && isAcceleration && isTransform {
		move := movable{
			Velocity:     velocity,
			Acceleration: acceleration,
			Transform:    transform,
		}
		m.entities[e.ID()] = move
	}
}

// RemoveEntity removes an Entity from the system.
func (m *movement) RemoveEntity(e entity.Entity) {
	delete(m.entities, e.ID())
}

// Process updates entities position based on their velocities.
func (m *movement) Process(elapsed float64) {
	e := float32(elapsed)
	for _, ent := range m.entities {
		// adjust the velocity based on the acceleration.
		updateVelocityUsingAcceleration(e, ent.Acceleration, ent.Velocity)

		// Apply velocity matrix to the transform
		ent.Transform.Translate(ent.Velocity.Translational().Mul(e))
		ent.Transform.Rotate(ent.Velocity.Rotational().Mul(e))
	}
}

// updateVelocityUsingAcceleration updates a velocity component by adding the acceleration component to the velocity component based on how much time has passed.
func updateVelocityUsingAcceleration(elapsed float32, a components.Acceleration, v components.Velocity) {
	accRot, accTrans := a.Rotational(), a.Translational()
	velRot, velTrans := v.Rotational(), v.Translational()

	velRot = velRot.Add(accRot.Mul(elapsed))
	velTrans = velTrans.Add(accTrans.Mul(elapsed))

	v.SetRotational(velRot)
	v.SetTranslational(velTrans)
}

type movable struct {
	Acceleration components.Acceleration
	Velocity     components.Velocity
	Transform    components.Transform
}
