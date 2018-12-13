package models

import (
	"log"

	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Astronaut struct {
	Body *univ.Body
	u    *univ.Universe
}

func NewAstronaut(u *univ.Universe) *Astronaut {

	tex, err := draw.NewTexture("models/pCylinder3Shape_color.gif")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/astronaut-animated.dae", u.Window.GetBoneProgram(), []*draw.Texture{tex})
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
