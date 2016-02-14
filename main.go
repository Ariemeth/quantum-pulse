package main

import (
	"runtime"

	"github.com/Ariemeth/frame-assault-2/engine"
	sm "github.com/Ariemeth/frame-assault-2/engine/shaderManager"
)

const vertShaderFile = "shaders/simple.vert"
const fragShaderFile = "shaders/simple.frag"
const shaderProgram = "simple"

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	engine := new(engine.Engine)

	engine.Init()
	shader := sm.Shader{VertSrcFile: vertShaderFile, FragSrcFile: fragShaderFile, Name: shaderProgram}

	engine.LoadShaders(shader, true)

	engine.Run()
}
