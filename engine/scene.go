package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Ariemeth/quantum-pulse/engine/components"
	"github.com/Ariemeth/quantum-pulse/engine/entity"
	sm "github.com/Ariemeth/quantum-pulse/engine/shaderManager"
	"github.com/Ariemeth/quantum-pulse/engine/systems"
	tm "github.com/Ariemeth/quantum-pulse/engine/textureManager"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// SceneSrcDir is the expected location of scenes
	SceneSrcDir = "assets/scenes/"
)

type scene struct {
	Renderer   systems.Renderer
	Animator   systems.Animator
	fileName   string
	shaders    sm.ShaderManager
	textures   tm.TextureManager
	camera     mgl32.Mat4
	projection mgl32.Mat4
}

// Scene represents a logical grouping of entities
type Scene interface {
	// Update is called to update all scene components.
	Update(float64)
	// Render will render each of it Renderable entities.
	Render()
	// ID retrieves the id of this scene.
	ID() string
}

// NewScene creates a new Scene
func NewScene(fileName string, shaders sm.ShaderManager, textures tm.TextureManager) Scene {
	scene := scene{
		fileName:   fileName,
		Renderer:   systems.NewRenderer(),
		Animator:   systems.NewAnimator(),
		shaders:    shaders,
		textures:   textures,
		camera:     mgl32.Ident4(),
		projection: mgl32.Ident4(),
	}

	scene.projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)

	scene.camera = mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	scene.loadSceneFile(fileName)

	return &scene
}

// ID returns the id of the scene which currently is the scene filename used to load the scene.
func (s *scene) ID() string {
	return s.fileName
}

// Update is called to update all scene components.
func (s *scene) Update(elapsed float64) {
	s.Animator.Process(elapsed)
}

// Render will render each of it Renderable entities.
func (s *scene) Render() {
	s.Renderer.Process()
}

func (s *scene) loadSceneFile(fileName string) {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", SceneSrcDir, fileName))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var sd sceneData
	json.Unmarshal(data, &sd)

	for _, modelFile := range sd.Models {

		// Load the mesh
		mesh := components.NewMesh()
		err := mesh.Load(modelFile.FileName)
		if err != nil {
			fmt.Printf("Unable to load mesh:%s", modelFile.FileName)
			continue
		}

		// Load a transform
		// TODO load the position into the transform
		t := components.NewTransform()

		// Load the shader
		md := mesh.Data()
		shaderName := fmt.Sprintf("%s:%s", md.VertShaderFile, md.FragShaderFile)
		program, err := s.shaders.LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, shaderName, false)
		if err != nil {
			fmt.Printf("Unable to load shader:%s", shaderName)
			continue
		}
		shader := components.NewShader(shaderName, program)
		shader.LoadMesh(mesh)
		shader.LoadTransform(t)

		// Load textures
		texture, err := s.textures.LoadTexture(md.TextureFile, md.TextureFile)
		if err != nil {
			fmt.Printf("Unable to load texture:%s", md.TextureFile)
			continue
		}
		shader.AddTexture(texture)

		//TODO this really only needs to be done once per shader
		gl.UniformMatrix4fv(shader.GetUniformLoc(components.ProjectionUniform), 1, false, &s.projection[0])

		//TODO this really only needs to be done once per shader
		gl.UniformMatrix4fv(shader.GetUniformLoc(components.CameraUniform), 1, false, &s.camera[0])

		ent := entity.NewEntity(modelFile.Name)
		ent.AddComponent(mesh)
		ent.AddComponent(t)
		ent.AddComponent(shader)

		s.Renderer.AddEntity(ent)

		// should be checking to see if the scene file indicates the model should be animated

		s.Animator.AddEntity(ent)

	}

}

type sceneData struct {
	Models []sceneModels `json:"models"`
}

type sceneModels struct {
	Name     string    `json:"name"`
	FileName string    `json:"fileName"`
	Position []float32 `json:"position"`
}
