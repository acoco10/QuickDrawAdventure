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
	player        *entities.Player
	enemies       []*entities.Enemy
	items         []*entities.Item
	tilemapJSON   *TilemapJSON
	tilemapeimage *ebiten.Image
	cam           *Camera
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

	g.cam.FollowTarget(g.player.X+16, g.player.Y+16, 640, 480)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		640,
		480,
	)

	return nil
}

// drawing screen + sprites
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}
	//map
	//loop through the tilemap

	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {

			//coordinates example 1%30=1 1/30=0 2%30=2 2/30 = 0 etc...
			x := index % layer.Width
			y := index / layer.Width

			//pixel position
			x *= 16
			y *= 16

			//tile location in asset image we subtract one becuase of json index
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			//pixel position of tile(each tile is a 16x16 square)
			srcX *= 16
			srcY *= 16

			//placing in correct coordinates
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.cam.X, g.cam.Y)

			screen.DrawImage(
				// beginning of tile = srcX, src Y
				//end of tile = srcX + 16, srcY +16
				g.tilemapeimage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	// player position variable

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

	// draw player
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
	ebiten.SetWindowSize(640, 480)
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

	tilemapeimage, _, err := ebitenutil.NewImageFromFile("assets/images/map/TilesetFloor.png")
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/images/map/level1.json")
	if err != nil {
		//handle error
		log.Fatal(err)
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
				Sprite: &entities.Sprite{Img: fishImg,
					X: 335.0,
					Y: 335.0,
				},
				AmtHeal: rand.Intn(4),
				Ifinv:   false,
			},
		},
		tilemapJSON:   tilemapJSON,
		tilemapeimage: tilemapeimage,
		cam:           NewCamera(0.0, 0.0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
