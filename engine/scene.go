package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
	"github.com/Ariemeth/quantum-pulse/resources"
	"github.com/Ariemeth/quantum-pulse/systems"
)

const (
	// SceneSrcDir is the expected location of scenes
	SceneSrcDir = "assets/scenes/"
)

type scene struct {
	Renderer systems.Renderer
	//	Animator systems.Animator
	Movement systems.Movement
	fileName string
	assets   *resources.Manager
}

// Scene represents a logical grouping of entities
type Scene interface {
	// ID retrieves the id of this scene.
	ID() string
	// Stop stops the components of the scene.
	Stop()
	// Start starts the scene components.
	Start()
	// Terminate stops the scene and frees any resources.
	Terminate()
}

// newScene creates a new Scene
func newScene(fileName string, assets *resources.Manager, window *glfw.Window, width, height int) (*scene, error) {
	scene := scene{
		fileName: fileName,
		Renderer: systems.NewRenderer(assets, runOnMain, window),
		//		Animator: systems.NewAnimator(),
		Movement: systems.NewMovement(),
		assets:   assets,
	}

	err := scene.loadSceneFile(fileName, width, height)
	if err != nil {
		return nil, err
	}

	return &scene, nil
}

// ID returns the id of the scene which currently is the scene filename used to load the scene.
func (s *scene) ID() string {
	return s.fileName
}

func (s *scene) Stop() {
	s.Renderer.Stop()
	s.Movement.Stop()
}

func (s *scene) Start() {
	s.Renderer.Start()
	s.Movement.Start()
}

func (s *scene) Terminate() {
	s.Renderer.Terminate()
	s.Movement.Terminate()
}

func (s *scene) loadSceneFile(fileName string, width, height int) error {

	data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", SceneSrcDir, fileName))
	if err != nil {
		return err
	}

	var sd sceneData
	err = json.Unmarshal(data, &sd)
	if err != nil {
		return err
	}

	// configure the camera
	cam := components.NewCamera()
	cam.SetView(sd.Camera.Position, sd.Camera.LookAt, sd.Camera.Up)
	cam.SetProjection(sd.Camera.FOVY, sd.Camera.NearPlane, sd.Camera.FarPlane, width, height)
	s.Renderer.LoadCamera(cam)

	// Load models
	for _, modelFile := range sd.Models {

		// Load the mesh component
		mesh := components.NewMesh()
		err := mesh.Load(modelFile.FileName)
		if err != nil {
			fmt.Printf("Unable to load mesh:%s", modelFile.FileName)
			continue
		}

		// Load a transform component
		t := components.NewTransform()
		pos := mgl32.Vec3{modelFile.Position[0], modelFile.Position[1], modelFile.Position[2]}
		t.Translate(pos)

		ent := entity.NewEntity(modelFile.Name)
		ent.AddComponent(mesh)
		ent.AddComponent(t)

		// Load the velocity component
		vel := components.NewVelocity()
		rotVel := mgl32.Vec3{modelFile.RotVelocity[0], modelFile.RotVelocity[1], modelFile.RotVelocity[2]}
		tranVel := mgl32.Vec3{modelFile.TransVelocity[0], modelFile.TransVelocity[1], modelFile.TransVelocity[2]}
		vel.Set(rotVel, tranVel)

		// Load the acceleration component
		accel := components.NewAcceleration()
		accRotVel := mgl32.Vec3{modelFile.RotAccel[0], modelFile.RotAccel[1], modelFile.RotAccel[2]}
		accTranVel := mgl32.Vec3{modelFile.TransAccel[0], modelFile.TransAccel[1], modelFile.TransAccel[2]}
		accel.Set(accRotVel, accTranVel)

		ent.AddComponent(vel)
		ent.AddComponent(accel)

		// Add the model to appropriate systems
		s.Renderer.AddEntity(ent)
		s.Movement.AddEntity(ent)
	}
	return nil
}

type sceneData struct {
	Camera sceneCamera   `json:"defaultCamera"`
	Models []sceneModels `json:"models"`
}

type sceneModels struct {
	Name          string     `json:"name"`
	FileName      string     `json:"fileName"`
	Position      [3]float32 `json:"position"`
	RotAccel      [3]float32 `json:"rotationalAcceleration"`
	TransAccel    [3]float32 `json:"translationalAcceleration"`
	RotVelocity   [3]float32 `json:"rotationalVelocity"`
	TransVelocity [3]float32 `json:"translationalVelocity"`
}

type sceneCamera struct {
	Position  [3]float32 `json:"position"`
	LookAt    [3]float32 `json:"lookat"`
	Up        [3]float32 `json:"up"`
	FOVY      float32    `json:"fovy"`
	NearPlane float32    `json:"nearPlane"`
	FarPlane  float32    `json:"farPlane"`
}
