package mapobjects

import "image"

type ObjectJSON struct {
	Name   string  `json:"name"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Class  string  `json:"class"`
	Type   string  `json:"type"`
}

func StoreMapObjects(tilemapJSON TilemapJSON) (colliders []image.Rectangle, entDoors map[string]Door, exDoors map[string]Door) {

	colliders = []image.Rectangle{}
	entDoors = make(map[string]Door)
	exDoors = make(map[string]Door)

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
			}
		}
	}

	return colliders, entDoors, exDoors
}
