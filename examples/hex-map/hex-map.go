package main

import (
	"github.com/Ariemeth/quantum-pulse/engine"
)

const (
	screenWidth  = 800
	screenHeight = 600
	windowTitle  = "hex map test"
)

func main() {

	// Create the engine
	e := engine.Engine{}

	// Initialize the engine.  This will create the window.
	err := e.Init(screenWidth, screenHeight, windowTitle)
	if err != nil {
		panic(err)
	}

	// Now load a json scene file.
	sceneID, err := e.LoadSceneFile("scene1.json")
	if err != nil {
		panic(err)
	}

	// Make the scene file that was just loaded the active scene.
	e.LoadScene(sceneID)

	// Start the main game loop.
	e.Run()
}
