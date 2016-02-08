package textureManager

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//TextureManager stores opengl textures
type TextureManager struct {
	textures    map[string]uint32
	textureLock sync.RWMutex
}

//NewTextureManager creates a new TextureManager
func NewTextureManager() *TextureManager {
	tm := TextureManager{textures: make(map[string]uint32)}
	return &tm
}

//LoadTexture loads a png file into an opengl texture
func (tm *TextureManager) LoadTexture(textureFile, key string) {
	texture, err := newTexture(textureFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	tm.textureLock.Lock()
	tm.textures[key] = texture
	tm.textureLock.Unlock()
}

//GetTexture returns a texture id if the texture was loaded, if it was not a 0 and
//false will be returned
func (tm *TextureManager) GetTexture(key string) (uint32, bool) {
	tm.textureLock.RLock()
	texture, status := tm.textures[key]
	tm.textureLock.RUnlock()
	return texture, status
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
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

	return texture, nil
}
