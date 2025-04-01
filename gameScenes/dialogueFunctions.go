package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

func CheckTrigger(player gameObjects.Character, x, y float64) bool {
	distance := DistanceEq(player.X, player.Y, x, y)
	if distance < 65 {
		return true
	}
	return false
}

func CheckDialoguePopup(player gameObjects.Character, npc map[string]*gameObjects.Character) gameObjects.Character {
	for _, char := range npc {
		if CheckTrigger(player, char.X, char.Y) {
			return *char
		}
	}
	return gameObjects.Character{}
}

func DrawPopUp(screen *ebiten.Image, x, y, width float64, camera *camera.Camera) {

	popupImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/popup.png")
	opts := ebiten.DrawImageOptions{}
	dopts := text.DrawOptions{}
	opts.GeoM.Scale(2, 2)
	dopts.GeoM.Translate(12, 20)
	opts.GeoM.Translate(x*4+width/2+camera.X*4, y*4-64+camera.Y*4)

	txt := "E"
	face, err := assetManagement.LoadFont(10, assetManagement.November)
	if err != nil {
		log.Fatal("err loading font")
	}
	text.Draw(popupImg, txt, face, &dopts)
	screen.DrawImage(popupImg, &opts)
}
