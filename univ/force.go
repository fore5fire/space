package univ

import (
	"log"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

// LinearForce is an observer that only applies the translation
// portion of a force to a body.
type LinearForce struct {
	vector       mgl32.Vec3
	velocityMut  sync.Mutex
	lastVelocity mgl32.Vec3
	velocity     mgl32.Vec3
	body         *Body

	velocityTicker *Ticker
	accelTicker    *Ticker
}

// NewLinearForce creates a new LinearForce.
func NewLinearForce(b *Body, vector mgl32.Vec3) *LinearForce {
	f := &LinearForce{
		vector: vector,
		body:   b,
	}
	f.accelTicker = NewTicker(DefaultRefreshRate, f.accelTick)
	f.velocityTicker = NewTicker(DefaultRefreshRate, f.velocityTick)
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
	f.velocityTicker.Close()
	f.accelTicker.Close()
}

func (f *LinearForce) velocityTick(elapsed float32) {
	f.velocityMut.Lock()
	f.body.Translate(f.velocity.Add(f.lastVelocity).Mul(elapsed / 2))
	f.lastVelocity = f.velocity
	f.velocityMut.Unlock()
}

func (f *LinearForce) accelTick(elapsed float32) {
	f.velocityMut.Lock()
	f.velocity = f.velocity.Add(f.vector.Mul(elapsed))
	f.velocityMut.Unlock()
}

// Torque is an observer that applies a torque to a body
type Torque struct {
	torque       mgl32.Quat
	velocityMut  sync.Mutex
	angularV     mgl32.Quat
	lastAngularV mgl32.Quat
	body         *Body

	velocityTicker *Ticker
	accelTicker    *Ticker
}

// NewTorque creates a new torque object, initially paused.
func NewTorque(b *Body, torque mgl32.Quat) *Torque {
	t := &Torque{
		torque:   torque,
		angularV: mgl32.QuatIdent(),
		body:     b,
	}
	t.velocityTicker = NewTicker(DefaultRefreshRate, t.velocityTick)
	t.accelTicker = NewTicker(DefaultRefreshRate, t.accelTick)
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

func (t *Torque) velocityTick(elapsed float32) {
	t.velocityMut.Lock()
	angularV := mgl32.QuatNlerp(t.angularV, t.lastAngularV, 0.5)
	t.body.Rotate(mgl32.QuatNlerp(mgl32.QuatIdent(), angularV, elapsed))
	t.lastAngularV = t.angularV
	t.velocityMut.Unlock()
}

func (t *Torque) accelTick(elapsed float32) {
	log.Println(elapsed)
	t.velocityMut.Lock()
	t.angularV = mgl32.QuatNlerp(t.angularV, t.angularV.Mul(t.torque), elapsed)
	t.velocityMut.Unlock()
}

// Destroy stops t and cleans up its resources. t should not be used after it is destroyed.
func (t *Torque) Destroy() {
	t.velocityTicker.Close()
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
