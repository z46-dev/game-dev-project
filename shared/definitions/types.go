package definitions

import "github.com/z46-dev/game-dev-project/util"

type (
	EntityType                uint8 // Represents the type of an entity
	ProjectileMovementPattern uint8 // Represents the movement pattern of a projectile
	HardpointDrawLayer        uint8 // Represents the draw layer of a weapon bank on a ship
	ShipClassification        uint8 // Represents the classification of a ship
	ShipID                    int   // Represents the key of a ship definition
	ProjectileID              int   // Represents the key of a projectile definition

	GenericExplosionCfg struct {
		ExplodesWhenOutOfRange bool    // Whether the projectile explodes when it reaches its maximum range or only on impact
		Radius                 float64 // The radius of the explosion
		EpicenterDamage        float64 // The damage dealt at the center of the explosion
		DamageFallsOff         bool    // Whether the damage falls off with distance from the epicenter f(x) = e/exp(x), 1 <= x <= radius
	}

	Projectile struct {
		ID                     ProjectileID              // The unique identifier for the projectile
		Name                   string                    // The name of the projectile
		Speed                  float64                   // The speed of the projectile
		Range                  float64                   // The maximum range of the projectile
		ImpactDamage           float64                   // The damage dealt on impact
		Explosion              *GenericExplosionCfg      // The explosion configuration, if any
		MovementPattern        ProjectileMovementPattern // The movement pattern of the projectile
		SineMovementAmplitude  float64                   // The amplitude of sine wave movement (if applicable)
		SineMovementFrequency  float64                   // The frequency of sine wave movement (if applicable)
		HomingMovementTurnRate float64                   // The turn rate for homing movement (if applicable)
		ReloadTicks            int                       // Interval in ticks between shots
	}

	RebirthConfig struct {
		RebirthDelayTicks  int     // The delay in ticks before the object or shield is re-build or regenerated
		RebirthHealthRatio float64 // The ratio of max health to restore when the object or shield is rebuilt or regenerated (0.0 - 1.0)
	}

	Hardpoint struct {
		Position      *util.Vector2D     // Position relative to the parent ship's center (normalized, -1 to 1)
		Size          float64            // Size of the hardpoint relative to the ship size (0.0 - 1.0)
		Direction     float64            // Facing direction in radians
		HullHealth    float64            // The health of the hardpoint's hull
		CanBeRepaired bool               // Whether the hardpoint can be repaired if damaged (different from if it can be re-built after being destroyed)
		Rebuild       *RebirthConfig     // The rebuild configuration, nil means no rebuild. If not nil, the hardpoint can be rebuilt after being destroyed after some time.
		DrawLayer     HardpointDrawLayer // The draw layer of the hardpoint
	}

	Gun struct {
		RelativePosition *util.Vector2D // Position relative to the parent object's center (normalized, -1 to 1)
		Direction        float64        // Facing direction in radians
		BarrelLength     float64        // The length of the barrel from the relative position relative to the parent object's size
		BarrelWidth      float64        // The width of the barrel relative to the parent object's size. Projectiles will be this real size when spawned.
	}

	Turret struct {
		Hardpoint                        // The hardpoint this weapon bank is mounted on
		Guns              []*Gun         // The weapons in this bank
		Projectile        *Projectile    // The projectile configuration
		Multishot         int            // Number of projectiles fired per shot
		MultishotInterval int            // Ticks between each projectile in a multishot
		Spread            *util.Vector2D // Spread in radians for random inaccuracy (x: horizontal, y: vertical). Will be calculated like so: rand(-spread value / 2, spread value / 2)
		EffectiveArc      *util.Vector2D // The effective firing arc (x: left limit, y: right limit) in radians relative to the facing direction. A nil value implies a full 360 degree arc.
		TraverseRate      float64        // The rate at which the turret can rotate in radians per tick
	}

	ShieldGenerator struct {
		Hardpoint
		ShieldRadius  float64        // The maximum radius of the shield
		ShieldHealth  float64        // The maximum shield health
		ShieldRegen   float64        // The amount of shield health regenerated per tick
		ShieldRebirth *RebirthConfig // The rebirth configuration for the shield, nil means no rebirth. If not nil, the shield will regenerate from 0 health after some time fully depleted.
	}

	Engine struct {
		Hardpoint
	}

	Ship struct {
		ID             ShipID             // The unique identifier for the ship
		Name           string             // The name of the ship
		Classification ShipClassification // The classification of the ship
		HullPath       []*util.Vector2D   // The polygonal hull of the ship (will be normalized, -1 to 1)
		Size           float64            // The size of the ship (used for scaling the hull and hardpoints)
		HullHealth     float64            // The health of the ship's hull
		Speed          float64            // The speed of the ship
		TurnSpeed      float64            // The maximum turn speed of the ship in radians per tick
		Shields        []*ShieldGenerator // The shield generators on the ship
		Engines        []*Engine          // The engines on the ship
		TurretBanks    []*Turret          // The turret weapon banks on the ship
	}
)
