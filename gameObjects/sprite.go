package gameObjects

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img          *ebiten.Image
	X, Y, Dx, Dy float64
	Visible      bool
	InAnimation  bool
	Shadow       bool
}

func (s *Sprite) EnterShadow() {
	s.Shadow = true
}
func (s *Sprite) ExitShadow() {
	s.Shadow = false
}
