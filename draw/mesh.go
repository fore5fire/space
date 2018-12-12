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
	rotation      mgl32.Quat
	position      mgl32.Vec3
	count         int32
	program       *Program
	texture       *Texture
}

type MeshFace [3]uint32

func (p *Program) NewMesh(vertexes []mgl32.Vec3, faces []MeshFace, uvCoords []mgl32.Vec2, normals []mgl32.Vec3) *Mesh {

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

func (m *Mesh) SetTexture(texture *Texture) {
	m.texture = texture
}

func (m *Mesh) Draw(state *GLState) {

	transform := m.rotation.Normalize().Mat4().Mul4(mgl32.Translate3D(m.position.Elem()))

	// Set Model Transform
	gl.UniformMatrix4fv(m.program.ModelID, 1, false, &transform[0])

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
