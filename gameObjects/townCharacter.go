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

type CharType uint8

const (
	Player CharType = iota
	NonPlayer
)

type Character struct {
	Name string
	*Sprite
	Animations  map[Direction]*animations.Animation
	SpriteSheet spritesheet.SpriteSheet
	CharType    CharType
	Spawned     bool
	Inventory   *Inventory
	BattleStats *battleStats.CharacterData
}

func (c *Character) ActiveAnimation(dX, dY float64) *animations.Animation {
	if c.CharType == NonPlayer {
		return c.Animations[MapIdle]
	}
	if dX > 0 {
		if dY > 0 {
			return c.Animations[Down]
		} else if dY < 0 {
			return c.Animations[Up]
		} else {
			return c.Animations[Right]
		}
	}
	if dX < 0 {
		if dY > 0 {
			return c.Animations[Down]
		} else if dY < 0 {
			return c.Animations[Up]
		} else {
			return c.Animations[Left]
		}
	}
	if dY > 0 {
		return c.Animations[Down]
	}
	if dY < 0 {
		return c.Animations[Up]
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
			Y:        spawnPoint.Y - 64,
			Z:        1,
			drawType: Char,
			Visible:  true,
		},

		SpriteSheet: sheet,
		Name:        spawnPoint.Name,
		CharType:    charType,
	}
	if charType == NonPlayer {
		character.Animations = map[Direction]*animations.Animation{
			MapIdle: animations.NewAnimation(0, 2, 1, 60),
		}
	}
	if charType == Player {
		character.Animations = map[Direction]*animations.Animation{
			Down:  animations.NewAnimation(0, 20, 4, 10),
			Up:    animations.NewAnimation(2, 22, 4, 12),
			Left:  animations.NewAnimation(3, 23, 4, 10.0),
			Right: animations.NewAnimation(1, 21, 4, 10.0),
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

func (c *Character) Draw(screen *ebiten.Image, cam camera.Camera) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(c.X, c.Y)
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Scale(4, 4)
	characterFrame := 0
	characterActiveAnimation := c.ActiveAnimation(c.Dx, c.Dy)
	if characterActiveAnimation != nil {
		characterFrame = characterActiveAnimation.Frame()
	}
	if c.Shadow {
		opts.ColorScale.Scale(0.5, 0.5, 0.5, 1)
	}
	if c.Visible {
		screen.DrawImage(
			//grab a subimage of the Spritesheet
			c.Img.SubImage(
				c.SpriteSheet.Rect(characterFrame)).(*ebiten.Image),
			opts,
		)
	}
	opts.GeoM.Reset()

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
}

func (c *Character) IncreaseZ() {
	c.Z += 1
}

func (c *Character) DecreaseZ() {
	c.Z -= 1
}
