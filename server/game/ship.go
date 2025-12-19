package game

import (
	"math/rand/v2"

	"github.com/z46-dev/game-dev-project/util"
)

func NewShip(g *Game, position *util.Vector2D) (s *Ship) {
	s = &Ship{}

	s.GenericObject = *NewGameObject(g, position)
	s.Polygon = util.NewPolygon(potentialShapes[rand.IntN(len(potentialShapes))], s.Position, s.Size/2, s.Rotation)
	return
}
