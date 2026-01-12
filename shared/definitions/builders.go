package definitions

import "github.com/z46-dev/game-dev-project/util"

// Ship Builder

func NewShip(id ShipID, name string, classification ShipClassification, hullPath []*util.Vector2D, size float64, assetName string) (s *Ship) {
	s = &Ship{
		ID:             id,
		Name:           name,
		Classification: classification,
		HullPath:       hullPath,
		Size:           size,
		AssetName:      assetName,
	}

	ShipConfigs[id] = s
	return
}

func (s *Ship) SetHullProps(health, speed, turnSpeed float64) (sh *Ship) {
	s.HullHealth = health
	s.Speed = speed / 2
	s.TurnSpeed = turnSpeed / 150000
	sh = s
	return
}

// Plane Builder

func NewPlane(id PlaneID, name string, size float64, assetName string) (p *Plane) {
	p = &Plane{
		ID:        id,
		Name:      name,
		AssetName: assetName,
		Size:      size,
	}

	PlaneConfigs[id] = p
	return
}

func (p *Plane) SetFlightProps(health, speed, turnSpeed float64) (pl *Plane) {
	p.Health = health
	p.Speed = speed
	p.TurnSpeed = turnSpeed
	pl = p
	return
}

// Ammo Builders

func NewDamageSource(fullDamage, penetration, fireChance float64) DamageSource {
	return DamageSource{
		FullDamage:  fullDamage,
		Penetration: penetration,
		FireChance:  fireChance,
	}
}

func NewEllipticalReticle(distance, width, height float64) EllipticalReticle {
	return EllipticalReticle{
		Distance: distance,
		Width:    width,
		Height:   height,
	}
}

func NewConeReticle(length, baseWidth, endWidth float64) ConeReticle {
	return ConeReticle{
		Length:    length,
		BaseWidth: baseWidth,
		EndWidth:  endWidth,
	}
}

func NewSkipReticle(length, baseWidth, endWidth float64, numSkips int) SkipReticle {
	return SkipReticle{
		Length:    length,
		BaseWidth: baseWidth,
		EndWidth:  endWidth,
		NumSkips:  numSkips,
	}
}

func NewPlaneAmmo(number int) (ammo *PlaneAmmo) {
	ammo = &PlaneAmmo{
		Number: number,
	}

	return
}

func (a *PlaneAmmo) WithRocket(rocket *PlaneAmmoRocket) *PlaneAmmo {
	a.Rocket = rocket
	return a
}

func (a *PlaneAmmo) WithTorpedo(torpedo *PlaneAmmoTorpedo) *PlaneAmmo {
	a.Torpedo = torpedo
	return a
}

func (a *PlaneAmmo) WithBomb(bomb *PlaneAmmoBomb) *PlaneAmmo {
	a.Bomb = bomb
	return a
}

func (a *PlaneAmmo) WithSkipBomb(skipBomb *PlaneAmmoSkipBomb) *PlaneAmmo {
	a.SkipBomb = skipBomb
	return a
}

func (a *PlaneAmmo) WithMine(mine *PlaneAmmoMine) *PlaneAmmo {
	a.Mine = mine
	return a
}

func NewPlaneAmmoRocket(damage DamageSource, reticle EllipticalReticle, speed float64) *PlaneAmmoRocket {
	return &PlaneAmmoRocket{
		DamageSource:      damage,
		EllipticalReticle: reticle,
		Speed:             speed,
	}
}

func NewPlaneAmmoTorpedo(damage DamageSource, reticle ConeReticle, speed, floodingChance float64) *PlaneAmmoTorpedo {
	return &PlaneAmmoTorpedo{
		DamageSource:   damage,
		ConeReticle:    reticle,
		Speed:          speed,
		FloodingChance: floodingChance,
	}
}

func NewPlaneAmmoBomb(damage DamageSource, reticle EllipticalReticle, fallTime float64) *PlaneAmmoBomb {
	return &PlaneAmmoBomb{
		DamageSource:      damage,
		EllipticalReticle: reticle,
		FallTime:          fallTime,
	}
}

func NewPlaneAmmoSkipBomb(damage DamageSource, reticle SkipReticle, fallTime float64) *PlaneAmmoSkipBomb {
	return &PlaneAmmoSkipBomb{
		DamageSource: damage,
		SkipReticle:  reticle,
		FallTime:     fallTime,
	}
}

func NewPlaneAmmoMine(damage DamageSource, reticle EllipticalReticle, activationDelay, duration int) *PlaneAmmoMine {
	return &PlaneAmmoMine{
		DamageSource:      damage,
		EllipticalReticle: reticle,
		ActivationDelay:   activationDelay,
		Duration:          duration,
	}
}

// Squadron Builder

func NewSquadron(plane *Plane, ammo *PlaneAmmo) (s *Squadron) {
	s = &Squadron{
		Plane: plane,
	}

	if ammo != nil {
		s.Ammo = *ammo
	}

	return
}

func (s *Squadron) WithAmmo(ammo *PlaneAmmo) *Squadron {
	if ammo != nil {
		s.Ammo = *ammo
		return s
	}

	s.Ammo = PlaneAmmo{}
	return s
}

func (s *Squadron) SetStrikeProps(isRTS, isTactical bool, squadronSize, attacksWith, cooldownBetweenStrikes int) *Squadron {
	s.IsRTS = isRTS
	s.IsTactical = isTactical
	s.SquadronSize = squadronSize
	s.AttacksWith = attacksWith
	s.CooldownBetweenStrikes = cooldownBetweenStrikes
	return s
}

func (s *Squadron) SetHangarProps(hangarSize, planePrepTime, planeLaunchTime, planeRecoveryTime, planeRegenerationTime int) *Squadron {
	s.HangarSize = hangarSize
	s.PlanePrepTime = planePrepTime
	s.PlaneLaunchTime = planeLaunchTime
	s.PlaneRecoveryTime = planeRecoveryTime
	s.PlaneRegenerationTime = planeRegenerationTime
	return s
}

func (s *Ship) AddSquadron(squadron *Squadron) *Ship {
	s.Squadrons = append(s.Squadrons, squadron)
	return s
}
