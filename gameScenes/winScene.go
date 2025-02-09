package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type WinScene struct {
	loaded             bool
	playerBattleSprite gameObjects.BattleSprite
	enemyBattleSprite  gameObjects.BattleSprite
	gameLog            *sceneManager.GameLog
}

func NewWinScene() *WinScene {
	return &WinScene{
		loaded: false,
	}
}

func (s *WinScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{74, 170, 40, 255})
	ebitenutil.DebugPrint(screen, "You Win! Press enter to battle again!")

	/*	titleScreenImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/titleScreen.png")
		if err != nil {
			log.Fatal(err)
		}
		screen.DrawImage(titleScreenImg, &ebiten.DrawImageOptions{})*/

	DrawBattleSprite(s.playerBattleSprite, screen, 8)

	dopts := text.DrawOptions{}
	face, err := assetManagement.LoadFont(60, assetManagement.November)

	if err != nil {
		log.Fatal(err)
	}

	dopts.GeoM.Translate(500, 450)
	text.Draw(screen, "You Win", face, &dopts)
	dopts.GeoM.Reset()
	dopts.GeoM.Translate(430, 550)
	text.Draw(screen, "Press Enter to Continue", face, &dopts)

}

func (s *WinScene) FirstLoad(gameLog *sceneManager.GameLog) {
	s.loaded = true
	s.playerBattleSprite = LoadPlayerBattleSprite()
	s.playerBattleSprite.CombatButtonAnimationTrigger("win")
	s.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
	s.gameLog = gameLog
}

func (s *WinScene) IsLoaded() bool {
	return s.loaded
}

func (s *WinScene) OnEnter() {
}

func (s *WinScene) OnExit() {

}

func (s *WinScene) Update() sceneManager.SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if s.gameLog.Mode == sceneManager.BattleTest {
			return sceneManager.StartSceneId
		} else {
			return sceneManager.TownSceneID
		}
	}

	s.playerBattleSprite.Update()

	return sceneManager.WinSceneID
}

var _ sceneManager.Scene = (*WinScene)(nil)
