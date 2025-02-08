package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type GameOverScene struct {
	loaded bool
}

func NewGameOverScene() *GameOverScene {

	return &GameOverScene{
		loaded: false,
	}
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
	screen.Fill(color.RGBA{255, 0, 0, 255})
	/*	titleScreenImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/titleScreen.png")
		if err != nil {
			log.Fatal(err)
		}
		screen.DrawImage(titleScreenImg, &ebiten.DrawImageOptions{})*/

	dopts := text.DrawOptions{}
	face, err := assetManagement.LoadFont(60, assetManagement.November)

	if err != nil {
		log.Fatal(err)
	}

	dopts.GeoM.Translate(450, 450)
	text.Draw(screen, "Game Over", face, &dopts)
	dopts.GeoM.Translate(430, 500)
	text.Draw(screen, "Press Enter to Continue", face, &dopts)
}

func (s *GameOverScene) FirstLoad(gameLog *sceneManager.GameLog) {
	s.loaded = true
}

func (s *GameOverScene) IsLoaded() bool {
	return s.loaded
}

func (s *GameOverScene) OnEnter() {
}

func (s *GameOverScene) OnExit() {
}

func (s *GameOverScene) Update() sceneManager.SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return sceneManager.StartSceneId
	}

	return sceneManager.GameOverSceneID
}

var _ sceneManager.Scene = (*GameOverScene)(nil)
