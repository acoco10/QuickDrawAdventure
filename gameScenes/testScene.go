package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type TestScene struct {
	loaded bool
}

func NewTestScene() *TestScene {
	return &TestScene{
		loaded: false,
	}
}

func (s *TestScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
	titleScreenImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/elyse/elyseBattleSprite.png")
	if err != nil {
		log.Fatal(err)
	}

	screen.DrawImage(titleScreenImg, &ebiten.DrawImageOptions{})

	dopts := text.DrawOptions{}
	face, err := assetManagement.LoadFont(60, assetManagement.November)
	if err != nil {
		log.Fatal(err)
	}
	dopts.GeoM.Translate(450, 450)
	text.Draw(screen, "Press Enter to Start", face, &dopts)

}

func (s *TestScene) FirstLoad(gameLog *sceneManager.GameLog) {
	s.loaded = true
}

func (s *TestScene) IsLoaded() bool {
	return s.loaded
}

func (s *TestScene) OnEnter() {

}

func (s *TestScene) OnExit() {
}

func (s *TestScene) Update() sceneManager.SceneId {
	return sceneManager.TestSceneID
}
