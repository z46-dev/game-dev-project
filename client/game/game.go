package game

import (
	"fmt"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/client/shaders"
	"github.com/z46-dev/game-dev-project/shared"
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/game-dev-project/util"
	"golang.org/x/image/colornames"
)

func NewGame() (g *Game) {
	g = &Game{
		Camera:      newCamera(),
		Ships:       make(map[uint64]*ClientShip),
		Projectiles: make(map[uint64]*ClientProjectile),
	}

	return
}

func (g *Game) Update() (err error) {
	g.LocalTime++
	var width, height int = ebiten.WindowSize()
	g.Camera.Width, g.Camera.Height = float64(width), float64(height)
	g.Camera.Update()

	if g.Socket != nil {
		var flags uint8
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			flags |= protocol.BITFLAG_INPUT_LEFT
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			flags |= protocol.BITFLAG_INPUT_RIGHT
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			flags |= protocol.BITFLAG_INPUT_UP
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			flags |= protocol.BITFLAG_INPUT_DOWN
		}
		if flags != g.lastInputFlags {
			var w *protocol.Writer = new(protocol.Writer)
			w.SetU8(protocol.PACKET_SERVERBOUND_INPUT)
			w.SetU8(flags)
			g.Socket.Write(w.GetBytes())
			g.lastInputFlags = flags
		}
	}

	return
}

