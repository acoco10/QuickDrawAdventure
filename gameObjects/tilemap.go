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
	Z          float64
	YSort      bool
}

type PropertiesJSON struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	// raw data for each tileset
	Tilesets []map[string]any `json:"tilesets"`
	RawData  json.RawMessage  `json:"-"`
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

func NewTilemapJSON(contents []byte) (*TilemapJSON, error) {

	var tilemapJSON TilemapJSON
	err := json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}
	for _, layer := range tilemapJSON.Layers {
		if layer.Class == "layer" {
			println("layer =", layer.Name)
			if len(layer.Properties) > 1 {
				layer.Z = layer.Properties[1].Value.(float64)
			} else {
				layer.Z = 1
			}

			println("z value for", layer.Name, "=", layer.Z)

			layer.YSort = layer.Properties[0].Value.(bool)

			println("y sort value for", layer.Name, "=", layer.YSort)
		}
	}
	return &tilemapJSON, nil
}
