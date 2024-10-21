package mapobjects

import "image"

type Door struct {
	Key   string
	Coord image.Rectangle
}

func NewDoor(obj ObjectJSON) (door Door) {
	door =
		Door{
			Key: obj.Name,
			Coord: image.Rect(
				int(obj.X),
				int(obj.Y)-10,
				int(obj.Width+obj.X),
				int(obj.Y)-20),
		}

	return door
}
