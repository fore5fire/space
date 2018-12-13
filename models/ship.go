package models

import (
	// "fmt"

	"log"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Ship struct {
	*univ.Body
	u      *univ.Universe
	ticker *univ.Ticker
}

func NewShip(u *univ.Universe) *Ship {

	body, err := draw.NewTexture("models/ship_body.png")
	if err != nil {
		log.Fatal(err)
	}
	// wings, err := draw.NewTexture("models/ship_wings.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	booster, err := draw.NewTexture("models/ship_boosters.png")
	if err != nil {
		log.Fatal(err)
	}

	b, err := u.NewBody("models/ship.dae", draw.ProgramTypeStandard, []*draw.Texture{booster, body})
	if err != nil {
		log.Fatal(err)
	}

	ship := &Ship{
		Body: b,
		u:    u,
	}

	ship.ticker = univ.NewTicker(univ.DefaultRefreshRate, ship.tick)
	ship.ticker.Start()

	return ship
}

func (ship *Ship) tick(elapsed float32) {
	offset := (float32((time.Now().UnixNano()/1000%3)-1) / 50.0) + 5
	loc := mgl32.Vec3{ship.Body.GetLocation().X(), float32(offset), ship.Body.GetLocation().Z()}
	ship.Body.SetLocation(loc)
}

func (r *Ship) Remove() {
	r.u.RemoveBody(r.Body)
}
