package gameObjects

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type Direction uint8

const (
	Down Direction = iota
	Up
	Left
	Right
	None
	MapIdle
)

type CharacterState uint8

const (
	Walking CharacterState = iota
	Idle
	Falling
)

type CharType uint8

const (
	Player CharType = iota
	NonPlayer
)

type Character struct {
	Name string
	*Sprite
	State             CharacterState
	WalkingAnimations map[Direction]*animations.Animation
	IdleAnimation     map[Direction]*animations.Animation
	SpriteSheet       spritesheet.SpriteSheet
	CharType          CharType
	Spawned           bool
	Inventory         *Inventory
	BattleStats       *battleStats.CharacterData
	funcQueue         []func()
	Lasso             *Lasso
	Direction         Direction
	Velocity          float32
}

func (c *Character) AddToFuncQueue(newFunc func()) {
	c.funcQueue = append(c.funcQueue, newFunc)
}

func (c *Character) ExecuteFuncQueue() {
	if len(c.funcQueue) > 0 {
		c.funcQueue[0]()
		c.funcQueue = c.funcQueue[1:]
	}
}

func (c *Character) ActiveAnimation() *animations.Animation {
	if c.CharType == NonPlayer {
		return c.WalkingAnimations[MapIdle]
	}
	if c.CharType == Player {
		switch c.State {
		case Walking:
			return c.WalkingAnimations[c.Direction]
		case Idle:
			return c.IdleAnimation[c.Direction]
		}
	}
	return nil
}
func (c *Character) Spawn() {
	c.Spawned = true
}
func (c *Character) DeSpawn() {
	c.Spawned = false
}

func NewCharacter(spawnPoint Spawn, sheet spritesheet.SpriteSheet, charType CharType) (*Character, error) {

	imgPath := "images/characters/npc/townFolk/" + spawnPoint.Name + ".png"
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		img = ebiten.NewImage(10, 10)
	}
	character := &Character{
		Sprite: &Sprite{
			Img:      img,
			X:        spawnPoint.X,
			Y:        spawnPoint.Y - float64(sheet.SpriteHeight),
			tileY:    int(spawnPoint.Y) / 16,
			tileX:    int(spawnPoint.X) % 16,
			Z:        1,
			PrevZ:    1,
			drawType: Char,
			Visible:  true,
		},

		SpriteSheet: sheet,
		Name:        spawnPoint.Name,
		CharType:    charType,
	}
	if charType == NonPlayer {
		character.WalkingAnimations = map[Direction]*animations.Animation{
			MapIdle: animations.NewAnimation(0, 2, 1, 60),
		}
	}
	if charType == Player {
		character.WalkingAnimations = map[Direction]*animations.Animation{
			Down:  animations.NewAnimation(0, 20, 4, 10),
			Up:    animations.NewAnimation(2, 22, 4, 12),
			Left:  animations.NewAnimation(3, 23, 4, 10.0),
			Right: animations.NewAnimation(1, 21, 4, 10.0),
		}
		character.IdleAnimation = map[Direction]*animations.Animation{
			Down:  animations.NewAnimation(0, 0, 1, 10),
			Up:    animations.NewAnimation(2, 2, 1, 12),
			Left:  animations.NewAnimation(3, 3, 1, 10.0),
			Right: animations.NewAnimation(1, 1, 1, 10.0),
		}
		inv := Inventory{
			items:          make([]Item, 0),
			ammo:           24,
			weaponEquipped: "Colt1851",
		}

		bloodyNote := Item{
			Name:         "BloodyNote",
			InventoryImg: *ebiten.NewImage(10, 10),
			Description:  "Your only clue to you Momma's whereabouts",
			ExamineImg:   *ebiten.NewImage(10, 10),
		}

		inv.items = append(inv.items, bloodyNote)
		character.Inventory = &inv
		character.Lasso = &Lasso{}
		character.Lasso.Triggered = false
		character.Lasso.chargeMax = 5
		charStats, err := battleStats.LoadSingleCharacter(spawnPoint.Name)

		if err != nil {
			return nil, err
		}
		character.BattleStats = &charStats

	}

	if spawnPoint.spawned {
		character.Spawned = true
	} else {
		character.Spawned = false
	}
	return character, nil
}
func (c *Character) CheckName() string {
	return c.Name
}

