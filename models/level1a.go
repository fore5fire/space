package models

import (
	"log"

	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Level1A struct {
	Body *univ.Body
	u    *univ.Universe
}

func NewLevel1A(u *univ.Universe) *Level1A {
	tex, err := draw.NewTexture("models/cement.jpg")
	if err != nil {
		log.Fatal(err)
	}

	metal, err := draw.NewTexture("models/level1a.png")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/game.dae", draw.ProgramTypeStandard, []*draw.Texture{tex, metal, metal, tex, metal})
	if err != nil {
		log.Fatal(err)
	}

	return &Level1A{
		Body: b,
		u:    u,
	}
}

func (l *Level1A) Remove() {
	l.u.RemoveBody(l.Body)
}
