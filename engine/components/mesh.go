package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
)

const (
	// ComponentTypeMesh represents a mesh component's type.
	ComponentTypeMesh = "mesh"
	// MeshSrcDir is the expected location of meshes
	MeshSrcDir = "assets/models/"
)

// Mesh represents a component that holds the data representing a mesh.
type Mesh interface {
	// Data retrieves the mesh data.
	Data() MeshData
	// Update updates the mesh data with the new data passed in as a parameter.
	Update(MeshData)
	// Load loads the mesh data from file.
	Load(string) error
}

// NewMesh creates a new Mesh component.
func NewMesh() Mesh {
	m := mesh{}
	return &m
}

type mesh struct {
	data     MeshData
	dataLock sync.RWMutex
	isLoaded bool
}

// ComponentType retrieves the type of this component.
func (m *mesh) ComponentType() string {
	return ComponentTypeMesh
}

// Data retrieves the mesh data from the component.
func (m *mesh) Data() MeshData {
	m.dataLock.RLock()
	defer m.dataLock.RUnlock()
	return m.data
}

// Update sets the MeshData to the new data.
func (m *mesh) Update(md MeshData) {
	m.dataLock.Lock()
	defer m.dataLock.Unlock()
	m.data = md
}

// Load loads a mesh from file and returns an error if an error occurs while loading the file.  If the mesh has already been loaded an "Already Loaded" error will be returned.
func (m *mesh) Load(fileName string) error {
	m.dataLock.Lock()
	defer m.dataLock.Unlock()

	if m.isLoaded {
		return errors.New("Already Loaded")
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", MeshSrcDir, fileName))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(data, &m.data)
	if err == nil {
		m.isLoaded = true
	}
	return err
}

// MeshData represents the the data needed to construct a 3d mesh.
type MeshData struct {
	Indexed        bool      `json:"indexed"`
	Verts          []float32 `json:"verts"`
	Indices        []uint32  `json:"indices"`
	VertSize       int32     `json:"vertSize"`
	TextureFile    string    `json:"textureFile"`
	FragShaderFile string    `json:"fragShaderFile"`
	VertShaderFile string    `json:"vertShaderFile"`
}