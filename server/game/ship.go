package game

import (
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/util"
)

func NewHealth(health float64, canBeRepaired bool) (hc *HealthComponent) {
	hc = &HealthComponent{
		MaxHealth:     health,
		Health:        health,
		CanBeRepaired: canBeRepaired,
	}
	return
}

func (hc *HealthComponent) IsAlive() (alive bool) {
	alive = hc.Health > 0
	return
}

func (hc *HealthComponent) Damage(amount float64) {
	hc.Health = max(0, min(hc.MaxHealth, hc.Health-amount))
}

func (hc *HealthComponent) Ratio() (ratio float64) {
	if hc.MaxHealth == 0 {
		ratio = 0
		return
	}

	ratio = hc.Health / hc.MaxHealth
	return
}

func NewShip(g *Game, position *util.Vector2D, def *definitions.Ship, f *Faction) (s *Ship) {
	s = &Ship{}
	s.GenericObject = *NewGameObject(g, position, f)
	s.Cfg = def
	s.Name = s.Cfg.Name
	s.Size = s.Cfg.Size
	s.Health = NewHealth(s.Cfg.HullHealth, true)
	s.Polygon = util.NewPolygon(s.Cfg.HullPath, s.Position, s.Size/2, s.Rotation)
	s.Control = NewControl(g, s)

	return
}

func (s *Ship) GetAABB() (aabb *util.AABB) {
	aabb = s.Polygon.AABB
	return
}

func (s *Ship) Insert() {
	s.Polygon.Transform(s.Position, s.Size/2, s.Rotation)
	s.Game.spatialHash.Insert(s)
	s.Faction.ShipsSpatialHash.Insert(s)
}

func (s *Ship) Update() {
	s.Control.Update()
	s.Position.Add(s.Velocity)
	s.Velocity.Scale(s.Friction)
	s.Insert()
}

func (s *Ship) Collide() {
	collideObjects(s.Game, s)
}

func (s *Ship) Think() {}
