package systems

import (
	"fmt"
	"sync"

	"time"

	am "github.com/Ariemeth/quantum-pulse/assets"
	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	// SystemTypeRenderer is the name of the renderer system.
	SystemTypeRenderer = "renderer"
)

// Renderer provides the interface needed to process the rendering of Entities.  Each time Process is called all Entities will be rendered.
type Renderer interface {
	System
	// Process renders all renderable entities.
	Process()
	// LoadCamera sets the camera to be used by the display.
	LoadCamera(camera components.Camera)
	Stop()
}

type renderer struct {
	entities     map[string]renderable
	assets       *am.AssetManager
	camera       components.Camera
	mainFunc     func(f func())
	window       *glfw.Window
	remove       chan entity.Entity
	add          chan entity.Entity
	quit         chan interface{}
	runningLock  sync.Mutex
	isRunning    bool
	requirements []string
}

// NewRenderer creates a new renderer system.  The renderer system handles rendering all renderable Entities to the screen.
func NewRenderer(assetManager *am.AssetManager, mainFunc func(f func()), window *glfw.Window) Renderer {
	r := renderer{
		entities:     make(map[string]renderable),
		assets:       assetManager,
		camera:       components.NewCamera(),
		mainFunc:     mainFunc,
		window:       window,
		remove:       make(chan entity.Entity, 0),
		add:          make(chan entity.Entity, 0),
		quit:         make(chan interface{}),
		requirements: []string{""}, //TODO eventually add the actual requirements
	}

	r.Start()

	return &r
}

// SystemType retrieves the type of system such as renderer, mover, etc.
func (r *renderer) SystemType() string {
	return SystemTypeRenderer
}

// AddEntity adds an entity to the renderer to be capable of being rendered to the screen.
func (r *renderer) AddEntity(e entity.Entity) {
	r.add <- e
}

// RemoveEntity removes an entity from the renderer.  Once removed the Entity will no longer be rendered to the screen unless it is added back to the renderer.
func (r *renderer) RemoveEntity(e entity.Entity) {
	r.remove <- e
}

// addEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (r *renderer) addEntity(e entity.Entity) {
	mesh, isMesh := e.Component(components.ComponentTypeMesh).(components.Mesh)
	transform, isTransform := e.Component(components.ComponentTypeTransform).(components.Transform)

	if isMesh && isTransform {
		rend := renderable{
			Mesh:      mesh,
			Transform: transform,
		}
		// Set up the shader
		md := mesh.Data()

		var shader am.Shader
		var err error
		r.mainFunc(func() {
			shader, err = r.assets.Shaders().LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, false)
		})
		if err != nil {
			fmt.Printf("Unable to load shaders %s,%s", md.VertShaderFile, md.FragShaderFile)
			return
		}

		r.mainFunc(func() {
			md.ProgramID = shader.ProgramID()
			md.VAO = shader.CreateVAO(mesh)
		})

		// Load and set the texture if it exists.
		var texture uint32
		r.mainFunc(func() {
			texture, err = r.assets.Textures().LoadTexture(md.TextureFile, md.TextureFile)
		})
		if err != nil {
			fmt.Printf("Unable to load texture:%s", md.TextureFile)
			return
		}

		md.TextureID = texture
		// Add the shader and texture ids to the mesh component.
		mesh.Set(md)

		r.entities[e.ID()] = rend
	}

}

// removeEntity removes an Entity from the system.
func (r *renderer) removeEntity(e entity.Entity) {
	delete(r.entities, e.ID())
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
		for {
			select {
			case ent := <-r.add:
				r.addEntity(ent)
				fmt.Printf("Adding %s to the renderer.\n", ent.ID())
			case ent := <-r.remove:
				r.removeEntity(ent)
			case <-r.quit:
				return
			case <-time.After(15 * time.Millisecond):
				r.Process()
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

func (r *renderer) Process() {

	r.mainFunc(func() {
		r.runningLock.Lock()
		defer r.runningLock.Unlock()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, ent := range r.entities {
			md := ent.Mesh.Data()
			shader, status := r.assets.Shaders().GetShaderProgram(md.ProgramID)
			if !status {
				// if the shader is not loaded there is no point in trying to render it to the screen.
				continue
			}
			gl.UseProgram(md.ProgramID)

			gl.BindVertexArray(md.VAO)
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

type renderable struct {
	Entity    entity.Entity
	Mesh      components.Mesh
	Transform components.Transform
}
