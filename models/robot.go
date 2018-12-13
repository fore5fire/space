package models

import (
	"log"

	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Robot struct {
	*univ.Body
	u *univ.Universe
}

func NewRobot(u *univ.Universe) *Robot {

	head, err := draw.NewTexture("models/Material Diffuse Color.png")
	if err != nil {
		log.Fatal(err)
	}
	body, err := draw.NewTexture("models/Material.001 Diffuse Color.png")
	if err != nil {
		log.Fatal(err)
	}
	leftfoot, err := draw.NewTexture("models/Material.002 Diffuse Color.png")
	if err != nil {
		log.Fatal(err)
	}
	rightfoot, err := draw.NewTexture("models/Material.003 Diffuse Color.png")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/robot.dae", u.Window.GetStandardProgram(), []*draw.Texture{head, head, leftfoot, rightfoot, body})
	if err != nil {
		log.Fatal(err)
	}

	return &Robot{
		Body: b,
		u:    u,
	}
}

func (r *Robot) Remove() {
	r.u.RemoveBody(r.Body)
}
