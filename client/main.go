package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/client/game"
	"github.com/z46-dev/game-dev-project/client/web"
	"github.com/z46-dev/game-dev-project/shared/protocol"
	"github.com/z46-dev/golog"
)

var (
	log *golog.Logger = golog.New().Prefix("[MAIN]", golog.BoldBlue).Timestamp()
	err error
)

func main() {
	log.Info("Starting...")

	ebiten.SetWindowSize(ebiten.Monitor().Size())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle("CS780 Project")
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowDecorated(false)

	var (
		g      *game.Game = game.NewGame()
		socket *web.Socket
	)

	if socket, err = web.Connect("ws://localhost:3000/ws?name=testuser"); err != nil {
		log.Panicf("Error connecting to server: %v", err)
	}

	defer socket.Close()

	socket.OnClose = func() {
		log.Error("Socket closed")
	}

	g.Socket = socket

	go socket.InitiateUpdateLoop(func(message []byte) {
		var (
			reader      *protocol.Reader = protocol.NewReader(message)
			messageType uint8            = reader.GetU8()
		)

		switch messageType {
		case protocol.PACKET_CLIENTBOUND_VIEW_UPDATE:
			g.ParseViewUpdate(reader)
		default:
			fmt.Printf("Unknown message type: %d\n", messageType)
		}
	})

	if err = ebiten.RunGame(g); err != nil {
		log.Panicf("Error running game: %v", err)
	}
}
