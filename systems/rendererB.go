package systems

import (
	"sync"

	"github.com/Ariemeth/quantum-pulse/engine/entity"
)

// RendererB provides the interface needed to process the renderin of Entities.  Each time Process is called all Entities will be rendered.
type RendererB interface {
	Process()
}

// NewRendererB creates a new rendererB system.  The rendererB system handles rendering all renderable Entities to the screen.
func NewRendererB() System {
	r := rendererB{
		entities:     make(map[string]entity.Entity),
		remove:       make(chan entity.Entity, 0),
		add:          make(chan entity.Entity, 0),
		quit:         make(chan interface{}),
		requirements: []string{""}, //TODO eventually add the actual requirements
	}

	return &r
}

// rendererB handles the rendering of renderable Entities.
type rendererB struct {
	entities     map[string]entity.Entity
	remove       chan entity.Entity
	add          chan entity.Entity
	quit         chan interface{}
	runningLock  sync.Mutex
	isRunning    bool
	requirements []string
}

// SystemType retrieves the name of rendererB system type.
func (r *rendererB) SystemType() string {
	return "rendererB"
}

// AddEntity adds an entity to the rendererB to be capable of being rendered to the screen.
func (r *rendererB) AddEntity(e entity.Entity) {
	r.add <- e
}

// RemoveEntity removes an entity from the rendererB.  Once removed the Entity will no longer be rendered to the screen unless it is added back to the rendererB.
func (r *rendererB) RemoveEntity(e entity.Entity) {
	r.remove <- e
}

// Start will begin attempting to Render Entities that have been added.  This routine is expected to wait for various channel inputs and act accordingly. The rendererB needs to run on the main thread.
func (r *rendererB) Start() {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()

	if r.isRunning {
		return
	}
	r.isRunning = true

	go func() {
		for {
			select {
			case ent := <-r.add:
				//TODO check to make are component requirements are met
				if _, ok := r.entities[ent.ID()]; !ok {
					r.entities[ent.ID()] = ent
				}
			case ent := <-r.remove:
				delete(r.entities, ent.ID())
			case <-r.quit:
				return
			}
		}
	}()
}

// Stop Will stop the rendererB from rendering any of its Entities.
func (r *rendererB) Stop() {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	if r.isRunning {
		r.quit <- true
		r.isRunning = false
	}
}

// IsRunning is useful to check if the rendererB is running its process routine.
func (r *rendererB) IsRunning() bool {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	return r.isRunning
}

// Process will render all Entities available to be rendered to the screen.
func (r *rendererB) Process() {
	//TODO process all Entities to render them on the screen.
}
