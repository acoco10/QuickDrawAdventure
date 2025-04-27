package gameObjects

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math"
)

type MapItem struct {
	State string
	Name  string
	*Sprite
	Scale float64
}

func (mi MapItem) Draw(screen *ebiten.Image, cam camera.Camera, player Character, debug bool) {
	if mi.Name == "bridgeLeft" || mi.Name == "bridgeRight" {
		overlap := false
		if player.Z == mi.Z {
			overlap = CheckOverlap(player.Sprite, mi)
		}
		if overlap {
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(mi.X, mi.Y+1)
			opts.GeoM.Translate(cam.X, cam.Y)
			opts.GeoM.Scale(4, 4)
			opts.GeoM.Rotate(1.5 * math.Pi / 360)
			screen.DrawImage(mi.Img, &opts)
			opts.GeoM.Reset()

		} else {
			/*	opts := ebiten.DrawImageOptions{}
				transformFactor := 5 - mi.Scale
				opts.GeoM.Translate(mi.X*transformFactor, mi.Y*transformFactor)
				opts.GeoM.Translate(cam.X*transformFactor, cam.Y*transformFactor)
				opts.GeoM.Scale(mi.Scale, mi.Scale)
				screen.DrawImage(mi.Img, &opts)*/
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Translate(mi.X, mi.Y)
			opts.GeoM.Translate(cam.X, cam.Y)
			opts.GeoM.Scale(mi.Scale, mi.Scale)
			screen.DrawImage(mi.Img, &opts)
		}
	} else {
		opts := ebiten.DrawImageOptions{}
		transformFactor := 5 - mi.Scale
		opts.GeoM.Translate(mi.X*transformFactor, mi.Y*transformFactor)
		opts.GeoM.Translate(cam.X*transformFactor, cam.Y*transformFactor)
		opts.GeoM.Scale(mi.Scale, mi.Scale)
		screen.DrawImage(mi.Img, &opts)
		opts.GeoM.Reset()
	}
}

func NewMapItem(object ObjectJSON) MapItem {

	imgPath := fmt.Sprintf("images/items/%s.png", object.Name)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)

	if err != nil {
		log.Fatal("couldnt load item img for:", object.Name)
	}

	mapItem := MapItem{}
	sprite := Sprite{}
	sprite.Img = img
	sprite.X = object.X
	sprite.Y = object.Y
	sprite.drawType = Obj
	mapItem.Name = object.Name

	for _, prop := range object.Properties {
		switch prop.Name {
		case "scale":
			scale, ok := prop.Value.(float64)
			if !ok {
				log.Fatal("could not convert item property scale to float64")
			}
			mapItem.Scale = scale
		case "z":
			z, ok := prop.Value.(float64)
			if !ok {
				log.Fatal("could not convert item property z to float64")
			}
			sprite.Z = z
		}
	}
	mapItem.Sprite = &sprite
	return mapItem
}
