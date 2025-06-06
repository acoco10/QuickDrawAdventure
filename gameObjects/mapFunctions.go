package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func DrawMapBelowPlayer(tileMapJson TilemapJSON, tilesets []Tileset, cam camera.Camera, screen *ebiten.Image, dark bool) {
	opts := ebiten.DrawImageOptions{}
	for _, layer := range tileMapJson.Layers {
		if layer.Type == "objectgroup" {
			continue
		}

		gids := make([]int, len(tilesets))

		for i := range gids {
			gids[i] = tilesets[i].Gid()
		}

		tileindex := DetermineTileSet(layer.Data, gids)

		if layer.Class == "above,below" {
			continue
		}
		if layer.Class == "eventTrigger" {
			continue
		}

		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			//coordinates example 1%30=1 1/30=0 2%30=2 2/30 = 0 etc...

			x := index % layer.Width
			y := index / layer.Width

			//pixel position
			x *= 16
			y *= 16

			img := tilesets[tileindex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))
			opts.GeoM.Translate(cam.X, cam.Y)
			opts.GeoM.Scale(4, 4)

			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()

		}
	}

}

func DrawMapAbovePlayer(tileMapJSON TilemapJSON, tilesets []Tileset, cam camera.Camera, screen *ebiten.Image, player Character, Triggers map[string]*Trigger, dark bool) {
	opts := ebiten.DrawImageOptions{}
	for _, layer := range tileMapJSON.Layers {

		if layer.Type == "objectgroup" {
			continue
		}

		if layer.Class == "trigger" {

			Check, ok := (layer.Properties[0].Value).(string)
			if !ok {
				log.Fatal("layer.properties value cannot be converted to string")
			}

			if Triggers[Check].Triggered {
				//in this case triggered being true means the player walked over this spot, meaning we want to draw the map below them but not above
				//so the trigger is true (on) but the layer is off
				continue
			}
		}

		if layer.Class == "eventTrigger" {
			Check, ok := (layer.Properties[0].Value).(string)
			if !ok {
				log.Fatal("layer.properties value cannot be converted to string")
			}

			if !Triggers[Check].Triggered {
				continue
			} else {
				println("drawing interaction:", Check)
			}
		}

		gids := make([]int, len(tilesets))
		for i := range gids {
			gids[i] = tilesets[i].Gid()
		}

		tileIndex := DetermineTileSet(layer.Data, gids)

		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			//coordinates example 1%30=1 1/30=0 2%30=2 2/30 = 0 etc...

			x := index % layer.Width
			y := index / layer.Width

			//pixel position
			x *= 16
			y *= 16

			if layer.Class == "above" || layer.Class == "trigger" {
				if int(player.Y)+48 < y {

					img := tilesets[tileIndex].Img(id)

					opts.GeoM.Translate(float64(x), float64(y))

					opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

					opts.GeoM.Translate(cam.X, cam.Y)
					opts.GeoM.Scale(4, 4)

					screen.DrawImage(img, &opts)

					// reset the opts for the next tile
					opts.GeoM.Reset()
				}
			}

			if layer.Class == "below" {
				if int(player.Y)-48 < y {

					img := tilesets[tileIndex].Img(id)

					opts.GeoM.Translate(float64(x), float64(y))

					opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

					opts.GeoM.Translate(cam.X, cam.Y)
					opts.GeoM.Scale(4, 4)

					screen.DrawImage(img, &opts)

					// reset the opts for the next tile
					opts.GeoM.Reset()
				}
			}
			if layer.Class == "above,below" || layer.Class == "eventTrigger" {
				img := tilesets[tileIndex].Img(id)

				opts.GeoM.Translate(float64(x), float64(y))

				opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

				opts.GeoM.Translate(cam.X, cam.Y)
				opts.GeoM.Scale(4, 4)

				screen.DrawImage(img, &opts)

				// reset the opts for the next tile
				opts.GeoM.Reset()
			}
		}

	}
}

func DrawEvent(tileMapJSON TilemapJSON, tilesets []Tileset, cam camera.Camera, screen *ebiten.Image, player Character, Triggers map[string]*Trigger) {

	opts := ebiten.DrawImageOptions{}
	for _, layer := range tileMapJSON.Layers {

		if layer.Class != "eventTrigger" {
			continue
		}

		if layer.Class == "eventTrigger" {
			Check, ok := (layer.Properties[0].Value).(string)
			if !ok {
				log.Fatal("layer.properties value cannot be converted to string")
			}

			if !Triggers[Check].Triggered {
				continue
			} else {
				println("drawing interaction:", Check)
			}
		}

		gids := make([]int, len(tilesets))
		for i := range gids {
			gids[i] = tilesets[i].Gid()
		}

		tileIndex := DetermineTileSet(layer.Data, gids)

		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			//coordinates example 1%30=1 1/30=0 2%30=2 2/30 = 0 etc...

			x := index % layer.Width
			y := index / layer.Width

			//pixel position
			x *= 16
			y *= 16

			if layer.Class == "eventTrigger" {
				img := tilesets[tileIndex].Img(id)

				opts.GeoM.Translate(float64(x), float64(y))

				opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

				opts.GeoM.Translate(cam.X, cam.Y)
				opts.GeoM.Scale(4, 4)

				screen.DrawImage(img, &opts)

				// reset the opts for the next tile
				opts.GeoM.Reset()
			}
		}

	}

}
