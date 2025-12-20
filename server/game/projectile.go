package game

import "github.com/z46-dev/game-dev-project/util"

func NewProjectile(g *Game, position *util.Vector2D) (p *Projectile) {
	p = &Projectile{}
	p.GenericObject = *NewGameObject(g, position)
	p.AABB = &util.AABB{}
	p.Size = 64
	p.Speed = 5
	p.Range = 300
	p.Friction = 1
	p.Damage = 1
	p.PrevPosition = position.Copy()
	return
}
