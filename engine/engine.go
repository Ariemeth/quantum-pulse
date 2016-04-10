// Package engine represents the core functionality that allow the engine to be run and process the various systems and components.
package engine

import (
	"fmt"
	"log"

	sm "github.com/Ariemeth/quantum-pulse/engine/shaderManager"
	tm "github.com/Ariemeth/quantum-pulse/engine/textureManager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const windowWidth = 800
const windowHeight = 600

// Engine constitutes the rendering engine which creates and initializes the rendering system.
type Engine struct {
	window       *glfw.Window
	shaders      sm.ShaderManager
	textures     tm.TextureManager
	scenes       map[string]Scene
	currentScene Scene
}

// Init is called to initialize glfw and opengl
func (e *Engine) Init() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	e.window = createWindow()

	initGL()

	e.shaders = sm.NewShaderManager()
	e.textures = tm.NewTextureManager()
	e.scenes = make(map[string]Scene)

}

// Run is runs the main engine loop
func (e *Engine) Run() {
	defer glfw.Terminate()

	previousTime := glfw.GetTime()

	for !e.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		e.currentScene.Update(elapsed)

		// Render
		e.currentScene.Render()

		// Maintenance
		e.window.SwapBuffers()
		glfw.PollEvents()
	}
}

// AddScene adds a scene into the engine
func (e *Engine) AddScene(scene Scene, name string) {
	e.scenes[name] = scene
}

// LoadScene loads a scene
func (e *Engine) LoadScene(name string) {
	scene, status := e.scenes[name]
	if status {
		e.currentScene = scene
	}
}

// LoadSceneFile loads a scene from a file but does not make it the current scene. You
// must still call LoadScene with the scene id to load it as the current scene.  LoadSceneFile
// should not be called before Init is called.
func (e *Engine) LoadSceneFile(fileName string) string {
	scene := NewScene(fileName, e.shaders, e.textures)
	e.AddScene(scene, scene.ID())
	return scene.ID()
}

func createWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Frame Assault", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(onKey)
	window.SetMouseButtonCallback(onMouseButton)
	window.SetCloseCallback(onClose)
	window.SetScrollCallback(onScroll)
	window.SetCursorPosCallback(onCursorPos)

	return window
}

func initGL() string {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	fmt.Println("OpenGl shading version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	fmt.Println("OpenGl renderer", gl.GoStr(gl.GetString(gl.RENDERER)))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)

	return version
}

func onKey(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}

	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	default:
		return
	}
}

func onMouseButton(window *glfw.Window, b glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}

	switch glfw.MouseButton(b) {
	case glfw.MouseButtonLeft:
		return
	case glfw.MouseButtonRight:
		return
	default:
		return
	}
}

func onClose(window *glfw.Window) {
	window.SetShouldClose(true)
}

func onScroll(window *glfw.Window, xoff float64, yoff float64) {

}

func onCursorPos(window *glfw.Window, xpos float64, ypos float64) {

}
