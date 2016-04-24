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
	name       string //shader name
	program    uint32
	vao        uint32
}

// Shader represents the behaviors needed to access a shader and its variables.
type Shader interface {
	// GetUniformLoc retrieves the shader location of the specified uniform.
	GetUniformLoc(name string) int32
	// GetAttribLoc retrieves the shader location of a the specified attribute.
	GetAttribLoc(name string) uint32
	// GetName retrieves the name of the shader program.
	GetName() string
	// ProgramID retrieves the program id of the shader program.
	ProgramID() uint32
	// LoadData loads the mesh data onto the gpu.
	LoadData(Mesh)
}

// NewShader creates a new shader program and populates the uniform and attribute layouts.
func NewShader(name string, shaderProgram uint32) Shader {
	s := shader{
		uniforms:   make(map[string]int32),
		attributes: make(map[string]uint32),
		name:       name,
		program:    shaderProgram,
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

// ProgramID retrieves the program id of the shader program.
func (s *shader) ProgramID() uint32 {
	return s.program
}

// LoadData loads the mesh data onto the gpu.
func (s *shader) LoadData(m Mesh) {

}

func (s *shader) storeLocations() {
	program := s.ProgramID()

	s.uniforms[ProjectionUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ProjectionUniform)))
	s.uniforms[CameraUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", CameraUniform)))
	s.uniforms[ModelUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ModelUniform)))
	s.uniforms[TextureUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", TextureUniform)))

	s.attributes[VertexAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexAttribute))))
	s.attributes[VertexTexCordAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexTexCordAttribute))))

	gl.BindFragDataLocation(program, 0, gl.Str(fmt.Sprintf("%s\x00", ShaderOutputColor)))
}
