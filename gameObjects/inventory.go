package gameObjects

import "github.com/hajimehoshi/ebiten/v2"

type Inventory struct {
	items          []Item
	ammo           int
	weaponEquipped string
}

type Item struct {
	Name         string
	InventoryImg ebiten.Image
	ExamineImg   ebiten.Image
	Description  string
}

func (i *Item) Examine() string {
	return i.Description
}
