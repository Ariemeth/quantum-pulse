package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sm "github.com/Ariemeth/frame-assault-2/engine/shaderManager"
	tm "github.com/Ariemeth/frame-assault-2/engine/textureManager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Model represents a physical entity
type Model struct {
	data              vertexData
	shaders           sm.ShaderManager
	textures          tm.TextureManager
	angle             float64
	model             mgl32.Mat4
	camera            mgl32.Mat4
	projection        mgl32.Mat4
	modelUniform      int32
	projectionUniform int32
	cameraUniform     int32
	textureUniform    int32
	vao               uint32
	currentProgram    uint32
	currentTexture    uint32
}

// NewModel creates a new model
func NewModel(id string, shaders sm.ShaderManager, textures tm.TextureManager, shader string) *Model {
	m := Model{
		shaders:    shaders,
		textures:   textures,
		angle:      0.0,
		model:      mgl32.Ident4(),
		camera:     mgl32.Ident4(),
		projection: mgl32.Ident4(),
	}

	program, status := shaders.GetShader(shader)
	if status {
		m.currentProgram = program
	} else {
		fmt.Println("unable to load program ", shader, ": ", program)
	}
	return &m
}

// Update updates the model
func (m *Model) Update(elapsed float64) {
	m.angle += elapsed
	m.model = mgl32.HomogRotate3D(float32(m.angle), mgl32.Vec3{0, 1, 0})
}

// Render renders the model
func (m *Model) Render() {

	gl.UseProgram(m.currentProgram)
	gl.UniformMatrix4fv(m.modelUniform, 1, false, &m.model[0])

	gl.BindVertexArray(m.vao)

	gl.ActiveTexture(gl.TEXTURE0)

	texture, isLoaded := m.textures.GetTexture("square")

	if isLoaded {
		gl.BindTexture(gl.TEXTURE_2D, texture)
	} else {
		fmt.Println("Unable to load texture")
	}

	if m.data.Indexed {
		gl.DrawElements(gl.TRIANGLE_FAN, int32(len(m.data.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.data.Verts))/5)
	}

	gl.BindVertexArray(0)
}

// Load loads and sets up the model
func (m *Model) Load(isIndexed bool) {
	m.LoadFile("hexagon.json")
	gl.UseProgram(m.currentProgram)

	m.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	m.projectionUniform = gl.GetUniformLocation(m.currentProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(m.projectionUniform, 1, false, &m.projection[0])

	m.camera = mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	m.cameraUniform = gl.GetUniformLocation(m.currentProgram, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(m.cameraUniform, 1, false, &m.camera[0])

	m.modelUniform = gl.GetUniformLocation(m.currentProgram, gl.Str("model\x00"))
	gl.UniformMatrix4fv(m.modelUniform, 1, false, &m.model[0])

	m.textureUniform = gl.GetUniformLocation(m.currentProgram, gl.Str("tex\x00"))
	gl.Uniform1i(m.textureUniform, 0)

	gl.BindFragDataLocation(m.currentProgram, 0, gl.Str("outputColor\x00"))

	// Load the texture
	m.textures.LoadTexture("assets/textures/square.png", "square")

	// Configure the vertex data
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.data.Verts)*4, gl.Ptr(m.data.Verts), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(m.currentProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, int32(m.data.VertSize)*4, gl.PtrOffset(0)) //5:number of values per vertex, 4:number of bytes in a float32

	texCoordAttrib := uint32(gl.GetAttribLocation(m.currentProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, true, int32(m.data.VertSize)*4, gl.PtrOffset(3*4)) //5:number of values per vertex, 4:number of bytes in a float32

	if m.data.Indexed {
		var indices uint32
		gl.GenBuffers(1, &indices)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indices)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.data.Indices)*4, gl.Ptr(m.data.Indices), gl.STATIC_DRAW)
	}

	gl.BindVertexArray(0)
}

// LoadFile loads a model's data from a json file.
func (m *Model) LoadFile(fileName string) {
	data, err := ioutil.ReadFile(fmt.Sprintf("assets/models/%s", fileName))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.Unmarshal(data, &m.data)
}

type vertexData struct {
	Indexed bool      `json:"indexed"`
	Verts   []float32 `json:"verts"`
	Indices []uint32  `json:"indices"`
	VertSize byte `json:"vertSize"`
}
