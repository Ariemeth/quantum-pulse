package engine

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/Ariemeth/quantum-pulse/resources"
)

var (
	mainQueue = make(chan func(), 15)
	stopQueue = make(chan interface{})
)

// Engine constitutes the rendering engine which creates and initializes the rendering system.
type Engine struct {
	assets       *resources.Manager
	currentScene Scene
	scenes       map[string]Scene
	window       *glfw.Window
	windowWidth  int
	windowHeight int
}

func init() {
	go processMainFuncs()
}

// processMainFuncs process functions that have been added to the mainQueue on the main os thread.
func processMainFuncs() {
	runtime.LockOSThread()

	for {
		select {
		case f := <-mainQueue:
			f()
		case <-stopQueue:
			return
		}
	}
}

// runOnMain takes a funcion adds it to a queue that will ensure it is run on the main thread.
func runOnMain(f func()) {
	done := make(chan bool, 1)
	mainQueue <- func() {
		f()
		done <- true
	}
	<-done
}

// Init is called to initialize glfw and opengl
func (e *Engine) Init(width, height int, title string) error {

	initError := make(chan error, 1)

	runOnMain(func() {
		if err := glfw.Init(); err != nil {
			log.Fatalln("failed to initialize glfw:", err)
		}
		e.windowWidth = width
		e.windowHeight = height
		window, err := createWindow(width, height, title)
		if err != nil {
			initError <- err
		}
		e.window = window

		if _, err := initGL(); err != nil {
			initError <- err
		}

		e.assets = resources.NewManager()
		e.scenes = make(map[string]Scene)
		initError <- nil
	})
	return <-initError
}

// Run starts and runs the main engine loop
func (e *Engine) Run() {
	defer runOnMain(func() { glfw.Terminate() })
	defer e.currentScene.Terminate()

	e.currentScene.Start()

	for !e.window.ShouldClose() {
	}

}

// AddScene adds a scene to the engine.
func (e *Engine) AddScene(scene Scene, name string) {
	e.scenes[name] = scene
}

// LoadScene switches the current scene to the one specificed if it has already been added.
func (e *Engine) LoadScene(name string) {

	if scene, status := e.scenes[name]; status {
		if e.currentScene != nil {
			e.currentScene.Stop()
		}
		e.currentScene = scene
	}
}

// LoadSceneFile loads a scene from a file but does not make it the current scene. You
// must still call LoadScene with the scene id to load it as the current scene.  LoadSceneFile
// should not be called before Init is called.
func (e *Engine) LoadSceneFile(fileName string) (string, error) {

	scene, err := newScene(fileName, e.assets, e.window, e.windowWidth, e.windowHeight)
	if err != nil {
		return "", err
	}

	if scene == nil {
		return "", fmt.Errorf("Unable to load nil scene")
	}

	e.AddScene(scene, scene.ID())
	return scene.ID(), nil
}

func createWindow(width, height int, title string) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(onKey)
	window.SetMouseButtonCallback(onMouseButton)
	window.SetCloseCallback(onClose)
	window.SetScrollCallback(onScroll)
	window.SetCursorPosCallback(onCursorPos)

	return window, nil
}

func initGL() (string, error) {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		return "", err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	fmt.Println("OpenGl shading version", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	fmt.Println("OpenGl renderer", gl.GoStr(gl.GetString(gl.RENDERER)))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 0.5, 1.0)

	return version, nil
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
