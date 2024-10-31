package entities

import (
	"ShootEmUpAdventure/animations"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	Health     uint
	Animations map[PlayerState]*animations.Animation
}

func (p *Player) ActiveAnimation(dX, dY int) *animations.Animation {
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

func NewPlayer(pImg *ebiten.Image, spawnX float64, spawnY float64) (*Player, error) {
	player := &Player{
		Sprite: &Sprite{
			Img:     pImg,
			X:       spawnX,
			Y:       spawnY,
			Visible: true,
		},
		Health: 10,
		Animations: map[PlayerState]*animations.Animation{
			Down:  animations.NewAnimation(0, 4, 4, 22.0),
			Up:    animations.NewAnimation(2, 6, 4, 22.0),
			Left:  animations.NewAnimation(1, 10, 4, 11.0),
			Right: animations.NewAnimation(3, 11, 4, 11.0),
		},
	}

	return player, nil
}
