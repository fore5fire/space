package main

import (
	"log"

	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Astronaut struct {
	Body *univ.Body
}

func NewAstronaut() *Astronaut {

	tex, err := draw.NewTexture("models/pCylinder3Shape_color.gif")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/astronaut-animated.dae", draw.ProgramTypeStandard, []*draw.Texture{tex})
	if err != nil {
		log.Fatal(err)
	}

	return &Astronaut{
		Body: b,
	}
}

func (m *Astronaut) Remove() {
	u.RemoveBody(m.Body)
}
