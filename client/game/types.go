package game

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/util"
)

type (
	GameObject interface {
		ID() uint64
		Update()
		Draw(screen *ebiten.Image)
		Destroy()
	}

	GenericObject struct {
		id                             uint64
		position, velocity             *util.Vector2D
		size, rotation, frictionFactor float64
		game                           *Game
		aabb                           *util.AABB
	}

	PlayerCamera struct {
		Position            *util.Vector2D
		Width, Height, Zoom float64

		realPosition *util.Vector2D
		realZoom     float64
	}

	SafeStorage[T GameObject] struct {
		storage           map[uint64]T
		mu                sync.RWMutex
		enqueuedAdditions []T
		enqueuedRemovals  []uint64
	}

	Game struct {
		time   int
		Camera *PlayerCamera
		nextID uint64

		genericObjects *SafeStorage[*GenericObject]
		PlayerObject   *GenericObject
		spatialHash    *util.SpatialHash[*GenericObject]
	}
)
