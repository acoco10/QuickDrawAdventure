package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type CharState uint8

const (
	Down CharState = iota
	Up
	Left
	Right
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
	Animations  map[CharState]*animations.Animation
	SpriteSheet spritesheet.SpriteSheet
	CharType    CharType
	Spawned     bool
	Inventory   *Inventory
	BattleStats *battleStats.CharacterData
}

func (p *Character) ActiveAnimation(dX, dY int) *animations.Animation {
	if p.CharType == NonPlayer {
		return p.Animations[MapIdle]
	}
	if dX > 0 {
		return p.Animations[Right]
	}
	if dX < 0 {
		return p.Animations[Left]
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
		return nil, err
	}
	character := &Character{
		Sprite: &Sprite{
			Img:       img,
			X:         float64(spawnPoint.X),
			Y:         float64(spawnPoint.Y - 64),
			Visible:   true,
			Direction: "Down",
		},

		SpriteSheet: sheet,
		Name:        spawnPoint.Name,
		CharType:    charType,
	}
	if charType == NonPlayer {
		character.Animations = map[CharState]*animations.Animation{
			MapIdle: animations.NewAnimation(0, 2, 1, 60),
		}
	}
	if charType == Player {
		character.Animations = map[CharState]*animations.Animation{
			Down:  animations.NewAnimation(0, 20, 4, 10),
			Up:    animations.NewAnimation(2, 22, 4, 12),
			Left:  animations.NewAnimation(3, 23, 4, 10.0),
			Right: animations.NewAnimation(1, 21, 4, 10.0),
		}
		inv := Inventory{
			items:          make([]Item, 0),
			ammo:           24,
			weaponEquipped: "Colt 1851",
		}
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
