package shaderManager

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

const (
	// ShaderSrcDir lists the expected location of shaders
	ShaderSrcDir = "assets/shaders/"
	windowWidth  = 800
	windowHeight = 600
)

// shaderManager stores shader programs
type shaderManager struct {
	programs      map[string]ShaderProgram
	programLock   sync.RWMutex
	DefaultShader string
}

// ShaderManager interface is used to interact with the shaderManager.
type ShaderManager interface {
	// LoadProgramFromFile creates a shader program from a vertex and fragment shader source files.
	LoadProgramFromFile(vertSrcFile string, fragSrcFile string, name string, shouldBeDefault bool) (ShaderProgram, error)
	// LoadProgramFromSrc creates a shader program from a vertex and fragment shader source strings.
	LoadProgramFromSrc(vertSrc string, fragSrc string, name string, shouldBeDefault bool) (ShaderProgram, error)
	// GetShader returns a program id if the shader program was loaded.
	GetShader(key string) (ShaderProgram, bool)
	// GetDefaultShader returns the name of the default shader.
	GetDefaultShader() string
}

// NewShaderManager creates a new ShaderManager
func NewShaderManager() ShaderManager {
	sm := shaderManager{programs: make(map[string]ShaderProgram)}
	return &sm
}

// LoadProgramFromFile creates a shader program from a vertex and fragment shader source files.
func (sm *shaderManager) LoadProgramFromFile(vertSrcFile string, fragSrcFile string, name string, shouldBeDefault bool) (ShaderProgram, error) {

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
func (sm *shaderManager) LoadProgramFromSrc(vertSrc string, fragSrc string, name string, shouldBeDefault bool) (ShaderProgram, error) {
	if program, alreadyLoaded := sm.GetShader(name); alreadyLoaded {
		return program, nil
	}

	if !strings.HasSuffix(vertSrc, "\x00") {
		vertSrc = vertSrc + "\x00"
	}

	if !strings.HasSuffix(fragSrc, "\x00") {
		fragSrc = fragSrc + "\x00"
	}

	//		program, err := newProgram(vertSrc, fragSrc)

	program, err := NewShaderProgram(name, vertSrc, fragSrc)
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

// GetShader returns a program id if the shader program was loaded, if it was not a 0
// false will be returned.
func (sm *shaderManager) GetShader(key string) (ShaderProgram, bool) {

	sm.programLock.RLock()
	defer sm.programLock.RUnlock()
	program, status := sm.programs[key]

	return program, status
}

// GetDefaultShader returns the name of the default shader.
func (sm *shaderManager) GetDefaultShader() string {
	return sm.DefaultShader
}

func loadShaderFile(fileName string) (string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(data) + "\x00", nil
}
