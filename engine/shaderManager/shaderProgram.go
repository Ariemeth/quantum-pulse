package shaderManager

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	CameraUniform          = "camera"
	ProjectionUniform      = "projection"
	ModelUniform           = "model"
	TextureUniform         = "tex"
	VertexAttribute        = "vert"
	VertexTexCordAttribute = "vertTexCoord"
)

// shaderProgram holds information about a shader program
type shaderProgram struct {
	uniforms   map[string]int32
	attributes map[string]uint32
	name       string //shader name
	id         uint32
}

// ShaderProgram represents the behaviors needed to access a shader and its variables.
type ShaderProgram interface {
	GetUniformLoc(name string) int32
	GetAttribLoc(name string) uint32
	GetName() string
	ProgramID() uint32
}

// NewShaderProgram creates a new shader program and populates the uniform and attribute layouts.
func NewShaderProgram(name string, vertSrc, fragSrc string) (ShaderProgram, error) {
	sp := shaderProgram{
		uniforms:   make(map[string]int32),
		attributes: make(map[string]uint32),
		name:       name,
	}

	program, err := newProgram(vertSrc, fragSrc)

	if err != nil {
		return nil, err
	}

	sp.id = program

	sp.populateLocations()

	return &sp, nil
}

func (sp *shaderProgram) GetUniformLoc(name string) int32 {
	u, ok := sp.uniforms[name]
	if !ok {
		return -1
	}
	return u
}

func (sp *shaderProgram) GetAttribLoc(name string) uint32 {
	a, ok := sp.attributes[name]
	if !ok {
		return 0
	}
	return a
}

func (sp *shaderProgram) GetName() string {
	return sp.name
}
func (sp *shaderProgram) ProgramID() uint32 {
	return sp.id
}

func (sp *shaderProgram) populateLocations() {
	program := sp.ProgramID()

	sp.uniforms[ProjectionUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ProjectionUniform)))
	sp.uniforms[CameraUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", CameraUniform)))
	sp.uniforms[ModelUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", ModelUniform)))
	sp.uniforms[TextureUniform] = gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", TextureUniform)))

	sp.attributes[VertexAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexAttribute))))
	sp.attributes[VertexTexCordAttribute] = uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", VertexTexCordAttribute))))

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
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