func (g *Game) Draw(screen *ebiten.Image) {
	var bounds = screen.Bounds()

	screen.DrawRectShader(bounds.Dx(), bounds.Dy(), shaders.BackgroundShader, &ebiten.DrawRectShaderOptions{
		GeoM: ebiten.GeoM{},
		Uniforms: map[string]any{
			"Time":       float32(g.LocalTime),
			"Camera":     []float32{float32(g.Camera.Position.X), float32(g.Camera.Position.Y), float32(g.Camera.Zoom)},
			"ScreenSize": []float32{float32(bounds.Dx()), float32(bounds.Dy())},
		},
	})

	// g.genericObjects.ForEach(func(o *GenericObject) {
	// 	o.Draw(screen)
	// })

	g.ShipsMu.RLock()
	var ships []*ClientShip = make([]*ClientShip, 0, len(g.Ships))
	for _, ship := range g.Ships {
		ships = append(ships, ship)
	}
	g.ShipsMu.RUnlock()
	sort.Slice(ships, func(i, j int) bool {
		if ships[i].Size == ships[j].Size {
			return ships[i].ID < ships[j].ID
		}
		return ships[i].Size > ships[j].Size
	})
	for _, ship := range ships {
		ship.Draw(g, screen)
	}

	g.ProjectilesMu.RLock()
	var projectiles []*ClientProjectile = make([]*ClientProjectile, 0, len(g.Projectiles))
	for _, projectile := range g.Projectiles {
		projectiles = append(projectiles, projectile)
	}
	g.ProjectilesMu.RUnlock()
	sort.Slice(projectiles, func(i, j int) bool {
		if projectiles[i].Size == projectiles[j].Size {
			return projectiles[i].ID < projectiles[j].ID
		}
		return projectiles[i].Size > projectiles[j].Size
	})
	for _, projectile := range projectiles {
		projectile.Draw(g, screen)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	w, h = ebiten.WindowSize()
	return
}

func (g *Game) ParseViewUpdate(reader *protocol.Reader) {
	g.ServerTime = int(reader.GetU32())
	g.Camera.RealPosition.X = float64(reader.GetF32())
	g.Camera.RealPosition.Y = float64(reader.GetF32())
	g.Camera.RealZoom = g.Camera.Width / float64(reader.GetF32())

	reader.GetU64()

	// Entities in View
	for {
		var (
			id         uint64
			entityType uint8
			isNew      bool
		)

		if id = reader.GetU64(); id == 0 {
			break
		}

		entityType = reader.GetU8()
		isNew = reader.GetU8() == 0

		switch entityType {
		case protocol.ENTITY_TYPE_SHIP:
			g.ParseIncomingShip(reader, id, isNew)
		case protocol.ENTITY_TYPE_PROJECTILE:
			g.ParseIncomingProjectile(reader, id, isNew)
		default:
			fmt.Printf("Unknown entity type: %d\n", entityType)
		}
	}

	// Deletes
	for {
		var id uint64 = reader.GetU64()
		if id == 0 {
			break
		}

		var entityType uint8 = reader.GetU8()
		switch entityType {
		case protocol.ENTITY_TYPE_SHIP:
			g.ShipsMu.Lock()
			delete(g.Ships, id)
			g.ShipsMu.Unlock()
		case protocol.ENTITY_TYPE_PROJECTILE:
			g.ProjectilesMu.Lock()
			delete(g.Projectiles, id)
			g.ProjectilesMu.Unlock()
		}
	}
}

func (g *Game) ParseIncomingShip(reader *protocol.Reader, id uint64, isNew bool) {
	if isNew {
		var ship *ClientShip = &ClientShip{
			ID:       id,
			Position: util.Vector(float64(reader.GetF32()), float64(reader.GetF32())),
			Size:     float64(reader.GetF32()),
			Rotation: float64(reader.GetF32()),
			Name:     reader.GetStringUTF8(),
		}

		ship.RealPosition = ship.Position.Copy()
		ship.RealSize = ship.Size
		ship.RealRotation = ship.Rotation

		var points []*util.Vector2D = make([]*util.Vector2D, 0, reader.GetU16())
		for range cap(points) {
			points = append(points, util.Vector(float64(reader.GetF32()), float64(reader.GetF32())))
		}

		ship.Definition = definitions.MustGetByKey(definitions.ShipConfigs, definitions.ShipID(reader.GetU8()))
		ship.asset = shared.CreateShipAsset(points, 1024, colornames.Lightsteelblue, colornames.Lightslategray)
		ship.HealthRatio = float64(reader.GetF32())

		var numEntries int = int(reader.GetU8())
		ship.Shields = make([][2]float64, numEntries)
		for i := range numEntries {
			ship.Shields[i][0] = float64(reader.GetF32())
			ship.Shields[i][1] = float64(reader.GetF32())
		}

		numEntries = int(reader.GetU8())
		ship.Engines = make([]float64, numEntries)
		for i := range numEntries {
			ship.Engines[i] = float64(reader.GetF32())
		}

		numEntries = int(reader.GetU8())
		ship.Turrets = make([][2]float64, numEntries)
		for i := range numEntries {
			ship.Turrets[i][0] = float64(reader.GetF32())
			ship.Turrets[i][1] = float64(reader.GetF32())
		}

		g.ShipsMu.Lock()
		g.Ships[id] = ship
		g.ShipsMu.Unlock()
	} else {
		g.ShipsMu.RLock()
		var ship *ClientShip = g.Ships[id]
		g.ShipsMu.RUnlock()

		if ship == nil {
			fmt.Printf("Received update for unknown ship ID: %d\n", id)
			return
		}

		var flags uint8 = reader.GetU8()

		if flags&(1<<0) != 0 {
			ship.RealPosition.X = float64(reader.GetF32())
			ship.RealPosition.Y = float64(reader.GetF32())
		}

		if flags&(1<<1) != 0 {
			ship.RealSize = float64(reader.GetF32())
		}

		if flags&(1<<2) != 0 {
			ship.RealRotation = float64(reader.GetF32())
		}

		if flags&(1<<3) != 0 {
			ship.HealthRatio = float64(reader.GetF32())
		}

		if flags&(1<<4) != 0 {
			var numEntries int = int(reader.GetU8())
			ship.Shields = make([][2]float64, numEntries)
			for i := range numEntries {
				ship.Shields[i][0] = float64(reader.GetF32())
				ship.Shields[i][1] = float64(reader.GetF32())
			}
		}

		if flags&(1<<5) != 0 {
			var numEntries int = int(reader.GetU8())
			ship.Engines = make([]float64, numEntries)
			for i := range numEntries {
				ship.Engines[i] = float64(reader.GetF32())
			}
		}

		if flags&(1<<6) != 0 {
			var numEntries int = int(reader.GetU8())
			ship.Turrets = make([][2]float64, numEntries)
			for i := range numEntries {
				ship.Turrets[i][0] = float64(reader.GetF32())
				ship.Turrets[i][1] = float64(reader.GetF32())
			}
		}
	}
}

func (g *Game) ParseIncomingProjectile(reader *protocol.Reader, id uint64, isNew bool) {
	if isNew {
		var projectile *ClientProjectile = &ClientProjectile{
			ID:       id,
			Position: util.Vector(float64(reader.GetF32()), float64(reader.GetF32())),
			Size:     float64(reader.GetF32()),
			Rotation: float64(reader.GetF32()),
		}

		projectile.RealPosition = projectile.Position.Copy()
		projectile.RealSize = projectile.Size
		projectile.RealRotation = projectile.Rotation
		projectile.asset = shared.CreateCircleAsset(projectile.Size*2, colornames.Lightsalmon)

		g.ProjectilesMu.Lock()
		g.Projectiles[id] = projectile
		g.ProjectilesMu.Unlock()
	} else {
		g.ProjectilesMu.RLock()
		var projectile *ClientProjectile = g.Projectiles[id]
		g.ProjectilesMu.RUnlock()

		if projectile == nil {
			fmt.Printf("Received update for unknown projectile ID: %d\n", id)
			return
		}

		var flags uint8 = reader.GetU8()
		if flags&(1<<0) != 0 {
			projectile.RealPosition.X = float64(reader.GetF32())
			projectile.RealPosition.Y = float64(reader.GetF32())
		}

		if flags&(1<<1) != 0 {
			projectile.RealSize = float64(reader.GetF32())
		}

		if flags&(1<<2) != 0 {
			projectile.RealRotation = float64(reader.GetF32())
		}
	}
}
