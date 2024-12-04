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
	TakingDamage
	BackFacing
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
	State                    SpriteBattleState
	Animations               map[SpriteBattleState]*animations.Animation
	CombatAnimations         map[CAnimation]*animations.Animation
	DialogueAnimations       map[DAnimation]*animations.Animation
}

func (bs *BattleSprite) GetAnimation() *animations.Animation {

	if bs.State == UsingCombatSkill {
		return bs.CombatAnimations[bs.CurrentCombatAnimation]
	}

	if bs.State == UsingDialogueSkill {
		return bs.DialogueAnimations[bs.CurrentDialogueAnimation]
	}

	return bs.Animations[bs.State]
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
		Animations: map[SpriteBattleState]*animations.Animation{},
	}

	return bSprite, nil
}
