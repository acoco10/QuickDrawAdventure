package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// every tileset must be able to give an image given an id
type Tileset interface {
	Img(id int) *ebiten.Image
}

// the tileset data deserialized from a standard, single-image tileset
type UniformTilesetJSON struct {
	Path string `json:"image"`
}

// struct for storing uniform tile sets ie 16 x 16 ground tiles
type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	// gets right sprite data based on starting point of tile set
	id -= u.gid
	srcX := id % 22
	srcY := id / 22
	//pixel position of tile(each tile is a 16x16 square)
	srcX *= 16
	srcY *= 16

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+16, srcY+16,
		),
	).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

// struct for tiles or objects of different sizes like buildings or fauna
type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {

	log.Printf("Before subtraction: id = %d, gid = %d", id, d.gid)
	id -= d.gid
	log.Printf("After subtraction: id = %d", id)
	img := d.imgs[id]

	if img == nil {
		log.Fatal("Error: img is nil")
	}
	return img
}

func NewTileSet(path string, gid int) (Tileset, error) {
	//read file contents
	contents, err := os.ReadFile(path)
	fmt.Println(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if strings.Contains(path, "Structures") || strings.Contains(path, "nature") {
		//return dyn tileset
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dyn tileset JSON: %w", err)
		}

		//create the tileset
		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0) //change back to ebiten image
		//loop over tile data and load image for each
		for _, tileJSON := range dynTilesetJSON.Tiles {

			// clean and convert tileset relative path to root relative path
			tileJSONpath := tileJSON.Path
			tileJSONpath = filepath.Clean(tileJSONpath)
			tileJSONpath = strings.ReplaceAll(tileJSONpath, "\\", "/")
			tileJSONpath = strings.TrimPrefix(tileJSONpath, "../")
			tileJSONpath = strings.TrimPrefix(tileJSONpath, "../")
			tileJSONpath = filepath.Clean(tileJSONpath)
			tileJSONpath = filepath.Join("assets", tileJSONpath)

			fmt.Printf("Loading dyn tileset image from: %s\n", tileJSONpath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONpath)

			if err != nil {
				return nil, err
			}

			dynTileset.imgs = append(dynTileset.imgs, img)

		}
		return &dynTileset, nil
	}
	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)

	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	//clean and convert tileset relative path to root relative path
	tileJSONpath := uniformTilesetJSON.Path
	tileJSONpath = filepath.Clean(tileJSONpath)
	tileJSONpath = strings.ReplaceAll(tileJSONpath, "\\", "/")
	tileJSONpath = strings.TrimPrefix(tileJSONpath, "../")
	tileJSONpath = strings.TrimPrefix(tileJSONpath, "../")
	tileJSONpath = filepath.Clean(tileJSONpath)
	tileJSONpath = filepath.Join("assets", tileJSONpath)

	fmt.Printf("Loading uniform tileset image from: %s\n", tileJSONpath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONpath)

	if err != nil {
		return nil, err
	}
	uniformTileset.img = img
	uniformTileset.gid = gid

	return &uniformTileset, nil
}
