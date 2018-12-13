package univ

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

// ChaseCam is a camera that keeps a relative location to a body.
type ChaseCam struct {
	locMut   sync.RWMutex
	rotMut   sync.RWMutex
	location mgl32.Vec3
	rotation mgl32.Quat
	target   *Body
	ticker   *Ticker
}

// NewChaseCam creates a new ChaseCam with the specified target and update interval.
func NewChaseCam(target *Body) *ChaseCam {
	cam := &ChaseCam{
		target: target,
	}
	cam.ticker = NewTicker(DefaultRefreshRate, cam.tick)
	cam.ticker.Start()
	return cam
}

// Remove stops the camera from updating. Always call Remove on ChaseCams that are no longer needed.
func (cam *ChaseCam) Remove() {
	cam.ticker.Close()
}

func (cam *ChaseCam) tick(elapsed float32) {
	// set the position of the camera and the look at point as relative positions to the direction and position of the target
	rot := cam.target.GetRotation()
	lookAtMat := mgl32.Translate3D(cam.target.GetLocation().Elem())
	lookAtMatRot := lookAtMat.Mul4(rot.Normalize().Mat4())
	lookAt := lookAtMatRot.Mul4(mgl32.Translate3D(0.0, 0.0, 5.0)).Col(3).Vec3()
	lookFrom := lookAtMatRot.Mul4(mgl32.Translate3D(cam.location.Elem())).Col(3).Vec3()

	transform := mgl32.LookAtV(lookFrom, lookAt, mgl32.Vec3{0, 1, 0})
	cam.target.program.SetView(transform)
}

// GetLocation returns the current location of b
func (cam *ChaseCam) GetLocation() mgl32.Vec3 {
	cam.locMut.RLock()
	defer cam.locMut.RUnlock()
	return cam.location
}

// Translate translates the location of b by offset
func (cam *ChaseCam) Translate(offset mgl32.Vec3) {
	cam.locMut.Lock()
	cam.location = cam.location.Add(offset)
	cam.locMut.Unlock()
}

// SetLocation sets the location of b
func (cam *ChaseCam) SetLocation(loc mgl32.Vec3) {
	cam.locMut.Lock()
	cam.location = loc
	cam.locMut.Unlock()
}

// GetRotation gets the rotation of b
func (cam *ChaseCam) GetRotation() mgl32.Quat {
	cam.rotMut.RLock()
	defer cam.rotMut.RUnlock()
	return cam.rotation
}

// Rotate rotates b by offset
func (cam *ChaseCam) Rotate(offset mgl32.Quat) {
	cam.rotMut.Lock()
	cam.rotation = cam.rotation.Mul(offset)
	cam.rotMut.Unlock()
}

// SetRotation sets the rotation of b to rot
func (cam *ChaseCam) SetRotation(rot mgl32.Quat) {
	cam.rotMut.Lock()
	cam.rotation = rot
	cam.rotMut.Unlock()
}