func (c *Character) noOverLapCurrentZ(zLayer []Drawable) bool {
	overLap := false
	for _, img := range zLayer {
		if CheckOverlap(c, img) {
			overLap = true
		}
	}
	return overLap
}

func (c *Character) Update(mapLayers map[float64][]Drawable) {
	if c.State != Falling {
		if c.Dx > 0 {
			c.State = Walking
			c.Direction = Right
		}
		if c.Dx < 0 {
			c.State = Walking
			c.Direction = Left
		}
		if c.Dy < 0 {
			c.State = Walking
			c.Direction = Up
		}
		if c.Dy > 0 {
			c.State = Walking
			c.Direction = Down
		}
		if c.Dy == 0 && c.Dx == 0 {
			c.State = Idle
		}
	}
	signal := c.Lasso.Update()
	if signal.value != nil {
		c.Velocity = signal.value.(float32) * 10
	}

	if c.Velocity > 0 {
		c.Velocity -= 1
		c.IncreaseX(2)
	}

	playerActiveAnimation := c.ActiveAnimation()
	if playerActiveAnimation != nil {
		playerActiveAnimation.Update()
	}
	c.Falling(mapLayers)
	c.CheckForFallingState(mapLayers[c.Z])

}

func (c *Character) CheckForFallingState(currentLayer []Drawable) {
	onZ := c.noOverLapCurrentZ(currentLayer)
	if !onZ && c.Z > 1 {
		if c.State != Falling {
			c.State = Falling
		}
	}
}

func (c *Character) CheckForNotFallingState(tile Drawable) {
	overlap := CheckOverlap(c, tile)
	_, _, z := tile.GetCoord()
	if overlap {
		if tile.GetType() == Map {
			println("player overlapping with tile checking for falling stop z= ", z)
			t, ok := tile.(Tile)
			if !ok {
				log.Fatal("could not convert map drawable to tile in char falling func")
			}
			if t.TileType == Surface {
				println(t.TileType, "found stopping fall")
				_, _, tileZ := tile.GetCoord()
				c.Z = tileZ
				c.State = Walking
			} else {
				println("Player not on surface, tile type =", t.TileType)
			}
		}
	}
}

func (c *Character) Falling(mapLayers map[float64][]Drawable) {
	if c.State == Falling {
		c.AddToFuncQueue(c.IncreaseY)
		for _, layer := range mapLayers {
			for _, tile := range layer {
				c.CheckForNotFallingState(tile)
			}
		}
	}
}

func (c *Character) Draw(screen *ebiten.Image, cam camera.Camera, Player Character, debugMode bool) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(c.X, c.Y)
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Scale(4, 4)
	characterFrame := 0
	characterActiveAnimation := c.ActiveAnimation()

	if characterActiveAnimation != nil {
		characterFrame = characterActiveAnimation.Frame()
	}

	if c.Shadow {
		opts.ColorScale.Scale(0.5, 0.5, 0.5, 1)
	}

	if c.Visible {
		screen.DrawImage(
			c.Img.SubImage(
				c.SpriteSheet.Rect(characterFrame)).(*ebiten.Image),
			opts,
		)
	}
	opts.GeoM.Reset()

	if debugMode {
		face, err := assetManagement.LoadFont(12, assetManagement.November)
		if err != nil {
			log.Fatal()
		}
		dopts := text.DrawOptions{}
		x := c.X*4 + cam.X*4
		y := c.Y*4 + cam.Y*4
		dopts.GeoM.Translate(x, y)
		coord := fmt.Sprintf("x = %f y = %f z = %f", c.X, c.Y, c.Z)
		text.Draw(screen, coord, face, &dopts)
		stateText := fmt.Sprintf("player state = %d", c.State)
		dopts.GeoM.Reset()
		y = y + 20
		dopts.GeoM.Translate(x, y)
		text.Draw(screen, stateText, face, &dopts)
	}
}

func (c *Character) IncreaseZ() {
	c.Z += 1
}

func (c *Character) DecreaseZ() {
	if c.Z > 1 {
		c.Z -= 1
	}
}

func (c *Character) IncreaseY() {
	c.Y += 6
}

func (c *Character) IncreaseX(amt float64) {
	c.X += amt
}

func (c *Character) GetSize() (w int, h int) {
	return c.SpriteSheet.SpriteWidth, c.SpriteSheet.SpriteHeight
}
