package components

import (
	"fmt"
	"strings"

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
	id         uint32
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
}

// NewShader creates a new shader program and populates the uniform and attribute layouts.
func NewShader(name string, vertSrc, fragSrc string) (Shader, error) {
	s := shader{
		uniforms:   make(map[string]int32),
		attributes: make(map[string]uint32),
		name:       name,
	}

	program, err := newProgram(vertSrc, fragSrc)

	if err != nil {
		return nil, err
	}

	s.id = program

	s.populateLocations()

	return &s, nil
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
	return s.id
}

func (s *shader) populateLocations() {
	program := s.ProgramID()

	s.uniforms[ProjectionUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ProjectionUniform)))
	s.uniforms[CameraUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", CameraUniform)))
	s.uniforms[ModelUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ModelUniform)))
	s.uniforms[TextureUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", TextureUniform)))

	s.attributes[VertexAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexAttribute))))
	s.attributes[VertexTexCordAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexTexCordAttribute))))

	gl.BindFragDataLocation(program, 0, gl.Str(fmt.Sprintf("%s\x00", ShaderOutputColor)))
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
