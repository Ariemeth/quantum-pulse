package main

import (
	"fmt"

	"github.com/Ariemeth/quantum-pulse/engine"
)

const (
	screenWidth  = 800
	screenHeight = 600
	windowTitle  = "hex map test"
)

func main() {

	e := engine.Engine{}

	err := e.Init(screenWidth, screenHeight, windowTitle)
	if err != nil {
		panic(err)
	}

	sceneID, err := e.LoadSceneFile("scene1.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	e.LoadScene(sceneID)

	e.Run()
}
