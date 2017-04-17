package components

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// TypeCamera represents a Camera component's type.
	TypeCamera = "Camera"
)

// Camera represents the behaviors of any camera.
type Camera interface {
	// SetView sets the view matrix.
	SetView(camEye, camLookAt, camUp [3]float32)
	// SetProjection sets the projection matrix.
	SetProjection(fovY, nearPlane, farPlane float32, windowWidth, windowHeight int)
	// Projection retrieves the projection matrix.
	Projection() mgl32.Mat4
	// View retrieves the view matrix.
	View() mgl32.Mat4
	// PositionVec3 retrieves the camera position.
	PositionVec3() mgl32.Vec3
	// LookAtVec3 retrieves the point the camera is looking.
	LookAtVec3() mgl32.Vec3
	// UpVec3 retrieves the camera's up vector.
	UpVec3() mgl32.Vec3
	// FOVy retrieves the camera's FOVY.
	FOVy() float32
	// NearPlane retrieves the cameras current near plane.
	NearPlane() float32
	// FarPlane retrieves the camera's current far plane.
	FarPlane() float32
	// WindowHeight retrieves the window height used to calculate the perspective.
	WindowHeight() int
	// WindowWidth retrieves the window width used to calculate the perspective.
	WindowWidth() int
}

type camera struct {
	position      mgl32.Vec3
	lookAt        mgl32.Vec3
	up            mgl32.Vec3
	fovY          float32
	nearPlane     float32
	farPlane      float32
	projectionMat mgl32.Mat4
	viewMat       mgl32.Mat4
	windowHeight  int
	windowWidth   int
}

// NewCamera creates a new Camera component.
func NewCamera() Camera {
	c := camera{
		projectionMat: mgl32.Ident4(),
		viewMat:       mgl32.Ident4(),
	}

	return &c
}

// Type retrieves the type name of this component.
func (c *camera) Type() string {
	return TypeCamera
}

// SetView sets the view matrix.
func (c *camera) SetView(camEye, camLookAt, camUp [3]float32) {
	c.position = mgl32.Vec3{camEye[0], camEye[1], camEye[2]}
	c.lookAt = mgl32.Vec3{camLookAt[0], camLookAt[1], camLookAt[2]}
	c.up = mgl32.Vec3{camUp[0], camUp[1], camUp[2]}
	c.viewMat = mgl32.LookAtV(c.position, c.lookAt, c.up)
}

// SetProjection sets the projection matrix.
func (c *camera) SetProjection(fovY, nearPlane, farPlane float32, windowWidth, windowHeight int) {
	// Set the perspective
	c.fovY = mgl32.DegToRad(fovY)
	c.nearPlane, c.farPlane = nearPlane, farPlane
	c.windowHeight, c.windowWidth = windowHeight, windowWidth

	c.projectionMat = mgl32.Perspective(c.fovY, float32(windowWidth)/float32(windowHeight), nearPlane, farPlane)
}

// Projection retrieves the projection matrix.
func (c camera) Projection() mgl32.Mat4 {
	return c.projectionMat
}

// View retrieves the view matrix.
func (c camera) View() mgl32.Mat4 {
	return c.viewMat
}

// PositionVec3 retrieves the camera position.
func (c camera) PositionVec3() mgl32.Vec3 {
	return c.position
}

// LookAtVec3 retrieves the point the camera is looking.
func (c camera) LookAtVec3() mgl32.Vec3 {
	return c.lookAt
}

// UpVec3 retrieves the camera's up vector.
func (c camera) UpVec3() mgl32.Vec3 {
	return c.up
}

// FOVy retrieves the camera's FOVY.
func (c camera) FOVy() float32 {
	return c.fovY
}

// NearPlane retrieves the cameras current near plane.
func (c camera) NearPlane() float32 {
	return c.nearPlane
}

// FarPlane retrieves the camera's current far plane.
func (c camera) FarPlane() float32 {
	return c.farPlane
}

// WindowHeight retrieves the window height used to calculate the perspective.
func (c camera) WindowHeight() int {
	return c.windowHeight
}

// WindowWidth retrieves the window width used to calculate the perspective.
func (c camera) WindowWidth() int {
	return c.windowWidth
}
