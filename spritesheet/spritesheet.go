package spritesheet

import (
	"image"
)

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	TileSize      int
	SpriteWidth   int
	SpriteHeight  int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := index % s.WidthInTiles * s.SpriteWidth
	y := index / s.WidthInTiles * s.SpriteHeight
	lowX := x + s.SpriteWidth
	lowY := y + s.SpriteHeight
	return image.Rect(x, y, lowX, lowY)
}

//playerSpriteSheet := spritesheet.NewSpritesheet(2, 5, 18, 18, 31)

func NewSpritesheet(widthTiles, heightTiles, tileSize, spriteWidth, spriteHeight int) *SpriteSheet {
	return &SpriteSheet{
		widthTiles, heightTiles, tileSize, spriteWidth, spriteHeight,
	}
}
