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

func CheckEntDoor(player *Character, entdoors map[string]Door, exdoors map[string]Door) bool {

	for _, door := range entdoors {
		if door.Coord.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			player.EnterShadow()
			return true
		}
	}
	return false
}

func CheckExDoor(player *Character, entdoors map[string]Door, exdoors map[string]Door) bool {
	for _, door := range exdoors {
		if door.Coord.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			//player.X = float64(entdoors[key].Coord.Dx()/2+entdoors[key].Coord.Min.X) - 8
			//player.Y = float64(entdoors[key].Coord.Min.Y + 20)

			player.ExitShadow()
			return true
		}
	}
	return false
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
