// Package systems provides the various systems used by the engine
package systems

import "github.com/Ariemeth/quantum-pulse/engine/entity"

// System represents the behavior all engine systems should implement.
type System interface {
	// SystemType retrieves the type of system such as renderer, mover, etc.
	SystemType() string
	// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
	AddEntity(e entity.Entity)
	// RemoveEntity removes an Entity from the system.
	RemoveEntity(e entity.Entity)
	// Start causes the system to begin processing all Entities which have been added.
	//	Start()
	// Stop stops the system.  Entities are not removed, but the the system will no perform any additional processing until it is started again.
	//	Stop()
	// Used to determine if the system is currently running and processing available Entities.
	//	IsRunning() bool
}
