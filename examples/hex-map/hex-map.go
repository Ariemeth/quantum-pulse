package main

import (
	"fmt"
	qp "github.com/Ariemeth/quantum-pulse"
)

func main() {

	engine := new(qp.Engine)

	engine.Init()

	sceneID, err := engine.LoadSceneFile("scene1.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	engine.LoadScene(sceneID)

	engine.Run()
}
