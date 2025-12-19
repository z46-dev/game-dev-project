package game

import (
	"sync"

	"github.com/z46-dev/game-dev-project/server/web"
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
		ShipCache                      map[uint64]*GenericObjectCache
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
		Socket     *web.Socket
		Body       *Ship
		Camera     *Camera
		InputFlags uint8
		InputMu    sync.RWMutex
	}

	// All game object should embed this either directly or through another embedded struct
	GenericObject struct {
		ID                                             uint64
		Game                                           *Game
		Position, Velocity                             *util.Vector2D
		Size, Rotation, Friction, Density, Pushability float64
	}

	CircularCollisionPlugin struct {
		GenericObject
		AABB *util.AABB
	}

	PolygonalCollisionPlugin struct {
		GenericObject
		Polygon *util.Polygon
	}

	Ship struct {
		PolygonalCollisionPlugin

		Name string
	}

	Projectile struct {
		CircularCollisionPlugin

		Speed float64
	}

	// Caches (Each renderable type should have a cache, using inheretence where possible)

	GenericObjectCache struct {
		AsOf                                int // Corresponds with Game.time
		New, Old                            *protocol.Writer
		ID                                  uint64
		X, Y, Size, Rotation                float64
		PosChanged, SizeChanged, RotChanged bool
	}
)
