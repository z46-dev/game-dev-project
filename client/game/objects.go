package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/z46-dev/game-dev-project/shared"
	"github.com/z46-dev/game-dev-project/util"
	"golang.org/x/image/colornames"
)

var (
	shieldAsset *ebiten.Image = shared.CreateCircleAsset(6, colornames.Lightcyan)
	engineAsset *ebiten.Image = shared.CreateRoundedRectAsset(14, 10, 3, colornames.Orange, colornames.Sienna, 1.5)
	turretBase  *ebiten.Image = shared.CreateRoundedRectAsset(16, 16, 4, colornames.Lightgray, colornames.Dimgray, 1.2)
	whitePixel  *ebiten.Image = func() *ebiten.Image {
		img := ebiten.NewImage(1, 1)
		img.Fill(colornames.White)
		return img
	}()
)

func (s *ClientShip) Draw(game *Game, screen *ebiten.Image) {
	s.Position.X = util.Lerp(s.Position.X, s.RealPosition.X, .3)
	s.Position.Y = util.Lerp(s.Position.Y, s.RealPosition.Y, .3)
	s.Size = util.Lerp(s.Size, s.RealSize, .3)
	s.Rotation = util.LerpAngle(s.Rotation, s.RealRotation, .3)

	var bounds image.Rectangle = s.asset.Bounds()
	var dx, dy float64 = float64(bounds.Dx()), float64(bounds.Dy())
	var width, height float64 = s.Size / dx, s.Size / dy

	var options *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}

	// Object transformations
	options.GeoM.Translate(-dx/2, -dy/2)
	options.GeoM.Scale(width, height)
	options.GeoM.Rotate(s.Rotation)
	options.GeoM.Translate(s.Position.X, s.Position.Y)

	// Camera transformations
	options.GeoM.Scale(game.Camera.Zoom, game.Camera.Zoom)
	options.GeoM.Translate(game.Camera.Width/2, game.Camera.Height/2)
	options.GeoM.Translate(-game.Camera.Position.X*game.Camera.Zoom, -game.Camera.Position.Y*game.Camera.Zoom)

	// Graphical improvements
	options.Filter = ebiten.FilterLinear
	options.DisableMipmaps = false

	screen.DrawImage(s.asset, options)
	s.drawHardpoints(game, screen)
	s.drawHealthBar(game, screen)
}

func (s *ClientShip) drawHealthBar(game *Game, screen *ebiten.Image) {
	if s.HealthRatio <= 0 {
		return
	}

	var barWidth float32 = float32(s.Size * 0.8 * game.Camera.Zoom)
	var barHeight float32 = float32(6 * game.Camera.Zoom)
	var x float32 = float32((s.Position.X-game.Camera.Position.X)*game.Camera.Zoom + game.Camera.Width/2 - float64(barWidth)/2)
	var y float32 = float32((s.Position.Y-game.Camera.Position.Y)*game.Camera.Zoom + game.Camera.Height/2 - s.Size*0.6*game.Camera.Zoom)

	vector.FillRect(screen, x-1, y-1, barWidth+2, barHeight+2, colornames.Black, false)
	vector.FillRect(screen, x, y, barWidth*float32(s.HealthRatio), barHeight, colornames.Limegreen, false)
}

