package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/graphicEffects"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
)

type SpriteBattleState uint8

const (
	Idle SpriteBattleState = iota
	CombatPhaseIdle
	UsingCombatSkill
	UsingDialogueSkill
	Dead
)

type SoundFxType uint8

const (
	Attack SoundFxType = iota
	GotHit
	Miss
	HitDialogue
	HitTarget
)

type CAnimation uint8

const (
	AttackOne CAnimation = iota
	AttackTwo
	AttackThree
	NoCombatSkill
	Win
	Reload
)

type DAnimation uint8

const (
	Insult DAnimation = iota
	Brag
	Intimidate
	NoDialogueSkill
	Draw
)

type CharEffect uint8

const (
	NoEffect CharEffect = iota
	Outline  CharEffect = iota
)

type CountDownEvent uint8

const (
	NoEvent CountDownEvent = iota
	TurnOffOutline
)

type BattleSprite struct {
	*Sprite
	Scale                    float64
	SpriteSheet              *spritesheet.SpriteSheet
	CurrentCombatAnimation   CAnimation
	CurrentDialogueAnimation DAnimation
	state                    SpriteBattleState
	HitCounter               int
	IdleAnimation            *animations.Animation
	CombatIdleAnimation      *animations.Animation
	CombatAnimations         map[CAnimation]*animations.CyclicAnimation
	DialogueAnimations       map[DAnimation]*animations.CyclicAnimation
	counter                  int
	countdown                int
	inAnimation              bool
	Effects                  map[graphicEffects.EffectType]graphicEffects.GraphicEffect
	SoundFX                  map[SoundFxType]resource.AudioID
	EffectApplied            CharEffect
	CountDownEvent           CountDownEvent
	Hit                      bool
}

func (bs *BattleSprite) GotHit(timer int) {
	bs.Hit = true
	bs.countdown = timer
}

func (bs *BattleSprite) UpdateDialogueAnimation(animation DAnimation) {
	bs.CurrentDialogueAnimation = animation
}

func (bs *BattleSprite) UpdateCombatAnimation(animation CAnimation) {
	bs.CurrentCombatAnimation = animation
}

func (bs *BattleSprite) UpdateState(state SpriteBattleState) {
	bs.state = state
}

func (bs *BattleSprite) GetAnimation() *animations.Animation {
	if bs.state == UsingCombatSkill && bs.CurrentCombatAnimation == NoCombatSkill {
		println("No combat skill but trying to use change state to combat")
		bs.state = CombatPhaseIdle
	}
	if bs.state == UsingCombatSkill {
		return bs.CombatAnimations[bs.CurrentCombatAnimation].Animation
	}

	if bs.state == UsingDialogueSkill && bs.CurrentDialogueAnimation == NoDialogueSkill {
		return bs.IdleAnimation
	}
	if bs.state == UsingDialogueSkill {
		return bs.DialogueAnimations[bs.CurrentDialogueAnimation].Animation
	}

	if bs.state == Idle {
		return bs.IdleAnimation
	}

	if bs.state == CombatPhaseIdle {
		return bs.CombatIdleAnimation
	}

	return nil
}
func (bs *BattleSprite) GetState() SpriteBattleState {
	return bs.state
}

func (bs *BattleSprite) GetCycles() int {
	if bs.state == UsingCombatSkill {
		return bs.CombatAnimations[bs.CurrentCombatAnimation].GetCycles()
	}
	if bs.state == UsingDialogueSkill {
		return bs.DialogueAnimations[bs.CurrentDialogueAnimation].GetCycles()
	}
	return 0
}

func (bs *BattleSprite) TriggerCountDownEvent() {
	if bs.CountDownEvent == TurnOffOutline {
		bs.EffectApplied = NoEffect
	}
}

func (bs *BattleSprite) Update() {
	if bs.countdown > 0 {
		bs.countdown--
		println("battleSprite counter", bs.countdown)
	}
	if bs.countdown == 0 {
		bs.Hit = false
	}

	bsAnimation := bs.GetAnimation()

	if (bs.state == UsingCombatSkill || bs.state == UsingDialogueSkill) && bsAnimation != nil {
		cycles := bs.GetCycles()
		if bsAnimation.Frame() == bsAnimation.LastF {
			bs.counter++
			if bs.counter > cycles*int(bsAnimation.SpeedInTPS) {
				println("battleSprite counter: ", bs.counter)
				bsAnimation.Reset()
				if bs.state == UsingDialogueSkill {
					bs.counter = 0
					if bs.CurrentDialogueAnimation == Draw {
						bs.CurrentDialogueAnimation = NoDialogueSkill
						bs.UpdateState(CombatPhaseIdle)
					} else {
						bs.CurrentDialogueAnimation = NoDialogueSkill
						bs.UpdateState(Idle)
					}

				}
				if bs.state == UsingCombatSkill {
					bs.counter = 0
					bs.CurrentCombatAnimation = NoCombatSkill
					bs.UpdateState(CombatPhaseIdle)

				}
				println("BattleSprite state after checks:", bs.state)

			}
		}
	}

	if bsAnimation != nil {
		bsAnimation.Update()
	}

}

func (bs *BattleSprite) CombatButtonAnimationTrigger(text string) {

	if text == "shoot" || text == "bite" {
		bs.CurrentCombatAnimation = AttackOne
	}
	if text == "focusedShot" {
		bs.CurrentCombatAnimation = AttackTwo
	}
	if text == "fanShot" {
		bs.CurrentCombatAnimation = AttackThree
	}
	if text == "win" {
		bs.CurrentCombatAnimation = Win
	}
	if text == "reload" {
		bs.CurrentCombatAnimation = Reload
	}

}

