package main

import (
	"fmt"

	"github.com/Ariemeth/quantum-pulse/engine"
)

func main() {

	e := new(engine.Engine)

	e.Init()

	sceneID, err := e.LoadSceneFile("scene1.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	e.LoadScene(sceneID)

	e.Run()
}
