package gameObjects

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
	Dead
)

type CAnimation uint8

const (
	FanShot CAnimation = iota
	Shot
	FocusedShot
	NoCombatSkill
	Win
	Reload
)

type DAnimation uint8

const (
	Insult DAnimation = iota
	Brag
	StareDown
	NoDialogueSkill
	Draw
)

type OTSCAnimation uint8

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

func (bs *BattleSprite) UpdateHitCounter(counter int) {
	bs.HitCounter = counter
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

	if text == "shoot" {
		bs.CurrentCombatAnimation = Shot
	}
	if text == "focused_shot" {
		bs.CurrentCombatAnimation = FocusedShot
	}
	if text == "fan_shot" {
		bs.CurrentCombatAnimation = FanShot
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
		bs.CurrentDialogueAnimation = StareDown
	}
	if text == "draw" {
		bs.CurrentDialogueAnimation = Draw
	}

}

func (bs *BattleSprite) UpdateScale(scale float64) {
	bs.Scale = scale
}

func NewBattleSprite(pImg *ebiten.Image, spriteSheet *spritesheet.SpriteSheet, x float64, y float64, scale float64) (*BattleSprite, error) {
	bSprite := &BattleSprite{
		Scale:       scale,
		SpriteSheet: spriteSheet,
		Sprite: &Sprite{
			Img:     pImg,
			X:       x,
			Y:       y,
			Visible: true,
		},
		CombatAnimations: map[CAnimation]*animations.CyclicAnimation{
			Shot:        animations.NewCyclicAnimation(5, 25, 10, 15, 1),
			FocusedShot: animations.NewCyclicAnimation(4, 34, 10, 15, 1),
			FanShot:     animations.NewCyclicAnimation(3, 23, 10, 15, 3),
			Win:         animations.NewCyclicAnimation(8, 68, 10, 15, 5),
			Reload:      animations.NewCyclicAnimation(9, 39, 10, 15, 1),
		},

		DialogueAnimations: map[DAnimation]*animations.CyclicAnimation{
			Insult:    animations.NewCyclicAnimation(1, 21, 10, 12, 3),
			Brag:      animations.NewCyclicAnimation(1, 21, 10, 12, 3),
			StareDown: animations.NewCyclicAnimation(1, 21, 10, 12, 3),
			Draw:      animations.NewCyclicAnimation(2, 22, 10, 7, 1),
		},

		IdleAnimation:            animations.NewAnimation(0, 0, 0, 10),
		CombatIdleAnimation:      animations.NewAnimation(12, 12, 0, 10),
		counter:                  0,
		state:                    Idle,
		CurrentDialogueAnimation: NoDialogueSkill,
		CurrentCombatAnimation:   NoCombatSkill,
		inAnimation:              false,
	}

	return bSprite, nil
}
