package draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	vao           uint32
	vertexes      []mgl32.Vec3
	vertexVBO     uint32
	indexBufferID uint32
	uvCoords      []mgl32.Vec2
	uvVBO         uint32
	normals       []mgl32.Vec3
	normalVBO     uint32
	boneIDs       []mgl32.Vec3
	boneIDsVBO    uint32
	bones         []mgl32.Mat4
	weights       []mgl32.Vec3
	weightsVBO    uint32
	rotation      mgl32.Quat
	position      mgl32.Vec3
	count         int32
	program       Program
	texture       *Texture
}

type MeshFace [3]uint32
type VertBone [4]int32

func (m *Mesh) SetTexture(texture *Texture) {
	m.texture = texture
}

func (m *Mesh) Draw(state *GLState) {

	transform := mgl32.Translate3D(m.position.Elem()).Mul4(m.rotation.Normalize().Mat4())

	// Set Model Transform
	gl.UniformMatrix4fv(m.program.GetModelID(), 1, false, &transform[0])

	if m.texture != nil {
		m.texture.Use(gl.TEXTURE0)
	}
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, m.count, gl.UNSIGNED_INT, nil)
}

func (m *Mesh) SetLocation(loc mgl32.Vec3) {
	m.position = loc
}

func (m *Mesh) SetRotation(rot mgl32.Quat) {
	m.rotation = rot
}
