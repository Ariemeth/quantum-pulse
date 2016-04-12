package systems

import (
	"sync"

	"github.com/Ariemeth/quantum-pulse/engine/entity"
)

// Renderer provides the interface needed to process the renderin of Entities.  Each time Process is called all Entities will be rendered.
type Renderer interface {
	Process()
}

// NewRenderer creates a new renderer system.  The renderer system handles rendering all renderable Entities to the screen.
func NewRenderer() System {
	r := renderer{
		entities:     make(map[string]entity.Entity),
		remove:       make(chan entity.Entity, 0),
		add:          make(chan entity.Entity, 0),
		quit:         make(chan interface{}),
		requirements: []string{""},
	}

	return &r
}

// renderer handles the rendering of renderable Entities.
type renderer struct {
	entities     map[string]entity.Entity
	remove       chan entity.Entity
	add          chan entity.Entity
	quit         chan interface{}
	runningLock  sync.Mutex
	isRunning    bool
	requirements []string
}

// SystemType retrieves the name of renderer system type.
func (r *renderer) SystemType() string {
	return "renderer"
}

// AddEntity adds an entity to the renderer to be capable of being rendered to the screen.
func (r *renderer) AddEntity(e entity.Entity) {
	r.add <- e
}

// RemoveEntity removes an entity from the renderer.  Once removed the Entity will no longer be rendered to the screen unless it is added back to the renderer.
func (r *renderer) RemoveEntity(e entity.Entity) {
	r.remove <- e
}

// Start will begin attempting to Render Entities that have been added.  This routine is expected to wait for various channel inputs and act acordingly. The renderer needs to run on the main thread.
func (r *renderer) Start() {
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

// Stop Will stop the renderer from rendering any of its Entities.
func (r *renderer) Stop() {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	if r.isRunning {
		r.quit <- true
		r.isRunning = false
	}
}

// IsRunning is useful to check if the renderer is running its process routine.
func (r *renderer) IsRunning() bool {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	return r.isRunning
}

// Process will render all Entities available to be rendered to the screen.
func (r *renderer) Process() {
	//TODO process all Entities to render them on the screen.
}
