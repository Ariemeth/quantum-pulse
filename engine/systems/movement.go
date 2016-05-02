package systems

import (
	"github.com/Ariemeth/quantum-pulse/engine/components"
	"github.com/Ariemeth/quantum-pulse/engine/entity"
	"github.com/go-gl/mathgl/mgl32"
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
		t := ent.Transform.Data()
		//		accRot := ent.Acceleration.Rotational()
		//		accTran := ent.Acceleration.Translational()
		//		vel := ent.Velocity.Data()
		//		rotX := mgl32.HomogRotate3D(e*accRot.X(), mgl32.Vec3{1, 0, 0})
		//		rotY := mgl32.HomogRotate3D(e*accRot.Y(), mgl32.Vec3{0, 1, 0})
		//		rotZ := mgl32.HomogRotate3D(e*accRot.Z(), mgl32.Vec3{0, 0, 1})
		//		accResult := rotX.Mul4(rotY.Mul4(rotZ))
		//		accResult = accResult.Mul4(mgl32.Translate3D(accTran.X(), accTran.Y(), accTran.Z()))

		// adjust the velocity based on the acceleration.
		//		vel = vel.Mul4(accResult)
		//		ent.Velocity.Set(vel)

		//		result := t.Mul4(vel)
		//		ent.Velocity.Set(result)
		velMatrix := mgl32.HomogRotate3D(float32(ent.Velocity.Rotational().Y()*e), mgl32.Vec3{0, 1, 0})
		ent.Transform.Set(t.Mul4(velMatrix))
		//	ent.Transform.Set(t.Mul4(result))
		//ent.Transform.Set(t.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0}))
	}
}

type movable struct {
	Acceleration components.Acceleration
	Velocity     components.Velocity
	Transform    components.Transform
}
