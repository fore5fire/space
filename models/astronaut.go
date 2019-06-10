package models

import (
	"log"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Astronaut struct {
	*univ.Body
	u            *univ.Universe
	walkingSound beep.StreamSeekCloser

	forward, back, left, right, up, down *univ.Acceleration
	rightroll, leftroll                  *univ.Acceleration
}

func NewAstronaut(u *univ.Universe) *Astronaut {

	tex, err := draw.NewTexture("models/astronaut.png")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/astronaut3.dae", u.Window.GetBoneProgram(), []*draw.Texture{tex})
	if err != nil {
		log.Fatal(err)
	}

	a := &Astronaut{
		Body:      b,
		u:         u,
		forward:   univ.NewLinearAcceleration(b, mgl32.Vec3{0, 0, 20}),
		back:      univ.NewLinearAcceleration(b, mgl32.Vec3{0, 0, -20}),
		left:      univ.NewLinearAcceleration(b, mgl32.Vec3{20, 0, 0}),
		right:     univ.NewLinearAcceleration(b, mgl32.Vec3{-20, 0, 0}),
		up:        univ.NewLinearAcceleration(b, mgl32.Vec3{0, 20, 0}),
		down:      univ.NewLinearAcceleration(b, mgl32.Vec3{0, -20, 0}),
		rightroll: univ.NewAngularAcceleration(b, mgl32.Vec3{0, -1.5, 0}),
		leftroll:  univ.NewAngularAcceleration(b, mgl32.Vec3{0, 1.5, 0}),
	}

	a.forward.Pause()
	a.back.Pause()
	a.left.Pause()
	a.right.Pause()
	a.up.Pause()
	a.down.Pause()
	a.rightroll.Pause()
	a.leftroll.Pause()

	return a
}

func (m *Astronaut) Remove() {
	m.u.RemoveBody(m.Body)
}

func (m *Astronaut) SetForward(enable bool) {
	f1, _ := os.Open("audio/walking.wav")
	s, _, _ := wav.Decode(f1)
	speaker.Play(s)
	if enable {
		m.forward.Start()
	} else {
		m.forward.Pause()
	}
}

func (m *Astronaut) SetBack(enable bool) {
	if enable {
		m.back.Start()
	} else {
		m.back.Pause()
	}
}

func (m *Astronaut) SetLeft(enable bool) {
	if enable {
		m.left.Start()
	} else {
		m.left.Pause()
	}
}

func (m *Astronaut) SetRight(enable bool) {
	if enable {
		m.right.Start()
	} else {
		m.right.Pause()
	}
}

func (m *Astronaut) SetDown(enable bool) {
	if enable {
		m.down.Start()
	} else {
		m.down.Pause()
	}
}

func (m *Astronaut) SetUp(enable bool) {
	if enable {
		m.up.Start()
	} else {
		m.up.Pause()
	}
}

func (m *Astronaut) SetRollRight(enable bool) {
	if enable {
		m.rightroll.Start()
	} else {
		m.rightroll.Pause()
	}
}

func (m *Astronaut) SetRollLeft(enable bool) {
	if enable {
		m.leftroll.Start()
	} else {
		m.leftroll.Pause()
	}
}
