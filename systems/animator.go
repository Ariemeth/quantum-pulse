package systems

const (
	// TypeAnimator is the name of the animator system.
	TypeAnimator = "animator"
)

/*
// Animator defines the behaviors expected of the animation system.
type Animator interface {
	System
	// Process performs handles animating all entities that have animations.
	Process(elapsed float64)
}

type animator struct {
	entities map[string]animatable
}

// NewAnimator creates a new Animator system.
func NewAnimator() Animator {
	a := animator{
		entities: make(map[string]animatable, 0),
	}
	return &a
}

func (a *animator) Start() {

}

func (a *animator) Stop() {

}

func (a *animator) IsRunning() bool {
	return false
}

// Type retrieves the type of system such as renderer, mover, etc.
func (a *animator) Type() string {
	return TypeAnimator
}

// AddEntity adds an Entity to the system.  Each system will have a component requirement that must be met before the Entity can be added.
func (a *animator) AddEntity(e entity.Entity) {
	transform, isTransform := e.Component(components.ComponentTypeTransform).(components.Transform)

	if isTransform {
		anim := animatable{
			Transform: transform,
			animation: make([]func(float64), 0),
		}
		// This is just temporary for testing.
		angle := 0.0
		rotate := func(elapsed float64) {
			angle += elapsed
			anim.Transform.Set(mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0}))
		}
		anim.animation = append(anim.animation, rotate)
		a.entities[e.ID()] = anim
	}
}

// RemoveEntity removes an Entity from the system.
func (a *animator) RemoveEntity(e entity.Entity) {
	delete(a.entities, e.ID())
}

// Process performs handles animating all entities that have animations.
func (a *animator) Process(elapsed float64) {
	for _, anims := range a.entities {
		for _, anim := range anims.animation {
			anim(elapsed)
		}
	}
}

type animatable struct {
	Entity    entity.Entity
	Transform components.Transform
	animation []func(float64)
}
*/
