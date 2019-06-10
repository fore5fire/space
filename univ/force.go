package univ

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
)

// Acceleration is an Observer that modifies a body's velocity and angular velocity over time.
type Acceleration struct {
	linearVector  mgl32.Vec3
	angularVector mgl32.Vec3
	body          *Body

	accelTicker *draw.Ticker
}

// NewAcceleration creates a new acceleration with both a linear and angular component.
//
// The new acceleration is initally paused, and will not be applied to the object until Start() is
// called on it.
func NewAcceleration(b *Body, linearVector, angularVector mgl32.Vec3) *Acceleration {
	a := &Acceleration{
		linearVector:  linearVector,
		angularVector: angularVector,
		body:          b,
	}
	a.accelTicker = draw.NewTicker(DefaultRefreshRate, a.accelTick)
	return a
}

// NewLinearAcceleration creates a new acceleration with only a linear component.
//
// The new acceleration is initally paused, and will not be applied to the object until Start() is
// called on it.
func NewLinearAcceleration(b *Body, vector mgl32.Vec3) *Acceleration {
	return NewAcceleration(b, vector, mgl32.Vec3{})
}

// NewAngularAcceleration creates a new acceleration with only an angular component.
//
// The new acceleration is initally paused, and will not be applied to the object until Start() is
// called on it.
func NewAngularAcceleration(b *Body, vector mgl32.Vec3) *Acceleration {
	return NewAcceleration(b, mgl32.Vec3{}, vector)
}

// Start starts this force accelerating its body.
func (a *Acceleration) Start() {
	a.accelTicker.Start()
}

// Pause stops this force from accelerating its body, but continues applying its velocity.
func (a *Acceleration) Pause() {
	a.accelTicker.Stop()
}

// Destroy stops f and cleans up its resources. f should not be used after it is destroyed.
func (a *Acceleration) Destroy() {
	a.accelTicker.Close()
}

func (a *Acceleration) accelTick(elapsed float32) {
	a.body.AddVelocity(a.body.Rotation().Rotate(a.linearVector.Mul(elapsed)))
	a.body.AddAngularV(a.body.Rotation().Rotate(a.angularVector.Mul(elapsed)))
}

// // Force is a force that applies a linear force and a torque, calculated
// // by the  application point of the force
// type Force struct {
// 	linearForce LinearForce
// 	torque      Torque
// }

// func NewForce(vector mgl32.Vec3, relativePosition mgl32.Vec3) *Force {
// 	return &Force{
// 		linearForce: LinearForce{},
// 		torque:      Torque{},
// 	}
// }
