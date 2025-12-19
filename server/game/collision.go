package game

import (
	"github.com/z46-dev/game-dev-project/util"
)

func applyElasticity(o, n *GenericObject, mtv *util.Vector2D) {
	var (
		imO, imN float64 = o.invMass(), n.invMass()
		totalM   float64 = imO + imN
	)

	if totalM == 0 {
		return
	}

	var (
		norm, rel      *util.Vector2D = mtv.Copy().Normalize(), o.Velocity.Copy().Subtract(n.Velocity)
		velAlongNormal float64        = rel.Dot(norm)
	)

	if velAlongNormal >= 0 {
		return
	}

	var (
		e       float64 = 0.2
		impulse float64 = -(1 + e) * velAlongNormal / totalM
	)

	o.Velocity.Add(norm.Copy().Scale(impulse * imO))
	n.Velocity.Subtract(norm.Copy().Scale(impulse * imN))
}

func applyMTVResolution(o, n *PolygonalCollisionPlugin) {
	var (
		wO, wN  float64 = o.PushWeight(), n.PushWeight()
		totalW  float64 = wO + wN
		lastMTV *util.Vector2D
	)

	if totalW == 0 {
		return
	}

	for range 4 {
		var resolution *util.Vector2D
		if resolution = util.ResolveTwoPolygons(o.Polygon, n.Polygon); resolution == nil {
			break
		}

		var mO, mC *util.Vector2D = resolution.Copy().Scale(wO / totalW), resolution.Copy().Scale(wN / totalW)
		o.Position.Add(mO)
		n.Position.Subtract(mC)
		o.Polygon.Transform(o.Position, o.Size/2, o.Rotation)
		n.Polygon.Transform(n.Position, n.Size/2, n.Rotation)

		lastMTV = resolution
		if !util.TwoPolygonsIntersect(o.Polygon, n.Polygon) {
			break
		}
	}

	if lastMTV != nil {
		applyElasticity(&o.GenericObject, &n.GenericObject, lastMTV)
	}
}

func polygonsIntersectAt(o, n *PolygonalCollisionPlugin, posO, posN *util.Vector2D) (intersect bool) {
	o.Polygon.Transform(posO, o.Size/2, o.Rotation)
	n.Polygon.Transform(posN, n.Size/2, n.Rotation)
	intersect = util.TwoPolygonsIntersect(o.Polygon, n.Polygon)
	return
}

func shipProjectileCollision(o *Ship, n *Projectile) {}

func shipShipCollision(o *Ship, n *Ship) {
	if o.Pushability == 0 && n.Pushability == 0 {
		return
	}

	var prevO, prevN *util.Vector2D = o.Position.Copy().Subtract(o.Velocity), n.Position.Copy().Subtract(n.Velocity)
	if o.Velocity.SquaredMagnitude() == 0 && n.Velocity.SquaredMagnitude() == 0 {
		var resolution *util.Vector2D
		if resolution = util.ResolveTwoPolygons(o.Polygon, n.Polygon); resolution != nil {
			var norm *util.Vector2D = resolution.Copy().Normalize()
			const nudge float64 = 0.05
			if o.Pushability > 0 {
				o.Velocity.Add(norm.Copy().Scale(nudge))
			}

			if n.Pushability > 0 {
				n.Velocity.Subtract(norm.Copy().Scale(nudge))
			}
		}

		for range 4 {
			if !util.TwoPolygonsIntersect(o.Polygon, n.Polygon) {
				break
			}

			applyMTVResolution(&o.PolygonalCollisionPlugin, &n.PolygonalCollisionPlugin)
		}

		return
	}

	if polygonsIntersectAt(&o.PolygonalCollisionPlugin, &n.PolygonalCollisionPlugin, prevO, prevN) {
		var curO, curN *util.Vector2D = o.Position.Copy(), n.Position.Copy()
		o.Polygon.Transform(curO, o.Size/2, o.Rotation)
		n.Polygon.Transform(curN, n.Size/2, n.Rotation)
		applyMTVResolution(&o.PolygonalCollisionPlugin, &n.PolygonalCollisionPlugin)
		return
	}

	var lo, hi float64 = 0, 1
	for range 12 {
		var mid float64 = (lo + hi) / 2
		var posO *util.Vector2D = prevO.Copy().Add(o.Velocity.Copy().Scale(mid))
		var posN *util.Vector2D = prevN.Copy().Add(n.Velocity.Copy().Scale(mid))

		if polygonsIntersectAt(&o.PolygonalCollisionPlugin, &n.PolygonalCollisionPlugin, posO, posN) {
			hi = mid
			continue
		}

		lo = mid
	}

	o.Position = prevO.Copy().Add(o.Velocity.Copy().Scale(hi))
	n.Position = prevN.Copy().Add(n.Velocity.Copy().Scale(hi))
	o.Polygon.Transform(o.Position, o.Size/2, o.Rotation)
	n.Polygon.Transform(n.Position, n.Size/2, n.Rotation)
	applyMTVResolution(&o.PolygonalCollisionPlugin, &n.PolygonalCollisionPlugin)
}

func projectileProjectileCollision(o *Projectile, n *Projectile) {}

func collideObjects(game *Game, self CollidableObject) {
	if game == nil {
		return
	}

	var myAABB *util.AABB = self.GetAABB()
	if myAABB == nil {
		return
	}

	var collisions []CollidableObject = game.spatialHash.Retrieve(myAABB)
	switch my := self.(type) {
	case *Ship:
		for _, c := range collisions {
			if c == self {
				continue
			}

			switch other := c.(type) {
			case *Ship:
				shipShipCollision(my, other)
			case *Projectile:
				shipProjectileCollision(my, other)
			}
		}
	case *Projectile:
		for _, c := range collisions {
			if c == self {
				continue
			}

			switch other := c.(type) {
			case *Ship:
				shipProjectileCollision(other, my)
			case *Projectile:
				projectileProjectileCollision(my, other)
			}
		}
	}
}

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

func (s *Ship) Update() {
	s.Position.Add(s.Velocity)
	s.Velocity.Scale(s.Friction)
	s.Insert()
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

func (p *Projectile) Update() {
	p.Position.Add(p.Velocity)
	p.Velocity.Scale(p.Friction)
	p.Insert()
}

func (p *Projectile) Collide() {
	collideObjects(p.Game, p)
}
