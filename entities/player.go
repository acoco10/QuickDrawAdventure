package entities

import "ShootEmUpAdventure/animations"

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
