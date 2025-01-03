package dialogue

import (
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/qdabattlesystem/gameScenes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
	"math"
)

func CheckDialogueTrigger(player gameObjects.Player, enemy gameObjects.NPC) bool {
	if math.Abs(player.X-enemy.X) < 5 && player.Y-enemy.Y < 5 {
		return true
	}
	return false
}

func DrawDialoguePopup(player gameObjects.Player, enemy gameObjects.NPC, screen *ebiten.Image) {

	dopts := text.DrawOptions{}
	dopts.GeoM.Translate(enemy.X, enemy.Y)
	txt := "press e to talk"
	face, err := gameScenes.LoadFont(10)

	if err != nil {
		log.Fatal("err loading font")
	}

	text.Draw(screen, txt, face, &dopts)
}
