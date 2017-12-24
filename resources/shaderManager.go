package resources

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
	windowWidth  = 800
	windowHeight = 600
)

// shaderManager stores shader programs
type shaderManager struct {
	shaders       map[string]Shader
	programLock   sync.RWMutex
	DefaultShader string
}

// ShaderManager interface is used to interact with the shaderManager.
type ShaderManager interface {
	// LoadProgramFromFile creates a shader program from a vertex and fragment shader source files.
	LoadProgramFromFile(vertSrcFile string, fragSrcFile string, shouldBeDefault bool) (Shader, error)
	// LoadProgramFromSrc creates a shader program from a vertex and fragment shader source strings.
	LoadProgramFromSrc(vertSrc string, fragSrc string, name string, shouldBeDefault bool) (Shader, error)
	// GetShader returns a program id if the shader program was loaded.
	GetShader(key string) (Shader, bool)
	// GetShaderProgram returns the Shader interface with information about the opengl shader.
	GetShaderProgram(id uint32) (Shader, bool)
	// GetDefaultShader returns the name of the default shader.
	GetDefaultShader() uint32
}

// newShaderManager creates a new ShaderManager
func newShaderManager() ShaderManager {
	sm := shaderManager{
		shaders: make(map[string]Shader),
	}
	return &sm
}

// LoadProgramFromFile creates a shader program from a vertex and fragment shader source files.
func (sm *shaderManager) LoadProgramFromFile(vertSrcFile string, fragSrcFile string, shouldBeDefault bool) (Shader, error) {
	name := fmt.Sprintf("%s:%s", vertSrcFile, fragSrcFile)
	simpleVert, err := loadShaderFile(fmt.Sprintf("%s%s", ShaderSrcDir, vertSrcFile))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	simpleFrag, err := loadShaderFile(fmt.Sprintf("%s%s", ShaderSrcDir, fragSrcFile))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return sm.LoadProgramFromSrc(simpleVert, simpleFrag, name, shouldBeDefault)
}

// LoadProgramFromSrc creates a shader program from a vertex and fragment shader source strings.
func (sm *shaderManager) LoadProgramFromSrc(vertSrc string, fragSrc string, name string, shouldBeDefault bool) (Shader, error) {
	if program, alreadyLoaded := sm.GetShader(name); alreadyLoaded {
		return program, nil
	}

	if !strings.HasSuffix(vertSrc, "\x00") {
		vertSrc = vertSrc + "\x00"
	}

	if !strings.HasSuffix(fragSrc, "\x00") {
		fragSrc = fragSrc + "\x00"
	}

	program, err := newProgram(vertSrc, fragSrc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sm.programLock.Lock()
	defer sm.programLock.Unlock()

	shader := newShader(name, program)
	sm.shaders[name] = shader

	if len(sm.shaders) == 1 || shouldBeDefault {
		sm.DefaultShader = name
	}

	return shader, nil
}

// GetShader returns a uint32 if the shader program was loaded, if it was not nil and
// false will be returned.
func (sm *shaderManager) GetShader(key string) (Shader, bool) {

	sm.programLock.RLock()
	defer sm.programLock.RUnlock()
	program, status := sm.shaders[key]

	return program, status
}

// GetShaderProgram returns the Shader interface with information about the opengl shader.
func (sm *shaderManager) GetShaderProgram(id uint32) (Shader, bool) {
	sm.programLock.RLock()
	defer sm.programLock.RUnlock()
	for _, shader := range sm.shaders {
		if shader.ProgramID() == id {
			return shader, true
		}
	}

	return nil, false
}

// GetDefaultShader returns the name of the default shader.
func (sm *shaderManager) GetDefaultShader() uint32 {
	s, _ := sm.GetShader(sm.DefaultShader)
	return s.ProgramID()
}

func loadShaderFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data) + "\x00", nil
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
