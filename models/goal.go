package models

import (
	// "fmt"

	"fmt"
	"log"
	"os"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/univ"
)

type Goal struct {
	*univ.Body
	u      *univ.Universe
	target *univ.Body
	ticker *draw.Ticker
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

func (goal *Goal) Pickup(t *univ.Body) {

	if t.Location().Sub(goal.Body.Location()).Len() < 10 {
		if goal.target == nil {
			f1, _ := os.Open("audio/pickup.wav")
			s, _, _ := wav.Decode(f1)
			speaker.Play(s)

			fmt.Println("Pick up")
			goal.target = t
			t.AddObserver(goal)
			goal.update()
		} else {
			f1, _ := os.Open("audio/putdown.wav")
			s, _, _ := wav.Decode(f1)
			speaker.Play(s)

			fmt.Println("Set down")
			goal.target.RemoveObserver(goal)
			goal.target = nil
		}
	} else {
		fmt.Println("Can't pick up")
	}
}

func (r *Goal) Remove() {
	r.u.RemoveBody(r.Body)
}

// BodyTranslated conforms to Observer.BodyTranslated and should not be called directly
func (goal *Goal) BodyTranslated(b *univ.Body) {
	goal.update()
}

// BodyRotated conforms to Observer.BodyRotated and should not be called directly
func (goal *Goal) BodyRotated(b *univ.Body) {
	goal.update()
}

func (goal *Goal) update() {
	rot := goal.target.Rotation()
	posMat := mgl32.Translate3D(goal.target.Location().Elem())
	posRot := posMat.Mul4(rot.Normalize().Mat4())
	pos := posRot.Mul4(mgl32.Translate3D(0.0, 0.0, -1.0)).Col(3).Vec3()

	goal.Body.SetRotation(rot)
	goal.Body.SetLocation(pos)
}