func (s *ClientShip) drawHardpoints(game *Game, screen *ebiten.Image) {
	if s.Definition == nil {
		return
	}

	drawPoint := func(asset *ebiten.Image, local *util.Vector2D, size float64, angle float64) {
		var world *util.Vector2D = local.Copy().Scale(s.Size / 2)
		world.Rotate(s.Rotation)
		world.Add(s.Position)

		var opts *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
		var bounds image.Rectangle = asset.Bounds()
		var dx, dy float64 = float64(bounds.Dx()), float64(bounds.Dy())
		var scale float64 = (s.Size * size) / math.Max(dx, dy)

		opts.GeoM.Translate(-dx/2, -dy/2)
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Rotate(angle)
		opts.GeoM.Translate(world.X, world.Y)
		opts.GeoM.Scale(game.Camera.Zoom, game.Camera.Zoom)
		opts.GeoM.Translate(game.Camera.Width/2, game.Camera.Height/2)
		opts.GeoM.Translate(-game.Camera.Position.X*game.Camera.Zoom, -game.Camera.Position.Y*game.Camera.Zoom)

		screen.DrawImage(asset, opts)
	}

	drawRect := func(local *util.Vector2D, width, height float64, angle float64, col color.Color) {
		var world *util.Vector2D = local.Copy().Scale(s.Size / 2)
		world.Rotate(s.Rotation)
		world.Add(s.Position)

		var opts *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
		opts.ColorScale.ScaleWithColor(col)
		opts.GeoM.Translate(-0.5, -0.5)
		opts.GeoM.Scale(width, height)
		opts.GeoM.Rotate(angle)
		opts.GeoM.Translate(world.X, world.Y)
		opts.GeoM.Scale(game.Camera.Zoom, game.Camera.Zoom)
		opts.GeoM.Translate(game.Camera.Width/2, game.Camera.Height/2)
		opts.GeoM.Translate(-game.Camera.Position.X*game.Camera.Zoom, -game.Camera.Position.Y*game.Camera.Zoom)

		screen.DrawImage(whitePixel, opts)
	}

	drawGun := func(local *util.Vector2D, width, length float64, angle float64, forward float64, col color.Color) {
		var offsetWorld *util.Vector2D = util.Vector(math.Cos(angle), math.Sin(angle)).Scale(forward)
		offsetWorld.Rotate(-s.Rotation)
		offsetWorld.Scale(2 / s.Size)
		drawRect(local.Copy().Add(offsetWorld), width, length, angle, col)
	}

	for _, shield := range s.Definition.Shields {
		drawPoint(shieldAsset, shield.Position, shield.Size, 0)
	}

	for _, engine := range s.Definition.Engines {
		drawPoint(engineAsset, engine.Position, engine.Size, engine.Direction+s.Rotation)
	}

	for i, turret := range s.Definition.TurretBanks {
		// var angle float64 = turret.Direction
		// if i < len(s.Turrets) {
		// 	angle = s.Hardpoints.TurretAngles[i]
		// }

		// for _, weapon := range turret.Weapons {
		// 	var w float64 = turret.Size * weapon.BarrelWidth * s.Size
		// 	var l float64 = turret.Size * weapon.BarrelLength * s.Size
		// 	drawGun(turret.Position, w, l, angle+s.Rotation+weapon.Direction, turret.Size*0.5*s.Size, colornames.Lightyellow)
		// }

		// drawPoint(turretBase, turret.Position, turret.Size, angle+s.Rotation)

		var angle float64 = s.Turrets[i][1] + s.Rotation

		for _, weapon := range turret.Guns {
			var w float64 = turret.Size * weapon.BarrelWidth * s.Size
			var l float64 = turret.Size * weapon.BarrelLength * s.Size
			drawGun(turret.Position, w, l, angle, turret.Size*0.5*s.Size, colornames.Lightyellow)
		}

		drawPoint(turretBase, turret.Position, turret.Size, angle)
	}
}

func (p *ClientProjectile) Draw(game *Game, screen *ebiten.Image) {
	p.Position.X = util.Lerp(p.Position.X, p.RealPosition.X, .4)
	p.Position.Y = util.Lerp(p.Position.Y, p.RealPosition.Y, .4)
	p.Size = util.Lerp(p.Size, p.RealSize, .4)
	p.Rotation = util.LerpAngle(p.Rotation, p.RealRotation, .4)

	var bounds image.Rectangle = p.asset.Bounds()
	var dx, dy float64 = float64(bounds.Dx()), float64(bounds.Dy())
	var width, height float64 = p.Size / dx, p.Size / dy

	var options *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}

	options.GeoM.Translate(-dx/2, -dy/2)
	options.GeoM.Scale(width, height)
	options.GeoM.Rotate(p.Rotation)
	options.GeoM.Translate(p.Position.X, p.Position.Y)

	options.GeoM.Scale(game.Camera.Zoom, game.Camera.Zoom)
	options.GeoM.Translate(game.Camera.Width/2, game.Camera.Height/2)
	options.GeoM.Translate(-game.Camera.Position.X*game.Camera.Zoom, -game.Camera.Position.Y*game.Camera.Zoom)

	options.Filter = ebiten.FilterLinear
	options.DisableMipmaps = false

	screen.DrawImage(p.asset, options)
}
