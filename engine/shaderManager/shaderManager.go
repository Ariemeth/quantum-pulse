package shaderManager

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/Ariemeth/quantum-pulse/engine/components"
)

const (
	// ShaderSrcDir lists the expected location of shaders
	ShaderSrcDir = "assets/shaders/"
	windowWidth  = 800
	windowHeight = 600
)

// shaderManager stores shader programs
type shaderManager struct {
	programs      map[string]components.Shader
	programLock   sync.RWMutex
	DefaultShader string
}

// ShaderManager interface is used to interact with the shaderManager.
type ShaderManager interface {
	// LoadProgramFromFile creates a shader program from a vertex and fragment shader source files.
	LoadProgramFromFile(vertSrcFile string, fragSrcFile string, name string, shouldBeDefault bool) (components.Shader, error)
	// LoadProgramFromSrc creates a shader program from a vertex and fragment shader source strings.
	LoadProgramFromSrc(vertSrc string, fragSrc string, name string, shouldBeDefault bool) (components.Shader, error)
	// GetShader returns a program id if the shader program was loaded.
	GetShader(key string) (components.Shader, bool)
	// GetDefaultShader returns the name of the default shader.
	GetDefaultShader() components.Shader
}

// NewShaderManager creates a new ShaderManager
func NewShaderManager() ShaderManager {
	sm := shaderManager{programs: make(map[string]components.Shader)}
	return &sm
}

// LoadProgramFromFile creates a shader program from a vertex and fragment shader source files.
func (sm *shaderManager) LoadProgramFromFile(vertSrcFile string, fragSrcFile string, name string, shouldBeDefault bool) (components.Shader, error) {

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
func (sm *shaderManager) LoadProgramFromSrc(vertSrc string, fragSrc string, name string, shouldBeDefault bool) (components.Shader, error) {
	if program, alreadyLoaded := sm.GetShader(name); alreadyLoaded {
		return program, nil
	}

	if !strings.HasSuffix(vertSrc, "\x00") {
		vertSrc = vertSrc + "\x00"
	}

	if !strings.HasSuffix(fragSrc, "\x00") {
		fragSrc = fragSrc + "\x00"
	}

	program, err := components.NewShader(name, vertSrc, fragSrc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sm.programLock.Lock()
	defer sm.programLock.Unlock()
	sm.programs[name] = program

	if len(sm.programs) == 1 || shouldBeDefault {
		sm.DefaultShader = name
	}
	return program, nil
}

// GetShader returns a components.Shader if the shader program was loaded, if it was not nil and
// false will be returned.
func (sm *shaderManager) GetShader(key string) (components.Shader, bool) {

	sm.programLock.RLock()
	defer sm.programLock.RUnlock()
	program, status := sm.programs[key]

	return program, status
}

// GetDefaultShader returns the name of the default shader.
func (sm *shaderManager) GetDefaultShader() components.Shader {
	s, _ := sm.GetShader(sm.DefaultShader)
	return s
}

func loadShaderFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data) + "\x00", nil
}
