package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img             *ebiten.Image
	X, Y, Dx, Dy, Z float64
	Visible         bool
	InAnimation     bool
	Shadow          bool
	drawType        DrawableType
}

func (s *Sprite) EnterShadow() {
	s.Shadow = true
}
func (s *Sprite) ExitShadow() {
	s.Shadow = false
}

func (s *Sprite) Draw(screen *ebiten.Image, cam camera.Camera) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(s.X, s.Y)
	opts.GeoM.Translate(0.0, -(float64(s.Img.Bounds().Dy()) + 16))
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Scale(4, 4)

	screen.DrawImage(s.Img, &opts)
}

func (s *Sprite) GetCoord() (x, y, z float64) {
	if s.drawType == Char {
		return s.X, s.Y + 48, s.Z

	}
	return s.X, s.Y, s.Z
}

func (s *Sprite) GetType() DrawableType {
	return s.drawType
}

func (s *Sprite) CheckYSort() bool {
	return true
}

func (s *Sprite) CheckName() string {
	return "generic Sprite"
}
