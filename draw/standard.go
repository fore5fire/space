package draw

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type StandardProgram struct {
	ID           uint32
	ProjectionID int32
	ModelID      int32
	CameraID     int32
	TextureID    int32
	BonesID      int32

	TextureLocID  uint32
	VertexID      uint32
	NormalID      uint32
	VertBonesID   uint32
	VertWeightsID uint32

	viewMut       sync.Mutex
	view          mgl32.Mat4
	projectionMut sync.Mutex
	projection    mgl32.Mat4
	drawablesMut  sync.Mutex
	drawables     map[Drawable]struct{}
}

func newStandardProgram(vertShaderPath, fragShaderPath string) *StandardProgram {

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
		panic("compile vertex shader: " + err.Error())
	}

	fragmentShader, err := compileShader(string(fragSource)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic("compile fragment shader: " + err.Error())
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

	p := &StandardProgram{
		ID:           id,
		ProjectionID: gl.GetUniformLocation(id, gl.Str("projection\x00")),
		CameraID:     gl.GetUniformLocation(id, gl.Str("camera\x00")),
		ModelID:      gl.GetUniformLocation(id, gl.Str("model\x00")),
		TextureID:    gl.GetUniformLocation(id, gl.Str("tex\x00")),
		BonesID:      gl.GetUniformLocation(id, gl.Str("bones\x00")),

		VertexID:      uint32(gl.GetAttribLocation(id, gl.Str("vert\x00"))),
		TextureLocID:  uint32(gl.GetAttribLocation(id, gl.Str("vertTexCoord\x00"))),
		NormalID:      uint32(gl.GetAttribLocation(id, gl.Str("vertNormal\x00"))),
		VertBonesID:   uint32(gl.GetAttribLocation(id, gl.Str("vertBones\x00"))),
		VertWeightsID: uint32(gl.GetAttribLocation(id, gl.Str("vertWeights\x00"))),

		projection: mgl32.Ident4(),
		view:       mgl32.Ident4(),
		drawables:  make(map[Drawable]struct{}),
	}

	gl.Uniform1i(p.TextureID, 0)
	gl.BindFragDataLocation(id, 0, gl.Str("outputColor\x00"))

	return p
}

func (p *StandardProgram) use() {
	gl.UseProgram(p.ID)
}

func (p *StandardProgram) setView(view mgl32.Mat4) {
	p.viewMut.Lock()
	p.view = view
	p.viewMut.Unlock()
}

func (p *StandardProgram) setProjection(projection mgl32.Mat4) {
	p.projectionMut.Lock()
	p.projection = projection
	p.projectionMut.Unlock()
}

func (p *StandardProgram) AddDrawable(d Drawable) {
	p.drawablesMut.Lock()
	p.drawables[d] = struct{}{}
	p.drawablesMut.Unlock()
}

func (p *StandardProgram) RemoveDrawable(d Drawable) {
	p.drawablesMut.Lock()
	delete(p.drawables, d)
	p.drawablesMut.Unlock()
}

func (p *StandardProgram) Draw(state *GLState) {

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

func (p *StandardProgram) GetModelID() int32 {
	return p.ModelID
}

func (p *StandardProgram) NewMesh(vertexes []mgl32.Vec3, faces []MeshFace, uvCoords []mgl32.Vec2, normals []mgl32.Vec3) *Mesh {

	mesh := &Mesh{
		count:    int32(len(faces) * 3),
		program:  p,
		rotation: mgl32.QuatIdent(),

		// Save references to data so it won't get garbage-collected prematurely
		vertexes: vertexes,
		uvCoords: uvCoords,
		normals:  normals,
	}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	gl.GenBuffers(1, &mesh.vertexVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vertexVBO)
	// buffer type - length in bytes - data pointer - draw type
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexes)*3*4, gl.Ptr(vertexes), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(p.VertexID)
	// attribute id - data type - transpose - stride - offset
	gl.VertexAttribPointer(p.VertexID, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	if len(uvCoords) > 0 {
		gl.GenBuffers(1, &mesh.uvVBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.uvVBO)
		// buffer type - length in bytes - data pointer - draw type
		gl.BufferData(gl.ARRAY_BUFFER, len(uvCoords)*2*4, gl.Ptr(uvCoords), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(p.TextureLocID)
		// attribute id - data type - transpose - stride - offset
		gl.VertexAttribPointer(p.TextureLocID, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	}

	gl.GenBuffers(1, &mesh.indexBufferID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.indexBufferID)
	// buffer type - length in bytes - data pointer - draw type
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(faces)*3*4, gl.Ptr(faces), gl.STATIC_DRAW)

	gl.GenBuffers(1, &mesh.normalVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.normalVBO)
	// buffer type - length in bytes - data pointer - draw type
	gl.BufferData(gl.ARRAY_BUFFER, len(normals)*3*4, gl.Ptr(normals), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(p.NormalID)
	gl.VertexAttribPointer(p.NormalID, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	return mesh
}
