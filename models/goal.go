package models

import (
	// "fmt"

	"fmt"
	"log"

	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Goal struct {
	*univ.Body
	u      *univ.Universe
	ticker *univ.Ticker
}

func NewGoal(u *univ.Universe) *Goal {

	goal, err := draw.NewTexture("models/goal.png")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/goal.dae", u.Window.GetStandardProgram(), []*draw.Texture{goal})
	if err != nil {
		log.Fatal(err)
	}

	return &Goal{
		Body: b,
		u:    u,
	}

}

func (r *Goal) Pickup(target *univ.Body) {
	fmt.Println(target.GetLocation())
	fmt.Println(r.Body.GetLocation())
}

func (r *Goal) Remove() {
	r.u.RemoveBody(r.Body)
}
