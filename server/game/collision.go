package game

import (
	"fmt"

	"github.com/z46-dev/game-dev-project/util"
)

// Circular Collision

func (ccp *CircularCollisionPlugin) GetAABB() (aabb *util.AABB) {
	aabb = ccp.AABB
	return
}

func (ccp *CircularCollisionPlugin) Insert() {
	ccp.AABB.X1, ccp.AABB.Y1 = ccp.Position.X-ccp.Size/2, ccp.Position.Y-ccp.Size/2
	ccp.AABB.X2, ccp.AABB.Y2 = ccp.Position.X+ccp.Size/2, ccp.Position.Y+ccp.Size/2
	ccp.Game.spatialHash.Insert(ccp)
}

// Polygonal Collision

func (pcp *PolygonalCollisionPlugin) GetAABB() (aabb *util.AABB) {
	aabb = pcp.Polygon.AABB
	return
}

func (pcp *PolygonalCollisionPlugin) Insert() {
	pcp.Polygon.Transform(pcp.Position, pcp.Size/2, pcp.Rotation)
	pcp.Game.spatialHash.Insert(pcp)
}

// Ship collision

func (s *Ship) GetAABB() (aabb *util.AABB) {
	aabb = s.Polygon.AABB
	return
}

func (s *Ship) Insert() {
	s.Polygon.Transform(s.Position, s.Size/2, s.Rotation)
	s.Game.spatialHash.Insert(s)
}

func (s *Ship) Collide() {
	collideObjects(s.Game, s)
}

// Projectile collision

func (p *Projectile) GetAABB() (aabb *util.AABB) {
	aabb = p.AABB
	return
}

func (p *Projectile) Insert() {
	p.AABB.X1, p.AABB.Y1 = p.Position.X-p.Size/2, p.Position.Y-p.Size/2
	p.AABB.X2, p.AABB.Y2 = p.Position.X+p.Size/2, p.Position.Y+p.Size/2
	p.Game.spatialHash.Insert(p)
}

func (p *Projectile) Collide() {
	collideObjects(p.Game, p)
}

func collideObjects(game *Game, self CollidableObject) {
	if game == nil {
		return
	}

	var myAABB *util.AABB = self.GetAABB()
	if myAABB == nil {
		fmt.Printf("No AABB for object %+v\n", self)
		return
	}

	var collisions []CollidableObject = game.spatialHash.Retrieve(myAABB)
	for _, c := range collisions {
		if c == self {
			continue
		}

		switch obj := c.(type) {
		case *Ship:
			fmt.Printf("Ship (id=%d)\n", obj.ID)
		case *Projectile:
			fmt.Printf("Projectile (id=%d)\n", obj.ID)
		default:
			fmt.Printf("Unknown collision type: %T\n", c)
		}
	}
}
