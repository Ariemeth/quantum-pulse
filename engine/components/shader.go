package components

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	// ComponentTypeShader represents a shaders components type.
	ComponentTypeShader = "shader"
	// CameraUniform is the expected name of view matrix uniform in the shader.
	CameraUniform = "camera"
	// ProjectionUniform is the expected name of the projection matrix uniform in the shader.
	ProjectionUniform = "projection"
	// ModelUniform is the expected name of the model matrix uniform in the shader.
	ModelUniform = "model"
	// TextureUniform is the expected name of the texture uniform in the shader.
	TextureUniform = "tex"
	// VertexAttribute is the expected name of the vertex data attribute in the shader.
	VertexAttribute = "vert"
	// VertexTexCordAttribute is the expected name of the vertex texture coordinates attribute in the shader.
	VertexTexCordAttribute = "vertTexCoord"
	// ShaderOutputColor is the expected name of the output color variable leaving the fragment shader.
	ShaderOutputColor = "outputColor"
)

// shader holds information about a shader program
type shader struct {
	uniforms   map[string]int32
	attributes map[string]uint32
	texture    uint32
	name       string //shader name
	program    uint32
	vao        uint32
}

// Shader represents the behaviors needed to access a shader and its variables.
type Shader interface {
	Component
	// GetUniformLoc retrieves the shader location of the specified uniform.
	GetUniformLoc(name string) int32
	// GetAttribLoc retrieves the shader location of a the specified attribute.
	GetAttribLoc(name string) uint32
	// GetName retrieves the name of the shader program.
	GetName() string
	// GetVAO retrieves the vao id being used by this shader.
	GetVAO() uint32
	// ProgramID retrieves the program id of the shader program.
	ProgramID() uint32
	// LoadTransform loads the transform matrix onto the gpu.
	LoadTransform(Transform)
	// LoadMesh loads the mesh data onto the gpu.
	LoadMesh(Mesh)
	// AddTexture adds the texture id to use with this shader.
	AddTexture(uint32)
	// GetTexture retrieves the texture id to use with this shader.
	GetTexture() (uint32, bool)
}

// NewShader creates a new shader program and populates the uniform and attribute layouts.
func NewShader(name string, shaderProgram uint32) Shader {
	s := shader{
		uniforms:   make(map[string]int32),
		attributes: make(map[string]uint32),
		name:       name,
		program:    shaderProgram,
		texture:    999999,
	}

	s.storeLocations()

	return &s
}

// ComponentType is expected to return a string representing the type of component.
func (s *shader) ComponentType() string {
	return ComponentTypeShader
}

// GetUniformLoc retrieves the shader location of the specified uniform.
func (s *shader) GetUniformLoc(name string) int32 {
	u, ok := s.uniforms[name]
	if !ok {
		return -1
	}
	return u
}

// GetAttribLoc retrieves the shader location of a the specified attribute.
func (s *shader) GetAttribLoc(name string) uint32 {
	a, ok := s.attributes[name]
	if !ok {
		return 0
	}
	return a
}

// GetName retrieves the name of the shader program.
func (s *shader) GetName() string {
	return s.name
}

// GetVAO retrieves the vao id being used by this shader.
func (s *shader) GetVAO() uint32 {
	return s.vao
}

// ProgramID retrieves the program id of the shader program.
func (s *shader) ProgramID() uint32 {
	return s.program
}

// LoadTransform loads the transform matrix onto the gpu.
func (s *shader) LoadTransform(t Transform) {

	td := t.Data()
	gl.UseProgram(s.program)
	gl.UniformMatrix4fv(s.GetUniformLoc(ModelUniform), 1, false, &td[0])
}

// LoadMesh loads the mesh data onto the gpu.  This will create a new VAO and should
// only be called once unless you need to reset the shader.
func (s *shader) LoadMesh(m Mesh) {

	md := m.Data()
	gl.UseProgram(s.program)

	// Configure vertex array object with the model's data
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(md.Verts)*4, gl.Ptr(md.Verts), gl.STATIC_DRAW)

	vertAttrib := s.GetAttribLoc(VertexAttribute)
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, md.VertSize*4, gl.PtrOffset(0)) // 4:number of bytes in a float32

	texCoordAttrib := s.GetAttribLoc(VertexTexCordAttribute)
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

func (s *shader) AddTexture(texture uint32) {
	s.texture = texture
}

func (s *shader) GetTexture() (uint32, bool) {
	if s.texture >= 999999 {
		return 0, false
	}
	return s.texture, true
}

func (s *shader) storeLocations() {
	program := s.ProgramID()
	gl.UseProgram(program)

	s.uniforms[ProjectionUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ProjectionUniform)))
	s.uniforms[CameraUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", CameraUniform)))
	s.uniforms[ModelUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ModelUniform)))
	s.uniforms[TextureUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", TextureUniform)))

	s.attributes[VertexAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexAttribute))))
	s.attributes[VertexTexCordAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexTexCordAttribute))))

	gl.BindFragDataLocation(program, 0, gl.Str(fmt.Sprintf("%s\x00", ShaderOutputColor)))
}
