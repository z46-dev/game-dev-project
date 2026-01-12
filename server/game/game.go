package game

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
)

func NewGame() (g *Game) {
	g = &Game{
		Ships:           util.NewSafeStorage[*Ship](),
		Planes:          util.NewSafeStorage[*Plane](),
		spatialHash:     util.NewSpatialHash[CollidableObject](),
		ShipCache:       make(map[uint64]*ShipCache),
		ProjectileCache: make(map[uint64]*GenericObjectCache),
		Players:         make(map[int]*Player),
		Factions:        make(map[uint64]*Faction),
	}

	return
}

func (g *Game) Init() {
	var botChoices []*definitions.Ship = []*definitions.Ship{
		definitions.ShipChkalov,
		definitions.ShipColossus,
		definitions.ShipEnterprise,
		definitions.ShipParseval,
	}

	var npcFaction *Faction = NewFaction(g, "NPCs")
	for range 5 {
		var ship *Ship = NewShip(g, util.RandomRadius(4096), botChoices[rand.IntN(len(botChoices))], npcFaction)
		g.Ships.Add(ship)
	}
}

func (g *Game) Update() {
	g.time++

	// Update spatial hash
	g.spatialHash.Clear()
	for _, f := range g.Factions {
		f.Update()
	}

	// Flush storages
	g.Ships.Flush()
	g.Planes.Flush()

	// Update ships & projectiles (Update & Insert phase)
	g.Ships.ForEach(func(s *Ship) {
		s.Update()
	})

	g.Planes.ForEach(func(p *Plane) {
		p.Update()
	})

	// Collision phase
	g.Ships.ForEach(func(s *Ship) {
		s.Collide()
		s.Think()
	})

	g.Planes.ForEach(func(p *Plane) {
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
