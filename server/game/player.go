package game

import (
	"github.com/z46-dev/game-dev-project/server/web"
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"github.com/z46-dev/game-dev-project/util"
)

func NewPlayer(game *Game, socket *web.Socket, name string) (p *Player) {
	p = &Player{
		Socket: socket,
		Body:   NewShip(game, util.RandomRadius(128), definitions.ShipTiger),
		Camera: NewCamera(3000),
	}

	p.Body.Name = name
	game.Ships.Add(p.Body)

	game.PlayersMu.Lock()
	game.Players[socket.ID] = p
	game.PlayersMu.Unlock()

	return
}

func (p *Player) SetInputFlags(flags uint8) {
	p.InputMu.Lock()
	p.InputFlags = flags
	p.InputMu.Unlock()
}

func (p *Player) GetInputFlags() (flags uint8) {
	p.InputMu.RLock()
	flags = p.InputFlags
	p.InputMu.RUnlock()
	return
}
