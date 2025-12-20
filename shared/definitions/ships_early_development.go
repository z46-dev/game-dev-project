package definitions

import (
	"math"

	"github.com/z46-dev/game-dev-project/util"
)

var ShipTiger *Ship = NewShip(SHIP_TIGER, "Tiger", ShipClassificationFrigate, []*util.Vector2D{
	util.Vector(1.0, 0.0), util.Vector(0.3, 0.85), util.Vector(-0.15, 0.6),
	util.Vector(-0.9, 0.8), util.Vector(-0.55, 0.0), util.Vector(-0.9, -0.8),
	util.Vector(-0.15, -0.6), util.Vector(0.3, -0.85),
}, 210).
	SetHullProps(3400, 2, 0.025).
	AddEngine(NewEngine(util.Vector(-1, 0), 0.05, math.Pi, 280)).
	AddShieldGenerator(NewShieldGenerator(util.Vector(0, 0), 0.035, 0, 140, 400, 3000, 0.1)).
	AddTurretWeaponBank(NewTurretWeaponBank(util.Vector(0.5, 0), 0.1, 0, 300, nil, 0.01, ProjLightLaser).
		AddWeapon(NewGun(util.Vector(0, 0), 0, 2, 0.5)))

var ShipHindenburg *Ship = NewShip(SHIP_HINDENBURG, "Hindenburg", ShipClassificationDestroyer, util.SVGPathToVector2DArray("M1-.05.35-.2.3-.15-.1-.35-.55-.35-.9-.2-.9-.15-1-.1-1 .1-.9.15-.9.2-.55.35-.1.35.3.15.35.2 1 .05Z"), 738).
	SetHullProps(3400, 1, 0.005).
	AddEngine(NewEngine(util.Vector(-1, 0), 0.05, math.Pi, 280)).
	AddShieldGenerator(NewShieldGenerator(util.Vector(0, 0), 0.035, 0, 140, 400, 3000, 0.1)).
	AddTurretWeaponBank(NewTurretWeaponBank(util.Vector(0.5, 0), 0.1, 0, 300, nil, 0.01, ProjHeavyLaser).
		AddWeapon(NewGun(util.Vector(0, 0), 0, 2, 0.5)))
