package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type StartScene struct {
	loaded            bool
	audioPlayer       *audioManagement.SFXAudioPlayer
	musicPlayer       *audioManagement.DJ
	playingIntroMusic bool
	gameLog           *sceneManager.GameLog
}

func NewStartScene() *StartScene {
	return &StartScene{
		loaded:            false,
		musicPlayer:       audioManagement.NewSongPlayer(audioManagement.IntroMusic),
		playingIntroMusic: false,
	}
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
	titleScreenImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/titleScreen.png")
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

func (s *StartScene) FirstLoad(gameLog *sceneManager.GameLog) {
	s.loaded = true
	s.gameLog = gameLog

}

func (s *StartScene) IsLoaded() bool {
	return s.loaded
}

func (s *StartScene) OnEnter() {

}

func (s *StartScene) OnExit() {
	s.musicPlayer.Stop()
}

func (s *StartScene) Update() sceneManager.SceneId {
	s.musicPlayer.Update()
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if s.gameLog.Mode == sceneManager.BattleTest {
			s.musicPlayer.Stop()
			println("entering battleTest mode")
			s.gameLog.EnemyEncountered = battleStats.Sheriff
			elyseStats, err := battleStats.LoadSingleCharacter("elyse")
			if err != nil {
				log.Fatal(err)
			}
			s.gameLog.PlayerStats = &elyseStats
			return sceneManager.BattleSceneId
		} else {
			return sceneManager.TownSceneID
		}

	}

	return sceneManager.StartSceneId
}

var _ sceneManager.Scene = (*StartScene)(nil)
