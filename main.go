package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/isksss/tap-jellyfish/jellyfish"
)

func main() {
	ebiten.SetWindowSize(jellyfish.ScreenWidth, jellyfish.ScreenHeight)
	ebiten.SetWindowTitle(jellyfish.Title)

	game := jellyfish.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
