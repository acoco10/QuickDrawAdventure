package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/graphicEffects"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type OnScreenStatsUI struct {
	ammoEffect   *graphicEffects.AnimatedEffect
	shots        int
	healthBar    *ebiten.Image
	tensionMeter *graphicEffects.AnimatedEffect
}

func (os *OnScreenStatsUI) LoadEffects() error {
	ammoImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/ammoTracker.png")
	if err != nil {
		return err
	}

	ammoSpriteSheet := spritesheet.NewSpritesheet(23, 1, 32, 32)

	os.ammoEffect = graphicEffects.NewEffect(ammoImg, ammoSpriteSheet, 100, 100, 23, 0, 1, 12, 4)

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
	width := float64(tensionMeterImg.Bounds().Dx()) / 10
	height := float64(tensionMeterImg.Bounds().Dy()) * 2
	xPosition := 1512*0.4 - width
	yPosition := 918*0.65 - height
	println("xPosition of tensionmeter =", xPosition)
	os.tensionMeter = graphicEffects.NewEffect(tensionMeterImg, tensionSpriteSheet, xPosition, yPosition, 10, 0, 1, 25, 2)
	return nil
}
func (os *OnScreenStatsUI) ProcessTension(bs BattleScene) {

	enemyTensionThreshHold := bs.battle.CharacterBattleData[battle.Enemy].DisplayStat(battleStats.TensionThreshold)

	if float32(bs.battle.Tension) >= float32(enemyTensionThreshHold/2) && bs.battle.Tension < enemyTensionThreshHold {
		if os.tensionMeter.Frame() == 0 {
			os.tensionMeter.Trigger()
			bs.audioPlayer.Play(audioManagement.TensionIncrease)
		}
	}

	if bs.battle.Tension > enemyTensionThreshHold && bs.battle.BattlePhase == battle.Dialogue {
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
		os.ammoEffect.Reset()
	}

	if os.shots > 0 {
		os.ammoEffect.Trigger()
	}

}

func (os *OnScreenStatsUI) TensionUpdate() {
	os.tensionMeter.Update()
	if os.tensionMeter.Frame() == 4 && os.tensionMeter.FrameCounter == 0 {
		os.tensionMeter.UnTrigger()
	}

	if os.tensionMeter.Frame() == os.tensionMeter.LastFrame-1 {
		os.tensionMeter.UnTrigger()
	}
}

func (os *OnScreenStatsUI) Update() error {
	os.TensionUpdate()
	os.ammoEffect.Update()

	if os.ammoEffect.Frame() > 0 && os.ammoEffect.Frame()%4 == 0 && os.ammoEffect.FrameCounter == 0 && os.shots > 0 {

		os.shots--
		os.ammoEffect.UnTrigger()

		if os.shots > 0 {
			os.ammoEffect.Trigger()
		}

	}

	if os.ammoEffect.Frame() == 22 {
		os.ammoEffect.UnTrigger()
	}
	return nil
}

func (os *OnScreenStatsUI) Draw(gameBattle battle.Battle, screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(100, 100)
	opts.GeoM.Scale(4, 4)
	os.ammoEffect.Draw(screen, 0)
	opts.GeoM.Reset()
	opts.GeoM.Scale(4, 4)
	opts.GeoM.Translate(100, 50)
	clr := color.RGBA{R: 10, G: 200, B: 100, A: 255}
	if gameBattle.BattlePhase == battle.Shooting {
		os.ammoEffect.MakeVisible()
		screen.DrawImage(os.healthBar, &opts)
		health := float32(gameBattle.CharacterBattleData[battle.Player].DisplayStat(battleStats.Health)) / float32(gameBattle.CharacterBattleData[battle.Player].DisplayBaselineStat(battleStats.Health))
		vector.DrawFilledRect(screen, float32(116), float32(75), float32(64*4*health), float32(15), clr, false)
	}

	if gameBattle.BattlePhase == battle.Dialogue {
		os.tensionMeter.Draw(screen, 0)
	}
}
