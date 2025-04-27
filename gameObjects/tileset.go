package gameObjects

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"image"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TileType uint8

const (
	Stair TileType = iota
	Surface
	Edge
	NotDefined
)

type Tileset interface {
	Img(id int) *ebiten.Image
	Gid() int
	Name() string
	Dimensions() (int, int)
	TileType(int) TileType
}

// UniformTilesetJSON the tileset data deserialized from a standard, single-image tileset
type UniformTilesetJSON struct {
	Path       string           `json:"image"`
	Width      int              `json:"columns"`
	Name       string           `json:"name"`
	Properties []PropertiesJSON `json:"properties"`
	TileProps  []TileProp       `json:"Tiles"`
}

// UniformTileset struct for storing uniform tile sets ie 16 x 16 ground tiles
type UniformTileset struct {
	img            *ebiten.Image
	tilesetWidth   int
	gid            int
	name           string
	TileProperties map[int]TileType
}

func MapTileType(tileProp string) (TileType, error) {
	switch tileProp {
	case "stair":
		return Stair, nil
	case "surface":
		return Surface, nil
	case "edge":
		return Edge, nil
	default:
		return NotDefined, errors.New("invalid string in tile type json value in tile set for: " + tileProp)
	}
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

func (u *UniformTileset) TileType(id int) TileType {
	println("tile tipe for id", id, "=", u.TileProperties[id-1])
	return u.TileProperties[id-1]
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

type TileProp struct {
	TileID int            `json:"id"`
	Prop   []TileProperty `json:"properties"`
}

type TileProperty struct {
	Name         string `json:"name"`
	PropertyType string `json:"propertytype"`
	Type         string `json:"type"`
	Value        string `json:"value"`
}

type TileJSON struct {
	Id        int      `json:"id"`
	Animation []string `json:"animation"`
	Path      string   `json:"image"`
	Width     int      `json:"imagewidth"`
	Height    int      `json:"imageheight"`
	Name      string   `json:"name"`
}

type AnimatedTileSet struct {
	Tiles            []*TileJSON `json:"tiles"`
	Name             string      `json:"name"`
	TileHeight       int         `json:"tileheight"`
	TileWidth        int         `json:"tilewidth"`
	Animation        animations.Animation
	OffsetX, OffsetY int
}

type DynTileSetJSON struct {
	Tiles      []*TileJSON      `json:"tiles"`
	Name       string           `json:"name"`
	TileHeight int              `json:"tileheight"`
	TileWidth  int              `json:"tilewidth"`
	Properties []PropertiesJSON `json:"properties"`
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

func (d *DynTileset) TileType(id int) TileType {
	return NotDefined
}

func (d *DynTileset) Name() string {
	return d.name
}

func (d *DynTileset) Img(id int) *ebiten.Image {

	id -= d.gid
	if id >= len(d.imgs) {
		println(d.name, "has error in id values for tiles")
	}
	img := d.imgs[id]

	if img == nil {
		log.Fatal("Error: img is nil")
	}
	return img
}

func NewDynamicTileSet(path string, gid int) (Tileset, error) {
	contents, err := assets.Map.ReadFile(path)
	fmt.Println(path)
	if err != nil {
		return nil, fmt.Errorf("failed to interactions file %s: %w", path, err)
	}
	var dynTileSetJSON DynTileSetJSON
	err = json.Unmarshal(contents, &dynTileSetJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dyn tileset JSON: %w", err)
	}

	//create the tileset
	dynTileset := DynTileset{}
	dynTileset.name = dynTileSetJSON.Name
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

		fmt.Printf("Loading dyn tileset image from: %s\n", tileJsonPath)

		img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, tileJsonPath)

		if err != nil {
			return nil, fmt.Errorf("failed to load dyntileset image from %s: %w", tileJsonPath, err)
		}

		dynTileset.imgs = append(dynTileset.imgs, img)

	}
	dynTileset.width = dynTileSetJSON.TileWidth
	dynTileset.height = dynTileSetJSON.TileHeight

	println(dynTileset.name)

	return &dynTileset, nil
}

func NewUniformTileSet(path string, gid int) (Tileset, error) {
	contents, err := assets.Map.ReadFile(path)
	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)

	if err != nil {
		return nil, err
	}

	uniformTileSet := UniformTileset{}

	//clean and convert tileset relative path to root relative path
	tileJsonPath := uniformTilesetJSON.Path
	tileJsonPath = filepath.Clean(tileJsonPath)
	tileJsonPath = strings.ReplaceAll(tileJsonPath, "\\", "/")
	tileJsonPath = strings.TrimPrefix(tileJsonPath, "../")
	tileJsonPath = strings.TrimPrefix(tileJsonPath, "../")
	tileJsonPath = filepath.Clean(tileJsonPath)

	fmt.Printf("Loading uniform tileset image from: %s %d\n", tileJsonPath, gid)

	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, tileJsonPath)
	if err != nil {
		return nil, err
	}

	uniformTileSet.img = img
	uniformTileSet.gid = gid
	uniformTileSet.name = uniformTilesetJSON.Name
	uniformTileSet.tilesetWidth = uniformTilesetJSON.Width

	propMap := make(map[int]TileType)
	for _, props := range uniformTilesetJSON.TileProps {
		prop := props.Prop[0]
		println("loading tile type", props.TileID, prop.Value)
		if props.TileID != 0 {
			tTypeIota, err := MapTileType(prop.Value)
			if err != nil {
				log.Fatal(err)
			}
			propMap[props.TileID] = tTypeIota

		}
	}
	uniformTileSet.TileProperties = propMap

	return &uniformTileSet, nil
}

func NewTileSet(path string, gid int) (Tileset, error) {
	//interactions file contents
	contents, err := assets.Map.ReadFile(path)
	fmt.Println(path)
	if err != nil {
		return nil, fmt.Errorf("failed to interactions file %s: %w", path, err)
	}

	var checkType map[string]interface{}

	err = json.Unmarshal(contents, &checkType)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dyn tileset JSON: %w", err)
	}

	properties, ok := checkType["properties"].([]interface{})
	if !ok {
		log.Fatal()
	}
	propValues, ok := properties[0].(map[string]interface{})
	if !ok {
		log.Fatal()
	}

	switch propValues["value"] {

	case "animated":
		println("skipping animated tile set")
		return nil, nil
	case "dynamic":
		tileSet, err := NewDynamicTileSet(path, gid)
		if err != nil {
			log.Fatal(err)
		}
		return tileSet, nil
	case "uniform":
		tileSet, err := NewUniformTileSet(path, gid)
		if err != nil {
			log.Fatal(err)
		}
		return tileSet, nil
	}

	return nil, nil
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
