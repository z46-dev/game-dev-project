package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/z46-dev/game-dev-project/util"
)

func TestEmbedding(t *testing.T) {
	var g *Game = NewGame()

	var s *Ship = NewShip(g, util.Vector(0, 0))
	var p *Projectile = NewProjectile(g, util.Vector(1, 1))

	// Quick asserts
	assert.Equal(t, g, s.Game)
	assert.Equal(t, g, p.Game)

	assert.Equal(t, s.Position, util.Vector(0, 0))
	s.Position.X = 1
	assert.Equal(t, s.Position, util.Vector(1, 0))

	g.Ships.Add(s)
	g.Projectiles.Add(p)

	s.Insert()
	p.Insert()

	s.Collide()
	p.Collide()
}
