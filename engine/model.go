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
	indices           []uint32
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
		vertices:   iHexVerts,
		indices:    hexIndices,
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
	
	if m.indices != nil {
		gl.DrawElements(gl.TRIANGLE_FAN, int32(len(m.indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.vertices))/5)
	}

	gl.BindVertexArray(0)
}

//Load loads and sets up the model
func (m *Model) Load(isIndexed bool) {
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
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertices)*4, gl.Ptr(m.vertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(m.currentProgram, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0)) //5:number of values per vertex, 4:number of bytes in a float32

	texCoordAttrib := uint32(gl.GetAttribLocation(m.currentProgram, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, true, 3*4, gl.PtrOffset(3*4)) //5:number of values per vertex, 4:number of bytes in a float32

	if isIndexed {
		var indices uint32
		gl.GenBuffers(1, &indices)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indices)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.indices)*4, gl.Ptr(m.indices), gl.STATIC_DRAW)
	}

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

type point struct {
	X float32
	Y float32
	Z float32
	U float32
	V float32
}

var (
	Center = point{0.0, 0.0, 0.0, 0.5, 0.5}
	P1     = point{0.0, 1.0, 0.0, 0.5, 0.0}
	P2     = point{-1.0, 0.5, 0.0, 0.0, 0.25}
	P3     = point{-1.0, -0.5, 0.0, 0.0, 0.75}
	P4     = point{0.0, -1.0, 0.0, 0.5, 1.0}
	P5     = point{1.0, -0.5, 0.0, 1.0, 0.75}
	P6     = point{1.0, 0.5, 0.0, 1.0, 0.25}
)

var iHexVerts = []float32{
	// X, Y, Z
	// Top
	Center.X, Center.Y, Center.Z,
	P1.X, P1.Y, P1.Z, 
	P2.X, P2.Y, P2.Z, 
	P3.X, P3.Y, P3.Z,
	P4.X, P4.Y, P4.Z,
	P5.X, P5.Y, P5.Z, 
	P6.X, P6.Y, P6.Z,
}

var iHexVerts2 = []float32{
	// X, Y, Z
	// Top
	Center.X, Center.Y, Center.Z,
	P1.X, P1.Y, P1.Z, P1.U, P1.V,
	P2.X, P2.Y, P2.Z, P2.U, P2.V,
	P3.X, P3.Y, P3.Z, P3.U, P3.V,
	P4.X, P4.Y, P4.Z, P4.U, P4.V,
	P5.X, P5.Y, P5.Z, P5.U, P5.V,
	P6.X, P6.Y, P6.Z, P6.U, P6.V,
}

var hexIndices = []uint32{
	0, 1, 2, 3, 4, 5, 6, 1,
}

var hexVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	Center.X, Center.Y, Center.Z, Center.U, Center.V,
	P1.X, P1.Y, P1.Z, P1.U, P1.V,
	P2.X, P2.Y, P2.Z, P2.U, P2.V,
	Center.X, Center.Y, Center.Z, Center.U, Center.V,
	P2.X, P2.Y, P2.Z, P2.U, P2.V,
	P3.X, P3.Y, P3.Z, P3.U, P3.V,
	Center.X, Center.Y, Center.Z, Center.U, Center.V,
	P3.X, P3.Y, P3.Z, P3.U, P3.V,
	P4.X, P4.Y, P4.Z, P4.U, P4.V,
	Center.X, Center.Y, Center.Z, Center.U, Center.V,
	P4.X, P4.Y, P4.Z, P4.U, P4.V,
	P5.X, P5.Y, P5.Z, P5.U, P5.V,
	Center.X, Center.Y, Center.Z, Center.U, Center.V,
	P5.X, P5.Y, P5.Z, P5.U, P5.V,
	P6.X, P6.Y, P6.Z, P6.U, P6.V,
	Center.X, Center.Y, Center.Z, Center.U, Center.V,
	P6.X, P6.Y, P6.Z, P6.U, P6.V,
	P1.X, P1.Y, P1.Z, P1.U, P1.V,

	// Top

	// Front

	// Back

	// Left

	// Right

}
