package game

import (
	"image/color"
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
		nextID, nextFactionID          uint64
		Factions                       map[uint64]*Faction
		FactionsMu                     sync.RWMutex
		Ships                          *util.SafeStorage[*Ship]
		Planes                         *util.SafeStorage[*Plane]
		spatialHash                    *util.SpatialHash[CollidableObject]
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

	Faction struct {
		ID               uint64
		Name             string
		Color            color.RGBA
		ShipsSpatialHash *util.SpatialHash[*Ship]
	}

	Player struct {
		Socket       *web.Socket
		Body         *Ship
		Camera       *Camera
		InputFlags   uint8
		LastFireTick int
		InputMu      sync.RWMutex
		Faction      *Faction
	}

	// All game object should embed this either directly or through another embedded struct
	GenericObject struct {
		ID                                             uint64
		Game                                           *Game
		Position, Velocity                             *util.Vector2D
		Size, Rotation, Friction, Density, Pushability float64
		Faction                                        *Faction
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
	}

	Ship struct {
		PolygonalCollisionPlugin
		Name    string
		Cfg     *definitions.Ship
		Health  *HealthComponent
		Control *Control
	}

	Plane struct {
		CircularCollisionPlugin
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
