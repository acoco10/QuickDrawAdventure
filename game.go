package main

import (
	"QuickDrawAdventure/entities"
	"QuickDrawAdventure/mapobjects"
	"QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

type Game struct {

	//game elements
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	//enemies           []*entities.Enemy
	tilemapJSON       *mapobjects.TilemapJSON
	tilesets          []mapobjects.Tileset
	cam               *Camera
	colliders         []image.Rectangle
	objects           *entities.Object
	objectSpriteSheet *spritesheet.SpriteSheet
	entranceDoors     map[string]mapobjects.Door
	exitDoors         map[string]mapobjects.Door
	action            bool
}

func (g *Game) NewGame(player *entities.Player, object *entities.Object, tilemapJSONpath string) *Game {

	tilemapJSON, err := mapobjects.NewTilemapJSON(tilemapJSONpath)
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	tileset, err := tilemapJSON.GenTileSets()
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	cam := NewCamera(0.0, 0.0)

	colliders, entranceDoors, exitDoors := mapobjects.StoreMapObjects(*tilemapJSON)

	playerSpriteSheet := spritesheet.NewSpritesheet(4, 4, 18, 18, 31)

	objectSpriteSheet := spritesheet.NewSpritesheet(2, 2, 20, 20, 21)

	return &Game{
		player,
		playerSpriteSheet,
		tilemapJSON,
		tileset,
		cam,
		colliders,
		object,
		objectSpriteSheet,
		entranceDoors,
		exitDoors,
		false,
	}
}

func (g *Game) Update() error {

	g.player.Dx = 0
	g.player.Dy = 0
	//react to key presses by adding directional velocity
	if !g.player.InAnimation {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.player.Dx = 1.5
			g.player.Direction = "L"
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.player.Dx = -1.5
			g.player.Direction = "R"
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.player.Dy = 1.5
			g.player.Direction = "U"
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.player.Dy = -1.5
			g.player.Direction = "D"
		}

		g.action = false
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			g.action = true
		}
	}

	//increase players position by their velocity every update
	g.player.X += g.player.Dx

	mapobjects.CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy

	mapobjects.CheckCollisionVertical(g.player.Sprite, g.colliders)

	playerOnEntDoor := false
	playerOnExDoor := false
	if !g.player.InAnimation {
		playerOnEntDoor = mapobjects.CheckEntDoor(g.player, g.entranceDoors, g.exitDoors)
	}
	if !g.player.InAnimation {
		playerOnExDoor = mapobjects.CheckExDoor(g.player, g.entranceDoors, g.exitDoors)
	}

	/* //for _, enemy := range g.enemies {

	enemy.Dx = 0.0

	enemy.Dy = 0.0

	if enemy.FollowsPlayer {
		if enemy.X < g.player.X {
			enemy.Dx += 1

		} else if enemy.X > g.player.X {
			enemy.Dx -= 1
		}
		if enemy.Y < g.player.Y {
			enemy.Dy += 1

		} else if enemy.Y > g.player.Y {
			enemy.Dy -= 1
		}
	}
	enemy.X += enemy.Dx */

	//mapobjects.CheckCollisionHorizontal(enemy.Sprite, g.colliders)

	//enemy.Y += enemy.Dy

	//mapobjects.CheckCollisionVertical(enemy.Sprite, g.colliders)

	//checking active player animation
	playerActiveAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if playerActiveAnimation != nil {
		playerActiveAnimation.Update()
	}

	//updating camera to player position
	g.cam.FollowTarget(g.player.X+16, g.player.Y+16, 320, 240)

	//when player hits the edge of the map the camera does not follow
	//need to update this logic for interiors, new map?
	g.cam.Constrain(
		//width of maps from map JSON * tile size
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		//screen resolution
		320,
		240,
	)

	//check if player has entered a door and update door object eventually this will need to be a loop for all object animations
	if playerOnExDoor && g.objects.Status == "" {
		g.player.InAnimation = true
		g.objects.Status = "leaving"
	}

	if playerOnEntDoor && g.objects.Status == "" {
		g.player.InAnimation = true
		g.objects.Status = "entering"
	}

	//custom script animation for tavern door (swings forward on entrance)
	objectAnimation := g.objects.ActiveAnimation(g.objects.Status)

	if objectAnimation != nil {
		if g.objects.Status == "entering" {
			objectAnimation.Update()

			if objectAnimation.Frame() == objectAnimation.LastF-3 {
				//remove sprite on last frame before they are shown inside the building
				g.player.Visible = false
				objectAnimation.Update()
			}

			if objectAnimation.Frame() == objectAnimation.LastF {
				g.player.Visible = true
				x, y := mapobjects.GetDoorCoord(g.exitDoors, "door1", "up")
				g.player.X = x
				g.player.Y = y
				objectAnimation.Update()
				g.objects.StopAnimation()
				objectAnimation.ResetFrame()
				g.player.InAnimation = false
			}
		}
		if g.objects.Status == "leaving" {

			if objectAnimation.Frame() == objectAnimation.FirstF {
				x, y := mapobjects.GetDoorCoord(g.entranceDoors, "door1", "down")
				g.player.X = x
				g.player.Y = y
				g.player.InAnimation = false
				objectAnimation.Update()

			} else if objectAnimation.Frame() == objectAnimation.LastF {
				g.objects.StopAnimation()
				objectAnimation.ResetFrame()

			} else {
				objectAnimation.Update()
			}
		}
	}

	return nil
}

