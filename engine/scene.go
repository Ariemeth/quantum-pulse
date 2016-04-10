package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sm "github.com/Ariemeth/quantum-pulse/engine/shaderManager"
	tm "github.com/Ariemeth/quantum-pulse/engine/textureManager"
)

const (
	// SceneSrcDir is the expected location of scenes
	SceneSrcDir = "assets/scenes/"
)

type scene struct {
	Entities    map[string]EntityOld
	Renderables map[string]Renderable
	fileName    string
	shaders     sm.ShaderManager
	textures    tm.TextureManager
}

// Scene represents a logical grouping of entities
type Scene interface {
	Renderable
	AddEntity(EntityOld)
}

// NewScene creates a new Scene
func NewScene(fileName string, shaders sm.ShaderManager, textures tm.TextureManager) Scene {
	scene := scene{
		fileName:    fileName,
		Entities:    make(map[string]EntityOld),
		Renderables: make(map[string]Renderable),
		shaders:     shaders,
		textures:    textures,
	}

	scene.loadSceneFile(fileName)

	return &scene
}

// ID returns the id of the scene which currently is the scene filename used to load the scene.
func (s *scene) ID() string {
	return s.fileName
}

// AddEntity adds an Entity to the scene.
func (s *scene) AddEntity(ent EntityOld) {
	s.Entities[ent.ID()] = ent

	// if it is also a Renderable add it to the map of Renderables
	if rend, ok := ent.(Renderable); ok {
		s.Renderables[rend.ID()] = rend
	}
}

// Update is called to update all scene components.
func (s *scene) Update(elapsed float64) {
	for _, r := range s.Renderables {
		r.Update(elapsed)
	}
}

// Render will render each of it Renderable entities.
func (s *scene) Render() {
	for _, r := range s.Renderables {
		r.Render()
	}
}

func (s *scene) loadSceneFile(fileName string) {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", SceneSrcDir, fileName))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var sd sceneData
	json.Unmarshal(data, &sd)

	for _, modelFile := range sd.Models {
		m := NewModel(modelFile.Name, s.shaders, s.textures)
		m.Load(modelFile.Name)
		s.AddEntity(m)
	}
}

type sceneData struct {
	Models []sceneModels `json:"models"`
}

type sceneModels struct {
	Name     string    `json:"name"`
	Position []float32 `json:"position"`
}
