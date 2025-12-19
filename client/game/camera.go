package game

import (
	"math"

	"github.com/z46-dev/game-dev-project/util"
)

const MIN_DRAW_SIZE = 3

func newCamera() *PlayerCamera {
	return &PlayerCamera{
		Position:     util.Vector(0, 0),
		RealPosition: util.Vector(0, 0),
		Zoom:         1,
		RealZoom:     4,
		Width:        1000,
		Height:       1000,
	}
}

func (c *PlayerCamera) Update() {
	c.RealZoom = math.Max(math.Min(c.RealZoom, 5), .05)

	c.Position.X = util.Lerp(c.Position.X, c.RealPosition.X, .1)
	c.Position.Y = util.Lerp(c.Position.Y, c.RealPosition.Y, .1)
	c.Zoom = util.Lerp(c.Zoom, c.RealZoom, .1)
}

func (c *PlayerCamera) IsInView(position *util.Vector2D, radius float64) (inView bool) {
	if radius*c.Zoom < MIN_DRAW_SIZE {
		inView = false
	} else {
		var x, y float64 = position.X - c.Position.X, position.Y - c.Position.Y
		inView = x >= -c.Width/2/c.Zoom-radius && x <= c.Width/2/c.Zoom+radius && y >= -c.Height/2/c.Zoom-radius && y <= c.Height/2/c.Zoom+radius
	}

	return
}
