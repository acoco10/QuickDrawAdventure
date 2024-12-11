package entities

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteBattleState uint8

const (
	Idle SpriteBattleState = iota
	CombatPhaseIdle
	UsingCombatSkill
	UsingDialogueSkill
)

type CAnimation uint8

const (
	FanShot CAnimation = iota
	Shot
	FocusedShot
	NoCombatSkill
	Draw
)

type DAnimation uint8

const (
	Insult DAnimation = iota
	Brag
	StareDown
	NoDialogueSkill
)

type BattleSprite struct {
	*Sprite
	SpriteSheet              *spritesheet.SpriteSheet
	CurrentCombatAnimation   CAnimation
	CurrentDialogueAnimation DAnimation
	state                    SpriteBattleState
	IdleAnimation            *animations.Animation
	CombatIdleAnimation      *animations.Animation
	CombatAnimations         map[CAnimation]*animations.CyclicAnimation
	DialogueAnimations       map[DAnimation]*animations.CyclicAnimation
	counter                  int
	countdown                int
	inAnimation              bool
}

func (bs *BattleSprite) changeEvent(state SpriteBattleState, timer int) {
	bs.state = state
	bs.counter = timer
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
		return nil
	}
	if bs.state == UsingCombatSkill {
		return bs.CombatAnimations[bs.CurrentCombatAnimation].Animation
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

func (bs *BattleSprite) GetCycles() int {
	if bs.state == UsingCombatSkill {
		return bs.CombatAnimations[bs.CurrentCombatAnimation].GetCycles()
	}
	if bs.state == UsingDialogueSkill {
		return bs.DialogueAnimations[bs.CurrentDialogueAnimation].GetCycles()
	}
	return 0
}

func (bs *BattleSprite) Update() {

	bsAnimation := bs.GetAnimation()

	if (bs.state == UsingCombatSkill || bs.state == UsingDialogueSkill) && bsAnimation != nil {
		cycles := bs.GetCycles()
		if bsAnimation.Frame() == bsAnimation.LastF {
			bs.counter++
			if bs.counter > cycles*int(bsAnimation.SpeedInTPS) {
				println("battleSprite counter: ", bs.counter)
				bsAnimation.Reset()
				if bs.state == UsingDialogueSkill {
					bs.UpdateState(Idle)
					bs.counter = 0
					bs.CurrentDialogueAnimation = NoDialogueSkill
				}
				if bs.state == UsingCombatSkill {
					bs.UpdateState(CombatPhaseIdle)
					bs.counter = 0
					bs.CurrentCombatAnimation = NoCombatSkill
					bs.inAnimation = false
				}
				println("battlesprite state after checks:", bs.state)

			}
		}
	}
	if bsAnimation != nil {
		bsAnimation.Update()
	}

}

func (bs *BattleSprite) CombatButtonAnimationTrigger(text string) {

	if text == "shoot" {
		bs.CurrentCombatAnimation = Shot
	}
	if text == "focused_shot" {
		bs.CurrentCombatAnimation = FocusedShot
	}
	if text == "fan_shot" {
		bs.CurrentCombatAnimation = FanShot
	}
	if text == "Draw" {
		bs.CurrentCombatAnimation = Draw
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
		bs.CurrentDialogueAnimation = StareDown
	}

}

func NewBattleSprite(pImg *ebiten.Image, spriteSheet *spritesheet.SpriteSheet, spawnX float64, spawnY float64) (*BattleSprite, error) {
	bSprite := &BattleSprite{
		SpriteSheet: spriteSheet,
		Sprite: &Sprite{
			Img:     pImg,
			X:       spawnX,
			Y:       spawnY,
			Visible: true,
		},
		CombatAnimations: map[CAnimation]*animations.CyclicAnimation{
			Shot:        animations.NewCyclicAnimation(5, 19, 7, 10, 1),
			FocusedShot: animations.NewCyclicAnimation(4, 25, 7, 14, 1),
			FanShot:     animations.NewCyclicAnimation(3, 17, 7, 12, 3),
			Draw:        animations.NewCyclicAnimation(2, 16, 7, 7, 1),
		},
		DialogueAnimations: map[DAnimation]*animations.CyclicAnimation{
			Insult:    animations.NewCyclicAnimation(1, 15, 7, 7, 3),
			Brag:      animations.NewCyclicAnimation(1, 15, 7, 7, 3),
			StareDown: animations.NewCyclicAnimation(0, 14, 7, 7, 3),
		},
		IdleAnimation:            animations.NewAnimation(0, 21, 7, 40),
		counter:                  0,
		state:                    Idle,
		CombatIdleAnimation:      animations.NewAnimation(9, 16, 7, 10),
		CurrentDialogueAnimation: NoDialogueSkill,
		CurrentCombatAnimation:   NoCombatSkill,
		inAnimation:              false,
	}

	return bSprite, nil
}