func (bs *BattleSprite) DialogueButtonAnimationTrigger(text string) {

	if text == "brag" {
		bs.CurrentDialogueAnimation = Brag
	}
	if text == "insult" {
		bs.CurrentDialogueAnimation = Insult
	}
	if text == "stare down" {
		bs.CurrentDialogueAnimation = Intimidate
	}
	if text == "draw" {
		bs.CurrentDialogueAnimation = Draw
	}

}

func (bs *BattleSprite) UpdateScale(scale float64) {
	bs.Scale = scale
}

func (bs *BattleSprite) UpdateCharEffect(effect CharEffect, countDown int) {
	bs.EffectApplied = effect
	bs.countdown = countDown
	bs.CountDownEvent = TurnOffOutline
}

func NewBattleSprite(pImg *ebiten.Image, spriteSheet *spritesheet.SpriteSheet, x float64, y float64, scale float64, canimations map[CAnimation]*animations.CyclicAnimation, cIdle *animations.Animation) (*BattleSprite, error) {
	bSprite := &BattleSprite{
		Scale:       scale,
		SpriteSheet: spriteSheet,
		Sprite: &Sprite{
			Img:     pImg,
			X:       x,
			Y:       y,
			Visible: true,
		},
		CombatAnimations: canimations,
		DialogueAnimations: map[DAnimation]*animations.CyclicAnimation{
			Insult:     animations.NewCyclicAnimation(1, 21, 10, 15, 3),
			Brag:       animations.NewCyclicAnimation(1, 21, 10, 15, 3),
			Intimidate: animations.NewCyclicAnimation(1, 21, 10, 15, 3),
			Draw:       animations.NewCyclicAnimation(2, 22, 10, 10, 1),
		},

		IdleAnimation:            animations.NewAnimation(0, 0, 0, 10),
		CombatIdleAnimation:      animations.NewAnimation(3, 3, 0, 10),
		counter:                  0,
		state:                    Idle,
		CurrentDialogueAnimation: NoDialogueSkill,
		CurrentCombatAnimation:   NoCombatSkill,
		inAnimation:              false,
		EffectApplied:            NoEffect,
	}
	if cIdle != nil {
		bSprite.CombatIdleAnimation = cIdle
	}

	return bSprite, nil
}

func (bs *BattleSprite) LoadEffect(char battleStats.CharacterData) {
	effects := make(map[graphicEffects.EffectType]graphicEffects.GraphicEffect, 0)
	basePath := "images/characters/battleSprites/" + char.Name + "/" + char.Name
	if len(char.DialogueSkills) > 0 {
		starePath := basePath + "StareEffect.png"
		stareSpriteSheet := spritesheet.NewSpritesheet(7, 1, 320, 180)
		stareDownImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, starePath)
		if err != nil {
			log.Fatal(err)
		}
		stareEffect := graphicEffects.NewEffect(stareDownImg, stareSpriteSheet, 0, 0, 3, 0, 1, 30, 5)
		successfulPath := basePath + "SuccessfulEffect.png"
		successfulImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, successfulPath)
		if err != nil {
			log.Fatal(err)
		}
		successfulEffect := graphicEffects.NewStaticEffect(successfulImg, 0, 0, 0, 5, graphicEffects.SuccessfulEffect)

		unSuccessfulPath := basePath + "UnsuccessfulEffect.png"
		unSuccessfulImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, unSuccessfulPath)
		unSuccessfulEffect := graphicEffects.NewStaticEffect(unSuccessfulImg, 0, 0, 0, 5, graphicEffects.UnsuccessfulEffect)
		if err != nil {
			log.Fatal(err)
		}

		effects[graphicEffects.StareEffect] = stareEffect
		effects[graphicEffects.SuccessfulEffect] = successfulEffect
		effects[graphicEffects.UnsuccessfulEffect] = unSuccessfulEffect

	}

	hitStaticEffect := graphicEffects.NewStaticEffect(ebiten.NewImage(1, 1), 0, 0, 30, 1, graphicEffects.TookDamageEffect)
	redClrMap := map[string]float32{
		"r": 1,
		"b": 0.2,
		"g": 0,
		"a": 0.5,
	}
	hitEffect := NewBattleSpriteEffect(*hitStaticEffect, redClrMap, *bs)
	effects[graphicEffects.TookDamageEffect] = hitEffect
	bs.Effects = effects
}

func (bs *BattleSprite) LoadSoundFx(char battleStats.CharacterData) {

}

type BattleSpriteEffect struct {
	battleSprite BattleSprite
	clr          map[string]float32
	graphicEffects.StaticEffect
}

func (bse *BattleSpriteEffect) Draw(screen *ebiten.Image, depth int) {
	if bse.StaticEffect.CheckState() == graphicEffects.Triggered {
		if bse.StaticEffect.Frame() < 20 {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(bse.battleSprite.Scale, bse.battleSprite.Scale)
			opts.GeoM.Translate(bse.battleSprite.X, bse.battleSprite.Y)
			opts.ColorScale.Scale(bse.clr["r"], bse.clr["g"], bse.clr["b"], bse.clr["a"])
			frame := bse.battleSprite.CombatIdleAnimation.FirstF
			img := bse.battleSprite.Img.SubImage(
				bse.battleSprite.SpriteSheet.Rect(frame),
			).(*ebiten.Image)
			screen.DrawImage(img, opts)
		}
	}
}

func NewBattleSpriteEffect(staticEffect graphicEffects.StaticEffect, clr map[string]float32, sprite BattleSprite) *BattleSpriteEffect {

	return &BattleSpriteEffect{
		battleSprite: sprite,
		StaticEffect: staticEffect,
		clr:          clr,
	}

}
