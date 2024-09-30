package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

type Player struct {
	*Sprite
	Health uint
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Item struct {
	*Sprite
	AmtHeal int
	Ifinv   bool
}

type Game struct {
	player        *Player
	enemies       []*Enemy
	items         []*Item
	tilemapJSON   *TilemapJSON
	tilemapeimage *ebiten.Image
}

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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	//loop through the tilemap

	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			//tile position
			fmt.Printf("index %v\n", index)
			x := index % layer.Width
			y := index / layer.Width

			fmt.Printf("tile position %v %v\n", x, y)

			//pixel position
			x *= 16
			y *= 16

			fmt.Printf("tile position %v %v\n", x, y)

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapeimage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	// player position variable

	opts.GeoM.Translate(g.player.X, g.player.Y)

	// draw player
	screen.DrawImage(
		//grab a subimage of the Spritesheet
		g.player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	// draw enemy sprites
	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)

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

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/ToughGuy.png")
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

	tilemapJSON, err := NewTilemapJSON("assets/images/map/startermap.json")
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	game := Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   50.0,
				Y:   50.0,
			},
			Health: 10,
		},
		enemies: []*Enemy{
			{
				&Sprite{Img: ghostImg,
					X: 100.0,
					Y: 100.0,
				},
				true,
			},
			{
				&Sprite{Img: ghostImg,
					X: 50.0,
					Y: 50.0,
				},
				false,
			},
			{
				&Sprite{Img: ghostImg,
					X: 100.0,
					Y: 100.0,
				},
				false,
			},
		},
		items: []*Item{
			{
				&Sprite{Img: fishImg,
					X: 335.0,
					Y: 335.0,
				},
				rand.Intn(4),
				false,
			},
		},
		tilemapJSON:   tilemapJSON,
		tilemapeimage: tilemapeimage,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
