package shaderManager

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	// ShaderSrcDir lists the expected location of shaders
	ShaderSrcDir = "assets/shaders/"
)

// Shader holds information about a shader program
type Shader struct {
	VertSrcFile string //file that contains the vertex shader
	FragSrcFile string //file that contains the fragment shader
	VertSrc     string //vertex shader source
	FragSrc     string //vertex shader source
	Name        string //shader name
}

// shaderManager stores shader programs
type shaderManager struct {
	programs      map[string]uint32
	programLock   sync.RWMutex
	DefaultShader string
}

// ShaderManager interface is used to interact with the shaderManager.
type ShaderManager interface {
	// LoadProgram creates a shader program from a vertex and fragment shader source files.
	LoadProgram(shader Shader, shouldBeDefault bool) (uint32, error)
	// GetShader returns a program id if the shader program was loaded.
	GetShader(key string) (uint32, bool)
	// GetDefaultShader returns the name of the default shader.
	GetDefaultShader() string
}

// NewShaderManager creates a new ShaderManager
func NewShaderManager() ShaderManager {
	sm := shaderManager{programs: make(map[string]uint32)}
	return &sm
}

// LoadProgram creates a shader program from a vertex and fragment shader source files.
func (sm *shaderManager) LoadProgram(shader Shader, shouldBeDefault bool) (uint32, error) {

	if program, alreadyLoaded := sm.GetShader(shader.Name); alreadyLoaded {
		return program, nil
	}

	if len(shader.VertSrc) < 1 {
		simpleVert, err := loadShader(fmt.Sprintf("%s%s", ShaderSrcDir, shader.VertSrcFile))
		shader.VertSrc = simpleVert
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	} else if !strings.HasSuffix(shader.VertSrc, "\x00") {
		shader.VertSrc = shader.VertSrc + "\x00"
	}

	if len(shader.FragSrc) < 1 {
		simpleFrag, err := loadShader(fmt.Sprintf("%s%s", ShaderSrcDir, shader.FragSrcFile))
		shader.FragSrc = simpleFrag
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	} else if !strings.HasSuffix(shader.FragSrc, "\x00") {
		shader.FragSrc = shader.FragSrc + "\x00"
	}

	program, err := newProgram(shader.VertSrc, shader.FragSrc)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	sm.programLock.Lock()
	defer sm.programLock.Unlock()
	sm.programs[shader.Name] = program

	if len(sm.programs) == 1 || shouldBeDefault {
		sm.DefaultShader = shader.Name
	}
	return program, nil
}

// GetShader returns a program id if the shader program was loaded, if it was not a 0
// false will be returned.
func (sm *shaderManager) GetShader(key string) (uint32, bool) {

	sm.programLock.RLock()
	defer sm.programLock.RUnlock()
	program, status := sm.programs[key]

	return program, status
}

// GetDefaultShader returns the name of the default shader.
func (sm *shaderManager) GetDefaultShader() string {
	return sm.DefaultShader
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

func loadShader(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data) + "\x00", nil
}
