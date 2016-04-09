package engine

// Renderable represents an object that can be rendered
type Renderable interface {
	Entity
	Update(elapsed float64)
	Render()
}
