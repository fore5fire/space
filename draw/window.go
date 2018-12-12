package draw

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	pause         chan bool
	close         chan struct{}
	width, height int
	window        *glfw.Window
	programs      map[ProgramType]*Program
}

func NewWindow(width, height int) *Window {

	if err := glfw.Init(); err != nil {
		log.Fatalf("failed to initialize glfw: %v", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Cube", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	w := &Window{
		pause:  make(chan bool),
		close:  make(chan struct{}),
		width:  width,
		height: height,
		window: window,
	}

	w.programs = map[ProgramType]*Program{
		ProgramTypeStandard: newProgram("shaders/shader.vert", "shaders/shader.frag"),
	}

	return w
}

func (w *Window) Start() {
	w.pause <- false
}

func (w *Window) Pause() {
	w.pause <- true
}

func (w *Window) Close() {
	defer recover()
	close(w.close)
}

func (w *Window) GetHeight() int {
	return w.height
}

func (w *Window) GetWidth() int {
	return w.width
}

func (w *Window) Loop(keyCallback glfw.KeyCallback, mouseButtonCallback glfw.MouseButtonCallback, cursorPosCallback glfw.CursorPosCallback) {

	defer glfw.Terminate()

	w.window.SetKeyCallback(keyCallback)
	w.window.SetMouseButtonCallback(mouseButtonCallback)
	w.window.SetCursorPosCallback(cursorPosCallback)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	var glState GLState

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(w.GetWidth())/float32(w.GetHeight()), 0.1, 10.0)

	for !w.window.ShouldClose() {

		w.waitIfPaused()
		if w.shouldClose() {
			break
		}

		// Clear buffer
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, p := range w.programs {
			p.SetProjection(projection)

			p.Draw(&glState)
		}

		// Maintenance
		w.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (w *Window) waitIfPaused() {
	select {
	case paused := <-w.pause:
		for paused {
			select {
			case paused = <-w.pause:
			case <-w.close:
				// Unpause if closed so we can exit
				break
			}
		}
	default:
	}
}

func (w *Window) shouldClose() bool {
	select {
	case <-w.close:
		return true
	default:
		return false
	}
}
