package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	am "github.com/Ariemeth/quantum-pulse/engine/assets"
	"github.com/Ariemeth/quantum-pulse/engine/components"
	"github.com/Ariemeth/quantum-pulse/engine/entity"
	"github.com/Ariemeth/quantum-pulse/engine/systems"
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
	Movement   systems.Movement
	fileName   string
	assets     *am.AssetManager
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
func NewScene(fileName string, assets *am.AssetManager) Scene {
	scene := scene{
		fileName:   fileName,
		Renderer:   systems.NewRenderer(),
		Animator:   systems.NewAnimator(),
		Movement:   systems.NewMovement(),
		assets:     assets,
		camera:     mgl32.Ident4(),
		projection: mgl32.Ident4(),
	}

	scene.loadSceneFile(fileName)

	return &scene
}

// ID returns the id of the scene which currently is the scene filename used to load the scene.
func (s *scene) ID() string {
	return s.fileName
}

// Update is called to update all scene components.
func (s *scene) Update(elapsed float64) {
	//s.Animator.Process(elapsed)
	s.Movement.Process(elapsed)
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
	err = json.Unmarshal(data, &sd)
	if err != nil {
		fmt.Println(err)
		return
	}

	// configure the camera
	camEye := mgl32.Vec3{sd.Camera.Position[0], sd.Camera.Position[1], sd.Camera.Position[2]}
	camAt := mgl32.Vec3{sd.Camera.LookAt[0], sd.Camera.LookAt[1], sd.Camera.LookAt[2]}
	camUp := mgl32.Vec3{sd.Camera.Up[0], sd.Camera.Up[1], sd.Camera.Up[2]}
	s.camera = mgl32.LookAtV(camEye, camAt, camUp)

	// Set the perspective
	fovy := mgl32.DegToRad(sd.Camera.FOVY)
	s.projection = mgl32.Perspective(fovy, float32(windowWidth)/windowHeight, sd.Camera.NearPlane, sd.Camera.FarPlane)

	// Load models
	for _, modelFile := range sd.Models {

		// Load the mesh
		mesh := components.NewMesh()
		err := mesh.Load(modelFile.FileName)
		if err != nil {
			fmt.Printf("Unable to load mesh:%s", modelFile.FileName)
			continue
		}

		// Load a transform
		t := components.NewTransform()
		pos := mgl32.Vec3{modelFile.Position[0], modelFile.Position[1], modelFile.Position[2]}
		t.Translate(pos)

		// Load the shader
		md := mesh.Data()
		shaderName := fmt.Sprintf("%s:%s", md.VertShaderFile, md.FragShaderFile)
		program, err := s.assets.Shaders().LoadProgramFromFile(md.VertShaderFile, md.FragShaderFile, shaderName, false)
		if err != nil {
			fmt.Printf("Unable to load shader:%s", shaderName)
			continue
		}
		shader := components.NewShader(shaderName, program)
		shader.CreateVAO(mesh)

		// Load textures
		texture, err := s.assets.Textures().LoadTexture(md.TextureFile, md.TextureFile)
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

		vel := components.NewVelocity()
		// set rotation velocity to rotate around the Y axis at 1 radian per second
		rotVel := mgl32.Vec3{modelFile.RotVelocity[0], modelFile.RotVelocity[1], modelFile.RotVelocity[2]}
		vel.SetRotational(rotVel)
		accel := components.NewAcceleration()
		accVel := mgl32.Vec3{modelFile.RotAccel[0], modelFile.RotAccel[1], modelFile.RotAccel[2]}
		accel.SetRotational(accVel)

		ent.AddComponent(vel)
		ent.AddComponent(accel)

		s.Renderer.AddEntity(ent)
		s.Movement.AddEntity(ent)
	}

}

type sceneData struct {
	Camera sceneCamera   `json:"defaultCamera"`
	Models []sceneModels `json:"models"`
}

type sceneModels struct {
	Name          string    `json:"name"`
	FileName      string    `json:"fileName"`
	Position      []float32 `json:"position"`
	RotAccel      []float32 `json:"rotationalAcceleration"`
	TransAccel    []float32 `json:"translationalAcceleration"`
	RotVelocity   []float32 `json:"rotationalVelocity"`
	TransVelocity []float32 `json:"translationalVelocity"`
}

type sceneCamera struct {
	Position  []float32 `json:"position"`
	LookAt    []float32 `json:"lookat"`
	Up        []float32 `json:"up"`
	FOVY      float32   `json:"fovy"`
	NearPlane float32   `json:"nearPlane"`
	FarPlane  float32   `json:"farPlane"`
}
