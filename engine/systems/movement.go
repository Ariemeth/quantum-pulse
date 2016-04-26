package systems

import "github.com/Ariemeth/quantum-pulse/engine/entity"

const (
	// SystemTypeMovement is the name of the movement system.
	SystemTypeMovement = "renderer"
)

// Movement represents a system that knows how to alter an Entity's position based on its velocities.
type Movement interface {
	System
	// Process updates entities position based on their velocities.
	Process()
}

type movement struct {
}

// NewMovement creates a new Movement system.
func NewMovement() Movement {
	m := movement{}

	return &m
}

// SystemType retrieves the type of system such as renderer, mover, etc.
func (m *movement) SystemType() string {
	return SystemTypeMovement
}

// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (m *movement) AddEntity(e entity.Entity) {

}

// RemoveEntity removes an Entity from the system.
func (m *movement) RemoveEntity(e entity.Entity) {

}

// Process updates entities position based on their velocities.
func (m *movement) Process() {

}
