package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/z46-dev/game-dev-project/assets"
	"github.com/z46-dev/game-dev-project/client/shaders"
	"github.com/z46-dev/game-dev-project/shared"
	"github.com/z46-dev/game-dev-project/util"
	"golang.org/x/image/colornames"
)

var (
	shieldAsset    *ebiten.Image = shared.CreateCircleAsset(6, colornames.Lightcyan)
	engineAsset    *ebiten.Image = shared.CreateRoundedRectAsset(14, 10, 3, colornames.Orange, colornames.Sienna, 1.5)
	turretBaseMask *ebiten.Image = shared.CreateAssetForPolygon(util.SVGPathToVector2DArray("M 1 -0.4 H 0.6 L 0.4 -0.8 H -0.8 L -1 -0.6 V 0.6 L -0.8 0.8 H 0.4 L 0.6 0.4 H 1 Z"), 128)
	turretBaseSets []*ebiten.Image = func() []*ebiten.Image {
		sheet := assets.NewSpriteSheet(assets.MustGet("assets/textures/turret-textures-2x2.png"), 2, 2)
		bounds := turretBaseMask.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()
		const turretTextureRepeats = 3
		tileWidth := width / turretTextureRepeats
		tileHeight := height / turretTextureRepeats
		if tileWidth < 1 {
			tileWidth = 1
		}
		if tileHeight < 1 {
			tileHeight = 1
		}
		sets := make([]*ebiten.Image, 0, len(sheet.Frames))
		for _, frame := range sheet.Frames {
			tile := ebiten.NewImage(width, height)
			tile.Fill(color.Transparent)
			frameBounds := frame.Bounds()
			scaleX := float64(tileWidth) / float64(frameBounds.Dx())
			scaleY := float64(tileHeight) / float64(frameBounds.Dy())
			stepsX := int(math.Ceil(float64(width) / float64(tileWidth)))
			stepsY := int(math.Ceil(float64(height) / float64(tileHeight)))
			for y := 0; y < stepsY; y++ {
				for x := 0; x < stepsX; x++ {
					tileOpts := &ebiten.DrawImageOptions{}
					tileOpts.GeoM.Scale(scaleX, scaleY)
					tileOpts.GeoM.Translate(float64(x*tileWidth), float64(y*tileHeight))
					tileOpts.Filter = ebiten.FilterLinear
					tile.DrawImage(frame, tileOpts)
				}
			}

			img := ebiten.NewImage(width, height)
			img.Fill(color.Transparent)
			opts := &ebiten.DrawRectShaderOptions{
				Images: [4]*ebiten.Image{turretBaseMask, tile},
				Uniforms: nil,
			}
			img.DrawRectShader(width, height, shaders.TurretTextureShader, opts)
			sets = append(sets, img)
		}
		return sets
	}()
	whitePixel *ebiten.Image = func() *ebiten.Image {
		img := ebiten.NewImage(1, 1)
		img.Fill(colornames.Beige)
		return img
	}()
)

func turretTextureIndex(shipID uint64, turretIndex int, count int) int {
	if count <= 0 {
		return 0
	}
	var v uint64 = shipID + uint64(turretIndex)*0x9e3779b97f4a7c15
	v ^= v >> 33
	v *= 0xff51afd7ed558ccd
	v ^= v >> 33
	return int(v % uint64(count))
}

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

	vector.FillRect(screen, x-1, y-1, barWidth+2, barHeight+2, colornames.Black, true)
	vector.FillRect(screen, x, y, barWidth*float32(s.HealthRatio), barHeight, colornames.Limegreen, true)
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
		opts.Filter = ebiten.FilterLinear

		screen.DrawImage(asset, opts)
	}

	for _, shield := range s.Definition.Shields {
		drawPoint(shieldAsset, shield.Position, shield.Size, 0)
	}

	for _, engine := range s.Definition.Engines {
		drawPoint(engineAsset, engine.Position, engine.Size, engine.Direction+s.Rotation)
	}

	var gunFillOpts *vector.FillOptions = &vector.FillOptions{}
	var gunDrawOpts *vector.DrawPathOptions = &vector.DrawPathOptions{AntiAlias: true}
	gunDrawOpts.ColorScale.ScaleWithColor(colornames.Dimgray)

	for i, turret := range s.Definition.TurretBanks {
		var angle float64 = s.Turrets[i][1] + s.Rotation
		var turretSize float64 = turret.Size * s.Size
		var halfTurretSize float64 = turretSize / 2
		var turretLocalX float64 = turret.Position.X * (s.Size / 2)
		var turretLocalY float64 = turret.Position.Y * (s.Size / 2)
		var cosShip float64 = math.Cos(s.Rotation)
		var sinShip float64 = math.Sin(s.Rotation)
		var turretWorldX float64 = turretLocalX*cosShip - turretLocalY*sinShip + s.Position.X
		var turretWorldY float64 = turretLocalX*sinShip + turretLocalY*cosShip + s.Position.Y

		for _, weapon := range turret.Guns {
			var gunLocalX float64 = weapon.RelativePosition.X * halfTurretSize
			var gunLocalY float64 = weapon.RelativePosition.Y * halfTurretSize
			var cosTurret float64 = math.Cos(angle)
			var sinTurret float64 = math.Sin(angle)
			var gunOffsetX float64 = gunLocalX*cosTurret - gunLocalY*sinTurret
			var gunOffsetY float64 = gunLocalX*sinTurret + gunLocalY*cosTurret
			var gunPosX float64 = turretWorldX + gunOffsetX
			var gunPosY float64 = turretWorldY + gunOffsetY

			var barrelAngle float64 = angle + weapon.Direction
			var cosBarrel float64 = math.Cos(barrelAngle)
			var sinBarrel float64 = math.Sin(barrelAngle)
			var barrelLen float64 = weapon.BarrelLength * turretSize
			var barrelEndX float64 = gunPosX + cosBarrel*barrelLen
			var barrelEndY float64 = gunPosY + sinBarrel*barrelLen
			var halfWidth float64 = weapon.BarrelWidth * turretSize / 2
			var perpX float64 = -sinBarrel * halfWidth
			var perpY float64 = cosBarrel * halfWidth

			var p1x float64 = gunPosX + perpX
			var p1y float64 = gunPosY + perpY
			var p2x float64 = barrelEndX + perpX
			var p2y float64 = barrelEndY + perpY
			var p3x float64 = barrelEndX - perpX
			var p3y float64 = barrelEndY - perpY
			var p4x float64 = gunPosX - perpX
			var p4y float64 = gunPosY - perpY

			var z float64 = game.Camera.Zoom
			var cx float64 = game.Camera.Width / 2
			var cy float64 = game.Camera.Height / 2
			var camX float64 = game.Camera.Position.X
			var camY float64 = game.Camera.Position.Y

			var path vector.Path
			path.MoveTo(float32((p1x-camX)*z+cx), float32((p1y-camY)*z+cy))
			path.LineTo(float32((p2x-camX)*z+cx), float32((p2y-camY)*z+cy))
			path.LineTo(float32((p3x-camX)*z+cx), float32((p3y-camY)*z+cy))
			path.LineTo(float32((p4x-camX)*z+cx), float32((p4y-camY)*z+cy))
			path.Close()

			vector.FillPath(screen, &path, gunFillOpts, gunDrawOpts)
		}

		var turretAsset *ebiten.Image = turretBaseMask
		if len(turretBaseSets) > 0 {
			turretAsset = turretBaseSets[turretTextureIndex(s.ID, i, len(turretBaseSets))]
		}
		drawPoint(turretAsset, turret.Position, turret.Size, angle)
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
