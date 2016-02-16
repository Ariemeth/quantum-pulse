package engine

import (
	"fmt"

	sm "github.com/Ariemeth/frame-assault-2/engine/shaderManager"
	tm "github.com/Ariemeth/frame-assault-2/engine/textureManager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//Model represents a physical entity
type Model struct {
	Entity
	vertices          []float32
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

//NewModel creates a new model
func NewModel(id string, shaders sm.ShaderManager, textures tm.TextureManager, shader string) *Model {
	m := Model{
		Entity:     NewEntity(id),
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

//Update updates the model
func (m *Model) Update(elapsed float64) {
	m.angle += elapsed
	m.model = mgl32.HomogRotate3D(float32(m.angle), mgl32.Vec3{0, 1, 0})
}

//Render renders the model
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

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

	gl.BindVertexArray(0)
}

func (m *Model) Load() {
	program, isLoaded := m.shaders.GetShader(m.shaders.GetDefaultShader())
	if isLoaded {
		gl.UseProgram(program)
	} else {
		fmt.Println("Unable to load default shader")
	}

	//////////////////////////////
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	m.textures.LoadTexture("assets/textures/square.png", "square")

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	m.currentProgram = program
	m.vao = vao
	m.model = model
	m.modelUniform = modelUniform

	gl.BindVertexArray(0)
}

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
