package mapobjects

import (
	"encoding/json"
	"os"
	"path"
)

type TilemapLayerJSON struct {
	Data    []int        `json:"data"`
	Width   int          `json:"width"`
	Height  int          `json:"height"`
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Objects []ObjectJSON `json:"objects"`
	Class   string       `json:"class"`
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	// raw data for each tileset
	Tilesets []map[string]any `json:"tilesets"`
}

func (t *TilemapJSON) GenTileSets() ([]Tileset, error) {

	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		tilesetpath := path.Join("assets/map/", tilesetData["source"].(string))
		tileset, err := NewTileSet(tilesetpath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}
		tilesets = append(tilesets, tileset)
	}
	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
