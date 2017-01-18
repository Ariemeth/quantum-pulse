package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	am "github.com/Ariemeth/quantum-pulse/assets"
	"github.com/Ariemeth/quantum-pulse/components"
	"github.com/Ariemeth/quantum-pulse/entity"
	"github.com/Ariemeth/quantum-pulse/systems"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// SceneSrcDir is the expected location of scenes
	SceneSrcDir = "assets/scenes/"
)

type scene struct {
	Renderer systems.Renderer
	Animator systems.Animator
	Movement systems.Movement
	fileName string
	assets   *am.AssetManager
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
		fileName: fileName,
		Renderer: systems.NewRenderer(assets),
		Animator: systems.NewAnimator(),
		Movement: systems.NewMovement(),
		assets:   assets,
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
	cam := components.NewCamera()
	cam.SetView(sd.Camera.Position, sd.Camera.LookAt, sd.Camera.Up)
	cam.SetProjection(sd.Camera.FOVY, sd.Camera.NearPlane, sd.Camera.FarPlane, windowWidth, windowHeight)
	s.Renderer.LoadCamera(cam)

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

		ent := entity.NewEntity(modelFile.Name)
		ent.AddComponent(mesh)
		ent.AddComponent(t)

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
