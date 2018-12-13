package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
	"github.com/lsmith130/space/models"
	"github.com/lsmith130/space/univ"
)

const windowWidth = 800
const windowHeight = 600

var u *univ.Universe

func init() {
	runtime.LockOSThread()
}

var bot *models.Robot
var cam *univ.ChaseCam
var man *models.Astronaut

func main() {
	window := draw.NewWindow(1000, 1000)
	u = univ.NewUniverse(window, time.Millisecond*10)
	man = models.NewAstronaut(u)
	man.SetLocation(mgl32.Vec3{27, 25, 119})
	man.SetRotation(mgl32.QuatRotate(20, mgl32.Vec3{0, 1, 0}))
	defer man.Remove()

	level1 := models.NewLevel1A(u)
	defer level1.Remove()

	cam = univ.NewChaseCam(man.Body)
	cam.SetLocation(mgl32.Vec3{0, 2, -10})
	defer cam.Remove()

	window.Loop(HandleKey, HandleMouseButton, HandleCursor)
}

func HandleKey(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) {
	switch key {
	case glfw.KeyLeft:
		man.Rotate(mgl32.QuatRotate(.1, mgl32.Vec3{0, 1, 0}))
	case glfw.KeyRight:
		man.Rotate(mgl32.QuatRotate(-.1, mgl32.Vec3{0, 1, 0}))
	case glfw.KeyUp:
		cam.Translate(mgl32.Vec3{0, 0.2, 0})
	case glfw.KeyDown:
		cam.Translate(mgl32.Vec3{0, -0.2, 0})

	case glfw.KeyA:
		man.StepLeft()
	case glfw.KeyD:
		man.StepRight()
	case glfw.KeyW:
		man.StepForward()
	case glfw.KeyS:
		man.StepBack()
	}
}
func HandleMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, modifier glfw.ModifierKey) {
	log.Println("Handle mouse button")
}
func HandleCursor(w *glfw.Window, xpos float64, ypos float64) {
	// cam.Translate(mgl32.)
}
