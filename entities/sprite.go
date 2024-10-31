package entities

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img          *ebiten.Image
	X, Y, Dx, Dy float64
	Direction    string
	Visible      bool
	InAnimation  bool
}

func NewSpriteImg(imgpath string) (*ebiten.Image, error) {
	img, _, err := ebitenutil.NewImageFromFile(imgpath)
	if err != nil {
		//handle error
		log.Fatal(err)
	}
	return img, err
}
