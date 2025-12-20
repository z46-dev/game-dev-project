package game

import (
	"math/rand/v2"

	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/util"
)

func NewHardpoint(parent *Ship, def definitions.Hardpoint) (hp *HardpointInstance) {
	hp = &HardpointInstance{
		Parent:           parent,
		Position:         util.Vector(0, 0),
		Size:             def.Size * parent.Size,
		FacingDir:        def.Direction,
		Health:           NewHealth(def.HullHealth, def.CanBeRepaired, def.Rebuild),
		RelativePosition: def.Position.Copy(),
	}

	hp.Position = parent.Position.Copy().Add(hp.RelativePosition.Copy().Rotate(parent.Rotation).Scale(parent.Size / 2))
	return
}

func (hp *HardpointInstance) GetID() (id uint64) {
	id = hp.Parent.ID
	return
}

func (hp *HardpointInstance) GetAABB() (aabb *util.AABB) {
	var halfSize float64 = hp.Size / 2
	aabb = &util.AABB{
		X1: hp.Position.X - halfSize,
		Y1: hp.Position.Y - halfSize,
		X2: hp.Position.X + halfSize,
		Y2: hp.Position.Y + halfSize,
	}

	return
}

func (hp *HardpointInstance) Update() (cont bool) {
	if !hp.Health.IsAlive() {
		cont = false
		return
	}

	hp.Position = hp.Parent.Position.Copy().Add(hp.RelativePosition.Copy().Rotate(hp.Parent.Rotation).Scale(hp.Parent.Size / 2))
	hp.Parent.Game.hardpointsSpatialHash.Insert(hp)
	cont = true
	return
}

// Turrets

func NewTurret(parent *Ship, def *definitions.Turret) (turret *TurretInstance) {
	turret = &TurretInstance{
		HardpointInstance: *NewHardpoint(parent, def.Hardpoint),
		Cfg:               def,
		Guns:              []*GunInstance{},
		Target:            nil,
		ReloadTick:        rand.IntN(def.Projectile.ReloadTicks),
	}

	for _, gunDef := range def.Guns {
		turret.Guns = append(turret.Guns, &GunInstance{
			RelativePosition: gunDef.RelativePosition.Copy(),
			Direction:        gunDef.Direction,
			RelLength:        gunDef.BarrelLength,
			RelWidth:         gunDef.BarrelWidth,
		})
	}

	return
}

func (t *TurretInstance) RealFacing() (angle float64) {
	angle = wrapAngle(t.FacingDir + t.Parent.Rotation)
	return
}

func (t *TurretInstance) Update() {
	if !t.HardpointInstance.Update() {
		return
	}

	var targetFacing float64 = t.Cfg.Direction
	if t.Target != nil {
		targetFacing = util.AngleBetween(t.Position, t.Target) - t.Parent.Rotation
	}

	t.FacingDir += min(t.Cfg.TraverseRate, max(-t.Cfg.TraverseRate, wrapAngle(targetFacing-t.FacingDir)))
	t.ReloadTick++

	if t.Target != nil && t.ReloadTick >= t.Cfg.Projectile.ReloadTicks {
		for _, gun := range t.Guns {
			var turretAngle float64 = t.RealFacing()
			var halfTurretSize float64 = t.Size / 2
			var gunPos *util.Vector2D = t.Position.Copy().Add(
				gun.RelativePosition.Copy().Rotate(turretAngle).Scale(halfTurretSize),
			)
			var barrelAngle float64 = wrapAngle(turretAngle + gun.Direction)
			var projectilePos *util.Vector2D = gunPos.Add(
				util.VectorFromAngle(barrelAngle, gun.RelLength*t.Size),
			)

			var p *Projectile = NewProjectile(t.Parent.Game, t.Parent, projectilePos)
			p.ProjectileID = t.Cfg.Projectile.ID
			p.Size = t.Size * (gun.RelWidth * 2)
			p.Speed = t.Cfg.Projectile.Speed
			p.Damage = t.Cfg.Projectile.ImpactDamage
			p.Range = t.Cfg.Projectile.Range
			p.Velocity = util.VectorFromAngle(barrelAngle, p.Speed)
			p.Rotation = barrelAngle
		}
		t.ReloadTick = 0
	}
}

// Engines

func NewEngine(parent *Ship, def *definitions.Engine) (engine *EngineInstance) {
	engine = &EngineInstance{
		HardpointInstance: *NewHardpoint(parent, def.Hardpoint),
	}

	return
}

// Shield Generators

func NewShieldGenerator(parent *Ship, def *definitions.ShieldGenerator) (sg *ShieldGenerator) {
	sg = &ShieldGenerator{
		HardpointInstance: *NewHardpoint(parent, def.Hardpoint),
		ShieldRadius:      def.ShieldRadius * parent.Size,
		ShieldHealth:      NewHealth(def.ShieldHealth, false, nil),
		ShieldRegen:       def.ShieldRegen,
	}

	return
}
