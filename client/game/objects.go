package game

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/util"
)

func (s *ClientShip) Draw(game *Game, screen *ebiten.Image) {
	s.Position.X = util.Lerp(s.Position.X, s.RealPosition.X, .1)
	s.Position.Y = util.Lerp(s.Position.Y, s.RealPosition.Y, .1)
	s.Size = util.Lerp(s.Size, s.RealSize, .1)
	s.Rotation = util.LerpAngle(s.Rotation, s.RealRotation, .1)

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
}
