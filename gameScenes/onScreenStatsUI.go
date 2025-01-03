package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type OnScreenStatsUI struct {
	ammoEffect  *AnimatedEffect
	ammoCounter int
}

func LoadAmmoEffect() (*AnimatedEffect, error) {
	ammoImg, _, err := ebitenutil.NewImageFromFile("assets/images/menuAssets/ammoTracker.png")
	if err != nil {
		return nil, err
	}

	ammoSpriteSheet := spritesheet.NewSpritesheet(30, 1, 32, 32)

	ammoEffect := NewEffect(ammoImg, ammoSpriteSheet, 100, 100, 29, 0, 1, 12)

	return ammoEffect, nil
}

func (os *OnScreenStatsUI) Update(turn battle.Turn) error {
	shots := 0
	for _, effect := range turn.PlayerSkillUsed.Effects {
		if effect.EffectType == "shot" {
			shots = effect.NShots
		}
	}
	os.ammoEffect.Update()
	if os.ammoEffect.state == Triggered {
		os.ammoCounter++
		println("ammo GUI Frame Count:", os.ammoCounter, "\n")
	}

	if os.ammoCounter > shots*48+12 {
		os.ammoEffect.state = NotTriggered
		os.ammoCounter = 12
	}

	return nil
}

func (os *OnScreenStatsUI) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(100, 100)
	os.ammoEffect.Draw(screen)
}
