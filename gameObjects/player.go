package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type CharState uint8

const (
	Down CharState = iota
	Up
	Left
	Right
	MapIdle
)

type CharType uint8

const (
	Player CharType = iota
	NonPlayer
)

type Character struct {
	Name string
	*Sprite
	Animations  map[CharState]*animations.Animation
	SpriteSheet spritesheet.SpriteSheet
	CharType    CharType
}

func (p *Character) ShowY() float64 {
	return p.Y
}

func (p *Character) ActiveAnimation(dX, dY int) *animations.Animation {
	if p.CharType == NonPlayer {
		return p.Animations[MapIdle]
	}
	if dX > 0 {
		return p.Animations[Right]
	}
	if dX < 0 {
		return p.Animations[Left]
	}
	if dY > 0 {
		return p.Animations[Down]
	}
	if dY < 0 {
		return p.Animations[Up]
	}

	return nil
}

func NewCharacter(img *ebiten.Image, spawnPoint Trigger, sheet spritesheet.SpriteSheet, charType CharType) (*Character, error) {
	character := &Character{
		Sprite: &Sprite{
			Img:       img,
			X:         float64(spawnPoint.Rect.Min.X),
			Y:         float64(spawnPoint.Rect.Min.Y - 64),
			Visible:   true,
			Direction: "Down",
		},

		SpriteSheet: sheet,
		Name:        spawnPoint.Name,
		CharType:    charType,
	}
	if charType == NonPlayer {
		character.Animations = map[CharState]*animations.Animation{
			MapIdle: animations.NewAnimation(0, 2, 1, 60),
		}
	}
	if charType == Player {
		character.Animations = map[CharState]*animations.Animation{
			Down:  animations.NewAnimation(0, 20, 4, 10),
			Up:    animations.NewAnimation(2, 22, 4, 12),
			Left:  animations.NewAnimation(3, 23, 4, 10.0),
			Right: animations.NewAnimation(1, 21, 4, 10.0),
		}

	}

	return character, nil
}
