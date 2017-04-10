package systems

import (
	"fmt"

	am "github.com/Ariemeth/quantum-pulse/assets"
	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
	"github.com/go-gl/gl/v4.1-core/gl"
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
}

type renderer struct {
	entities map[string]renderable
	assets   *am.AssetManager
	camera   components.Camera
	mainFunc func(f func())
}

// NewRenderer creates a new rendererB system.  The rendererB system handles rendering all renderable Entities to the screen.
func NewRenderer(assetManager *am.AssetManager, mainFunc func(f func())) Renderer {
	r := renderer{
		entities: make(map[string]renderable),
		assets:   assetManager,
		camera:   components.NewCamera(),
		mainFunc: mainFunc,
	}

	return &r
}

// SystemType retrieves the type of system such as renderer, mover, etc.
func (r *renderer) SystemType() string {
	return SystemTypeRenderer
}

// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (r *renderer) AddEntity(e entity.Entity) {
	mesh, isMesh := e.Component(components.ComponentTypeMesh).(components.Mesh)
	transform, isTransform := e.Component(components.ComponentTypeTransform).(components.Transform)

	if isMesh && isTransform {
		rend := renderable{
			Mesh:      mesh,
			Transform: transform,
		}
		// Set up the shader
		md := mesh.Data()
		shader, err := r.assets.Shaders().LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, false)
		if err != nil {
			fmt.Printf("Unable to load shaders %s,%s", md.VertShaderFile, md.FragShaderFile)
			return
		}

		md.ProgramID = shader.ProgramID()
		md.VAO = shader.CreateVAO(mesh)

		// Load and set the texture if it exists.
		texture, err := r.assets.Textures().LoadTexture(md.TextureFile, md.TextureFile)
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

// RemoveEntity removes an Entity from the system.
func (r *renderer) RemoveEntity(e entity.Entity) {
	delete(r.entities, e.ID())
}

func (r *renderer) Process() {

	r.mainFunc(func() {
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
