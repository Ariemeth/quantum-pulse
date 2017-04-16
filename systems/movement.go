package systems

import (
	"fmt"
	"sync"
	"time"

	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
)

const (
	// SystemTypeMovement is the name of the movement system.
	SystemTypeMovement = "renderer"
)

// Movement represents a system that knows how to alter an Entity's position based on its velocities.
type Movement interface {
	System
	// Process updates entities position based on their velocities.
	Process(elapsed float32)
}

type movement struct {
	entities     map[string]movable
	remove       chan entity.Entity
	add          chan entity.Entity
	quit         chan interface{}
	runningLock  sync.Mutex
	isRunning    bool
	requirements []string
	interval     time.Duration
}

// NewMovement creates a new Movement system.
func NewMovement() Movement {
	m := movement{
		entities:     make(map[string]movable, 0),
		remove:       make(chan entity.Entity, 0),
		add:          make(chan entity.Entity, 0),
		quit:         make(chan interface{}),
		requirements: []string{""}, //TODO eventually add the actual requirements
		interval:     (1000 / 144) * time.Millisecond,
	}

	m.Start()

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

// Start will begin attempting to Render Entities that have been added.  This routine is expected to wait for various channel inputs and act accordingly. The renderer needs to run on the main thread.
func (m *movement) Start() {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()

	if m.isRunning {
		return
	}
	m.isRunning = true

	go func() {

		tr := time.NewTimer(m.interval)
		previousTime := time.Now().UnixNano()

		for {
			select {
			case ent := <-m.add:
				fmt.Printf("Adding %s to the movement system.\n", ent.ID())
				m.addEntity(ent)
			case ent := <-m.remove:
				fmt.Printf("Removing %s from the movement system.\n", ent.ID())
				m.removeEntity(ent)
			case <-m.quit:
				m.isRunning = false
				return
			case since := <-tr.C:
				current := since.UnixNano()
				// Get the elapsed time in seconds
				elapsed := float32((current - previousTime)) / 1000000000.0
				previousTime = current
				m.Process(elapsed)
				fmt.Printf("Unix interval %f\n", elapsed)
				tr.Reset(m.interval)
			}
		}
	}()
}

// Stop Will stop the renderer from rendering any of its Entities.
func (m *movement) Stop() {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()
	if m.isRunning {
		m.quit <- true
	}
}

// IsRunning is useful to check if the renderer is running its process routine.
func (m *movement) IsRunning() bool {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()
	return m.isRunning
}

// Process updates entities position based on their velocities.
func (m *movement) Process(elapsed float32) {
	for _, ent := range m.entities {
		// adjust the velocity based on the acceleration.
		updateVelocityUsingAcceleration(elapsed, ent.Acceleration, ent.Velocity)

		// Apply velocity matrix to the transform
		ent.Transform.Update(ent.Velocity.Translational().Mul(elapsed), ent.Velocity.Rotational().Mul(elapsed))
	}
}

// addEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (m *movement) addEntity(e entity.Entity) {
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

// removeEntity removes an Entity from the system.
func (m *movement) removeEntity(e entity.Entity) {
	delete(m.entities, e.ID())
}

// updateVelocityUsingAcceleration updates a velocity component by adding the acceleration component to the velocity component based on how much time has passed.
func updateVelocityUsingAcceleration(elapsed float32, a components.Acceleration, v components.Velocity) {
	accRot, accTrans := a.Rotational(), a.Translational()
	velRot, velTrans := v.Rotational(), v.Translational()

	velRot = velRot.Add(accRot.Mul(elapsed))
	velTrans = velTrans.Add(accTrans.Mul(elapsed))

	v.Update(velRot, velTrans)
}

type movable struct {
	Acceleration components.Acceleration
	Velocity     components.Velocity
	Transform    components.Transform
}
