package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

const (
	// TypeMesh represents a mesh component's type.
	TypeMesh = "mesh"
	// MeshSrcDir is the expected location of meshes
	MeshSrcDir = "assets/models/"
)

// Mesh represents a component that holds the data representing a mesh.
type Mesh interface {
	Component
	// Data retrieves the mesh data.
	Data() MeshData
	// Set updates the mesh data with the new data passed in as a parameter.
	Set(MeshData)
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

// Type retrieves the type of this component.
func (m *mesh) Type() string {
	return TypeMesh
}

// Data retrieves the mesh data from the component.
func (m *mesh) Data() MeshData {
	m.dataLock.RLock()
	defer m.dataLock.RUnlock()
	return m.data
}

// Set sets the MeshData to the new data.
func (m *mesh) Set(md MeshData) {
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

	fullFileName := ""
	if strings.Contains(fileName, MeshSrcDir) {
		fullFileName = fileName
	} else {
		fullFileName = fmt.Sprintf("%s%s", MeshSrcDir, fileName)
	}

	data, err := ioutil.ReadFile(fullFileName)
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
