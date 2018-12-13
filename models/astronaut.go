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
}

func NewAstronaut(u *univ.Universe) *Astronaut {

	tex, err := draw.NewTexture("models/astronaut.png")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/astronaut.fbx", u.Window.GetBoneProgram(), []*draw.Texture{tex})
	if err != nil {
		log.Fatal(err)
	}

	return &Astronaut{
		Body: b,
		u:    u,
	}
}

func (m *Astronaut) Remove() {
	m.u.RemoveBody(m.Body)
}

func (m *Astronaut) StepForward() {
	f1, _ := os.Open("audio/walking.wav")
	s, _, _ := wav.Decode(f1)
	speaker.Play(s)
	rot := m.Body.GetRotation().Rotate(mgl32.Vec3{0, 0, 1})
	m.Translate(mgl32.Vec3{rot.X(), 0, rot.Z()})
}

func (m *Astronaut) StepBack() {
	rot := m.Body.GetRotation().Rotate(mgl32.Vec3{0, 0, -1})
	m.Translate(mgl32.Vec3{rot.X(), 0, rot.Z()})
}

func (m *Astronaut) StepLeft() {
	rot := m.Body.GetRotation().Rotate(mgl32.Vec3{1, 0, 0})
	m.Translate(mgl32.Vec3{rot.X(), 0, rot.Z()})
}

func (m *Astronaut) StepRight() {
	rot := m.Body.GetRotation().Rotate(mgl32.Vec3{-1, 0, 0})
	m.Translate(mgl32.Vec3{rot.X(), 0, rot.Z()})
}
