package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

type StartScene struct {
	loaded             bool
	audioPlayer        *audioManagement.SFXAudioPlayer
	musicPlayer        *audioManagement.DJ
	playingIntroMusic  bool
	gameLog            *sceneManager.GameLog
	titleAnimation     *animations.Animation
	sImg               *ebiten.Image
	sheer              *spritesheet.SpriteSheet
	animationTriggered bool
}

func NewStartScene() *StartScene {
	return &StartScene{
		loaded:            false,
		musicPlayer:       audioManagement.NewSongPlayer(audioManagement.IntroMusic),
		playingIntroMusic: false,
	}
}

func (s *StartScene) FirstLoad(gameLog *sceneManager.GameLog) {
	s.loaded = true
	s.gameLog = gameLog
	startImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/titleScreen-sheet.png")
	if err != nil {
		log.Fatal(err)
	}
	s.sImg = startImg
	startSheet := spritesheet.NewSpritesheet(15, 1, 378, 228)

	s.sheer = startSheet
	s.audioPlayer = audioManagement.NewAudioPlayer()
	startAnimation := animations.NewAnimation(0, 14, 1, 5)

	s.titleAnimation = startAnimation

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
		s.audioPlayer.Play(audioManagement.DrawButton)
		s.animationTriggered = true
	}
	if s.animationTriggered {
		s.titleAnimation.Update()
		if s.titleAnimation.Frame() == 7 {
			s.audioPlayer.Play(audioManagement.HorseWinny)
		}
		if s.titleAnimation.Frame() == s.titleAnimation.LastF {
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
	}
	return sceneManager.StartSceneId
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(4, 4)
	frame := s.titleAnimation.Frame()
	frameRect := s.sheer.Rect(frame)
	screen.DrawImage(s.sImg.SubImage(frameRect).(*ebiten.Image), opts)

}

var _ sceneManager.Scene = (*StartScene)(nil)
