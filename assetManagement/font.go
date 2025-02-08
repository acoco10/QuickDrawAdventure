package assetManagement

import (
	"bytes"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type FontType uint8

const (
	November FontType = iota
	Lady
	NovemberOutline
)

func LoadFont(size float64, font FontType) (text.Face, error) {
	//reading tff file
	LoadedFont, err := assets.Fonts.ReadFile("fonts/novem.ttf")
	if err != nil {
		return nil, err
	}
	if font == Lady {
		LoadedFont, err = assets.Fonts.ReadFile("fonts/KOMIKZBA.ttf")
		if err != nil {
			return nil, err
		}
	}
	if font == NovemberOutline {
		LoadedFont, err = assets.Fonts.ReadFile("fonts/novemmOutline.ttf")
		if err != nil {
			return nil, err
		}
	}

	//extrapolating bytes to new reader object
	s, err := text.NewGoTextFaceSource(bytes.NewReader(LoadedFont))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	face := text.GoTextFace{
		Source: s,    //source font from tff file
		Size:   size, //input by function
	}

	return &face, nil
}
