package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type ObjectState uint8

const (
	Entering ObjectState = iota
	Leaving
)

type Object struct {
	*Sprite
	Animation         map[ObjectState]*animations.Animation
	AnimationActive   bool
	AnimationComplete bool
	Status            string
	SpriteSheet       spritesheet.SpriteSheet
	Name              string
}

func (o *Object) ActiveAnimation(door string) *animations.Animation {
	if door == "entering" {
		return o.Animation[Entering]
	}
	if door == "leaving" {
		return o.Animation[Leaving]
	}
	return nil
}

func (o *Object) PlayAnimation() {
	o.AnimationActive = true
}

func (o *Object) StopAnimation() {
	o.AnimationActive = false
	o.Status = ""
}

func NewObject(pImg *ebiten.Image, locationX float64, locationY float64, sheet spritesheet.SpriteSheet, enteringAnimation *animations.Animation, leavingAnimation *animations.Animation, name string) (*Object, error) {

	object := &Object{
		Sprite: &Sprite{
			Img: pImg,
			X:   locationX - 1,
			Y:   locationY - 32,
		},
		Animation: map[ObjectState]*animations.Animation{
			Entering: enteringAnimation,
			Leaving:  leavingAnimation,
		},
		AnimationActive:   false,
		AnimationComplete: false,
		Status:            "",
		SpriteSheet:       sheet,
		Name:              name,
	}
	return object, nil
}
