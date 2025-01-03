package main

import (
	"github.com/acoco10/qdabattlesystem/gameManager"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(1512, 982)
	ebiten.SetWindowTitle("Quick Draw Adventure")

	game := gameManager.NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
