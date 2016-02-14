package engine

import (
	"fmt"
	"log"

	sm "github.com/Ariemeth/frame-assault-2/engine/shaderManager"
	tm "github.com/Ariemeth/frame-assault-2/engine/textureManager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const windowWidth = 800
const windowHeight = 600

//Engine constitutes the rendering engine
type Engine struct {
	window   *glfw.Window
	shaders  sm.ShaderManager
	textures tm.TextureManager
}

//Init is called to initialize glfw and opengl
func (e *Engine) Init() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	e.window = createWindow()

	initGL()

	e.shaders = sm.NewShaderManager()
	e.textures = tm.NewTextureManager()
}

//Run is runs the main engine loop
func (e *Engine) Run() {
	defer glfw.Terminate()

	program, isLoaded := e.shaders.GetShader(e.shaders.GetDefaultShader())
	if isLoaded {
		gl.UseProgram(program)
	} else {
		fmt.Println("Unable to load default shader")
	}

	for !e.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//Update

		//Render

		// Maintenance
		e.window.SwapBuffers()
		glfw.PollEvents()
	}
}

//LoadShaders loads a vertex and fragment shader as a shader program
func (e *Engine) LoadShaders(shader sm.Shader, shouldBeDefault bool) {
	e.shaders.LoadProgram(shader, shouldBeDefault)
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
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

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
