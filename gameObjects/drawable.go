package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
)

type DrawableType uint8

const (
	Map = iota
	Obj
	Char
)

type Drawable interface {
	GetCoord() (x, y, z float64)
	Draw(screen *ebiten.Image, cam camera.Camera)
	GetType() DrawableType
	CheckYSort() bool
	CheckName() string
}

func ZSortDrawables(drawables []Drawable) []Drawable {
	sort.Slice(drawables, func(i, j int) bool {
		_, y1, z1 := drawables[i].GetCoord()
		_, y2, z2 := drawables[j].GetCoord()
		ySort1 := drawables[i].CheckYSort()
		ySort2 := drawables[j].CheckYSort()
		type1 := drawables[i].GetType()
		type2 := drawables[j].GetType()
		if z1 != z2 {
			return z1 < z2
		}
		if !ySort1 {
			return ySort2
		}
		if !ySort2 {
			return ySort1
		}
		if y1 != y2 {
			return y1 < y2
		}
		return type1 < type2
	})
	/*for _, drawable := range drawables {
		x, y, z := drawable.GetCoord()
		fmt.Printf(" X: %f, Y: %f, Z: %f, Type: %v, Name: %s,\n", x, y, z, drawable.GetType(), drawable.CheckName())
	}*/
	return drawables
}

func DrawGameObjects(drawables []Drawable, screen *ebiten.Image, cam camera.Camera) {
	for _, drawable := range drawables {
		drawable.Draw(screen, cam)
	}
}
