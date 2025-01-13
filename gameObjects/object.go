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
	Update()
}

type DoorObject struct {
	*Sprite
	Trigger
	Animation         map[ObjectState]*animations.Animation
	AnimationActive   bool
	AnimationComplete bool
	Status            ObjectState
	SpriteSheet       spritesheet.SpriteSheet
	DrawAbovePlayer   bool
}

func (do *DoorObject) ActiveAnimation(objectState ObjectState) *animations.Animation {
	if do.Type == EntryDoor || do.Type == ExitDoor {
		return do.Animation[objectState]
	}
	if do.Type == ContextualObject {
		return do.Animation[objectState]
	}
	return nil
}

func (do *DoorObject) PlayAnimation() {
	do.AnimationActive = true
}

func (do *DoorObject) StopAnimation() {
	do.AnimationActive = false
	do.Status = NotTriggered
}

func NewObject(pImg *ebiten.Image, sheet spritesheet.SpriteSheet, enteringAnimation *animations.Animation, leavingAnimation *animations.Animation, trigger Trigger) (*DoorObject, error) {
	objHeight := sheet.SpriteHeight

	object := &DoorObject{
		Sprite: &Sprite{
			Img: pImg,
			X:   float64(trigger.Rect.Min.X),
			Y:   float64(trigger.Rect.Max.Y - objHeight),
		},
		Animation: map[ObjectState]*animations.Animation{
			Entering: enteringAnimation,
			Leaving:  leavingAnimation,
			On:       animations.NewAnimation(2, 2, 1, 10),
		},
		AnimationActive:   false,
		AnimationComplete: false,
		Status:            NotTriggered,
		SpriteSheet:       sheet,
		Trigger:           trigger,
		DrawAbovePlayer:   false,
	}
	return object, nil
}
