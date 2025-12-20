package game

import (
	"math"
	"time"

	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
)

func NewGame() (g *Game) {
	g = &Game{
		Ships:                 util.NewSafeStorage[*Ship](),
		Projectiles:           util.NewSafeStorage[*Projectile](),
		spatialHash:           util.NewSpatialHash[CollidableObject](),
		hardpointsSpatialHash: util.NewSpatialHash[*HardpointInstance](),
		ShipCache:             make(map[uint64]*ShipCache),
		ProjectileCache:       make(map[uint64]*GenericObjectCache),
		Players:               make(map[int]*Player),
	}

	return
}

func (g *Game) Init() {
	for range 3 {
		var ship *Ship = NewShip(g, util.RandomRadius(4096), definitions.ShipHindenburg)
		g.Ships.Add(ship)
	}
}

func (g *Game) Update() {
	g.time++

	// Update spatial hash
	g.spatialHash.Clear()
	g.hardpointsSpatialHash.Clear()

	// Flush storages
	g.Ships.Flush()
	g.Projectiles.Flush()

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

	for _, player := range g.Players {
		var w *protocol.Writer = new(protocol.Writer)
		w.SetU8(protocol.PACKET_CLIENTBOUND_VIEW_UPDATE)
		player.Camera.See(g, player, w)
		player.Socket.Write(w.GetBytes())
	}
}

func wrapAngle(angle float64) float64 {
	for angle > math.Pi {
		angle -= math.Pi * 2
	}

	for angle < -math.Pi {
		angle += math.Pi * 2
	}

	return angle
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

func RemovePlayer(g *Game, socketID int) {
	g.PlayersMu.Lock()
	player := g.Players[socketID]
	delete(g.Players, socketID)
	g.PlayersMu.Unlock()

	if player == nil {
		return
	}

	if player.Body != nil {
		g.Ships.Remove(player.Body)
	}
}
