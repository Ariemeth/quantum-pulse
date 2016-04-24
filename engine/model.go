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
	vao        uint32
	//currentProgram uint32
	//	program        components.Shader
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

	gl.BindVertexArray(m.vao)

	gl.ActiveTexture(gl.TEXTURE0)

	md := m.meshComp.Data()
	texture, isLoaded := m.textures.GetTexture(md.TextureFile)

	if isLoaded {
		gl.BindTexture(gl.TEXTURE_2D, texture)
	} else {
		if md.TextureFile != "" {
			go fmt.Printf("Unable to load texture %s", md.TextureFile)
		}
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

	shader, err := m.shaders.LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, fmt.Sprintf("%s:%s", md.VertShaderFile, md.FragShaderFile), false)
	if err != nil {
		return
	}
	m.shader = shader

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

	// Load the model vertex data.
	td := m.transform.Data()
	gl.UniformMatrix4fv(m.shader.GetUniformLoc(components.ModelUniform), 1, false, &td[0])

	gl.Uniform1i(m.shader.GetUniformLoc(components.TextureUniform), 0)

	// Configure vertex array object with the model's data
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(md.Verts)*4, gl.Ptr(md.Verts), gl.STATIC_DRAW)

	vertAttrib := m.shader.GetAttribLoc(components.VertexAttribute)
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, md.VertSize*4, gl.PtrOffset(0)) // 4:number of bytes in a float32

	texCoordAttrib := m.shader.GetAttribLoc(components.VertexTexCordAttribute)
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, true, md.VertSize*4, gl.PtrOffset(3*4)) // 4:number of bytes in a float32

	if md.Indexed {
		var indices uint32
		gl.GenBuffers(1, &indices)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indices)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(md.Indices)*4, gl.Ptr(md.Indices), gl.STATIC_DRAW)
	}

	gl.BindVertexArray(0)
}
