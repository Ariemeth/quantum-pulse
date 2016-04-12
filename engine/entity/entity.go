package entity

import (
	"sync"

	"github.com/Ariemeth/quantum-pulse/engine/components"
)

// Entity represents the basic behaviors all engine entities possess.
type Entity interface {
	// ID retrieves the id of this entity.
	ID() string
	// AddComponent adds a component to the entity.
	AddComponent(components.Component)
	// RemoveComponent removes a component from the entity.
	RemoveComponent(components.Component)
	// Component accepts a ComponentRequest and returns a component to the channel in the request.  If the Component is not loaded on the entity the receiving channel will receive a nil value.
	Component(ComponentRequest)
	// Start will start main go routine for the entity.  This routine is expected to wait for various channel inputs and act acordingly.
	Start()
	// Stop Will stop the entities main processing routine.
	Stop()
	// IsRunning is useful to check if the entity is running its process routine.
	IsRunning() bool
}

// NewEntity creates a new entity with the given ID.
func NewEntity(id string) Entity {
	ent := entity{
		id:          id,
		components:  make(map[string]components.Component),
		compRequest: make(chan ComponentRequest, 0),
		remove:      make(chan components.Component, 0),
		add:         make(chan components.Component, 0),
		replace:     make(chan components.Component, 0),
		quit:        make(chan interface{}),
	}
	return &ent
}

type entity struct {
	id          string
	components  map[string]components.Component
	compRequest chan ComponentRequest
	remove      chan components.Component
	add         chan components.Component
	replace     chan components.Component
	quit        chan interface{}
	runningLock sync.Mutex
	isRunning   bool
}

// ID retrieves the id of this entity.
func (e *entity) ID() string {
	return e.id
}

// AddComponent adds a component to the entity.
func (e *entity) AddComponent(c components.Component) {
	e.add <- c
}

// RemoveComponent removes a component from the entity.
func (e *entity) RemoveComponent(c components.Component) {
	e.remove <- c
}

// Component accepts a ComponentRequest and returns a component to the channel in the request.  If the Component is not loaded on the entity the receiving channel will receive a nil value.
func (e *entity) ReplaceComponent(c components.Component) {
	e.replace <- c
}

// Start will start main go routine for the entity.  This routine is expected to wait for various channel inputs and act acordingly.
func (e *entity) Component(request ComponentRequest) {
	e.compRequest <- request
}

func (e *entity) ComponentTypes() []string {
	t := make([]string, len(e.components))
	for k := range e.components {
		t = append(t, k)
	}
	return t
}

// Start will start main go routine for the entity.  This routine is expected to wait for various channel inputs and act acordingly.
func (e *entity) Start() {
	defer e.runningLock.Unlock()
	e.runningLock.Lock()

	if e.isRunning {
		return
	}
	e.isRunning = true

	go func() {
		for {
			select {
			case req := <-e.compRequest:
				req.Channel() <- e.components[req.Name()]
			case comp := <-e.add:
				if _, ok := e.components[comp.ComponentType()]; !ok {
					e.components[comp.ComponentType()] = comp
				}
			case comp := <-e.remove:
				delete(e.components, comp.ComponentType())
			case comp := <-e.replace:
				e.components[comp.ComponentType()] = comp
			case <-e.quit:
				return
			}
		}
	}()
}

// Stop Will stop the entities main processing routine.
func (e *entity) Stop() {
	defer e.runningLock.Unlock()
	e.runningLock.Lock()
	if e.isRunning {
		e.quit <- true
		e.isRunning = false
	}
}

// IsRunning is useful to check if the entity is running its process routine.
func (e *entity) IsRunning() bool {
	defer e.runningLock.Unlock()
	e.runningLock.Lock()
	return e.isRunning
}
