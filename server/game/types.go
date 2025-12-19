package game

import "github.com/z46-dev/game-dev-project/util"

type (
	CollidableObject interface {
		util.Identifiable
		GetAABB() *util.AABB
		Insert()
		Collide()
	}

	Game struct {
		time        int
		nextID      uint64
		Ships       *util.SafeStorage[*Ship]
		Projectiles *util.SafeStorage[*Projectile]
		spatialHash *util.SpatialHash[CollidableObject]
	}

	// All game object should embed this either directly or through another embedded struct
	GenericObjectTemplate struct {
		ID                                             uint64
		Game                                           *Game
		Position, Velocity                             *util.Vector2D
		Size, Rotation, Friction, Density, Pushability float64
	}

	CircularCollisionPlugin struct {
		GenericObjectTemplate
		AABB *util.AABB
	}

	PolygonalCollisionPlugin struct {
		GenericObjectTemplate
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
)
