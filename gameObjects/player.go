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
)

type Character struct {
	Name string
	*Sprite
	Animations  map[CharState]*animations.Animation
	SpriteSheet spritesheet.SpriteSheet
}

func (p *Character) ShowY() float64 {
	return p.Y
}

func (p *Character) ActiveAnimation(dX, dY int) *animations.Animation {
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

func NewCharacter(img *ebiten.Image, spawnX float64, spawnY float64, sheet spritesheet.SpriteSheet, name string) (*Character, error) {
	character := &Character{
		Sprite: &Sprite{
			Img:       img,
			X:         spawnX,
			Y:         spawnY,
			Visible:   true,
			Direction: "Down",
		},
		Animations: map[CharState]*animations.Animation{
			Down:  animations.NewAnimation(0, 20, 4, 10),
			Up:    animations.NewAnimation(2, 22, 4, 12),
			Left:  animations.NewAnimation(3, 23, 4, 10.0),
			Right: animations.NewAnimation(1, 21, 4, 10.0),
		},
		SpriteSheet: sheet,
		Name:        name,
	}

	return character, nil
}
