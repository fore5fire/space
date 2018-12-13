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

func main() {
	window := draw.NewWindow(1000, 1000)
	u = univ.NewUniverse(window, time.Millisecond*10)
	// man := NewAstronaut()
	// defer man.Remove()

	// level1 := NewLevel1A()
	// defer level1.Remove()

	bot = models.NewRobot(u)
	defer bot.Remove()

	cam = univ.NewChaseCam(bot.Body, u.Window)
	cam.SetLocation(mgl32.Vec3{5, 5, 5})

	// force := univ.NewLinearForce(bot2.Body, mgl32.Vec3{0.5, 0.5, 0.5})
	// defer force.Destroy()
	// force.Start()

	// torque := univ.NewTorque(bot.Body, mgl32.QuatRotate(3, mgl32.Vec3{0, 1, 0}))
	// defer torque.Destroy()
	// bot.Rotate(mgl32.QuatRotate(1, mgl32.Vec3{0, 1, 0}))

	window.Loop(HandleKey, HandleMouseButton, HandleCursor)
}

func HandleKey(w *glfw.Window, key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) {
	switch key {
	case glfw.KeyLeft:
		cam.Translate(mgl32.Vec3{0.1, 0, 0})
	case glfw.KeyRight:
		cam.Translate(mgl32.Vec3{-0.1, 0, 0})
	case glfw.KeyUp:
		cam.Translate(mgl32.Vec3{0, 0, 0.1})
	case glfw.KeyDown:
		cam.Translate(mgl32.Vec3{0, 0, -0.1})

	case glfw.KeyA:
		bot.Translate(mgl32.Vec3{0.1, 0, 0})
	case glfw.KeyD:
		bot.Translate(mgl32.Vec3{-0.1, 0, 0})
	case glfw.KeyW:
		bot.Translate(mgl32.Vec3{0, 0, 0.1})
	case glfw.KeyS:
		bot.Translate(mgl32.Vec3{0, 0, -0.1})
	}
}
func HandleMouseButton(w *glfw.Window, button glfw.MouseButton, action glfw.Action, modifier glfw.ModifierKey) {
	log.Println("Handle mouse button")
}
func HandleCursor(w *glfw.Window, xpos float64, ypos float64) {
	// cam.Translate(mgl32.)
}
