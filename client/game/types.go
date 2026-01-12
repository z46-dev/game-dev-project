package game

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/client/web"
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/util"
)

type (
	GameObject interface {
		Draw(screen *ebiten.Image)
	}

	PlayerCamera struct {
		Position            *util.Vector2D
		Width, Height, Zoom float64

		RealPosition *util.Vector2D
		RealZoom     float64
	}

	ClientShip struct {
		ID                                     uint64
		Position, RealPosition                 *util.Vector2D
		Size, RealSize, Rotation, RealRotation float64
		asset                                  *ebiten.Image
		Name                                   string
		Color                                  color.Color
		Outline                                color.Color
		Definition                             *definitions.Ship
		HealthRatio                            float64
	}

	Game struct {
		ServerTime, LocalTime int
		Camera                *PlayerCamera
		PlayerID              uint64
		Socket                *web.Socket
		lastInputFlags        uint8

		Ships         map[uint64]*ClientShip
		ShipsMu       sync.RWMutex

		MousePosition *util.Vector2D
	}
)
