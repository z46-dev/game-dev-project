package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/z46-dev/game-dev-project/client/shaders"
	"github.com/z46-dev/game-dev-project/util"
)

func NewGame() (g *Game) {
	g = &Game{
		Camera:         newCamera(),
		genericObjects: newSafeStorage[*GenericObject](),
		spatialHash:    util.NewSpatialHash[*GenericObject](),
	}

	return
}

func (g *Game) next() (next uint64) {
	next = g.nextID
	g.nextID++
	return
}

func (g *Game) Update() (err error) {
	g.time++

	var width, height int = ebiten.WindowSize()
	g.Camera.Width, g.Camera.Height = float64(width), float64(height)

	if g.PlayerObject != nil {
		g.Camera.realPosition = g.PlayerObject.position
		g.Camera.realZoom = 128 / g.PlayerObject.size

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.PlayerObject.velocity.Y -= 0.5
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.PlayerObject.velocity.Y += 0.5
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.PlayerObject.velocity.X -= 0.5
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.PlayerObject.velocity.X += 0.5
		}
	}

	g.Camera.Update()

	g.spatialHash.Clear()

	g.genericObjects.Flush()

	g.genericObjects.ForEach(func(o *GenericObject) {
		o.Update()
	})

	g.genericObjects.ForEach(func(o *GenericObject) {
		o.Collide()
	})

	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	var bounds = screen.Bounds()

	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), shaders.BackgroundShader, &ebiten.DrawRectShaderOptions{
		GeoM: ebiten.GeoM{},
		Uniforms: map[string]any{
			"Time":       float32(g.time),
			"Camera":     []float32{float32(g.Camera.Position.X), float32(g.Camera.Position.Y), float32(g.Camera.Zoom)},
			"ScreenSize": []float32{float32(bounds.Dx()), float32(bounds.Dy())},
		},
	})

	g.genericObjects.ForEach(func(o *GenericObject) {
		o.Draw(screen)
	})
}

func (g *Game) Layout(_, _ int) (w, h int) {
	w, h = ebiten.WindowSize()
	return
}

func (g *Game) Init() {
	g.PlayerObject = newGenericObject(g, util.Vector(0, 0))
	g.genericObjects.Add(g.PlayerObject)

	for i := 0; i < 10; i++ {
		g.genericObjects.Add(newGenericObject(g, util.RandomRadius(1024)))
	}
}
