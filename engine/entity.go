package engine

import (
	"github.com/Ariemeth/frame-assault-2/engine/components"
)

// Entity represents the basic behaviors all engine entities possess.
type Entity interface {
	// ID retrieves the id of this entity.
	ID() string
	// AddComponent adds a component to the entity.
	AddComponent(components.Component)
	// RemoveComponent removes a component from the entity.
	RemoveComponent(components.Component)
	// Components retrievesa list of components attached to the entity.
	Components() []components.Component
}

// EntityOld is the original interface used for entities.  This is being added to allow
// the current scene and model classes to keep functioning as the system is being redone
// to be more like an ECS.
type EntityOld interface {
	// ID retrieves the id of this entity.
	ID() string
}
