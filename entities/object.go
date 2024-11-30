package entities

import (
	"github.com/acoco10/QuickDrawAdventure/animations"

	"github.com/hajimehoshi/ebiten/v2"
)

type ObjectState uint8

const (
	Entering ObjectState = iota
	Leaving
)

type Object struct {
	*Sprite
	Animation          map[ObjectState]*animations.Animation
	Animation_active   bool
	Animation_complete bool
	Status             string
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
	o.Animation_active = true
}

func (o *Object) StopAnimation() {
	o.Animation_active = false
	o.Status = ""
}

func NewObject(pImg *ebiten.Image, locationX float64, locationY float64) (*Object, error) {
	object := &Object{
		Sprite: &Sprite{
			Img: pImg,
			X:   locationX - 1,
			Y:   locationY - 32,
		},
		Animation: map[ObjectState]*animations.Animation{
			Entering: animations.NewAnimation(0, 6, 1, 10.0),
			Leaving:  animations.NewAnimation(0, 6, 1, 10.0),
		},
		Animation_active:   false,
		Animation_complete: false,
		Status:             "",
	}
	return object, nil
}
