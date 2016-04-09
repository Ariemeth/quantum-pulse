package engine

// Renderable represents an object that can be rendered
type Renderable interface {
	EntityOld
	Update(elapsed float64)
	Render()
}
