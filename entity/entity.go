package entity

import (
	"errors"
	"sync"

	"github.com/Ariemeth/quantum-pulse/components"
)

// Entity represents the basic behaviors all engine entities possess.
type Entity interface {
	// ID retrieves the id of this entity.
	ID() string
	// AddComponent adds a component to the entity.  Returns an error if the component type has already been added.
	AddComponent(components.Component) error
	// RemoveComponent removes a component from the entity.  If the component type is not already attached an error will be returned.
	RemoveComponent(string) error
	// ReplaceComponent replaces an existing component with a new one of the same type.
	ReplaceComponent(c components.Component) error
	// Component retrieves a component based on its component type.
	Component(string) components.Component
	// Retrieves a list of component types that have been attached to this entity.
	ComponentTypes() []string
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
	id         string
	components map[string]components.Component
	lock       sync.RWMutex
}

// ID retrieves the id of this entity.
func (e *entity) ID() string {
	return e.id
}

// AddComponent adds a component to the entity.  Returns an error if the component type has already been added.
func (e *entity) AddComponent(c components.Component) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if _, ok := e.components[c.Type()]; !ok {
		e.components[c.Type()] = c
		return nil
	}
	return errors.New("component type already attached")
}

// RemoveComponent removes a component from the entity.
func (e *entity) RemoveComponent(compType string) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if _, ok := e.components[compType]; ok {
		delete(e.components, compType)
		return nil
	}
	return errors.New("component type not attached")
}

// ReplaceComponent replaces an existing component with a new one of the same type.  If the component type has not already been attached returns an error.
func (e *entity) ReplaceComponent(c components.Component) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if _, ok := e.components[c.Type()]; ok {
		e.components[c.Type()] = c
		return nil
	}
	return errors.New("component type not already attached")
}

// Component retrieves a component based on its component type.
func (e *entity) Component(compType string) components.Component {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.components[compType]
}

// Retrieves a list of component types that have been attached to this entity.
func (e *entity) ComponentTypes() []string {
	e.lock.RLock()
	defer e.lock.RUnlock()
	t := make([]string, len(e.components))
	for k := range e.components {
		t = append(t, k)
	}
	return t
}
