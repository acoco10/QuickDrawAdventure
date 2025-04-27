package gameObjects

import (
	"encoding/json"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"path"
)

type TileMapLayer struct {
	Data             []int `json:"data"`
	TileSet          Tileset
	Tiles            []Tile
	HorizontalOffset int              `json:"offsetx"`
	VerticalOffset   int              `json:"offsety"`
	Width            int              `json:"width"`
	Height           int              `json:"height"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	Objects          []ObjectJSON     `json:"objects"`
	Class            string           `json:"class"`
	Properties       []PropertiesJSON `json:"properties"`
	Z                float64
	YSort            bool
	Layers           []*TileMapLayer `json:"layers"`
}
type TileKey struct {
	X, Y int
}
type PropertiesJSON struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// TilemapJSON all layers in a tilemap
type TilemapJSON struct {
	Layers []*TileMapLayer `json:"layers"`
	// raw data for each tileset
	Tilesets []map[string]any `json:"tilesets"`
	RawData  json.RawMessage  `json:"-"`
}

type Tile struct {
	Img      *ebiten.Image
	X, Y, Z  float64
	YSort    bool
	Layer    string
	TileType TileType
}

func (t Tile) GetType() DrawableType {
	return Map
}

func (t Tile) GetSize() (h int, w int) {
	return t.Img.Bounds().Dx(), t.Img.Bounds().Dy()
}

func (t Tile) CheckName() string {
	return t.Layer
}

func (t Tile) GetCoord() (x, y, z float64) {
	return t.X, t.Y, t.Z
}

func (t Tile) CheckYSort() bool {
	return t.YSort
}

func (t Tile) Draw(screen *ebiten.Image, cam camera.Camera, player Character, debugMode bool) {
	opts := ebiten.DrawImageOptions{}
	if debugMode {
		switch t.Z {
		case 2:
			opts.ColorScale.Scale(0.7, 0.8, 0.7, 1)
		case 3:
			opts.ColorScale.Scale(0.8, 0.5, 0.6, 1)
		case 4:
			opts.ColorScale.Scale(0.6, 0.5, 0.8, 1)
		}
	}
	opts.GeoM.Translate(t.X, t.Y)
	opts.GeoM.Translate(0.0, -(float64(t.Img.Bounds().Dy()) + 16))
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Scale(4, 4)
	screen.DrawImage(t.Img, &opts)
	/*if t.YSort {
	#very non performant debug mode
		face, err := assetManagement.LoadFont(12, assetManagement.November)
		if err != nil {
			log.Fatal()
		}
		dopts := text.DrawOptions{}
		x := t.X*4 + cam.X*4
		y := t.Y*4 + cam.Y*4 - 64
		dopts.GeoM.Translate(x, y)
		coord := fmt.Sprintf("x = %f y = %f", t.X, t.Y)
		text.Draw(screen, coord, face, &dopts)
	}*/

}

func (t *TilemapJSON) GenTileSets() ([]Tileset, error) {

	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		tilesetpath := path.Join("map/", tilesetData["source"].(string))
		tileset, err := NewTileSet(tilesetpath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}
		if tileset != nil {
			tilesets = append(tilesets, tileset)
		}
	}
	return tilesets, nil
}

func (t *TilemapJSON) GenTileSetMap() (map[string]Tileset, error) {

	tilesets := make(map[string]Tileset)

	for _, tilesetData := range t.Tilesets {
		tilesetpath := path.Join("map/", tilesetData["source"].(string))
		tileset, err := NewTileSet(tilesetpath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}
		if tileset != nil {
			tilesets[tileset.Name()] = tileset
		}
	}
	return tilesets, nil
}

func FlattenLayers(layers []*TileMapLayer) []*TileMapLayer {
	var result []*TileMapLayer
	for _, layer := range layers {
		if len(layer.Layers) > 0 {
			// it's a group layer
			result = append(result, FlattenLayers(layer.Layers)...)
		} else {
			// it's a normal tile/object/image layer
			result = append(result, layer)
		}
	}
	return result
}

func NewTilemapJSON(contents []byte) (*TilemapJSON, error) {

	var tilemapJSON TilemapJSON
	err := json.Unmarshal(contents, &tilemapJSON)

	if err != nil {
		return nil, err
	}
	tilemapJSON.Layers = FlattenLayers(tilemapJSON.Layers)

	tilesets, err := tilemapJSON.GenTileSetMap()
	if err != nil {
		log.Fatal(err)
	}
	for _, layer := range tilemapJSON.Layers {
		if layer.Class == "layer" {
			println("layer =", layer.Name)
			if len(layer.Properties) > 1 {
				layer.Z = layer.Properties[2].Value.(float64)
			} else {
				layer.Z = 1
			}

			println("z value for", layer.Name, "=", layer.Z)

			layer.YSort = layer.Properties[1].Value.(bool)

			println("y sort value for", layer.Name, "=", layer.YSort)

			tileSetName := layer.Properties[0].Value.(string)
			layer.TileSet = tilesets[tileSetName]

			println("layer:", layer.Name, "tileSet =", layer.TileSet.Name())

			var tiles []Tile
			for index, id := range layer.Data {

				if id == 0 {
					continue
				}

				x := index % layer.Width
				y := index / layer.Width

				x = x * 16
				y = (y + 2) * 16

				x = x + layer.HorizontalOffset
				y = y + layer.VerticalOffset

				tile := Tile{
					X:        float64(x),
					Y:        float64(y),
					Z:        layer.Z,
					Img:      layer.TileSet.Img(id),
					YSort:    layer.YSort,
					Layer:    layer.Name,
					TileType: layer.TileSet.TileType(id),
				}

				tiles = append(tiles, tile)

			}
			layer.Tiles = tiles
			println("len of tiles for layer", layer.Name, "=", len(layer.Tiles))

		}
	}
	return &tilemapJSON, nil
}
