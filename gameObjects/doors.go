package gameObjects

import (
	"image"
)

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
				int(obj.Y)-20,
				int(obj.Width+obj.X),
				int(obj.Y)-10,
			),
		}

	return door
}

func CheckEntDoor(player *Character, entdoors map[string]Door, exdoors map[string]Door) map[string]bool {
	returnMap := make(map[string]bool)
	for _, door := range entdoors {
		if door.Coord.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			player.EnterShadow()
			returnMap[door.Key] = true
		} else {
			returnMap[door.Key] = false
		}
	}
	return returnMap
}

func GetDoorCoord(doors map[string]Door, key string, direction string) (float64, float64) {
	x := 0.0
	y := 0.0
	if direction == "up" {
		x = float64(doors[key].Coord.Dx()/2+doors[key].Coord.Min.X) - 8
		y = float64(doors[key].Coord.Min.Y) - 32
	}
	if direction == "down" {
		x = float64(doors[key].Coord.Dx()/2+doors[key].Coord.Min.X) - 8
		y = float64(doors[key].Coord.Max.Y) - 16
	}

	return x, y

}
