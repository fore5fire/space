package univ

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
)

// LinearForce is an observer that only applies the translation
// portion of a force to a body.
type LinearForce struct {
	vector mgl32.Vec3
	body   *Body

	accelTicker *draw.Ticker
}

// NewLinearForce creates a new LinearForce.
func NewLinearForce(b *Body, vector mgl32.Vec3) *LinearForce {
	f := &LinearForce{
		vector: vector,
		body:   b,
	}
	f.accelTicker = draw.NewTicker(DefaultRefreshRate, f.accelTick)
	return f
}

// Start starts this force accelerating its body.
func (f *LinearForce) Start() {
	f.accelTicker.Start()
}

// Pause stops this force from accelerating its body, but continues applying its velocity.
func (f *LinearForce) Pause() {
	f.accelTicker.Stop()
}

// Destroy stops f and cleans up its resources. f should not be used after it is destroyed.
func (f *LinearForce) Destroy() {
	f.accelTicker.Close()
}

func (f *LinearForce) accelTick(elapsed float32) {
	f.body.SetVelocity(f.body.GetVelocity().Add(f.body.GetRotation().Rotate(f.vector.Mul(elapsed))))
}

// Torque is an observer that applies a torque to a body
type Torque struct {
	torque mgl32.Quat
	body   *Body

	accelTicker *draw.Ticker
}

// NewTorque creates a new torque object, initially paused.
func NewTorque(b *Body, torque mgl32.Quat) *Torque {
	t := &Torque{
		torque: torque,
		body:   b,
	}
	t.accelTicker = draw.NewTicker(DefaultRefreshRate, t.accelTick)
	return t
}

// Start makes this force accelerate its body until `Pause` is called.
func (t *Torque) Start() {
	t.accelTicker.Start()
}

// Pause stops this force from accelerating its body, but continues applying its velocity.
func (t *Torque) Pause() {
	t.accelTicker.Stop()
}

func (t *Torque) accelTick(elapsed float32) {
	v := t.body.GetAngularV()
	t.body.SetAngularV(mgl32.QuatNlerp(v, v.Mul(t.torque), elapsed))
}

// Destroy stops t and cleans up its resources. t should not be used after it is destroyed.
func (t *Torque) Destroy() {
	t.accelTicker.Close()
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
