package main

import (
	"ShootEmUpAdventure/entities"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// structs
type Game struct {
	//game elements
	player       *entities.Player
	enemies      []*entities.Enemy
	items        []*entities.Item
	tilemapJSON  *TilemapJSON
	tilesets     []Tileset
	tilemapimage *ebiten.Image
	cam          *Camera
}

// game update function
func (g *Game) Update() error {
	//react to key presses
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += 1
			} else if enemy.X > g.player.X {
				enemy.X -= 1
			}
			if enemy.Y < g.player.Y {
				enemy.Y += 1
			} else if enemy.Y > g.player.Y {
				enemy.Y -= 1
			}
		}
	}
	for _, item := range g.items {
		if math.Abs(item.X-g.player.X) <= 2 && math.Abs(item.Y-g.player.Y) <= 2 {
			g.player.Health += uint(item.AmtHeal)
			item.Ifinv = true
			fmt.Printf("Picked up an item! Health: %d\n", g.player.Health)
		}

	}

	g.cam.FollowTarget(g.player.X+16, g.player.Y+16, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		320,
		240,
	)

	return nil
}

// drawing screen + sprites
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	//map
	//loop through the tilemap
	for layerIndex, layer := range g.tilemapJSON.Layers {

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

			img := g.tilesets[layerIndex].Img(id)
			fmt.Println("Loading image from path:", g.tilesets)
			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()

		}
	}

	//draw player

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	screen.DrawImage(
		//grab a subimage of the Spritesheet
		g.player.Img.SubImage(
			image.Rect(0, 0, 32, 32),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	// draw enemy sprites
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()

	}

	//draw item
	for _, sprite := range g.items {
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
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(320, 240)
	ebiten.SetWindowTitle("Quick Draw Adventure")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/Elyse.png")
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	ghostImg, _, err := ebitenutil.NewImageFromFile("assets/images//enemies/Ghost.png")
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	fishImg, _, err := ebitenutil.NewImageFromFile("assets/images//items/Fish.png")
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	tilemapimage, _, err := ebitenutil.NewImageFromFile("assets/images/terrain/TilesetFloor.png")
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/map/demoMap.json")
	if err != nil {
		//handle error
		log.Fatal(err)
	}
	tilesets, err := tilemapJSON.GenTileSets()

	if err != nil {

		log.Fatalf("Failed to generate tilesets: %v", err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   75.0,
				Y:   75.0,
			},
			Health: 10,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: ghostImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: ghostImg,
					X:   50.0,
					Y:   50.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: ghostImg,
					X:   100.0,
					Y:   100.0,
				},
				FollowsPlayer: false,
			},
		},
		items: []*entities.Item{
			{
				Sprite: &entities.Sprite{
					Img: fishImg,
					X:   335.0,
					Y:   335.0,
				},
				AmtHeal: rand.Intn(4),
				Ifinv:   false,
			},
		},
		tilemapJSON:  tilemapJSON,
		tilemapimage: tilemapimage,
		tilesets:     tilesets,
		cam:          NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
