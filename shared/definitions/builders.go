package definitions

import "github.com/z46-dev/game-dev-project/util"

// Projectile Builders

func NewProjectile(id ProjectileID, name string, speed, range_, damage float64, reloadTicks int) (p *Projectile) {
	p = &Projectile{
		ID:              id,
		Name:            name,
		Speed:           speed,
		Range:           range_,
		ImpactDamage:    damage,
		Explosion:       nil,
		MovementPattern: MovementTypeLinear,
		ReloadTicks:     reloadTicks,
	}

	ProjectileConfigs[id] = p
	return
}

func (pr *Projectile) SetExplosion(explodesWhenOutOfRange bool, radius, epicenterDamage float64, damageFallsOff bool) (p *Projectile) {
	pr.Explosion = &GenericExplosionCfg{
		ExplodesWhenOutOfRange: explodesWhenOutOfRange,
		Radius:                 radius,
		EpicenterDamage:        epicenterDamage,
		DamageFallsOff:         damageFallsOff,
	}

	p = pr
	return
}

func (pr *Projectile) SetSineMovement(amplitude, frequency float64) (p *Projectile) {
	pr.MovementPattern = MovementTypeSineWave
	pr.SineMovementAmplitude = amplitude
	pr.SineMovementFrequency = frequency

	p = pr
	return
}

func (pr *Projectile) SetHomingMovement(turnRate float64) (p *Projectile) {
	pr.MovementPattern = MovementTypeHoming
	pr.HomingMovementTurnRate = turnRate

	p = pr
	return
}

// Weapon Deployer Builder

func NewGun(relativePosition *util.Vector2D, direction float64, barrelLength, barrelWidth float64) (wd *Gun) {
	wd = &Gun{
		RelativePosition: relativePosition,
		Direction:        direction,
		BarrelLength:     barrelLength,
		BarrelWidth:      barrelWidth,
	}

	return
}

// Generic Hardpoint embedded funcs

func (hp *Hardpoint) SetRepairable(reparable bool) (h *Hardpoint) {
	hp.CanBeRepaired = reparable
	h = hp
	return
}

func (hp *Hardpoint) SetRebuildConfig(rebuildDelayTicks int, rebuildHealthRatio float64) (h *Hardpoint) {
	hp.Rebuild = &RebirthConfig{
		RebirthDelayTicks:  rebuildDelayTicks,
		RebirthHealthRatio: rebuildHealthRatio,
	}

	h = hp
	return
}

func (hp *Hardpoint) SetDrawLayer(layer HardpointDrawLayer) (h *Hardpoint) {
	hp.DrawLayer = layer
	h = hp
	return
}

// Turret Weapon Bank Builder

func NewTurretWeaponBank(position *util.Vector2D, size, direction, hullHealth float64, effectiveArc *util.Vector2D, traverseRate float64, projectile *Projectile) (twb *Turret) {
	twb = &Turret{}
	twb.Position = position
	twb.Size = size
	twb.Direction = direction
	twb.HullHealth = hullHealth
	twb.Guns = []*Gun{}
	twb.EffectiveArc = effectiveArc
	twb.TraverseRate = traverseRate
	twb.DrawLayer = HardpointDrawLayerAboveHull
	twb.Projectile = projectile
	twb.Multishot = 1
	twb.MultishotInterval = 0
	twb.Spread = util.Vector(0, 0)

	return
}

func (twb *Turret) AddWeapon(deployer ...*Gun) (t *Turret) {
	twb.Guns = append(twb.Guns, deployer...)
	t = twb
	return
}

// Shield Generator Builder

func NewShieldGenerator(position *util.Vector2D, size, direction, hullHealth, shieldRadius, shieldHealth, shieldRegen float64) (sg *ShieldGenerator) {
	sg = &ShieldGenerator{}
	sg.Position = position
	sg.Size = size
	sg.Direction = direction
	sg.HullHealth = hullHealth
	sg.ShieldRadius = shieldRadius
	sg.ShieldHealth = shieldHealth
	sg.ShieldRegen = shieldRegen
	sg.ShieldRebirth = nil
	sg.DrawLayer = HardpointDrawLayerAboveHull

	return
}

func (sg *ShieldGenerator) SetShieldRebirthConfig(rebirthDelayTicks int, rebirthHealthRatio float64) (s *ShieldGenerator) {
	sg.ShieldRebirth = &RebirthConfig{
		RebirthDelayTicks:  rebirthDelayTicks,
		RebirthHealthRatio: rebirthHealthRatio,
	}

	s = sg
	return
}

// Engine Builder

func NewEngine(position *util.Vector2D, size, direction, hullHealth float64) (e *Engine) {
	e = &Engine{}
	e.Position = position
	e.Size = size
	e.Direction = direction
	e.HullHealth = hullHealth
	e.DrawLayer = HardpointDrawLayerAboveHull

	return
}

// Ship Builder

func NewShip(id ShipID, name string, classification ShipClassification, hullPath []*util.Vector2D, size float64) (s *Ship) {
	s = &Ship{
		ID:             id,
		Name:           name,
		Classification: classification,
		HullPath:       hullPath,
		Size:           size,
		Shields:        []*ShieldGenerator{},
		Engines:        []*Engine{},
		TurretBanks:    []*Turret{},
	}

	ShipConfigs[id] = s
	return
}

func (s *Ship) SetHullProps(health, speed, turnSpeed float64) (sh *Ship) {
	s.HullHealth = health
	s.Speed = speed
	s.TurnSpeed = turnSpeed
	sh = s
	return
}

func (s *Ship) AddShieldGenerator(shield ...*ShieldGenerator) (sh *Ship) {
	s.Shields = append(s.Shields, shield...)
	sh = s
	return
}

func (s *Ship) AddEngine(engine ...*Engine) (sh *Ship) {
	s.Engines = append(s.Engines, engine...)
	sh = s
	return
}

func (s *Ship) AddTurretWeaponBank(bank ...*Turret) (sh *Ship) {
	s.TurretBanks = append(s.TurretBanks, bank...)
	sh = s
	return
}
