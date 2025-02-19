package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"code.rocketnine.space/tslocum/brownboxbatman/world"

	"code.rocketnine.space/tslocum/brownboxbatman/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("Brown Box Bat Man")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetMaxTPS(144)
	ebiten.SetRunnableOnUnfocused(true) // Note - this currently does nothing in ebiten
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	parseFlags()

	if world.World.Debug == 0 {
		world.SetMessage("MOVE: ARROW KEYS\nFIRE: Z KEY\nMUTE: M KEY", 144*4)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigc

		g.Exit()
	}()

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
