package game

import (
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/util"
)

func NewHealth(health float64, canBeRepaired bool, rebuild *definitions.RebirthConfig) (hc *HealthComponent) {
	hc = &HealthComponent{
		MaxHealth:     health,
		Health:        health,
		CanBeRepaired: canBeRepaired,
		Rebuild:       rebuild,
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

func NewShip(g *Game, position *util.Vector2D, def *definitions.Ship) (s *Ship) {
	s = &Ship{}
	s.GenericObject = *NewGameObject(g, position)
	s.Cfg = def
	s.Name = s.Cfg.Name
	s.Size = s.Cfg.Size
	s.Health = NewHealth(s.Cfg.HullHealth, true, nil)
	s.Polygon = util.NewPolygon(s.Cfg.HullPath, s.Position, s.Size/2, s.Rotation)
	s.Control = NewControl(g, s)

	for _, shieldDef := range def.Shields {
		s.Shields = append(s.Shields, NewShieldGenerator(s, shieldDef))
	}

	for _, engineDef := range def.Engines {
		s.Engines = append(s.Engines, NewEngine(s, engineDef))
	}

	for _, turretDef := range def.TurretBanks {
		s.TurretBanks = append(s.TurretBanks, NewTurret(s, turretDef))
	}

	return
}

func (s *Ship) GetAABB() (aabb *util.AABB) {
	aabb = s.Polygon.AABB
	return
}

func (s *Ship) Insert() {
	s.Polygon.Transform(s.Position, s.Size/2, s.Rotation)
	s.Game.spatialHash.Insert(s)
}

func (s *Ship) Update() {
	s.Control.Update()
	s.Position.Add(s.Velocity)
	s.Velocity.Scale(s.Friction)
	s.Insert()

	for _, turret := range s.TurretBanks {
		turret.Update()
	}

	for _, shield := range s.Shields {
		shield.Update()
	}

	for _, engine := range s.Engines {
		engine.Update()
	}
}

func (s *Ship) Collide() {
	collideObjects(s.Game, s)
}
