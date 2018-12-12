package draw

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ProgramType is used to specify which shader program to use when rendering a body
type ProgramType int

const (
	// ProgramTypeStandard is the standard shader
	ProgramTypeStandard ProgramType = iota
)

// GetProgram returns the program of the requested type
func (w *Window) GetProgram(ptype ProgramType) *Program {
	return w.programs[ptype]
}

type Program struct {
	ID           uint32
	ProjectionID int32
	ModelID      int32
	CameraID     int32
	TextureID    int32
	TextureLocID uint32
	VertexID     uint32

	viewMut       sync.Mutex
	view          mgl32.Mat4
	projectionMut sync.Mutex
	projection    mgl32.Mat4
	drawablesMut  sync.Mutex
	drawables     map[Drawable]struct{}
}

func newProgram(vertShaderPath, fragShaderPath string) *Program {

	vertSource, err := ioutil.ReadFile(vertShaderPath)
	if err != nil {
		panic(fmt.Sprintf("open vertex shader %s: %v", vertShaderPath, err))
	}

	fragSource, err := ioutil.ReadFile(fragShaderPath)
	if err != nil {

		panic(fmt.Sprintf("open fragment shader %s: %v", fragShaderPath, err))
	}

	vertexShader, err := compileShader(string(vertSource)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err.Error())
	}

	fragmentShader, err := compileShader(string(fragSource)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err.Error())
	}

	id := gl.CreateProgram()

	gl.AttachShader(id, vertexShader)
	gl.AttachShader(id, fragmentShader)
	gl.LinkProgram(id)

	var status int32
	gl.GetProgramiv(id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(id, logLength, nil, gl.Str(log))

		panic(fmt.Sprintf("link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	p := &Program{
		ID:           id,
		ProjectionID: gl.GetUniformLocation(id, gl.Str("projection\x00")),
		CameraID:     gl.GetUniformLocation(id, gl.Str("camera\x00")),
		ModelID:      gl.GetUniformLocation(id, gl.Str("model\x00")),
		TextureID:    gl.GetUniformLocation(id, gl.Str("tex\x00")),
		VertexID:     uint32(gl.GetAttribLocation(id, gl.Str("vert\x00"))),
		TextureLocID: uint32(gl.GetAttribLocation(id, gl.Str("vertTexCoord\x00"))),

		projection: mgl32.Ident4(),
		view:       mgl32.Ident4(),
		drawables:  make(map[Drawable]struct{}),
	}

	gl.Uniform1i(p.TextureID, 0)
	gl.BindFragDataLocation(id, 0, gl.Str("outputColor\x00"))

	return p
}

func (p *Program) use() {
	gl.UseProgram(p.ID)
}

func (p *Program) SetView(view mgl32.Mat4) {
	p.viewMut.Lock()
	p.view = view
	p.viewMut.Unlock()
}

func (p *Program) SetProjection(projection mgl32.Mat4) {
	p.projectionMut.Lock()
	p.projection = projection
	p.projectionMut.Unlock()
}

func (p *Program) AddDrawable(d Drawable) {
	p.drawablesMut.Lock()
	p.drawables[d] = struct{}{}
	p.drawablesMut.Unlock()
}

func (p *Program) RemoveDrawable(d Drawable) {
	p.drawablesMut.Lock()
	delete(p.drawables, d)
	p.drawablesMut.Unlock()
}

func (p *Program) Draw(state *GLState) {

	p.use()

	p.viewMut.Lock()
	gl.UniformMatrix4fv(p.CameraID, 1, false, &p.view[0])
	p.viewMut.Unlock()

	p.projectionMut.Lock()
	gl.UniformMatrix4fv(p.ProjectionID, 1, false, &p.projection[0])
	p.projectionMut.Unlock()

	p.drawablesMut.Lock()
	for drawable := range p.drawables {
		drawable.Draw(state)
	}
	p.drawablesMut.Unlock()
}

func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	defer free()
	gl.ShaderSource(shader, 1, csources, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logs))

		return 0, errors.New(logs)
	}

	return shader, nil
}
