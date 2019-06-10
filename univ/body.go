package univ

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
)

// Body is the atomic unit of the univ package. A Body can be drawn in the universe at a location and orientation,
// and can be updated by observers
//
// All body functions are safe to use concurrently.
type Body struct {
	meshes  []*draw.Mesh
	program draw.Program
	// observerMut sync.RWMutex
	// observers   map[Observer]struct{}
	locMut   sync.RWMutex
	rotMut   sync.RWMutex
	location mgl32.Vec3
	rotation mgl32.Quat

	observerMut sync.RWMutex
	observers   map[Observer]struct{}

	velocityMut  sync.Mutex
	lastVelocity mgl32.Vec3
	velocity     mgl32.Vec3
	angularV     mgl32.Vec3
	lastAngularV mgl32.Vec3
	ticker       *draw.Ticker
	animators    []*draw.Animator
}

// AddObserver adds an observer to b. o.BodyUpdated will be called whenever b is updated.
// If an observer is added to a body that it is already observing, AddObserver has no effect.
func (b *Body) AddObserver(o Observer) {
	b.observerMut.Lock()
	b.observers[o] = struct{}{}
	b.observerMut.Unlock()
}

// RemoveObserver removes an observer from b, such that o no longer recieves updates from b.
// If an observer is removed from a body that it is not observing, RemoveObserver has no effect.
func (b *Body) RemoveObserver(o Observer) {
	b.observerMut.Lock()
	delete(b.observers, o)
	b.observerMut.Unlock()
}

// GetLocation returns the current location of b
func (b *Body) Location() mgl32.Vec3 {
	b.locMut.RLock()
	defer b.locMut.RUnlock()
	return b.location
}

// Translate translates the location of b by offset
func (b *Body) Translate(offset mgl32.Vec3) {
	b.locMut.Lock()

	b.location = b.location.Add(offset)
	for _, m := range b.meshes {
		m.SetLocation(b.location)
	}

	b.locMut.Unlock()
	b.notifyTranslation()
}

// SetLocation sets the location of b
func (b *Body) SetLocation(loc mgl32.Vec3) {
	b.locMut.Lock()
	b.location = loc
	for _, m := range b.meshes {
		m.SetLocation(loc)
	}
	b.locMut.Unlock()
	b.notifyTranslation()
}

// Rotation gets the rotation of b
func (b *Body) Rotation() mgl32.Quat {
	b.rotMut.RLock()
	defer b.rotMut.RUnlock()
	return b.rotation
}

// Rotate rotates b by offset
func (b *Body) Rotate(offset mgl32.Quat) {
	b.rotMut.Lock()
	b.rotation = b.rotation.Mul(offset)
	for _, m := range b.meshes {
		m.SetRotation(b.rotation)
	}
	b.rotMut.Unlock()
	b.notifyRotation()
}

// SetRotation sets the rotation of b to rot
func (b *Body) SetRotation(rot mgl32.Quat) {
	b.rotMut.Lock()
	b.rotation = rot
	for _, m := range b.meshes {
		m.SetRotation(rot)
	}
	b.rotMut.Unlock()
	b.notifyRotation()
}

// Draw draws b's meshes at the current location and rotation.
//
// Draw allows b to conform to draw.Drawable, and should not usually be called directly
func (b *Body) Draw(state *draw.GLState) {
	for _, mesh := range b.meshes {
		mesh.Draw(state)
	}
}

func (b *Body) notifyTranslation() {
	b.observerMut.Lock()
	observers := make([]Observer, len(b.observers))
	i := 0
	for o := range b.observers {
		observers[i] = o
		i++
	}
	b.observerMut.Unlock()

	for o := range b.observers {
		o.BodyTranslated(b)
	}
}

func (b *Body) notifyRotation() {
	b.observerMut.Lock()
	observers := make([]Observer, len(b.observers))
	i := 0
	for o := range b.observers {
		observers[i] = o
		i++
	}
	b.observerMut.Unlock()

	for o := range b.observers {
		o.BodyRotated(b)
	}
}

func (b *Body) SetVelocity(velocity mgl32.Vec3) {
	b.velocityMut.Lock()
	b.velocity = velocity
	b.velocityMut.Unlock()
}

func (b *Body) Velocity() mgl32.Vec3 {
	b.velocityMut.Lock()
	defer b.velocityMut.Unlock()
	return b.velocity
}

func (b *Body) AddVelocity(deltaV mgl32.Vec3) {
	b.velocityMut.Lock()
	defer b.velocityMut.Unlock()
	b.velocity = b.velocity.Add(deltaV)
}

func (b *Body) SetAngularV(angularV mgl32.Vec3) {
	b.velocityMut.Lock()
	b.angularV = angularV
	b.velocityMut.Unlock()
}

func (b *Body) AngularV() mgl32.Vec3 {
	b.velocityMut.Lock()
	defer b.velocityMut.Unlock()
	return b.angularV
}

func (b *Body) AddAngularV(deltaAngularV mgl32.Vec3) {
	b.velocityMut.Lock()
	defer b.velocityMut.Unlock()
	b.angularV = b.angularV.Add(deltaAngularV)
}

func (b *Body) velocityTick(elapsed float32) {
	b.velocityMut.Lock()
	defer b.velocityMut.Unlock()

	b.Translate(b.velocity.Add(b.lastVelocity).Mul(elapsed / 2))
	b.lastVelocity = b.velocity

	angularV := b.angularV.Add(b.lastAngularV).Mul(0.5)
	// TODO: verify the use of elapsed as the angle (assuming angularV includes the magnitude of rotation per second) correct.
	deltaRotation := mgl32.QuatRotate(elapsed, angularV)
	b.Rotate(deltaRotation)
	b.lastAngularV = b.angularV
}
