package systems

import (
	"github.com/Ariemeth/quantum-pulse/engine/components"
	"github.com/Ariemeth/quantum-pulse/engine/entity"
	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	SystemTypeRenderer = "renderer"
)

// Renderer provides the interface needed to process the rendering of Entities.  Each time Process is called all Entities will be rendered.
type Renderer interface {
	System
	Process()
}

type renderer struct {
	entities map[string]renderable
}

// NewRenderer creates a new rendererB system.  The rendererB system handles rendering all renderable Entities to the screen.
func NewRenderer() Renderer {
	r := renderer{
		entities: make(map[string]renderable),
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
	shader, isShader := e.Component(components.ComponentTypeShader).(components.Shader)
	transform, isTransform := e.Component(components.ComponentTypeTransform).(components.Transform)

	if isMesh && isShader && isTransform {
		rend := renderable{
			Mesh:      mesh,
			Shader:    shader,
			Transform: transform,
		}
		r.entities[e.ID()] = rend
	}

}

// RemoveEntity removes an Entity from the system.
func (r *renderer) RemoveEntity(e entity.Entity) {
	delete(r.entities, e.ID())
}

func (r *renderer) Process() {

	for _, ent := range r.entities {
		gl.UseProgram(ent.Shader.ProgramID())

		gl.BindVertexArray(ent.Shader.GetVAO())
		ent.Shader.LoadTransform(ent.Transform)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.Uniform1i(ent.Shader.GetUniformLoc(components.TextureUniform), 0)

		texture, isLoaded := ent.Shader.GetTexture()
		if isLoaded {
			gl.BindTexture(gl.TEXTURE_2D, texture)
		}

		md := ent.Mesh.Data()
		if md.Indexed {
			gl.DrawElements(gl.TRIANGLE_FAN, int32(len(md.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
		} else {
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(md.Verts))/md.VertSize)
		}

		gl.BindVertexArray(0)
	}
}

type renderable struct {
	Entity    entity.Entity
	Mesh      components.Mesh
	Shader    components.Shader
	Transform components.Transform
}
