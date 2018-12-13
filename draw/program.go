package draw

import (
	"errors"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ProgramType is used to specify which shader program to use when rendering a body
type ProgramType int

const (
	// ProgramTypeStandard is the standard shader program
	ProgramTypeStandard ProgramType = iota
	// ProgramTypeBoned is a shader program for bone animated models
	ProgramTypeBoned = iota
)

type GLState struct {
}

// Program is a shader program.
type Program interface {
	setView(view mgl32.Mat4, camPosition mgl32.Vec3)
	setProjection(projection mgl32.Mat4)
	RemoveMesh(m *Mesh)
	Draw(state *GLState)
	GetModelID() int32
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
