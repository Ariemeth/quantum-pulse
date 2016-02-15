package engine

type scene struct {
	entities map[string]Renderable
}

//Scene represents a logical grouping of entities
type Scene interface {
	Renderable
}

//NewScene creates a new Scene
func NewScene() Scene {
	scene := new(scene)
	return scene
}

//Update updates all entities in the scene
func (s *scene) Update() {

}

//Render renders all objects in the scene
func (s *scene) Render() {

}
