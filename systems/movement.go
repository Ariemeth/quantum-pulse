package systems

import (
	"fmt"
	"sync"
	"time"

	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
)

const (
	// TypeMovement is the name of the movement system.
	TypeMovement = "mover"
)

// Movement represents a system that knows how to alter an Entity's position based on its velocities.
type Movement interface {
	System
	// Process updates entities position based on their velocities.
	Process(elapsed float32)
}

type movement struct {
	entities       map[string]movable
	remove         chan entity.Entity
	add            chan entity.Entity
	quit           chan interface{}
	quitProcessing chan interface{}
	runningLock    sync.Mutex
	requirements   []string
	interval       time.Duration
	isRunning      bool
}

// NewMovement creates a new Movement system.
func NewMovement() Movement {
	m := movement{
		entities:       make(map[string]movable, 0),
		remove:         make(chan entity.Entity, 0),
		add:            make(chan entity.Entity, 0),
		quit:           make(chan interface{}),
		quitProcessing: make(chan interface{}),
		requirements: []string{components.TypeTransform,
			components.TypeAcceleration,
			components.TypeVelocity},
		interval:  (1000 / 144) * time.Millisecond,
		isRunning: false,
	}

	go func() {
		for {
			select {
			case ent := <-m.add:
				fmt.Printf("Adding %s to the movement system.\n", ent.ID())
				m.addEntity(ent)
			case ent := <-m.remove:
				fmt.Printf("Removing %s from the movement system.\n", ent.ID())
				m.removeEntity(ent)
			case <-m.quit:
				return
			}
		}
	}()

	return &m
}

// Type retrieves the type of system such as renderer, mover, etc.
func (m *movement) Type() string {
	return TypeMovement
}

// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (m *movement) AddEntity(e entity.Entity) {
	m.add <- e
}

// RemoveEntity removes an Entity from the system.
func (m *movement) RemoveEntity(e entity.Entity) {
	m.remove <- e
}

// IsRunning is useful to check if the movement system is processing entities.
func (m *movement) IsRunning() bool {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()
	return m.isRunning
}

// Start will begin updating entities transform based on their velocity and acceleration.
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
			case <-m.quitProcessing:
				m.isRunning = false
				return
			case <-tr.C:
				current := time.Now().UnixNano()
				// Get the elapsed time in seconds
				elapsed := float32((current - previousTime)) / 1000000000.0
				previousTime = current

				m.Process(elapsed)
				processingTime := time.Now().UnixNano() - current

				if processingTime > 0 {
					tr.Reset(m.interval - time.Duration(processingTime))
				} else {
					tr.Reset(time.Nanosecond)
				}
			}
		}
	}()
}

// Stop Will stop the movement system from moving any of its Entities.
func (m *movement) Stop() {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()
	m.quitProcessing <- true
}

// Terminate stops the movement system and releases all resources.  Once Terminate has been called, the system cannot be reused.
func (m *movement) Terminate() {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()
	m.quitProcessing <- true
	m.quit <- true
}

// Process updates entities position based on their velocities.
func (m *movement) Process(elapsed float32) {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()

	for _, ent := range m.entities {
		// adjust the velocity based on the acceleration.
		updateVelocityUsingAcceleration(elapsed, ent.Acceleration, ent.Velocity)

		// Apply velocity matrix to the transform
		ent.Transform.Update(ent.Velocity.Translational().Mul(elapsed), ent.Velocity.Rotational().Mul(elapsed))
	}
}

// addEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (m *movement) addEntity(e entity.Entity) {
	velocity, isVelocity := e.Component(components.TypeVelocity).(components.Velocity)
	acceleration, isAcceleration := e.Component(components.TypeAcceleration).(components.Acceleration)
	transform, isTransform := e.Component(components.TypeTransform).(components.Transform)

	if isVelocity && isAcceleration && isTransform {
		move := movable{
			Velocity:     velocity,
			Acceleration: acceleration,
			Transform:    transform,
		}
		defer m.runningLock.Unlock()
		m.runningLock.Lock()
		m.entities[e.ID()] = move
	}
}

// removeEntity removes an Entity from the system.
func (m *movement) removeEntity(e entity.Entity) {
	defer m.runningLock.Unlock()
	m.runningLock.Lock()

	delete(m.entities, e.ID())
}

// updateVelocityUsingAcceleration updates a velocity component by adding the acceleration component to the velocity component based on how much time has passed.
func updateVelocityUsingAcceleration(elapsed float32, a components.Acceleration, v components.Velocity) {
	accRot, accTrans := a.Rotational(), a.Translational()
	velRot, velTrans := v.Rotational(), v.Translational()

	velRot = velRot.Add(accRot.Mul(elapsed))
	velTrans = velTrans.Add(accTrans.Mul(elapsed))

	v.Set(velRot, velTrans)
}

type movable struct {
	Acceleration components.Acceleration
	Velocity     components.Velocity
	Transform    components.Transform
}
