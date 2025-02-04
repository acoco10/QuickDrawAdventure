package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type OnScreenStatsUI struct {
	ammoEffect *AnimatedEffect
	shots      int
}

func LoadAmmoEffect() (*AnimatedEffect, error) {
	ammoImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/ammoTracker.png")
	if err != nil {
		return nil, err
	}

	ammoSpriteSheet := spritesheet.NewSpritesheet(23, 1, 32, 32)

	ammoEffect := NewEffect(ammoImg, ammoSpriteSheet, 100, 100, 23, 0, 1, 12, 4)

	return ammoEffect, nil
}

func (os *OnScreenStatsUI) ProcessTurn(damage []int, skillName string) {
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

func (os *OnScreenStatsUI) Update() error {

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

func (os *OnScreenStatsUI) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(100, 100)
	opts.GeoM.Scale(4, 4)
	os.ammoEffect.Draw(screen)
}
