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
	shieldAsset    *ebiten.Image   = shared.CreateCircleAsset(6, colornames.Lightcyan)
	engineAsset    *ebiten.Image   = shared.CreateRoundedRectAsset(14, 10, 3, colornames.Orange, colornames.Sienna, 1.5)
	turretBaseMask *ebiten.Image   = shared.CreateAssetForPolygon(util.SVGPathToVector2DArray("M 1 -0.4 H 0.6 L 0.4 -0.8 H -0.8 L -1 -0.6 V 0.6 L -0.8 0.8 H 0.4 L 0.6 0.4 H 1 Z"), 128)
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
				Images:   [4]*ebiten.Image{turretBaseMask, tile},
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
	colossus *ebiten.Image = assets.MustGet("assets/ships/colossus.png")
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