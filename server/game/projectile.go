package game

import "github.com/z46-dev/game-dev-project/util"

func NewProjectile(g *Game, position *util.Vector2D) (p *Projectile) {
	p = &Projectile{}
	p.GenericObject = *NewGameObject(g, position)
	p.AABB = &util.AABB{}
	return
}
