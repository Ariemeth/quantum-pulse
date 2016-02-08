package main

import (
	"runtime"
	"github.com/Ariemeth/frame-assault-2/engine"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	engine := new(engine.Engine)

	engine.Init()

	engine.Run()
}
