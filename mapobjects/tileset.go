package mapobjects

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

// every tileset must be able to give an image given an id
type Tileset interface {
	Img(id int) *ebiten.Image
	Gid() int
}

// the tileset data deserialized from a standard, single-image tileset
type UniformTilesetJSON struct {
	Path  string `json:"image"`
	Width int    `json:"columns"`
}

// struct for storing uniform tile sets ie 16 x 16 ground tiles
type UniformTileset struct {
	img          *ebiten.Image
	tilesetWidth int
	gid          int
}

func (u *UniformTileset) Gid() int {
	return u.gid
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
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

// struct for tiles or objects of different sizes like buildings or fauna
type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Gid() int {
	return d.gid
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
	//read file contents
	contents, err := os.ReadFile(path)
	fmt.Println(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var dynTilesetJSON DynTilesetJSON
	err = json.Unmarshal(contents, &dynTilesetJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dyn tileset JSON: %w", err)
	}

	if len(dynTilesetJSON.Tiles) > 0 {
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

	fmt.Printf("Loading uniform tileset image from: %s %d\n", tileJSONpath, gid)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONpath)

	if err != nil {
		return nil, err
	}
	uniformTileset.img = img
	uniformTileset.gid = gid
	uniformTileset.tilesetWidth = uniformTilesetJSON.Width

	return &uniformTileset, nil
}

func DetermineTileSet(tiles []int, tilesetgids []int) int {

	maxid := slices.Max(tiles)
	tileindex := 0

	if maxid >= tilesetgids[len(tilesetgids)-1] {
		//id is greater then last gid make tileindex highest index
		tileindex = len(tilesetgids) - 1
		fmt.Printf("first print: maxid:%d tileindex:%d\n", maxid, tileindex)
	} else if tilesetgids[tileindex] <= maxid && tilesetgids[tileindex+1] > maxid {
		fmt.Printf("secondprintp1: maxid:%d tilesetindex+1gid:%d\n", maxid, tilesetgids[tileindex+1])
		//goldilocks zone just return 0
		fmt.Printf("secondprintp2: maxid:%d tileindex:%d\n", maxid, tileindex)
	} else {
		//tileindex gid is too small loop through and increase tileindex until max tile is smaller then next tileset gid
		fmt.Printf("last print p1: maxid:%d tileindex+1gid:%d %t\n", maxid, tilesetgids[tileindex+1], tilesetgids[tileindex+1] > maxid)
		for tilesetgids[tileindex+1] < maxid+1 {
			fmt.Printf("last print tileindex:%d\n", tileindex)
			tileindex += 1
		}
		fmt.Printf("last print p2: maxid:%d tileindex+1gid:%d\n", maxid, tilesetgids[tileindex+1])
	}

	return tileindex

}
