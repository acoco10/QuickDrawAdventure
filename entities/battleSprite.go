package entities

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteBattleState uint8

const (
	Idle SpriteBattleState = iota
	UsingCombatSkill
	UsingDialogueSkill
)

type CAnimation uint8

const (
	FanShot CAnimation = iota
	Shot
	FocusedShot
)

type DAnimation uint8

const (
	Insult DAnimation = iota
	Brag
	StareDown
	Draw
)

type BattleSprite struct {
	*Sprite
	SpriteSheet              *spritesheet.SpriteSheet
	CurrentCombatAnimation   CAnimation
	CurrentDialogueAnimation DAnimation
	state                    SpriteBattleState
	IdleAnimation            *animations.Animation
	CombatAnimations         map[CAnimation]*animations.CyclicAnimation
	DialogueAnimations       map[DAnimation]*animations.CyclicAnimation
	counter                  int
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

	if bs.state == UsingCombatSkill {
		return bs.CombatAnimations[bs.CurrentCombatAnimation].Animation
	}

	if bs.state == UsingDialogueSkill {
		return bs.DialogueAnimations[bs.CurrentDialogueAnimation].Animation
	}
	if bs.state == Idle {
		return bs.IdleAnimation
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
	if bs.state == UsingCombatSkill || bs.state == UsingDialogueSkill {

		cycles := bs.GetCycles()

		if bsAnimation.Frame() == bsAnimation.LastF {
			println("counter:", bs.counter, "cycles:", cycles)
			bs.counter++
			if bs.counter > cycles {
				bs.counter = 0
				bsAnimation.ResetFrame()
				bs.UpdateState(Idle)
			}
		}
	}

	bsAnimation.Update()

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
			Shot:        animations.NewCyclicAnimation(5, 12, 7, 7, 1),
			FocusedShot: animations.NewCyclicAnimation(4, 18, 7, 20, 1),
			FanShot:     animations.NewCyclicAnimation(3, 17, 7, 10, 3),
		},
		DialogueAnimations: map[DAnimation]*animations.CyclicAnimation{
			Insult:    animations.NewCyclicAnimation(0, 8, 7, 7, 3),
			Brag:      animations.NewCyclicAnimation(0, 8, 7, 7, 3),
			StareDown: animations.NewCyclicAnimation(0, 8, 7, 7, 3),
		},
		IdleAnimation: animations.NewAnimation(0, 2, 1, 10),
		counter:       0,
		state:         Idle,
	}

	return bSprite, nil
}
