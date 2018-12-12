package univ

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
)

// FreeCam is a camera that moves independantly of any body.
type FreeCam struct {
	mut      sync.RWMutex
	location mgl32.Vec3
	rotation mgl32.Quat
	ticker   *Ticker

	program *draw.Program
}

// NewFreeCam creates a new FreeCam with the specified target and update interval.
func NewFreeCam(u *Universe, programType draw.ProgramType) *FreeCam {
	cam := &FreeCam{
		program: u.window.GetProgram(programType),
	}
	cam.ticker = NewTicker(DefaultRefreshRate, cam.tick)
	cam.ticker.Start()
	return cam
}

// Destroy stops the camera from updating. Always call Destroy on ChaseCams that are no longer needed.
func (cam *FreeCam) Destroy() {
	cam.ticker.Close()
}

func (cam *FreeCam) tick(elapsed float32) {
	cam.mut.Lock()
	transform := cam.rotation.Normalize().Mat4().Mul4(mgl32.Translate3D(cam.location.Elem()))
	cam.program.SetView(transform)
	cam.mut.Unlock()
}

// GetLocation returns the current location of b
func (cam *FreeCam) GetLocation() mgl32.Vec3 {
	cam.mut.RLock()
	defer cam.mut.RUnlock()
	return cam.location
}

// Translate translates the location of b by offset
func (cam *FreeCam) Translate(offset mgl32.Vec3) {
	cam.mut.Lock()
	cam.location = cam.location.Add(offset)
	cam.mut.Unlock()
}

// SetLocation sets the location of b
func (cam *FreeCam) SetLocation(loc mgl32.Vec3) {
	cam.mut.Lock()
	cam.location = loc
	cam.mut.Unlock()
}

// GetRotation gets the rotation of cam
func (cam *FreeCam) GetRotation() mgl32.Quat {
	cam.mut.RLock()
	defer cam.mut.RUnlock()
	return cam.rotation
}

// Rotate rotates cam by offset
func (cam *FreeCam) Rotate(offset mgl32.Quat) {
	cam.mut.Lock()
	cam.rotation = cam.rotation.Mul(offset)
	cam.mut.Unlock()
}

// SetRotation sets the rotation of cam to rot
func (cam *FreeCam) SetRotation(rot mgl32.Quat) {
	cam.mut.Lock()
	cam.rotation = rot
	cam.mut.Unlock()
}
