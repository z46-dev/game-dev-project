package game

import (
	"sync"

	"github.com/z46-dev/game-dev-project/server/web"
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
)

type (
	CollidableObject interface {
		util.Identifiable
		GetAABB() *util.AABB
		Insert()
		Collide()
	}

	Game struct {
		time                           int
		nextID                         uint64
		Ships                          *util.SafeStorage[*Ship]
		Projectiles                    *util.SafeStorage[*Projectile]
		spatialHash                    *util.SpatialHash[CollidableObject]
		hardpointsSpatialHash          *util.SpatialHash[*HardpointInstance]
		ShipCache                      map[uint64]*ShipCache
		ProjectileCache                map[uint64]*GenericObjectCache
		ShipCacheMu, ProjectileCacheMu sync.RWMutex
		Players                        map[int]*Player
		PlayersMu                      sync.RWMutex
	}

	Camera struct {
		Position        *util.Vector2D
		FOV             float64
		ShipsSeen       map[uint64]bool
		ProjectilesSeen map[uint64]bool
	}

	Player struct {
		Socket       *web.Socket
		Body         *Ship
		Camera       *Camera
		InputFlags   uint8
		LastFireTick int
		InputMu      sync.RWMutex
	}

	// All game object should embed this either directly or through another embedded struct
	GenericObject struct {
		ID                                             uint64
		Game                                           *Game
		Position, Velocity                             *util.Vector2D
		Size, Rotation, Friction, Density, Pushability float64
		Team                                           int
	}

	CircularCollisionPlugin struct {
		GenericObject
		AABB *util.AABB
	}

	PolygonalCollisionPlugin struct {
		GenericObject
		Polygon *util.Polygon
	}

	HealthComponent struct {
		Health, MaxHealth float64
		CanBeRepaired     bool
		Rebuild           *definitions.RebirthConfig
	}

	HardpointInstance struct {
		Parent           *Ship
		Position         *util.Vector2D
		Size             float64
		FacingDir        float64
		Health           *HealthComponent
		RelativePosition *util.Vector2D
	}

	GunInstance struct {
		RelativePosition    *util.Vector2D
		Direction           float64
		RelLength, RelWidth float64
	}

	TurretInstance struct {
		HardpointInstance
		Cfg        *definitions.Turret
		Guns       []*GunInstance
		TargetShip *Ship
		Target     *util.Vector2D
		ReloadTick int
	}

	ShieldGenerator struct {
		HardpointInstance
		ShieldRadius float64
		ShieldHealth *HealthComponent
		ShieldRegen  float64
	}

	EngineInstance struct {
		HardpointInstance
	}

	Ship struct {
		PolygonalCollisionPlugin
		Name        string
		Cfg         *definitions.Ship
		Health      *HealthComponent
		Shields     []*ShieldGenerator
		Engines     []*EngineInstance
		TurretBanks []*TurretInstance
		Control     *Control
	}

	Projectile struct {
		CircularCollisionPlugin
		ProjectileID definitions.ProjectileID
		Speed        float64
		Range        float64
		Damage       float64
		PrevPosition *util.Vector2D
		Parent       *Ship
	}

	// Caches (Each renderable type should have a cache, using inheretence where possible)

	GenericObjectCache struct {
		AsOf                                int // Corresponds with Game.time
		New, Old                            *protocol.Writer
		ID                                  uint64
		X, Y, Size, Rotation                float64
		PosChanged, SizeChanged, RotChanged bool
	}

	ShipCache struct {
		GenericObjectCache
		Health                                                        float64
		Shields                                                       [][2]float64 // [][hardpoint health ratio, shield health ratio]
		Engines                                                       []float64    // []engine health ratio
		Turrets                                                       [][2]float64 // [][hardpoint health ratio, turret facing (absolute)]
		HealthChanged, ShieldsChanged, EnginesChanged, TurretsChanged bool
	}
)
