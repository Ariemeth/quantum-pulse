package components

// ComponentTypeModel represents a model component's type.
const ComponentTypeModel = "model"

// Model represents a component that holds the data representing a model.
type Model interface {
	Component
	Data(chan VertexData)
}

// NewModel creates a new Model component.
func NewModel() Model {
	m := model{}
	return &m
}

type model struct {
	data     VertexData
	dataChan chan VertexData
}

// ComponentType retrieves the type of this component.
func (m *model) ComponentType() string {
	return ComponentTypeModel
}

// Data is retrieves the model matrix from the component
func (m *model) Data(dc chan VertexData) {
	dc <- m.data
}

func (m *model) UpdateDataChannel() {
	//TODO how to update the data without creating another go routine in all componets.
}

// VertexData represents the the data needed to construct a 3d object.
type VertexData struct {
	Indexed  bool      `json:"indexed"`
	Verts    []float32 `json:"verts"`
	Indices  []uint32  `json:"indices"`
	VertSize int32     `json:"vertSize"`
}
