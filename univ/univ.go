package univ

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lsmith130/space/draw"
	assimp "github.com/tbogdala/assimp-go"
)

// DefaultRefreshRate is the default refresh rate
const DefaultRefreshRate = time.Millisecond * 16

// Universe is a group of Bodies drawn and updated together. It is the base object of
// the univ package, and all Bodies are created in a Universe.
type Universe struct {
	// bodies is a set of bodies
	bodies map[*Body]struct{}
	Window *draw.Window
}

// NewUniverse constructs a new empty Universe
func NewUniverse(window *draw.Window, updateRate time.Duration) *Universe {

	u := &Universe{
		bodies: make(map[*Body]struct{}),
		Window: window,
	}

	return u
}

// NewBody constructs a new body in u with a given model and shader
func (u *Universe) NewBody(modelPath string, program draw.Program, textures []*draw.Texture) (*Body, error) {

	meshes, err := assimp.ParseFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("load model %s: %v", modelPath, err)
	}
	if len(textures) > 0 && len(textures) != len(meshes) {
		return nil, fmt.Errorf("%d textures dosen't match %d meshes", len(textures), len(meshes))
	}

	body := &Body{
		meshes:    make([]*draw.Mesh, len(meshes)),
		animators: make([]*draw.Animator, len(meshes)),
		rotation:  mgl32.QuatIdent(),
		program:   program,
		observers: make(map[Observer]struct{}),
		angularV:  mgl32.QuatIdent(),
	}

	switch program := program.(type) {
	case *draw.BoneProgram:
		for i, mesh := range meshes {

			bones := make([]mgl32.Mat4, mesh.BoneCount)
			for _, bone := range mesh.Bones {
				bones[bone.Id] = mgl32.Ident4()
			}

			vertBones := make([]draw.VertBone, len(mesh.VertexWeightIds))
			for i, ids := range mesh.VertexWeightIds {
				vertBones[i] = draw.VertBone{int32(ids.X()), int32(ids.Y()), int32(ids.Z()), int32(ids.W())}
			}

			faces := *(*[]draw.MeshFace)(unsafe.Pointer(&mesh.Faces))
			body.meshes[i] = program.NewMesh(mesh.Vertices, faces, mesh.UVChannels[0], mesh.Normals, vertBones, mesh.VertexWeights, bones)
			body.meshes[i].SetTexture(textures[i])

			body.animators[i] = draw.NewAnimator(mesh.Bones, mesh.Animations, body.meshes[i])
		}
	case *draw.StandardProgram:
		for i, mesh := range meshes {
			faces := *(*[]draw.MeshFace)(unsafe.Pointer(&mesh.Faces))
			body.meshes[i] = program.NewMesh(mesh.Vertices, faces, mesh.UVChannels[0], mesh.Normals)
			body.meshes[i].SetTexture(textures[i])
		}
	}

	u.bodies[body] = struct{}{}
	body.ticker = draw.NewTicker(DefaultRefreshRate, body.velocityTick)
	return body, nil
}

// RemoveBody removes a body from u, such that it will no longer be drawn or recieve updates
func (u *Universe) RemoveBody(body *Body) {
	for _, mesh := range body.meshes {
		body.program.RemoveMesh(mesh)
	}
	delete(u.bodies, body)
}
