package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img                    *ebiten.Image
	X, Y, Dx, Dy, Z, PrevZ float64
	tileX, tileY           int
	Visible                bool
	InAnimation            bool
	Shadow                 bool
	drawType               DrawableType
}

func (s *Sprite) EnterShadow() {
	s.Shadow = true
}
func (s *Sprite) ExitShadow() {
	s.Shadow = false
}
func (s *Sprite) GetSize() (w int, h int) {
	return s.Img.Bounds().Dx(), s.Img.Bounds().Dy()
}

func (s *Sprite) Draw(screen *ebiten.Image, cam camera.Camera, player Character, debugMode bool) {
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

	} //change back to +48 for y sorting
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
