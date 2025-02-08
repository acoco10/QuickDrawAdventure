package gameObjects

import (
	"encoding/json"
	"path"
)

// go: embed
type TilemapLayerJSON struct {
	Data       []int            `json:"data"`
	Width      int              `json:"width"`
	Height     int              `json:"height"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Objects    []ObjectJSON     `json:"objects"`
	Class      string           `json:"class"`
	Properties []PropertiesJSON `json:"properties"`
}

type PropertiesJSON struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
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
		tilesetpath := path.Join("map/", tilesetData["source"].(string))
		tileset, err := NewTileSet(tilesetpath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}
		tilesets = append(tilesets, tileset)
	}
	return tilesets, nil
}

func NewTilemapJSON(contents []byte) (*TilemapJSON, error) {

	var tilemapJSON TilemapJSON
	err := json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
