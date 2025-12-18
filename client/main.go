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

	if err = ebiten.RunGame(g); err != nil {
		log.Panicf("Error running game: %v", err)
	}
}
