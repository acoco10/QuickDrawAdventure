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
	NotTriggered
	On
)

type Object interface {
	ActiveAnimation() *animations.Animation
	PlayAnimation()
	StopAnimation()
}

type DoorObject struct {
	*Sprite
	*Trigger
	Animation         map[ObjectState]*animations.Animation
	AnimationActive   bool
	AnimationComplete bool
	State             ObjectState
	SpriteSheet       spritesheet.SpriteSheet
	DrawAbovePlayer   bool
}

func (do *DoorObject) ActiveAnimation() *animations.Animation {
	if do.State == NotTriggered {
		return nil
	}
	return do.Animation[do.State]
}

func (do *DoorObject) PlayAnimation() {
	do.AnimationActive = true
}

func (do *DoorObject) StopAnimation() {
	do.AnimationActive = false
	do.State = NotTriggered
}

func NewObject(pImg *ebiten.Image, sheet spritesheet.SpriteSheet, enteringAnimation *animations.Animation, leavingAnimation *animations.Animation, trigger *Trigger) (*DoorObject, error) {
	objHeight := sheet.SpriteHeight

	object := &DoorObject{
		Sprite: &Sprite{
			Img:      pImg,
			X:        float64(trigger.Rect.Min.X),
			Y:        float64(trigger.Rect.Max.Y - objHeight),
			drawType: Obj,
		},
		Animation: map[ObjectState]*animations.Animation{
			Entering: enteringAnimation,
			Leaving:  leavingAnimation,
			On:       animations.NewAnimation(2, 2, 1, 10),
		},
		AnimationActive:   false,
		AnimationComplete: false,
		State:             NotTriggered,
		SpriteSheet:       sheet,
		Trigger:           trigger,
		DrawAbovePlayer:   false,
	}
	return object, nil
}
