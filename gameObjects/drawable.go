package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
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
	Draw(screen *ebiten.Image, cam camera.Camera, player Character, debug bool)
	GetType() DrawableType
	CheckYSort() bool
	CheckName() string
	GetSize() (w int, h int)
}

func SortDrawables(drawables []Drawable) []Drawable {
	sort.Slice(drawables, func(i, j int) bool {
		_, y1, _ := drawables[i].GetCoord()
		_, y2, _ := drawables[j].GetCoord()
		ySort1 := drawables[i].CheckYSort()
		ySort2 := drawables[j].CheckYSort()
		type1 := drawables[i].GetType()
		type2 := drawables[j].GetType()

		switch {
		case ySort1 && !ySort2:
			// j is non-ySorted → it goes first
			return false
		case !ySort1 && ySort2:
			// i is non-ySorted → it goes first
			return true
		case ySort1 && ySort2:
			// Both are ySorted → sort by y
			if y1 != y2 {
				return y1 < y2
			}
		}

		// Fallback for both non-ySorted or equal y → sort by type enum
		return type1 < type2

	})

	return drawables
}

func DrawGameObjects(drawables []Drawable, screen *ebiten.Image, cam camera.Camera, player Character, debugMode bool) {
	for _, drawable := range drawables {
		/*if drawable.CheckName() == "" {
			for z := i; z < i+20; z++ {
				forwardCheck := drawables[i+10].CheckName()
				println(forwardCheck)
			}
		}*/
		drawable.Draw(screen, cam, player, debugMode)
	}
}

func CheckOverlap(draw1, draw2 Drawable) bool {
	x1, y1, _ := draw1.GetCoord()
	x2, y2, _ := draw2.GetCoord()

	w1, h1 := draw1.GetSize()
	w2, h2 := draw2.GetSize()

	if draw1.GetType() == Map {
		y1 = y1 - float64(h1)
	}
	if draw2.GetType() == Obj {
		y2 = y2 + 16
	}

	if draw2.GetType() == Map {
		y2 = y2 - float64(h2)
	}
	if draw1.GetType() == Char {
		x1 = x1 + 4
		w1 = 4
		h1 = 4
	}
	
	if draw2.GetType() == Char {
		x2 = x2 + 4
		w2 = w2 - 10
		h2 = h2 - 48
	}

	rect := image.Rect(int(x1), int(y1), int(x1)+w1, int(y1)+h1)

	rect2 := image.Rect(int(x2), int(y2), int(x2)+w2, int(y2)+h2)

	if rect.Overlaps(rect2) {

		//fmt.Printf("draw 1 %s intersects draw 2 %s\n", draw1.CheckName(), draw2.CheckName())
		//fmt.Printf("draw 1 min (x,y) = %f,%f, draw1 w,h = %d,%d\n", x1, y1, w1, h1)
		//fmt.Printf("draw 2 min (x,y) = %f,%f, draw2.max w,h = %d,%d\n", x2, y2, w2, h2)

		return true
	}
	return false
}
