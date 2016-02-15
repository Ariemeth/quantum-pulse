package engine

type scene struct {
	Entities map[string]*Entity
}

//Scene represents a logical grouping of entities
type Scene interface {
}

//NewScene creates a new Scene
func NewScene() Scene {
	scene := new(scene)
	return scene
}

//AddEntity adds an Entity to the scene
func (s *scene) AddEntity(ent *Entity){
	s.Entities[ent.ID()] = ent
}
