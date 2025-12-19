package game

import (
	"math"
	"time"

	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
)

func genPolySides(n int) (sides []*util.Vector2D) {
	for i := range n {
		var angle float64 = 2 * math.Pi / float64(n) * float64(i)
		sides = append(sides, util.Vector(math.Cos(angle), math.Sin(angle)))
	}

	return
}

func genStarSides(n int, radMul float64) (sides []*util.Vector2D) {
	n *= 2
	for i := range n {
		var angle float64 = 2 * math.Pi / float64(n) * float64(i)
		var radius float64 = 1
		if i%2 == 0 {
			radius *= radMul
		}

		sides = append(sides, util.Vector(math.Cos(angle)*radius, math.Sin(angle)*radius))
	}

	return
}

var potentialShapes = [][]*util.Vector2D{
	genPolySides(3),
	genPolySides(4),
	genPolySides(5),
	genPolySides(6),
	genStarSides(3, 0.25),
	genStarSides(4, 0.5),
	genStarSides(5, 0.5),
}

func NewGame() (g *Game) {
	g = &Game{
		Ships:           util.NewSafeStorage[*Ship](),
		Projectiles:     util.NewSafeStorage[*Projectile](),
		spatialHash:     util.NewSpatialHash[CollidableObject](),
		ShipCache:       make(map[uint64]*GenericObjectCache),
		ProjectileCache: make(map[uint64]*GenericObjectCache),
		Players:         make(map[int]*Player),
	}

	return
}

func (g *Game) Init() {
	var rectObj *Ship = NewShip(g, util.Vector(0, 0))
	rectObj.Size = 256
	rectObj.Pushability = 0
	rectObj.Polygon = util.NewPolygon([]*util.Vector2D{
		util.Vector(-1, 0.1),
		util.Vector(1, 0.1),
		util.Vector(1, -0.1),
		util.Vector(-1, -0.1),
	}, rectObj.Position, rectObj.Size/2, rectObj.Rotation)
	g.Ships.Add(rectObj)

	for range 16 {
		var ship *Ship = NewShip(g, util.RandomRadius(1024))
		ship.Size = 64
		g.Ships.Add(ship)
	}
}

func (g *Game) Update() {
	g.time++

	// Update spatial hash
	g.spatialHash.Clear()

	// Flush storages
	g.Ships.Flush()
	g.Projectiles.Flush()

	g.applyPlayerInput()

	// Update ships & projectiles (Update & Insert phase)
	g.Ships.ForEach(func(s *Ship) {
		s.Update()
	})

	g.Projectiles.ForEach(func(p *Projectile) {
		p.Update()
	})

	// Collision phase
	g.Ships.ForEach(func(s *Ship) {
		s.Collide()
	})

	g.Projectiles.ForEach(func(p *Projectile) {
		p.Collide()
	})

	if g.time%3 == 0 {
		for _, player := range g.Players {
			var w *protocol.Writer = new(protocol.Writer)
			w.SetU8(protocol.PACKET_CLIENTBOUND_VIEW_UPDATE)
			player.Camera.See(g, player, w)
			player.Socket.Write(w.GetBytes())
		}
	}
}

func (g *Game) applyPlayerInput() {
	const speed float64 = 4

	g.PlayersMu.RLock()
	defer g.PlayersMu.RUnlock()

	for _, player := range g.Players {
		if player.Body == nil {
			continue
		}

		var flags uint8 = player.GetInputFlags()
		var vx, vy float64

		if flags&protocol.BITFLAG_INPUT_LEFT != 0 {
			vx -= 1
		}
		if flags&protocol.BITFLAG_INPUT_RIGHT != 0 {
			vx += 1
		}
		if flags&protocol.BITFLAG_INPUT_UP != 0 {
			vy -= 1
		}
		if flags&protocol.BITFLAG_INPUT_DOWN != 0 {
			vy += 1
		}

		if vx == 0 && vy == 0 {
			player.Body.Velocity.X = 0
			player.Body.Velocity.Y = 0
		} else {
			var mag float64 = math.Hypot(vx, vy)
			player.Body.Velocity.X = (vx / mag) * speed
			player.Body.Velocity.Y = (vy / mag) * speed
		}

		if player.Camera != nil {
			player.Camera.Position = player.Body.Position
		}
	}
}

func (g *Game) BeginUpdateLoop(tps int) {
	var ticker *time.Ticker = time.NewTicker(time.Duration(1000/tps) * time.Millisecond)
	for range ticker.C {
		g.Update()
	}
}

func SetPlayerInput(g *Game, socketID int, flags uint8) {
	g.PlayersMu.RLock()
	player := g.Players[socketID]
	g.PlayersMu.RUnlock()
	if player == nil {
		return
	}

	player.SetInputFlags(flags)
}
