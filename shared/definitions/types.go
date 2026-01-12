package definitions

import "github.com/z46-dev/game-dev-project/util"

type (
	EntityType         uint8 // Represents the type of an entity
	ShipClassification uint8 // Represents the classification of a ship
	SquadronType       uint8 // Represents the type of squadron
	ShipID             int   // Represents the key of a ship definition
	PlaneID            int   // Represents the key of a plane definition

	Ship struct {
		ID             ShipID             // The unique identifier for the ship
		Name           string             // The name of the ship
		Classification ShipClassification // The classification of the ship
		AssetName      string             // The name of the asset used for rendering the ship
		HullPath       []*util.Vector2D   // The polygonal hull of the ship (will be normalized, -1 to 1)
		Size           float64            // The size of the ship (used for scaling the hull and hardpoints)
		HullHealth     float64            // The health of the ship's hull
		Speed          float64            // The speed of the ship
		TurnSpeed      float64            // The maximum turn speed of the ship in radians per tick
		Squadrons      []*Squadron        // The squadrons carried by the ship
	}

	EllipticalReticle struct {
		Distance float64 // The distance from the planes to the center of the reticle
		Width    float64 // The width of the elliptical reticle
		Height   float64 // The height of the elliptical reticle
	}

	ConeReticle struct {
		Length    float64 // The length of the cone reticle
		BaseWidth float64 // The base width of the cone reticle
		EndWidth  float64 // The end width of the cone reticle
	}

	SkipReticle struct {
		Length    float64 // The length of the skip reticle
		BaseWidth float64 // The base width of the skip reticle
		EndWidth  float64 // The end width of the skip reticle
		NumSkips  int     // The number of skips the bomb will make
	}

	DamageSource struct {
		FullDamage  float64 // The full damage
		Penetration float64 // The penetration value
		FireChance  float64 // The chance to start a fire (0.0 - 1.0)
	}

	PlaneAmmoRocket struct {
		DamageSource
		EllipticalReticle
		Speed float64 // The speed of the rocket
	}

	PlaneAmmoTorpedo struct {
		DamageSource
		ConeReticle
		Speed          float64 // The speed of the torpedo
		FloodingChance float64 // The chance to cause flooding (0.0 - 1.0)
	}

	PlaneAmmoBomb struct {
		DamageSource
		EllipticalReticle
		FallTime float64 // The time it takes for the bomb to fall to the target
	}

	PlaneAmmoSkipBomb struct {
		DamageSource
		SkipReticle
		FallTime float64 // The time it takes for the bomb to fall to the target
	}

	PlaneAmmoMine struct {
		DamageSource
		EllipticalReticle
		ActivationDelay int // The delay in ticks before the mine becomes active
		Duration        int // The duration in ticks the mine remains active
	}

	PlaneAmmo struct {
		Number   int // The number of ammo available
		Rocket   *PlaneAmmoRocket
		Torpedo  *PlaneAmmoTorpedo
		Bomb     *PlaneAmmoBomb
		SkipBomb *PlaneAmmoSkipBomb
		Mine     *PlaneAmmoMine
	}

	Plane struct {
		ID        PlaneID // The unique identifier for the plane
		Name      string  // The name of the plane
		AssetName string  // The name of the asset used for rendering the plane
		Size      float64 // The size of the plane (used for scaling the hull and hardpoints)
		Health    float64 // The health of the plane
		Speed     float64 // The speed of the plane
		TurnSpeed float64 // The maximum turn speed of the plane in radians per tick
	}

	Squadron struct {
		Plane                  *Plane // The plane type used in the squadron
		Ammo                   PlaneAmmo
		HangarSize             int  // The number of planes in reserves at the start of the ship's lifecycle
		IsRTS                  bool // Whether the squadron is RTS controlled or manually controlled
		IsTactical             bool // Do the planes not return to the carrier?
		SquadronSize           int  // Number of planes launched at once
		AttacksWith            int  // Number of planes that attack a target at once (only applicable for manual control)
		CooldownBetweenStrikes int  // Cooldown in ticks between strikes (only applicable for manual control)
		PlanePrepTime          int  // Cooldown in ticks per plane on the carrier deck before a squad can be launched
		PlaneLaunchTime        int  // Time in ticks it takes to launch each plane
		PlaneRecoveryTime      int  // Time in ticks it takes to recover each plane
		PlaneRegenerationTime  int  // Time in ticks it takes to regenerate a single plane in the hangar (if set 0, no regeneration occurs)
	}
)
