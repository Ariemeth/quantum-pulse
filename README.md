# Quantum Pulse game engine

[![Build Status](https://travis-ci.org/Ariemeth/quantum-pulse.svg?branch=master)](https://travis-ci.org/Ariemeth/quantum-pulse)
[![Go Report Card](https://goreportcard.com/badge/github.com/ariemeth/quantum-pulse)](https://goreportcard.com/report/github.com/ariemeth/quantum-pulse)
[![GoDoc](https://godoc.org/github.com/Ariemeth/quantum-pulse?status.svg)](https://godoc.org/github.com/Ariemeth/quantum-pulse)

## What is Quantum Pulse

Quantum Pulse is a learning project looking into building an entity component system based game engine using go

## Requirements

- Go 1.8+
- OpenGL 4.1+
- [go-gl/glfw](https://github.com/go-gl/glfw)

## Getting started

The hex map example

```go
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
```
