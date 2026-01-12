package game

import (
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
)

func NewCamera(fov float32) (c *Camera) {
	c = &Camera{
		Position:        util.Vector(0, 0),
		FOV:             float64(fov),
		ShipsSeen:       make(map[uint64]bool),
		ProjectilesSeen: make(map[uint64]bool),
	}

	return
}

func (c *Camera) IsInView(aabb *util.AABB) (inside bool) {
	var halfFOV float64 = c.FOV / 2
	inside = !(aabb.X1 > c.Position.X+halfFOV ||
		aabb.X2 < c.Position.X-halfFOV ||
		aabb.Y1 > c.Position.Y+halfFOV ||
		aabb.Y2 < c.Position.Y-halfFOV)

	return
}

func (c *Camera) SeeShip(w *protocol.Writer, o *Ship) {
	var cache *ShipCache
	o.Game.ShipCacheMu.RLock()
	cache = o.Game.ShipCache[o.ID]
	o.Game.ShipCacheMu.RUnlock()

	if cache == nil {
		cache = &ShipCache{}
		o.Game.ShipCacheMu.Lock()
		o.Game.ShipCache[o.ID] = cache
		o.Game.ShipCacheMu.Unlock()
	}

	if cache.AsOf != o.Game.time {
		if cache.X != o.Position.X || cache.Y != o.Position.Y {
			cache.X = o.Position.X
			cache.Y = o.Position.Y
			cache.PosChanged = true
		}

		if cache.Size != o.Size {
			cache.Size = o.Size
			cache.SizeChanged = true
		}

		if cache.Rotation != o.Rotation {
			cache.Rotation = o.Rotation
			cache.RotChanged = true
		}

		if cache.Health != o.Health.Ratio() {
			cache.Health = o.Health.Ratio()
			cache.HealthChanged = true
		}

		cache.AsOf = o.Game.time

		// Clear the buffers
		cache.New = nil
		cache.Old = nil
	}

	// Are we new?
	if _, seen := c.ShipsSeen[o.ID]; !seen {
		c.ShipsSeen[o.ID] = true

		// Build new buffer if needed
		if cache.New == nil {
			cache.New = new(protocol.Writer)

			cache.New.SetU8(0)
			cache.New.SetF32(float32(o.Position.X))
			cache.New.SetF32(float32(o.Position.Y))
			cache.New.SetF32(float32(o.Size))
			cache.New.SetF32(float32(o.Rotation))
			cache.New.SetStringUTF8(o.Name)

			cache.New.SetU16(uint16(len(o.Polygon.Reference)))
			for _, p := range o.Polygon.Reference {
				cache.New.SetF32(float32(p.X))
				cache.New.SetF32(float32(p.Y))
			}

			cache.New.SetU8(uint8(o.Cfg.ID))
			cache.New.SetF32(float32(o.Health.Ratio()))
		}

		// Send new buffer
		w.Append(cache.New)
	} else {
		// Build old buffer if needed
		if cache.Old == nil {
			cache.Old = new(protocol.Writer)
			cache.Old.SetU8(1)

			var flags uint8 = 0
			if cache.PosChanged {
				flags |= 1 << 0
			}

			if cache.SizeChanged {
				flags |= 1 << 1
			}

			if cache.RotChanged {
				flags |= 1 << 2
			}

			if cache.HealthChanged {
				flags |= 1 << 3
			}

			cache.Old.SetU8(flags)

			if cache.PosChanged {
				cache.Old.SetF32(float32(o.Position.X))
				cache.Old.SetF32(float32(o.Position.Y))
			}

			if cache.SizeChanged {
				cache.Old.SetF32(float32(o.Size))
			}

			if cache.RotChanged {
				cache.Old.SetF32(float32(o.Rotation))
			}

			if cache.HealthChanged {
				cache.Old.SetF32(float32(o.Health.Ratio()))
			}
		}

		// Send old buffer
		w.Append(cache.Old)
	}
}

func (c *Camera) See(g *Game, player *Player, w *protocol.Writer) {
	w.SetU32(uint32(g.time))
	w.SetF32(float32(c.Position.X))
	w.SetF32(float32(c.Position.Y))
	w.SetF32(float32(c.FOV))

	if player.Body == nil {
		w.SetU64(0)
	} else {
		c.Position = player.Body.Position
		w.SetU64(player.Body.ID)
	}

	// Entities in View
	var (
		shipsSeenNow       = make(map[uint64]bool)
	)

	for _, something := range g.spatialHash.Retrieve(&util.AABB{
		X1: c.Position.X - c.FOV/2,
		Y1: c.Position.Y - c.FOV/2,
		X2: c.Position.X + c.FOV/2,
		Y2: c.Position.Y + c.FOV/2,
	}) {
		switch o := something.(type) {
		case *Ship:
			shipsSeenNow[o.ID] = true
			w.SetU64(o.ID)
			w.SetU8(protocol.ENTITY_TYPE_SHIP)
			c.SeeShip(w, o)
		}
	}

	// Say we're done
	w.SetU64(0)

	// Deletes
	for id := range c.ShipsSeen {
		if _, stillSeen := shipsSeenNow[id]; !stillSeen {
			w.SetU64(id)
			w.SetU8(protocol.ENTITY_TYPE_SHIP)
			delete(c.ShipsSeen, id)
		}
	}

	// Say we're done with deletes
	w.SetU64(0)
}
