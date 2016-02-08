package main

import (
	"runtime"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	engine := new(Engine)

	engine.Init()

	engine.Run()
}
