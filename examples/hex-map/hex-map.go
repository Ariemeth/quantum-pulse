package main

import (
	"fmt"

	"github.com/Ariemeth/quantum-pulse/engine"
)

func main() {

	e := engine.Engine{}

	e.Init(800, 600,"hex map test")

	sceneID, err := e.LoadSceneFile("scene1.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	e.LoadScene(sceneID)

	e.Run()
}
