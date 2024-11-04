package main

import (
	"ShootEmUpAdventure/entities"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	playerimg, _ := entities.NewSpriteImg("assets/images/characters/elyseSpriteSheet.png")
	player, _ := entities.NewPlayer(playerimg, 75, 75)

	doorimg, _ := entities.NewSpriteImg("assets/images/buildings/tavernDoorSpriteSheet.png")
	doorObject, _ := entities.NewObject(doorimg, 167.18, 158.76)

	game := &Game{}
	game = game.NewGame(player, doorObject, "assets/map/town1Map.json")

	ebiten.SetWindowSize(320, 240)
	ebiten.SetWindowTitle("Quick Draw Adventure")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
