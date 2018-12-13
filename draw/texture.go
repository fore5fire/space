package draw

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	ID uint32
}

func NewTexture(file string) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open %s: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture Texture
	gl.GenTextures(1, &texture.ID)
	texture.Use(gl.TEXTURE0)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	// When MAGnifying the image (no bigger mipmap available), use LINEAR filtering
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// When MINifying the image, use a LINEAR blend of two mipmaps, each filtered LINEARLY too
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	// Generate mipmaps, by the way.
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return &texture, nil
}

func (t *Texture) Use(textureSlot uint32) {
	gl.ActiveTexture(textureSlot)
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}
