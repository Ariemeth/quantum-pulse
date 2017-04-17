package systems

import "github.com/Ariemeth/quantum-pulse/entity"

// System represents the behavior all engine systems should implement.
type System interface {
	// Type retrieves the type of system such as renderer, mover, etc.
	Type() string
	// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
	AddEntity(e entity.Entity)
	// RemoveEntity removes an Entity from the system.
	RemoveEntity(e entity.Entity)
	// IsRunning returns true if the system is currently processing entities, false otherwise.
	IsRunning() bool
	// Start causes the system to begin processing all Entities which have been added.
	Start()
	// Stop stops the system.  Entities are not removed, but the the system will not perform any additional processing until it is started again.
	Stop()
	// Terminate stops the system and cleans up all resources.
	Terminate()
}
