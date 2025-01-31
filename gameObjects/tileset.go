package gameObjects

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Tileset every tileset must be able to give an image given an id
type Tileset interface {
	Img(id int) *ebiten.Image
	Gid() int
	Name() string
	Dimensions() (int, int)
}

// UniformTilesetJSON the tileset data deserialized from a standard, single-image tileset
type UniformTilesetJSON struct {
	Path  string `json:"image"`
	Width int    `json:"columns"`
	Name  string `json:"name"`
}

// UniformTileset struct for storing uniform tile sets ie 16 x 16 ground tiles
type UniformTileset struct {
	img          *ebiten.Image
	tilesetWidth int
	gid          int
	name         string
}

func (u *UniformTileset) Gid() int {
	return u.gid
}

func (u *UniformTileset) Dimensions() (int, int) {
	//need to figure out this calc for height
	return 0, 0
}

func (u *UniformTileset) Name() string {
	return u.name
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	// gets right sprite data based on starting point of tile set
	id -= u.gid
	srcX := id % u.tilesetWidth
	srcY := id / u.tilesetWidth
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
	Name   string `json:"name"`
}

type DynTilesetJSON struct {
	Tiles      []*TileJSON `json:"tiles"`
	Name       string      `json:"name"`
	TileHeight int         `json:"tileheight"`
	TileWidth  int         `json:"tilewidth"`
}

// DynTileset struct for tiles or objects of different sizes like buildings or fauna
type DynTileset struct {
	imgs   []*ebiten.Image
	gid    int
	name   string
	width  int
	height int
}

func (d *DynTileset) Dimensions() (int, int) {
	return d.width, d.height
}

func (d *DynTileset) Gid() int {
	return d.gid
}

func (d *DynTileset) Name() string {
	return d.name
}

func (d *DynTileset) Img(id int) *ebiten.Image {

	id -= d.gid
	img := d.imgs[id]

	if img == nil {
		log.Fatal("Error: img is nil")
	}
	return img
}

func NewTileSet(path string, gid int) (Tileset, error) {
	//interactions file contents
	contents, err := os.ReadFile(path)
	fmt.Println(path)
	if err != nil {
		return nil, fmt.Errorf("failed to interactions file %s: %w", path, err)
	}

	var checkDynTileSetJSON DynTilesetJSON
	err = json.Unmarshal(contents, &checkDynTileSetJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dyn tileset JSON: %w", err)
	}

	if len(checkDynTileSetJSON.Tiles) > 0 {
		//return dyn tileset
		var dynTileSetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTileSetJSON)
		if err != nil {

			return nil, fmt.Errorf("failed to unmarshal dyn tileset JSON: %w", err)
		}

		//create the tileset
		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)
		//change back to ebiten image
		//loop over tile data and load image for each
		for _, tileJSON := range dynTileSetJSON.Tiles {

			// clean and convert tileset relative path to root relative path
			tileJsonPath := tileJSON.Path
			tileJsonPath = filepath.Clean(tileJsonPath)
			tileJsonPath = strings.ReplaceAll(tileJsonPath, "\\", "/")
			tileJsonPath = strings.TrimPrefix(tileJsonPath, "../")
			tileJsonPath = strings.TrimPrefix(tileJsonPath, "../")
			tileJsonPath = filepath.Clean(tileJsonPath)
			tileJsonPath = filepath.Join("assets", tileJsonPath)

			fmt.Printf("Loading dyn tileset image from: %s\n", tileJsonPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJsonPath)

			if err != nil {
				return nil, fmt.Errorf("failed to load dyntileset image from %s: %w", tileJsonPath, err)
			}

			dynTileset.imgs = append(dynTileset.imgs, img)

		}
		dynTileset.width = dynTileSetJSON.TileWidth
		dynTileset.height = dynTileSetJSON.TileHeight
		return &dynTileset, nil
	}
	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)

	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}

	//clean and convert tileset relative path to root relative path
	tileJsonPath := uniformTilesetJSON.Path
	tileJsonPath = filepath.Clean(tileJsonPath)
	tileJsonPath = strings.ReplaceAll(tileJsonPath, "\\", "/")
	tileJsonPath = strings.TrimPrefix(tileJsonPath, "../")
	tileJsonPath = strings.TrimPrefix(tileJsonPath, "../")
	tileJsonPath = filepath.Clean(tileJsonPath)
	tileJsonPath = filepath.Join("assets", tileJsonPath)

	fmt.Printf("Loading uniform tileset image from: %s %d\n", tileJsonPath, gid)

	img, _, err := ebitenutil.NewImageFromFile(tileJsonPath)

	if err != nil {
		return nil, err
	}
	uniformTileset.img = img
	uniformTileset.gid = gid
	uniformTileset.name = uniformTilesetJSON.Name
	uniformTileset.tilesetWidth = uniformTilesetJSON.Width

	return &uniformTileset, nil
}

func DetermineTileSet(tiles []int, tilesetgids []int) int {

	maxid := slices.Max(tiles)
	tileindex := 0

	if maxid >= tilesetgids[len(tilesetgids)-1] {
		tileindex = len(tilesetgids) - 1
	} else if tilesetgids[tileindex] <= maxid && tilesetgids[tileindex+1] > maxid {
	} else {
		for tilesetgids[tileindex+1] < maxid+1 {
			tileindex += 1
		}
	}

	return tileindex

}