// Draw screen + sprites
func (g *Game) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}

	//map
	//loop through the tile map

	for _, layer := range g.tilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			continue
		}
		gids := make([]int, len(g.tilesets))
		for i := range gids {
			gids[i] = g.tilesets[i].Gid()
		}
		tileindex := mapobjects.DetermineTileSet(layer.Data, gids)

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

			img := g.tilesets[tileindex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()

		}
	}

	//draw player

	opts.GeoM.Translate(g.objects.X, g.objects.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	objectFrame := 0
	objectAnimation := g.objects.ActiveAnimation(g.objects.Status)

	if objectAnimation != nil {
		objectFrame = objectAnimation.Frame()
	}

	screen.DrawImage(
		g.objects.Img.SubImage(
			g.objectSpriteSheet.Rect(objectFrame),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	playerFrame := 0
	playerActiveAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if playerActiveAnimation != nil {

		playerFrame = playerActiveAnimation.Frame()

	} else {
		if g.player.Direction == "U" {
			playerFrame = g.player.Animations[0].FirstF
		}
		if g.player.Direction == "D" {
			playerFrame = g.player.Animations[1].FirstF
		}
		if g.player.Direction == "R" {
			playerFrame = g.player.Animations[2].FirstF
		}
		if g.player.Direction == "L" {
			playerFrame = g.player.Animations[3].FirstF
		}

	}

	if g.player.Visible {
		screen.DrawImage(
			//grab a subimage of the Spritesheet
			g.player.Img.SubImage(
				g.playerSpriteSheet.Rect(playerFrame),
			).(*ebiten.Image),
			&opts,
		)
	}

	opts.GeoM.Reset()

	//make list of all gids in the games tilesets

	for _, layer := range g.tilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			continue
		}

		gids := make([]int, len(g.tilesets))
		for i := range gids {
			gids[i] = g.tilesets[i].Gid()
		}
		tileIndex := mapobjects.DetermineTileSet(layer.Data, gids)

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

			if layer.Class == "above" || layer.Class == "above,below" {
				if int(g.player.Y)+48 < y {

					img := g.tilesets[tileIndex].Img(id)

					opts.GeoM.Translate(float64(x), float64(y))

					opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

					opts.GeoM.Translate(g.cam.X, g.cam.Y)

					screen.DrawImage(img, &opts)

					// reset the opts for the next tile
					opts.GeoM.Reset()
				}
			}
			if layer.Class == "below" || layer.Class == "above,below" {
				if int(g.player.Y)-48 < y {

					img := g.tilesets[tileIndex].Img(id)

					opts.GeoM.Translate(float64(x), float64(y))

					opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

					opts.GeoM.Translate(g.cam.X, g.cam.Y)

					screen.DrawImage(img, &opts)

					// reset the opts for the next tile
					opts.GeoM.Reset()
				}
			}

		}
	}

	opts.GeoM.Reset()

	// draw enemy sprites
	/* for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()

	} */

	//draw item
	/* for _, sprite := range g.items {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		if !sprite.Ifinv {
			screen.DrawImage(
				sprite.Img.SubImage(
					image.Rect(0, 0, 16, 16),
				).(*ebiten.Image),
				&opts,
			)
		}

		opts.GeoM.Reset()
	} */

	/*for _, sprite := range g.objectSprites {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
	}*/

	//TESTING drawing colliders for testing

	/* for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.cam.X),
			float32(collider.Min.Y)+float32(g.cam.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	} */

	//drawing doors for testing
	/* for _, door := range g.entranceDoors {
		vector.StrokeRect(
			screen,
			float32(door.Coord.Min.X)+float32(g.cam.X),
			float32(door.Coord.Min.Y)+float32(g.cam.Y),
			float32(door.Coord.Dx()),
			float32(door.Coord.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	} */

	/* for _, door := range g.exitDoors {
		vector.StrokeRect(
			screen,
			float32(door.Coord.Min.X)+float32(g.cam.X),
			float32(door.Coord.Min.Y)+float32(g.cam.Y),
			float32(door.Coord.Dx()),
			float32(door.Coord.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)

	}*/

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}
