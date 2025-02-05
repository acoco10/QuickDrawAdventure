package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type OnScreenStatsUI struct {
	ammoEffect   *AnimatedEffect
	shots        int
	healthBar    *ebiten.Image
	tensionMeter *AnimatedEffect
}

func (os *OnScreenStatsUI) LoadEffects() error {
	ammoImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/ammoTracker.png")
	if err != nil {
		return err
	}

	ammoSpriteSheet := spritesheet.NewSpritesheet(23, 1, 32, 32)

	os.ammoEffect = NewEffect(ammoImg, ammoSpriteSheet, 100, 100, 23, 0, 1, 12, 4)

	healthImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/healthBar.png")

	if err != nil {
		return err
	}

	os.healthBar = healthImg

	tensionMeterImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/tensionMeter.png")

	if err != nil {
		return err
	}
	tensionSpriteSheet := spritesheet.NewSpritesheet(10, 1, 93, 54)

	os.tensionMeter = NewEffect(tensionMeterImg, tensionSpriteSheet, 20, 300, 10, 0, 1, 25, 3)
	return nil
}
func (os *OnScreenStatsUI) ProcessTension(bs BattleScene) {

	enemyTensionThreshHold := bs.battle.Enemy.DisplayStat(battleStats.TensionThreshold)

	if float32(bs.battle.Tension) >= float32(enemyTensionThreshHold/2) && bs.battle.Tension < enemyTensionThreshHold {
		if os.tensionMeter.frame == 0 {
			os.tensionMeter.Trigger()
			bs.audioPlayer.Play(audioManagement.TensionIncrease)
		}
	}

	if bs.battle.Tension > enemyTensionThreshHold {
		os.tensionMeter.Trigger()
		bs.audioPlayer.Play(audioManagement.TensionIncrease)
	}

}
func (os *OnScreenStatsUI) ProcessTurn(bs BattleScene, damage []int, skillName string) {
	os.ProcessTension(bs)
	for _, shot := range damage {
		if shot != -1 {
			os.shots++
		}
	}

	if skillName == "reload" {
		os.shots = 0
		os.ammoEffect.frame = 0
	}

	if os.shots > 0 {
		os.ammoEffect.Trigger()
	}

}

func (os *OnScreenStatsUI) TensionUpdate() {
	os.tensionMeter.Update()
	if os.tensionMeter.frame == 4 && os.tensionMeter.frameCounter == 0 {
		os.tensionMeter.UnTrigger()
	}

	if os.tensionMeter.frame == os.tensionMeter.lastFrame-1 {
		os.tensionMeter.UnTrigger()
	}
}

func (os *OnScreenStatsUI) Update() error {
	os.TensionUpdate()
	os.ammoEffect.Update()

	if os.ammoEffect.frame > 0 && os.ammoEffect.frame%4 == 0 && os.ammoEffect.frameCounter == 0 && os.shots > 0 {

		os.shots--
		os.ammoEffect.UnTrigger()

		if os.shots > 0 {
			os.ammoEffect.Trigger()
		}

	}

	if os.ammoEffect.frame == 22 {
		os.ammoEffect.UnTrigger()
	}
	return nil
}

func (os *OnScreenStatsUI) Draw(gameBattle battle.Battle, screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(100, 100)
	opts.GeoM.Scale(4, 4)
	os.ammoEffect.Draw(screen)
	opts.GeoM.Reset()
	opts.GeoM.Scale(4, 4)
	opts.GeoM.Translate(100, 50)
	clr := color.RGBA{R: 10, G: 200, B: 100, A: 255}
	if gameBattle.BattlePhase == battle.Shooting {
		os.ammoEffect.visible = true
		screen.DrawImage(os.healthBar, &opts)
		health := float32(gameBattle.Player.DisplayStat(battleStats.Health)) / float32(gameBattle.Player.DisplayBaselineStat(battleStats.Health))
		println(health)
		println("baseHealth = ", gameBattle.Player.DisplayBaselineStat(battleStats.Health), "health = ", gameBattle.Player.DisplayBaselineStat(battleStats.Health))
		vector.DrawFilledRect(screen, float32(116), float32(75), float32(64*4*health), float32(15), clr, false)
	}

	if gameBattle.BattlePhase == battle.Dialogue {
		os.tensionMeter.visible = true
		os.tensionMeter.Draw(screen)
	}
}
