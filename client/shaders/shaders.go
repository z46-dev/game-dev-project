package shaders

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

func InitShader(input []byte) (shader *ebiten.Shader) {
	var err error
	if shader, err = ebiten.NewShader(input); err != nil {
		panic(err)
	}

	return
}

var (
	//go:embed background.kage.go
	backgroundShaderSource []byte
	BackgroundShader       *ebiten.Shader = InitShader(backgroundShaderSource)

	//go:embed turret_texture.kage.go
	turretTextureShaderSource []byte
	TurretTextureShader       *ebiten.Shader = InitShader(turretTextureShaderSource)
)
