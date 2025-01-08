package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
	"math"
)

func CheckDialogueTrigger(player gameObjects.Character, npc gameObjects.Character) bool {
	if math.Abs(player.X-npc.X) < 20 && player.Y-npc.Y < 20 {
		return true
	}
	return false
}

func DrawDialoguePopup(player gameObjects.Character, npc gameObjects.Character, screen *ebiten.Image, camera camera.Camera) {
	if CheckDialogueTrigger(player, npc) {
		popupImg, _, err := ebitenutil.NewImageFromFile("assets/images/menuAssets/popup.png")
		opts := ebiten.DrawImageOptions{}
		dopts := text.DrawOptions{}
		opts.GeoM.Scale(2, 2)
		dopts.GeoM.Translate(12, 20)
		opts.GeoM.Translate(npc.X*4+8+camera.X*4, npc.Y*4-64+camera.Y*4)

		txt := "E"
		face, err := LoadFont(10, November)
		if err != nil {
			log.Fatal("err loading font")
		}
		text.Draw(popupImg, txt, face, &dopts)
		screen.DrawImage(popupImg, &opts)
	}
}
