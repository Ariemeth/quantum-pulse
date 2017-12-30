package systems

import (
	"log"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
	am "github.com/Ariemeth/quantum-pulse/resources"
)

const (
	// TypeRenderer is the name of the renderer system.
	TypeRenderer = "renderer"
)

// Renderer provides the interface needed to process the rendering of Entities.  Each time Process is called all Entities will be rendered.
type Renderer interface {
	System
	// Process renders all renderable entities.
	Process()
	// LoadCamera sets the camera to be used by the display.
	LoadCamera(camera components.Camera)
}

type renderer struct {
	entities       map[string]renderable
	assets         *am.Manager
	camera         components.Camera
	mainFunc       func(f func())
	window         *glfw.Window
	remove         chan entity.Entity
	add            chan entity.Entity
	quit           chan interface{}
	quitProcessing chan interface{}
	runningLock    sync.Mutex
	isRunning      bool
	requirements   []string
	interval       time.Duration
}

// NewRenderer creates a new renderer system.  The renderer system handles rendering all renderable Entities to the screen.
func NewRenderer(assetManager *am.Manager, mainFunc func(f func()), window *glfw.Window) Renderer {
	r := renderer{
		entities:       make(map[string]renderable),
		assets:         assetManager,
		camera:         components.NewCamera(),
		mainFunc:       mainFunc,
		window:         window,
		remove:         make(chan entity.Entity, 0),
		add:            make(chan entity.Entity, 0),
		quit:           make(chan interface{}),
		quitProcessing: make(chan interface{}),
		requirements:   []string{components.TypeTransform, components.TypeMesh},
		interval:       (1000 / 144) * time.Millisecond,
	}

	go func() {
		for {
			select {
			case ent := <-r.add:
				log.Printf("Adding %s to the rendering system.\n", ent.ID())
				r.addEntity(ent)
			case ent := <-r.remove:
				log.Printf("Removing %s to the rendering system.\n", ent.ID())
				r.removeEntity(ent)
			case <-r.quit:
				return
			}
		}
	}()

	return &r
}

// Type retrieves the type of system such as renderer, mover, etc.
func (r *renderer) Type() string {
	return TypeRenderer
}

// AddEntity adds an entity to the renderer to be capable of being rendered to the screen.
func (r *renderer) AddEntity(e entity.Entity) {
	r.add <- e
}

// RemoveEntity removes an entity from the renderer.  Once removed the Entity will no longer be rendered to the screen unless it is added back to the renderer.
func (r *renderer) RemoveEntity(e entity.Entity) {
	r.remove <- e
}

// IsRunning is useful to check if the renderer is running its process routine.
func (r *renderer) IsRunning() bool {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	return r.isRunning
}

// Start will begin attempting to Render Entities that have been added.  This routine is expected to wait for various channel inputs and act accordingly. The renderer needs to run on the main thread.
func (r *renderer) Start() {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()

	if r.isRunning {
		return
	}
	r.isRunning = true

	go func() {

		tr := time.NewTimer(r.interval)

		for {
			select {
			case <-r.quitProcessing:
				r.isRunning = false
				return
			case <-tr.C:
				current := time.Now().UnixNano()

				r.Process()
				processingTime := time.Now().UnixNano() - current

				if processingTime > 0 {
					tr.Reset(r.interval - time.Duration(processingTime))
				} else {
					tr.Reset(time.Nanosecond)
				}
			}
		}
	}()
}

// Stop Will stop the renderer from rendering any of its Entities.
func (r *renderer) Stop() {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	if r.isRunning {
		r.quitProcessing <- true
	}
}

func (r *renderer) Terminate() {
	defer r.runningLock.Unlock()
	r.runningLock.Lock()
	r.quitProcessing <- true
	r.quit <- true
}

func (r *renderer) Process() {

	r.mainFunc(func() {
		r.runningLock.Lock()
		defer r.runningLock.Unlock()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, ent := range r.entities {
			md := ent.Mesh.Data()
			shader, status := r.assets.Shaders().GetShaderProgram(ent.ProgramID)
			if !status {
				// if the shader is not loaded there is no point in trying to render it to the screen.
				continue
			}
			gl.UseProgram(ent.ProgramID)

			gl.BindVertexArray(ent.VAO)
			// set the camera uniform.  TODO this only needs to be done once per shader
			projection := r.camera.Projection()
			gl.UniformMatrix4fv(shader.GetUniformLoc(am.ProjectionUniform), 1, false, &projection[0])
			// set the view uniform. TODO this only needs to be done once per shader
			view := r.camera.View()
			gl.UniformMatrix4fv(shader.GetUniformLoc(am.CameraUniform), 1, false, &view[0])

			td := ent.Transform.Data()
			gl.UniformMatrix4fv(shader.GetUniformLoc(am.ModelUniform), 1, false, &td[0])
			gl.ActiveTexture(gl.TEXTURE0)
			gl.Uniform1i(shader.GetUniformLoc(am.TextureUniform), 0)

			texture, isLoaded := r.assets.Textures().GetTexture(md.TextureFile)
			if isLoaded {
				gl.BindTexture(gl.TEXTURE_2D, texture)
			}

			if md.Indexed {
				gl.DrawElements(gl.TRIANGLE_FAN, int32(len(md.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
			} else {
				gl.DrawArrays(gl.TRIANGLES, 0, int32(len(md.Verts))/md.VertSize)
			}

			gl.BindVertexArray(0)
		}

		// Maintenance
		r.window.SwapBuffers()
		glfw.PollEvents()

	})
}

// LoadCamera sets the camera to be used by the display.
func (r *renderer) LoadCamera(camera components.Camera) {
	r.camera = camera
}

// addEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (r *renderer) addEntity(e entity.Entity) {
	mesh, isMesh := e.Component(components.TypeMesh).(components.Mesh)
	transform, isTransform := e.Component(components.TypeTransform).(components.Transform)

	if !isMesh || !isTransform {
		log.Printf("cannot load entity %s as a renderable", e.ID())
		return
	}

	rend := renderable{
		Mesh:      mesh,
		Transform: transform,
	}
	// Set up the shader
	md := mesh.Data()

	done := make(chan bool, 1)

	r.mainFunc(func() {
		shader, err := r.assets.Shaders().LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, false)

		if err != nil {
			log.Printf("Unable to load shaders %s,%s", md.VertShaderFile, md.FragShaderFile)
			done <- false
			return
		}

		rend.ProgramID = shader.ProgramID()
		rend.VAO = shader.CreateVAO(mesh)

		// Load and set the texture if it exists.
		if md.TextureFile != "" {
			texture, err := r.assets.Textures().LoadTexture(md.TextureFile, md.TextureFile)
			if err != nil {
				log.Printf("Unable to load texture:%s", md.TextureFile)
				done <- false
				return
			}

			rend.TextureID = texture
		}
		done <- true
	})
	// Add the shader and texture ids to the mesh component.
	if <-done {
		r.entities[e.ID()] = rend
	}
}

// removeEntity removes an Entity from the system.
func (r *renderer) removeEntity(e entity.Entity) {
	delete(r.entities, e.ID())
}

type renderable struct {
	Entity    entity.Entity
	Mesh      components.Mesh
	Transform components.Transform
	VAO       uint32
	ProgramID uint32
	TextureID uint32
}
