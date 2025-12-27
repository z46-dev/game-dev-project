package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/golog"
)

func InitShader(input []byte) (shader *ebiten.Shader) {
	var err error
	if shader, err = ebiten.NewShader(input); err != nil {
		panic(err)
	}

	return
}

var (
	//go:embed star.kage.go
	starShaderSource []byte
	starShader       *ebiten.Shader = InitShader(starShaderSource)
)

type Game struct {
	Time int
}

func (g *Game) Update() error {
	g.Time++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var bounds = screen.Bounds()

	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), starShader, &ebiten.DrawRectShaderOptions{
		GeoM: ebiten.GeoM{},
		Uniforms: map[string]any{
			"Time":          float32(g.Time) * 0.01,
			"Camera":        []float32{0.0, 0.0, 1.0},
			"ScreenSize":    []float32{float32(bounds.Dx()), float32(bounds.Dy())},
			"StarCenter":    []float32{float32(bounds.Dx()) / 2.0, float32(bounds.Dy()) / 2.0},
			"StarRadius":    float32(220.0),
			"StarIntensity": float32(1),
			"StarColor":     []float32{0.3, 0.6, 1.0},
			"StarPulse":     float32(5),
			"StarDetail":    float32(2),
		},
	})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ebiten.WindowSize()
}

var (
	log *golog.Logger = golog.New().Prefix("[MAIN]", golog.BoldBlue).Timestamp()
	err error
)

func main() {
	log.Info("Starting...")

	ebiten.SetWindowTitle("CS780 Project")
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &Game{}

	if err = ebiten.RunGame(g); err != nil {
		log.Panicf("Error running game: %v", err)
	}
}
