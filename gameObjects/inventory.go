package gameObjects

import "github.com/hajimehoshi/ebiten/v2"

type Inventory struct {
	items          []Item
	ammo           int
	weaponEquipped string
}

type Item struct {
	name         string
	inventoryImg ebiten.Image
	examineImg   ebiten.Image
	description  string
}

func (i *Item) Examine() string {
	return i.description
}
