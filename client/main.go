package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/client/game"
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
		g       *game.Game     = game.NewGame()
		spinner *golog.Spinner = log.Spinner("Game is running...", golog.SpinnerRunner, 5)
	)

	g.Init()

	spinner.Start()
	defer spinner.Stop()

	// var socket *web.Socket
	// if socket, err = web.Connect("ws://localhost:3000/ws?name=testuser"); err != nil {
	// 	log.Panicf("Error connecting to server: %v", err)
	// }

	// defer socket.Close()

	// socket.OnClose = func() {
	// 	socket.Logger.Error("Socket closed")
	// }

	// go socket.InitiateUpdateLoop(func(message []byte) {})

	if err = ebiten.RunGame(g); err != nil {
		log.Panicf("Error running game: %v", err)
	}
}
