package univ

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

// ChaseCam is a camera that keeps a relative location to a body.
type ChaseCam struct {
	locMut   sync.RWMutex
	location mgl32.Vec3
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

// Destroy stops the camera from updating. Always call Destroy on ChaseCams that are no longer needed.
func (cam *ChaseCam) Destroy() {
	cam.ticker.Close()
}

func (cam *ChaseCam) tick(elapsed float32) {
	loc := cam.target.GetLocation().Add(cam.location)
	// rot := b.GetRotation().Mul(cam.rotation)
	transform := mgl32.LookAtV(loc, cam.target.GetLocation(), mgl32.Vec3{0, 1, 0})
	// transform := mgl32.Translate3D(loc.Elem()).Mul4(rot.Normalize().Mat4())
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

// // GetRotation gets the rotation of b
// func (cam *ChaseCam) GetRotation() mgl32.Quat {
// 	cam.rotMut.RLock()
// 	defer cam.rotMut.RUnlock()
// 	return cam.rotation
// }

// // Rotate rotates b by offset
// func (cam *ChaseCam) Rotate(offset mgl32.Quat) {
// 	cam.rotMut.Lock()
// 	cam.rotation = cam.rotation.Mul(offset)
// 	cam.rotMut.Unlock()
// }

// // SetRotation sets the rotation of b to rot
// func (cam *ChaseCam) SetRotation(rot mgl32.Quat) {
// 	cam.rotMut.Lock()
// 	cam.rotation = rot
// 	cam.rotMut.Unlock()
// }
