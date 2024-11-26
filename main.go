package main

import (
	"ShootEmUpAdventure/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {

	playerimg, _ := entities.NewSpriteImg("assets/images/characters/Zephrynthes.png")
	player, _ := entities.NewPlayer(playerimg, 75, 75)

	doorimg, _ := entities.NewSpriteImg("assets/images/buildings/tavernDoorSpriteSheet.png")
	doorObject, _ := entities.NewObject(doorimg, 167.18, 158.76)

	game := &Game{}
	game = game.NewGame(player, doorObject, "assets/map/expirementmap.json")

	ebiten.SetWindowSize(480, 270)
	ebiten.SetWindowTitle("Quick Draw Adventure")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
