package engine

// EntityOld is the original interface used for entities.  This is being added to allow
// the current scene and model classes to keep functioning as the system is being redone
// to be more like an ECS.
type EntityOld interface {
	// ID retrieves the id of this entity.
	ID() string
}
