package engine

import "github.com/Ariemeth/frame-assault-2/engine/components"

// EntityOld is the original interface used for entities.  This is being added to allow
// the current scene and model classes to keep functioning as the system is being redone
// to be more like an ECS.
type EntityOld interface {
	// ID retrieves the id of this entity.
	ID() string
}

// Entity represents the basic behaviors all engine entities possess.
type Entity interface {
	// ID retrieves the id of this entity.
	ID() string
	// AddComponent adds a component to the entity.
	AddComponent(components.Component)
	// RemoveComponent removes a component from the entity.
	RemoveComponent(components.Component)
	// Components retrievesa list of components attached to the entity.
	Components() map[string]components.Component
}

// NewEntity creates a new entity with the given ID.
func NewEntity(id string) Entity {
	ent := entity{
		id:         id,
		components: make(map[string]components.Component),
	}
	return &ent
}

type entity struct {
	id          string
	components  map[string]components.Component
	compRequest chan string
}

func (e entity) ID() string {
	return e.id
}

func (e *entity) AddComponent(components.Component) {

}

func (e *entity) RemoveComponent(components.Component) {

}

func (e *entity) Components() map[string]components.Component {
	return e.components
}

func (e *entity) GetComponent(cType string) components.Component {
	return e.components[cType]
}

func (e *entity) ComponentTypes() []string {
	t := make([]string, len(e.components))
	for k := range e.components {
		t = append(t, k)
	}
	return t
}

func (e *entity) Run(quit chan interface{}) {
	for {
		select {
		case <-quit:
			return
		}
	}
}
