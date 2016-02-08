package shaderManager

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//ShaderManager stores shader programs
type ShaderManager struct {
	programs map[string]uint32
}

//NewShaderManager creates a new ShaderManager
func NewShaderManager() *ShaderManager{
	sm := ShaderManager{programs: make(map[string]uint32)}
	return &sm
}

//LoadProgram creates a shader program from a vertex and fragment shader source files.
func (sm *ShaderManager) LoadProgram(vertexSourceFile, fragmentSourceFile, key string) {
	simpleVert, err := loadShader(vertexSourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	simpleFrag, err := loadShader(fragmentSourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}
		
	program, err := newProgram(simpleVert, simpleFrag)
	if(err != nil){
		fmt.Println(err)
		return
	}
	
	sm.programs[key] = program	
}

//GetProgram returns a program id if the shader program was loaded, if it was not a 0 
//false will be returned
func (sm *ShaderManager) GetProgram(key string) (uint32, bool){
	program, status := sm.programs[key]
	return program, status
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

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
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
