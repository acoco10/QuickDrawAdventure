package gameObjects

import (
	"image"
)

type ObjectJSON struct {
	Name   string  `json:"name"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	X      float64 `json:"x"`
	Y      float64 `json:"Y"`
	Class  string  `json:"class"`
	Type   string  `json:"type"`
}

type NPCspawn struct {
	Name string
	X, Y float64
}

func NewNPCspawn(json ObjectJSON) NPCspawn {
	println("loading npc spawn:", json.Name, "\n")
	npcspawn := new(NPCspawn)
	npcspawn.Name = json.Name
	npcspawn.X = json.X
	npcspawn.Y = json.Y
	return *npcspawn
}

func StoreMapObjects(tilemapJSON TilemapJSON) (colliders []image.Rectangle, entDoors map[string]Door, exDoors map[string]Door, NPCSpawns map[string]NPCspawn) {
	entDoors = make(map[string]Door)
	exDoors = make(map[string]Door)
	NPCSpawns = make(map[string]NPCspawn)
	for _, layer := range tilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			for _, object := range layer.Objects {
				if object.Type == "RectCollider" {
					img := image.Rect(
						int(object.X),
						int(object.Y)-32,
						int(object.Width+object.X),
						int(object.Y+object.Height)-32,
					)
					colliders = append(colliders, img)
				}
				if object.Type == "entranceDoor" {
					door := NewDoor(object)
					entDoors[object.Name] = door
				}
				if object.Type == "exitDoor" {
					door := NewDoor(object)
					exDoors[object.Name] = door

				}
				if object.Type == "NPCSpawn" {
					npcspawn := NewNPCspawn(object)
					NPCSpawns[object.Name] = npcspawn

				}
			}
		}
	}

	return colliders, entDoors, exDoors, NPCSpawns
}
