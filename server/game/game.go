package game

import (
	"math"

	"github.com/z46-dev/game-dev-project/util"
)

func genPolySides(n int) (sides []*util.Vector2D) {
	for i := range n {
		var angle float64 = 2 * math.Pi / float64(n) * float64(i)
		sides = append(sides, util.Vector(math.Cos(angle), math.Sin(angle)))
	}

	return
}

func genStarSides(n int, radMul float64) (sides []*util.Vector2D) {
	n *= 2
	for i := range n {
		var angle float64 = 2 * math.Pi / float64(n) * float64(i)
		var radius float64 = 1
		if i%2 == 0 {
			radius *= radMul
		}

		sides = append(sides, util.Vector(math.Cos(angle)*radius, math.Sin(angle)*radius))
	}

	return
}

var npcShapes [][]*util.Vector2D = [][]*util.Vector2D{
	genPolySides(3),
	genPolySides(4),
	genPolySides(5),
	genPolySides(6),
	genStarSides(3, .25),
	genStarSides(4, .5),
	genStarSides(5, .75),
}

func NewGame() (g *Game) {
	g = &Game{
		Ships:       util.NewSafeStorage[*Ship](),
		Projectiles: util.NewSafeStorage[*Projectile](),
		spatialHash: util.NewSpatialHash[CollidableObject](),
	}

	return
}

func NewGameObject(game *Game, position *util.Vector2D) (o *GenericObjectTemplate) {
	o = &GenericObjectTemplate{}
	o.ID = game.nextID
	game.nextID++
	o.Game = game
	o.Position = position
	o.Size = 32
	o.Rotation = 0
	return
}

func NewShip(g *Game, position *util.Vector2D) (s *Ship) {
	s = &Ship{}

	s.GenericObjectTemplate = *NewGameObject(g, position)
	s.Polygon = util.NewPolygon(genPolySides(4), s.Position, s.Size, s.Rotation)
	return
}

func NewProjectile(g *Game, position *util.Vector2D) (p *Projectile) {
	p = &Projectile{}
	p.GenericObjectTemplate = *NewGameObject(g, position)
	p.AABB = &util.AABB{}
	return
}
