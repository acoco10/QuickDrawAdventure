package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type Line struct {
	X   float32
	toX float32
	Y   float32
}

func Init(x float32, toX float32, y float32) *Line {

	line := Line{
		x,
		toX,
		y,
	}
	return &line
}

func (s *Line) Update() {
	if s.X == 0 {
		s.toX += 15
	}

	if s.X > 0 {
		s.toX -= 15
	}

}

func (s *Line) Draw(screen *ebiten.Image) {
	var c color.Color
	if s.X == 0 {
		c = color.RGBA{
			0,
			2,
			220,
			100,
		}
	} else {
		c = color.RGBA{
			255,
			255,
			255,
			220,
		}
	}

	vector.StrokeLine(screen, s.X, s.Y, s.toX, s.Y, 10, c, false)
}
