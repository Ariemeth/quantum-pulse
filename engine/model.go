package engine

import (
	"fmt"

	"github.com/Ariemeth/quantum-pulse/engine/components"
	sm "github.com/Ariemeth/quantum-pulse/engine/shaderManager"
	tm "github.com/Ariemeth/quantum-pulse/engine/textureManager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// ModelSrcDir is the expected location of models
	ModelSrcDir = "assets/models/"
)

// Model represents a physical entity
type Model struct {
	id         string
	meshComp   components.Mesh
	transform  components.Transform
	shader     components.Shader
	shaders    sm.ShaderManager
	textures   tm.TextureManager
	angle      float64
	camera     mgl32.Mat4
	projection mgl32.Mat4
}

// NewModel creates a new model.
func NewModel(id string, shaders sm.ShaderManager, textures tm.TextureManager) *Model {
	m := Model{
		id:         id,
		meshComp:   components.NewMesh(),
		transform:  components.NewTransform(),
		shaders:    shaders,
		textures:   textures,
		angle:      0.0,
		camera:     mgl32.Ident4(),
		projection: mgl32.Ident4(),
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

	md := m.meshComp.Data()
	texture, isLoaded := m.textures.GetTexture(md.TextureFile)

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

// Load loads and sets up the model
func (m *Model) Load(fileName string) {

	m.meshComp.Load(fileName)
	md := m.meshComp.Data()

	shaderName := fmt.Sprintf("%s:%s", md.VertShaderFile, md.FragShaderFile)
	program, err := m.shaders.LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, shaderName, false)
	if err != nil {
		return
	}
	m.shader = components.NewShader(shaderName, program)

	// Load the texture
	m.textures.LoadTexture(md.TextureFile, md.TextureFile)

	// Load the uniforms
	gl.UseProgram(m.shader.ProgramID())

	// Load the perspective matrix into a shader uniform
	// TODO this should be pulled out of the model Load method
	m.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	gl.UniformMatrix4fv(m.shader.GetUniformLoc(components.ProjectionUniform), 1, false, &m.projection[0])

	// Load the camera matrix into a shader
	// TODO this should be pulled out of the model Load method
	m.camera = mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(m.shader.GetUniformLoc(components.CameraUniform), 1, false, &m.camera[0])

	m.shader.LoadTransform(m.transform)
	m.shader.LoadMesh(m.meshComp)
	gl.Uniform1i(m.shader.GetUniformLoc(components.TextureUniform), 0)
}
