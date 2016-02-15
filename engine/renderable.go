package engine

//Renderable represents an object that can be rendered
type Renderable interface {
	Update(elapsed float64)
	Render()
}
