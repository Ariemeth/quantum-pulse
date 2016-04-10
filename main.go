package main

import (
	"runtime"

	"github.com/Ariemeth/quantum-pulse/engine"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	engine := new(engine.Engine)

	engine.Init()

	engine.LoadScene(engine.LoadSceneFile("scene1.json"))

	engine.Run()
}
