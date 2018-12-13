package draw

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/gombz"
)

type Animator struct {
	bones            []gombz.Bone
	animations       []gombz.Animation
	channels         map[int32]gombz.AnimationChannel
	currentAnimation int
	mut              sync.Mutex
	mesh             *Mesh
	ticker           *Ticker
	startTime        time.Time
}

func NewAnimator(bones []gombz.Bone, animations []gombz.Animation, mesh *Mesh) *Animator {
	a := &Animator{
		bones:      bones,
		mesh:       mesh,
		animations: animations,
		startTime:  time.Now(),
		channels:   make(map[int32]gombz.AnimationChannel, len(bones)),
	}
	log.Println(len(animations))

	for _, bone := range bones {
		for _, channel := range animations[a.currentAnimation].Channels {
			if channel.BoneId == bone.Id {
				a.channels[bone.Id] = channel
			}
		}
	}

	a.ticker = NewTicker(time.Millisecond*16, a.tick)

	return a
}

func (a *Animator) tick(elapsed float32) {
	anim := a.animations[a.currentAnimation]
	animTime := float64(time.Now().Sub(a.startTime)) / float64(time.Second)
	duration := anim.Duration / 4
	ticksPerSecond := float64(len(anim.Channels[0].PositionKeys)) / float64(duration)
	animOffset := float32(math.Mod(animTime, float64(duration)))
	current := animOffset * float32(ticksPerSecond)
	last := int(current)

	// next := int(current + 1)
	// interpolationFactor := current - float32(last)

	a.mesh.bonesMut.Lock()
	for _, bone := range a.bones {
		if bone.Parent == -1 {
			a.setBone(bone, last, anim)
		}
	}
	a.mesh.bonesMut.Unlock()
}

func (a *Animator) setBone(bone gombz.Bone, last int, anim gombz.Animation) {
	a.mesh.bones[bone.Id], _ = a.calcBone(bone, last, anim)
	for _, b := range a.bones {
		if b.Parent == bone.Id {
			a.setBone(b, last, anim)
		}
	}
}

func (a *Animator) calcBone(bone gombz.Bone, last int, anim gombz.Animation) (mgl32.Mat4, mgl32.Mat4) {

	parent := mgl32.Ident4()
	if bone.Parent >= 0 {
		_, parent = a.calcBone(a.bones[bone.Parent], last, anim)
	}

	local := bone.Transform
	channel, ok := a.channels[bone.Id]
	if ok {
		// rot := mgl32.QuatNlerp(channel.RotationKeys[last].Key, channel.RotationKeys[next].Key, interpolationFactor)
		// pos := channel.PositionKeys[last].Key.Mul(interpolationFactor).Add(channel.PositionKeys[next].Key.Mul(1 - interpolationFactor))
		// scale := channel.ScaleKeys[last].Key.Mul(interpolationFactor).Add(channel.ScaleKeys[next].Key.Mul(1 - interpolationFactor))
		rot := mgl32.QuatIdent()
		if len(channel.RotationKeys) > last {
			rot = channel.RotationKeys[last].Key
		} else {
			// log.Println(bone.Name, "missing rot key", last)
		}
		var pos, scale mgl32.Vec3
		if len(channel.PositionKeys) > last {
			pos = channel.PositionKeys[last].Key
		} else {
			// log.Println(bone.Name, "missing pos key", last)
		}
		if len(channel.ScaleKeys) > last {
			scale = channel.ScaleKeys[last].Key
		} else {
			// log.Println(bone.Name, "missing scale key", last)
		}

		local = mgl32.Translate3D(pos.Elem()).
			Mul4(rot.Mat4()).
			Mul4(mgl32.Scale3D(scale.Elem()))
	}

	global := parent.Mul4(local)

	return global.Mul4(bone.Offset), global

}
