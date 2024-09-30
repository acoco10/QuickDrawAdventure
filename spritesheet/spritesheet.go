package spritesheet

type SpriteSheet struct {
	WidthInTiles int 
	HeightInTiles int
	TileSize int
}

func NewSpritesheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		w, h, t,
	}
}