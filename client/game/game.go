package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/z46-dev/game-dev-project/client/shaders"
	"github.com/z46-dev/game-dev-project/shared"
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
			g.PlayerObject.velocity.Y -= 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.PlayerObject.velocity.Y += 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.PlayerObject.velocity.X -= 1
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.PlayerObject.velocity.X += 1
		}

		var _, wheelY float64 = ebiten.Wheel()
		g.PlayerObject.rotation += wheelY * .1
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
	var rectObj *GenericObject = newGenericObject(g).Spawn(util.Vector(0, 0))
	rectObj.size = 256
	rectObj.pushability = 0
	rectObj.polygon = util.NewPolygon([]*util.Vector2D{
		util.Vector(-1, 0.1),
		util.Vector(1, 0.1),
		util.Vector(1, -0.1),
		util.Vector(-1, -0.1),
	}, rectObj.position, rectObj.rotation, rectObj.size)
	rectObj.asset = shared.CreateAssetForPolygon(rectObj.polygon, 1024)
	g.genericObjects.Add(rectObj)

	g.PlayerObject = newGenericObject(g).SafelySpawn(func() *util.Vector2D {
		return util.RandomRadius(1024)
	}, 16)

	g.genericObjects.Add(g.PlayerObject)

	for i := 0; i < 64; i++ {
		g.genericObjects.Add(newGenericObject(g).SafelySpawn(func() *util.Vector2D {
			return util.RandomRadius(1024)
		}, 16))
	}
}
