package main

import (
	"log"

	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Robot struct {
	*univ.Body
}

func NewRobot() *Robot {

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

	b, err := u.NewBody("models/robot.dae", draw.ProgramTypeStandard, []*draw.Texture{head, head, leftfoot, rightfoot, body})
	if err != nil {
		log.Fatal(err)
	}

	return &Robot{
		Body: b,
	}
}

func (r *Robot) Remove() {
	u.RemoveBody(r.Body)
}
