package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	Z           float64
}

func (p *Character) ActiveAnimation(dX, dY float64) *animations.Animation {
	if p.CharType == NonPlayer {
		return p.Animations[MapIdle]
	}
	if dX > 0 {
		if dY > 0 {
			return p.Animations[Down]
		} else if dY < 0 {
			return p.Animations[Up]
		} else {
			return p.Animations[Right]
		}
	}
	if dX < 0 {
		if dY > 0 {
			return p.Animations[Down]
		} else if dY < 0 {
			return p.Animations[Up]
		} else {
			return p.Animations[Left]
		}
	}
	if dY > 0 {
		return p.Animations[Down]
	}
	if dY < 0 {
		return p.Animations[Up]
	}

	return nil
}
func (p *Character) Spawn() {
	p.Spawned = true
}
func (p *Character) DeSpawn() {
	p.Spawned = false
}

func NewCharacter(spawnPoint Spawn, sheet spritesheet.SpriteSheet, charType CharType) (*Character, error) {

	imgPath := "images/characters/npc/townFolk/" + spawnPoint.Name + ".png"
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		img = ebiten.NewImage(10, 10)
	}
	character := &Character{
		Sprite: &Sprite{
			Img:     img,
			X:       float64(spawnPoint.X),
			Y:       float64(spawnPoint.Y - 64),
			Visible: true,
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
