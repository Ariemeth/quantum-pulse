package engine

import (
	"github.com/Ariemeth/quantum-pulse/engine/components"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// ModelSrcDir is the expected location of models
	ModelSrcDir = "assets/models/"
)

// Model represents a physical entity
type Model struct {
	id        string
	meshComp  components.Mesh
	transform components.Transform
	shader    components.Shader
	angle     float64
}

func (m *Model) AddMesh(mesh components.Mesh) {
	m.meshComp = mesh
}
func (m *Model) AddShader(shader components.Shader) {
	m.shader = shader
}

func (m *Model) AddTransform(transform components.Transform) {
	m.transform = transform
}

// NewModel creates a new model.
func NewModel(id string) *Model {
	m := Model{
		id:    id,
		angle: 0.0,
	}

	return &m
}

// ID returns the id of the model.
func (m Model) ID() string {
	return m.id
}

// Update updates the model
func (m *Model) Update(elapsed float64) {
	m.angle += elapsed
	m.transform.Set(mgl32.HomogRotate3D(float32(m.angle), mgl32.Vec3{0, 1, 0}))
}

// Render renders the model
func (m *Model) Render() {

	gl.UseProgram(m.shader.ProgramID())
	td := m.transform.Data()
	gl.UniformMatrix4fv(m.shader.GetUniformLoc(components.ModelUniform), 1, false, &td[0])

	gl.BindVertexArray(m.shader.GetVAO())

	gl.ActiveTexture(gl.TEXTURE0)
	gl.Uniform1i(m.shader.GetUniformLoc(components.TextureUniform), 0)

	md := m.meshComp.Data()
	texture, isLoaded := m.shader.GetTexture()

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
